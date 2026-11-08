package main

import (
	"fmt"
	"log"
	"strings"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

const (
	maxNameLen = 24
	maxNoteLen = 200
)

var (
	titleStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#b3d4c3"))
	mutedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#6f7974"))
	labelStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#9ea7a2")).Width(6)
)

type model struct {
	name       textinput.Model
	note       textinput.Model
	focused    int
	submitted  bool
	remoteAddr string
}

func newModel(remoteAddr string) model {
	name := textinput.New()
	name.Placeholder = "your name"
	name.CharLimit = maxNameLen
	name.Focus()

	note := textinput.New()
	note.Placeholder = "say something to the world!"
	note.CharLimit = maxNoteLen

	return model{name: name, note: note, remoteAddr: remoteAddr}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.submitted {
		if _, ok := msg.(tea.KeyPressMsg); ok {
			return m, tea.Quit
		}
		return m, nil
	}

	if key, ok := msg.(tea.KeyPressMsg); ok {
		switch key.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit

		case "tab", "shift+tab", "up", "down":
			return m.switchFocus(), nil

		case "enter":
			if m.focused == 0 {
				return m.switchFocus(), nil
			}
			name := strings.TrimSpace(m.name.Value())
			note := strings.TrimSpace(m.note.Value())
			if name == "" || note == "" {
				return m, nil
			}
			m.submitted = true
			// TODO add saving of notes
			log.Printf("note from %s: name=%q note=%q", m.remoteAddr, name, note)
			return m, nil
		}
	}

	var cmd tea.Cmd
	if m.focused == 0 {
		m.name, cmd = m.name.Update(msg)
	} else {
		m.note, cmd = m.note.Update(msg)
	}
	return m, cmd
}

func (m model) switchFocus() model {
	if m.focused == 0 {
		m.focused = 1
		m.name.Blur()
		m.note.Focus()
	} else {
		m.focused = 0
		m.note.Blur()
		m.name.Focus()
	}
	return m
}

func (m model) View() tea.View {
	s := titleStyle.Render("guestbook") + "\n"
	s += mutedStyle.Render("leave a note on tosterabgx.dev") + "\n\n"

	if m.submitted {
		s += "thanks! your note will show up on my site.\n\n"
		s += mutedStyle.Render("press any key to leave")
		return tea.NewView(s)
	}

	s += labelStyle.Render("name") + m.name.View() + "\n"
	s += labelStyle.Render("note") + m.note.View() + "\n"

	count := len([]rune(m.note.Value()))
	s += mutedStyle.Render(fmt.Sprintf("    %d / %d", count, maxNoteLen)) + "\n\n"
	s += mutedStyle.Render("tab next · enter send · esc quit")

	return tea.NewView(s)
}
