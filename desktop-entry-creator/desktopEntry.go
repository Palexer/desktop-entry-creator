package main

import "fmt"

type desktopEntry struct {
	Name       string
	Icon       string
	Type       string
	Command    string
	Categories string
	Comment    string
}

func (e *desktopEntry) string() string {
	if e.Type == "nil" {
		e.Type = "Application"
	}
	return fmt.Sprintf("[Desktop Entry]\nName=%s\nIcon=%s\nType=%s\nExec=%s\nTerminal=false\nCategories=%v\nComment=%s\nVersion=1.0\n",
		e.Name,
		e.Icon,
		e.Type,
		e.Command,
		e.Categories,
		e.Comment,
	)
}
