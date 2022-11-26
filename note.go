package main

import (
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
)

const (
	noteCountKey  = "noteCount"
	noteKeyFormat = "note%d"
)

type Note struct {
	content string
}

func (n *Note) title() string {
	if strings.TrimSpace(n.content) == "" {
		return "Untitled"
	}
	return strings.SplitN(n.content, "\n", 2)[0]
}

type NoteList struct {
	list  []*Note
	prefs fyne.Preferences
}

func (l *NoteList) add() *Note {
	n := &Note{}
	l.list = append([]*Note{n}, l.list...)
	return n
}

func (l *NoteList) loadNotesFromPreferences() {
	total := l.prefs.Int(noteCountKey)
	if total == 0 {
		tutorialNode := &Note{
			content: `Tutorial Note
			- Add a new note by pressing the  +  symbol in the toolbar.
			- Remove the selected note by pressing the  -  symbol in the toolbar.
			- The title of the note will be the first line in the edit window.
			- Notes are saved automatically.`,
		}
		l.list = append(l.list, tutorialNode)
		return
	}

	for i := 0; i < total; i++ {
		key := fmt.Sprintf(noteKeyFormat, i)
		content := l.prefs.String(key)

		l.list = append(l.list, &Note{content: content})
	}
}

func (l *NoteList) remove(n *Note) int {
	if len(l.list) == 0 {
		return -1
	}

	for i, note := range l.list {
		if note != n {
			continue
		}

		if i == len(l.list)-1 {
			l.list = l.list[:i]
			return i
		} else {
			l.list = append(l.list[:i], l.list[i+1:]...)
			return i
		}
	}

	return -1
}

func (l *NoteList) saveNotes() {
	for i, n := range l.list {
		key := fmt.Sprintf(noteKeyFormat, i)
		l.prefs.SetString(key, n.content)
	}
	l.prefs.SetInt(noteCountKey, len(l.list))
}
