package main

import (
	"fmt"
	"goKeeperViewer/internal/settings"
	"net/url"
	"runtime"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/lang"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
)

func actionHelpAbout() {
	dialogTitle := fmt.Sprintf(
		"%s %s",
		a.Metadata().Name,
		a.Metadata().Version,
	)
	urlEmailTitle := "moose@ylsoftware.com"
	urlEmail, _ := url.Parse(fmt.Sprintf("mailto:%s", urlEmailTitle))
	urlWSTitle := "https://github.com/moose-kazan/goKeeperViewer"
	urlWS, _ := url.Parse("https://github.com/moose-kazan/goKeeperViewer")
	aboutLayout := container.NewVBox(
		widget.NewLabelWithStyle(dialogTitle, fyne.TextAlignCenter, fyne.TextStyle{Bold: true}),
		widget.NewForm(
			widget.NewFormItem(lang.L("Author"), widget.NewLabel("Vadim Kalinnikov")),
			widget.NewFormItem(lang.L("E-Mail"), widget.NewHyperlink(urlEmailTitle, urlEmail)),
			widget.NewFormItem(lang.L("Website"), widget.NewHyperlink(urlWSTitle, urlWS)),
			widget.NewFormItem("", widget.NewLabel(lang.L("Simple viewer for KDBX (KeePass) files."))),
			widget.NewFormItem(lang.L("OS"), widget.NewLabel(fmt.Sprintf("%s %s", runtime.GOOS, runtime.GOARCH))),
		),
	)
	d := dialog.NewCustom(
		lang.L("About"),
		lang.L("OK"),
		aboutLayout,
		w,
	)
	d.Show()
}

func actionMenuOpen() {
	d := dialog.NewFileOpen(func(u fyne.URIReadCloser, e error) {
		if e != nil {
			dialog.NewError(e, w).Show()
			return
		}
		if u != nil {
			loadFile(u.URI())
		}

	}, w)
	d.SetFilter(storage.NewExtensionFileFilter([]string{".kdbx"}))
	d.Show()
}

func actionSettings() {
	var startLoadVariants = []string{
		lang.L("None"),
		lang.L("Last File"),
	}

	loadOnStart := widget.NewRadioGroup(
		startLoadVariants,
		func(s string) {

		},
	)
	loadOnStart.Selected = loadOnStart.Options[settings.New(a.Preferences()).GetStartLoadOption()]

	confirmExit := widget.NewCheck("", func(b bool) {})
	confirmExit.SetChecked(settings.New(a.Preferences()).GetConfirmExit())

	dialog.NewForm(
		lang.L("Settings"),
		lang.L("OK"),
		lang.L("Cancel"),
		[]*widget.FormItem{
			widget.NewFormItem(
				lang.L("Load on start"),
				loadOnStart,
			),
			widget.NewFormItem(
				lang.L("Confirm exit"),
				confirmExit,
			),
		},
		func(b bool) {
			if !b {
				return
			}
			s := settings.New(a.Preferences())
			for k, v := range startLoadVariants {
				if v == loadOnStart.Selected {
					s.SetStartLoadOption(k)
					break
				}
			}
			s.SetConfirmExit(confirmExit.Checked)
		},
		w,
	).Show()
}

func actionWindowClose() {
	if settings.New(a.Preferences()).GetConfirmExit() {
		dialog.NewConfirm(
			lang.L("Confirm"),
			lang.L("Are you want to close app?"),
			func(b bool) {
				if b {
					w.Close()
				}
			},
			w,
		).Show()
	} else {
		a.Quit()
	}
}
