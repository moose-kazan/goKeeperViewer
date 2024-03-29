package main

import (
	"goKeeperViewer/internal/fynefilechooser"
	"goKeeperViewer/internal/fynetheme"
	"goKeeperViewer/internal/kdb"
	"goKeeperViewer/internal/settings"
	"os"
	"path/filepath"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

var (
	a               fyne.App
	w               fyne.Window
	passwordTree    *widget.Tree
	passwordDetails *widget.Form
	db              *kdb.KDB
)

func loadFile(fileName fyne.URI) {
	pwdEntry := widget.NewPasswordEntry()
	keyFileChooser := fynefilechooser.NewFileChooser(w, storage.NewExtensionFileFilter([]string{".keyx", ".key"}))
	d := dialog.NewForm(
		"Enter password",
		"OK",
		"Cancel",
		[]*widget.FormItem{
			widget.NewFormItem("File Name", widget.NewLabel(filepath.Base(fileName.Path()))),
			widget.NewFormItem("Password", pwdEntry),
			widget.NewFormItem("Key File", keyFileChooser),
		},
		func(b bool) {
			if !b {
				return
			}

			tmpDb := kdb.New()
			err := tmpDb.Load(fileName, pwdEntry.Text, keyFileChooser.GetURI())

			if err != nil {
				dialog.NewError(err, w).Show()
				return
			}

			_ = tmpDb.Tree()

			db = tmpDb
			//db.SetDebug(true)
			passwordTree.Root = "/"
			passwordTree.Refresh()

			settings.New(a.Preferences()).SetLastFile(fileName.String())
		},
		w,
	)
	d.Show()
}

func buildPasswordDetails() *widget.Form {
	formUserName := widget.NewEntry()
	formUserName.Disable()
	formNotes := widget.NewEntry()
	formNotes.Disable()
	formNotes.MultiLine = true
	formPassword := widget.NewPasswordEntry()
	formPassword.Password = true
	formPassword.Disable()

	passwordDetails = widget.NewForm(
		widget.NewFormItem("Title", widget.NewLabel("")),
		widget.NewFormItem("URL", widget.NewHyperlink("", nil)),
		widget.NewFormItem("UserName", formUserName),
		widget.NewFormItem("Password", formPassword),
		widget.NewFormItem("Notes", formNotes),
	)
	return passwordDetails
}

func main() {
	os.Setenv("FYNE_THEME", "light")
	a = app.NewWithID("goKeeperViewer")
	a.Settings().SetTheme(fynetheme.New())
	w = a.NewWindow("goKeeperViewer")
	w.Resize(fyne.NewSize(640, 480))
	w.SetCloseIntercept(actionWindowClose)

	w.SetMainMenu(BuildMenu())

	toolbar := BuildToolbar()

	passwordTree = widget.NewTree(
		func(s string) []string {
			return db.GetChildIDs(s)
		},
		func(s string) bool {
			if db == nil {
				return false
			}
			d := db.IsBranch(s)
			return d
		},
		func(b bool) fyne.CanvasObject {
			if b {
				return widget.NewLabel("")
			}
			return widget.NewLabel("")
		},
		func(s string, b bool, co fyne.CanvasObject) {
			item := db.GetItemByID(s)
			co.(*widget.Label).SetText(item.Title)
		},
	)
	passwordTree.OnSelected = func(uid widget.TreeNodeID) {
		item := db.GetItemByID(uid)
		if item.Entry == nil {
			return
		}
		passwordDetails.Show()
		// TODO: Process all Entry fields dinamicaly
		// TODO: Do something with internal binaries, like ssh-keys
		for _, v := range passwordDetails.Items {
			if v.Text == "Title" {
				v.Widget.(*widget.Label).SetText(item.Entry.GetTitle())
			} else if v.Text == "Password" {
				v.Widget.(*widget.Entry).SetText(item.Entry.GetPassword())
			} else if v.Text == "URL" {
				v.Widget.(*widget.Hyperlink).SetURLFromString(item.Entry.GetContent("URL"))
				v.Widget.(*widget.Hyperlink).SetText(item.Entry.GetContent("URL"))
			} else if v.Text == "UserName" {
				v.Widget.(*widget.Entry).SetText(item.Entry.GetContent("UserName"))
			} else if v.Text == "Notes" {
				v.Widget.(*widget.Entry).SetText(item.Entry.GetContent("Notes"))
			}
		}
	}

	passwordDetails = buildPasswordDetails()
	passwordDetails.Hide()

	content := container.NewBorder(
		toolbar,
		nil,
		nil,
		nil,
		container.NewGridWithColumns(
			2,
			passwordTree,
			passwordDetails,
		),
	)

	w.SetContent(content)

	if len(os.Args) > 1 {
		loadFile(storage.NewFileURI(os.Args[1]))
	} else if settings.New(a.Preferences()).GetStartLoadOption() == settings.START_LOAD_LAST {
		var fileName = settings.New(a.Preferences()).GetLastFile()
		if fileName != "" {
			loadFile(storage.NewURI(fileName))
		}
	}

	w.ShowAndRun()
}
