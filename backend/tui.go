package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"
	"time"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

const (
	maxNameLen = 24
	maxNoteLen = 200
	greeting   = "You've reached Arsenii. I can't come to the terminal right now.\nPlease leave your name and a message after the beep."
)

var (
	textStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#eceee9"))
	beepStyle  = lipgloss.NewStyle().Foreground(lipgloss.Color("#b3d4c3"))
	mutedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#6f7974"))
	errStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#e08a7a"))
)

type stage int

const (
	stageRinging stage = iota
	stageGreeting
	stageBeep
	stageName
	stageNote
	stageDone
)

type (
	advanceMsg struct{}
	typeMsg    struct{}
	hangUpMsg  struct{}
)

func after(d time.Duration, msg tea.Msg) tea.Cmd {
	return tea.Tick(d, func(time.Time) tea.Msg { return msg })
}

type model struct {
	db         *sql.DB
	remoteAddr string

	stage stage
	shown int // for typing effect

	name textinput.Model
	note textinput.Model
	err  error
}

func newModel(db *sql.DB, remoteAddr string) model {
	name := textinput.New()
	name.SetWidth(30)
	name.Placeholder = "your name"
	name.Prompt = ""
	name.CharLimit = maxNameLen

	note := textinput.New()
	note.SetWidth(50)
	note.Placeholder = "your message"
	note.Prompt = ""
	note.CharLimit = maxNoteLen

	return model{db: db, remoteAddr: remoteAddr, name: name, note: note}
}

func (m model) Init() tea.Cmd {
	return after(900*time.Millisecond, advanceMsg{})
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case advanceMsg:
		switch m.stage {
		case stageRinging:
			m.stage = stageGreeting
			return m, after(35*time.Millisecond, typeMsg{})
		case stageGreeting:
			m.stage = stageBeep
			return m, after(700*time.Millisecond, advanceMsg{})
		case stageBeep:
			return m.startName(), nil
		}
		return m, nil

	case typeMsg:
		if m.stage != stageGreeting {
			return m, nil
		}
		if m.shown < len([]rune(greeting)) {
			m.shown++
			return m, after(35*time.Millisecond, typeMsg{})
		}
		return m, after(time.Second, advanceMsg{})

	case hangUpMsg:
		return m, tea.Quit

	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "esc":
			return m, tea.Quit
		}

		switch m.stage {
		case stageName:
			if msg.String() == "enter" {
				if strings.TrimSpace(m.name.Value()) == "" {
					return m, nil
				}
				m.stage = stageNote
				m.name.Blur()
				return m, m.note.Focus()
			}
		case stageNote:
			if msg.String() == "enter" {
				return m.save()
			}

		case stageDone:
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd
	switch m.stage {
	case stageName:
		m.name, cmd = m.name.Update(msg)
	case stageNote:
		m.note, cmd = m.note.Update(msg)
	}
	return m, cmd
}

func (m model) startName() model {
	m.stage = stageName
	m.shown = len([]rune(greeting))
	m.name.Focus()
	return m
}

func (m model) save() (tea.Model, tea.Cmd) {
	name := strings.TrimSpace(m.name.Value())
	note := strings.TrimSpace(m.note.Value())
	if note == "" {
		return m, nil
	}
	if err := insertNote(m.db, name, note, m.remoteAddr); err != nil {
		log.Printf("saving note from %s: %v", m.remoteAddr, err)
		m.err = err
	}
	m.note.Blur()
	m.stage = stageDone
	return m, after(3*time.Second, hangUpMsg{})
}

func (m model) View() tea.View {
	s := mutedStyle.Render("*bzzzz...*") + "\n\n"

	if m.stage >= stageGreeting {
		visible := string([]rune(greeting)[:m.shown])
		s += textStyle.Render(visible) + "\n\n"
	}

	if m.stage >= stageBeep {
		s += beepStyle.Render("*BEEEEP*") + "\n\n"
	}

	if m.stage == stageName {
		s += mutedStyle.Render("name: ")
		s += m.name.View() + "\n\n"
	}

	if m.stage >= stageNote {
		s += mutedStyle.Render("name: ")
		s += textStyle.Render(m.name.Value()) + "\n\n"
	}

	if m.stage == stageNote {
		s += mutedStyle.Render("message:") + "\n"
		s += m.note.View() + "\n"
		count := len([]rune(m.note.Value()))
		s += mutedStyle.Render(fmt.Sprintf("%d / %d", count, maxNoteLen)) + "\n"
	}

	if m.stage == stageDone {
		if m.err != nil {
			s += errStyle.Render("the tape got stuck. try again later.") + "\n\n"
		} else {
			s += textStyle.Render("message recorded. thanks!") + "\n\n"
		}
		s += mutedStyle.Render("*click*") + "\n"
	}

	return tea.NewView(s)
}
