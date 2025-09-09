package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/lang"
	"fyne.io/fyne/v2/theme"
)

func newMenuItem(label string, action func(), Icon fyne.Resource, Shortcut fyne.Shortcut) *fyne.MenuItem {
	m := fyne.NewMenuItem(label, action)
	m.Icon = Icon
	m.Shortcut = Shortcut
	return m
}

func BuildMenu() *fyne.MainMenu {
	return fyne.NewMainMenu(
		fyne.NewMenu(
			lang.L("File"),
			newMenuItem(lang.L("Open"), actionMenuOpen, theme.DocumentIcon(), nil),
			newMenuItem(lang.L("Quit"), actionWindowClose, theme.LogoutIcon(), nil),
		),
		fyne.NewMenu(
			lang.L("Settings"),
			newMenuItem(lang.L("Settings"), actionSettings, theme.SettingsIcon(), nil),
		),
		fyne.NewMenu(
			lang.L("Help"),
			newMenuItem(lang.L("About"), actionHelpAbout, theme.InfoIcon(), nil),
		),
	)
}
