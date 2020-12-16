package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/dialog"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/theme"
	"fyne.io/fyne/widget"
)

type ui struct {
	mainWin fyne.Window
	app     fyne.App
	entry   desktopEntry

	nameEntry       *widget.Entry
	categoriesEntry *widget.Entry
	typeEntry       *widget.Entry
	commandEntry    *widget.Entry
	terminalSelect  *widget.Select
	iconPathLabel   *widget.Label
}

func (u *ui) loadMainUI() fyne.CanvasObject {
	u.entry = desktopEntry{}

	u.nameEntry = widget.NewEntry()
	u.iconPathLabel = widget.NewLabel("")
	u.commandEntry = widget.NewEntry()
	u.categoriesEntry = widget.NewEntry()
	u.typeEntry = widget.NewEntry()
	u.typeEntry.SetText("Application")
	u.terminalSelect = widget.NewSelect([]string{"true", "false"}, func(s string) {
		if s == "true" {
			u.entry.Terminal = true
		} else {
			u.entry.Terminal = false
		}
	})
	u.terminalSelect.SetSelected("false")

	form := widget.NewForm(
		widget.NewFormItem("Name", u.nameEntry),
		widget.NewFormItem("Icon", widget.NewButtonWithIcon("Open", theme.FolderOpenIcon(), u.openIconDialog)),
		widget.NewFormItem("Icon Path: ", u.iconPathLabel),
		widget.NewFormItem("Command", u.commandEntry),
		widget.NewFormItem("Terminal", u.terminalSelect),
		widget.NewFormItem("Type", u.typeEntry),
		widget.NewFormItem("Categories", u.categoriesEntry),
	)
	return widget.NewVBox(
		form,
		layout.NewSpacer(),
		widget.NewHBox(
			layout.NewSpacer(),
			widget.NewButtonWithIcon("Save", theme.DocumentSaveIcon(), u.saveEntry),
		),
	)
}

func (u *ui) openIconDialog() {
	dialog.ShowFileOpen(func(f fyne.URIReadCloser, e error) {
		if f == nil {
			return
		}
		if e != nil {
			dialog.ShowError(e, u.mainWin)
			return
		}
		if f.URI().Extension() != ".png" {
			dialog.ShowError(fmt.Errorf("unsupported file extension: %s", f.URI().Extension()), u.mainWin)
			return
		}
		u.entry.Icon = f.URI().String()[7:]
		u.iconPathLabel.SetText(u.entry.Icon)
	}, u.mainWin)
}

func (u *ui) saveEntry() {
	u.entry.Name = u.nameEntry.Text
	u.entry.Command = u.commandEntry.Text
	u.entry.Type = u.typeEntry.Text
	u.entry.Categories = u.categoriesEntry.Text

	err := u.saveFile()
	if err != nil {
		dialog.ShowError(err, u.mainWin)
	}
}

func (u *ui) saveFile() error {
	if getProcessOwner() != "root\n" {
		return fmt.Errorf("can't save desktop entry: \nyou need to run this application with root privileges")
	}

	if u.entry.Name == "" {
		return fmt.Errorf("can't save desktop entry: missing file name")
	}
	if u.entry.Command == "" {
		return fmt.Errorf("can't save desktop entry: missing command")
	}

	fileName := "/usr/share/applications/" + u.nameEntry.Text + ".desktop"

	file, err := os.Create(fileName)
	if err != nil {
		return fmt.Errorf("failed to create file: %s", err)
	}

	_, err = file.WriteString(u.entry.string())
	if err != nil {
		return fmt.Errorf("failed to write file: %s", err)
	}
	return nil
}

func getProcessOwner() string {
	stdout, err := exec.Command("ps", "-o", "user=", "-p", strconv.Itoa(os.Getpid())).Output()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return string(stdout)
}

func main() {
	a := app.New()
	a.SetIcon(resourceIconPng)
	w := a.NewWindow("Desktop Entry Creator")
	w.SetIcon(resourceIconPng)
	w.Resize(fyne.NewSize(500, 380))
	app := &ui{app: a, mainWin: w}
	w.SetContent(app.loadMainUI())
	w.ShowAndRun()
}
