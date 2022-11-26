package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

type Ui struct {
	selectedNote *Note
	notes        *NoteList

	content *widget.Entry
	list    *fyne.Container
}

func (u *Ui) addNote() {
	newNote := u.notes.add()
	u.setNote(newNote)
	u.notes.saveNotes()
}

func (u *Ui) setNote(n *Note) {
	u.selectedNote = n
	u.content.SetText(n.content)
	u.refreshList()
}

func (u *Ui) refreshList() {
	u.list.Objects = nil
	for _, n := range u.notes.list {
		thisNote := n

		button := widget.NewButton(n.title(), func() {
			u.setNote(thisNote)
		})

		if n == u.selectedNote {
			button.Importance = widget.HighImportance
		}
		u.list.Add(button)
	}
}

func (u *Ui) removeSelectedNote() {
	idx := u.notes.remove(u.selectedNote)

	if idx >= 0 && idx < len(u.notes.list) {
		u.setNote(u.notes.list[idx])
	} else {
		if idx > 0 {
			u.setNote(u.notes.list[idx-1])
		} else {
			u.content.SetText("")
			u.content.SetPlaceHolder("...")
		}
	}

	u.refreshList()
	u.notes.saveNotes()
}

func (u *Ui) loadUI() fyne.CanvasObject {
	u.content = widget.NewMultiLineEntry()
	u.list = container.NewVBox()
	u.refreshList()

	if len(u.notes.list) > 0 {
		u.setNote(u.notes.list[0])
	} else {
		u.content.SetText("")
		u.content.SetPlaceHolder("...")
	}

	u.content.OnChanged = func(content string) {
		if u.selectedNote == nil {
			return
		}
		u.selectedNote.content = content
		u.refreshList()
		u.notes.saveNotes()
	}

	tbar := widget.NewToolbar(
		widget.NewToolbarAction(theme.ContentAddIcon(), func() {
			u.addNote()
		}),
		widget.NewToolbarAction(theme.ContentRemoveIcon(), func() {
			u.removeSelectedNote()
		}),
	)

	scroll := container.NewScroll(u.list)

	side := container.New(layout.NewBorderLayout(tbar, nil, nil, nil), tbar, scroll)

	split := container.NewHSplit(side, u.content)
	split.Offset = 0.25
	return split
}

func main() {
	a := app.NewWithID("xyz.nono.notes")
	w := a.NewWindow("Nono Notes")

	list := &NoteList{prefs: a.Preferences()}
	list.loadNotesFromPreferences()
	notesUi := &Ui{notes: list}
	w.SetContent(notesUi.loadUI())
	w.Resize(fyne.Size{Width: 800, Height: 600})
	w.ShowAndRun()
}
