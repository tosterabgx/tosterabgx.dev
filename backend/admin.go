package main

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

const adminVisible = 15

var selectedStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#b3d4c3"))

type adminModel struct {
	db         *sql.DB
	notes      []AdminNote
	cursor     int
	confirming bool
	status     string
}

func newAdminModel(db *sql.DB) adminModel {
	return adminModel{db: db}.reload()
}

func (m adminModel) reload() adminModel {
	notes, err := listNotesAdmin(m.db, 500)
	if err != nil {
		m.status = "error loading notes: " + err.Error()
		return m
	}
	m.notes = notes
	if m.cursor >= len(m.notes) {
		m.cursor = max(0, len(m.notes)-1)
	}
	return m
}

func (m adminModel) Init() tea.Cmd { return nil }

func (m adminModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	key, ok := msg.(tea.KeyPressMsg)
	if !ok {
		return m, nil
	}

	if m.confirming {
		m.confirming = false
		if strings.ToLower(key.String()) != "y" {
			m.status = "cancelled"
			return m, nil
		}
		n := m.notes[m.cursor]
		if err := deleteNote(m.db, n.ID); err != nil {
			m.status = "delete failed: " + err.Error()
			return m, nil
		}
		log.Printf("admin deleted note %d (%q from %s)", n.ID, n.Name, n.IP)
		m.status = fmt.Sprintf("deleted #%d", n.ID)
		return m.reload(), nil
	}

	switch key.String() {
	case "q", "esc", "ctrl+c":
		return m, tea.Quit
	case "up":
		if m.cursor > 0 {
			m.cursor--
		}
	case "down":
		if m.cursor < len(m.notes)-1 {
			m.cursor++
		}
	case "d":
		if len(m.notes) > 0 {
			m.confirming = true
		}
	case "r":
		m = m.reload()
		m.status = "reloaded"
	}
	return m, nil
}

func (m adminModel) View() tea.View {
	s := beepStyle.Render(fmt.Sprintf("%d notes", len(m.notes))) + "\n\n"

	if len(m.notes) == 0 {
		s += mutedStyle.Render("no notes yet") + "\n"
	}

	start := max(0, m.cursor-adminVisible/2)
	end := min(len(m.notes), start+adminVisible)

	for i := start; i < end; i++ {
		n := m.notes[i]
		line := fmt.Sprintf("#%-4d %-15s %-12s %s",
			n.ID, n.IP, truncate(n.Name, 12), truncate(n.Message, 50))
		if i == m.cursor {
			s += selectedStyle.Render("> "+line) + "\n"
		} else {
			s += mutedStyle.Render("  "+line) + "\n"
		}
	}

	s += "\n"
	if m.confirming {
		n := m.notes[m.cursor]
		s += errStyle.Render(fmt.Sprintf("delete #%d from %s? y/n", n.ID, n.Name)) + "\n"
	} else if m.status != "" {
		s += textStyle.Render(m.status) + "\n"
	}

	return tea.NewView(s)
}

func truncate(s string, max int) string {
	r := []rune(s)
	if len(r) <= max {
		return s
	}
	return string(r[:max-1]) + "..."
}
