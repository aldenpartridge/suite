package ui

import (
	"suite/pkg/styles"

	tea "github.com/charmbracelet/bubbletea"
)

type state int

const (
	stateCategories state = iota
	stateTools
)

// Category represents a group of related tools
type Category struct {
	Name  string
	Tools []string
}

// Model represents the main application UI model
type Model struct {
	categories     []Category
	currentTool    string
	cursor         int
	width          int
	height         int
	state          state
	activeCategory int
}

// New creates a new UI model
func New() Model {
	return Model{
		categories: []Category{
			{
				Name: "Network Analysis",
				Tools: []string{
					"Port Scanner",
					"Network Mapper",
				},
			},
			{
				Name: "Web Security",
				Tools: []string{
					"Subdomain Enumerator",
					"Header Scanner",
				},
			},
			{
				Name: "Cryptography",
				Tools: []string{
					"Hash Cracker",
				},
			},
			{
				Name: "System",
				Tools: []string{
					"Quit",
				},
			},
		},
		state: stateCategories,
	}
}

// CurrentTool returns the currently selected tool
func (m Model) CurrentTool() string {
	return m.currentTool
}

// Init implements tea.Model
func (m Model) Init() tea.Cmd {
	return nil
}

// Update implements tea.Model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.state == stateCategories {
				if m.cursor < len(m.categories)-1 {
					m.cursor++
				}
			} else {
				if m.cursor < len(m.categories[m.activeCategory].Tools)-1 {
					m.cursor++
				}
			}
		case "enter":
			if m.state == stateCategories {
				if m.categories[m.cursor].Name == "System" {
					return m, tea.Quit
				}
				m.state = stateTools
				m.activeCategory = m.cursor
				m.cursor = 0
			} else {
				tool := m.categories[m.activeCategory].Tools[m.cursor]
				if tool == "Quit" {
					return m, tea.Quit
				}
				m.currentTool = tool
				return m, tea.Quit
			}
		case "esc", "backspace":
			if m.state == stateTools {
				m.state = stateCategories
				m.cursor = m.activeCategory
			}
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
	}
	return m, nil
}

// View implements tea.Model
func (m Model) View() string {
	var s string

	if m.state == stateCategories {
		s = styles.Title.Render("🛡️  Cybersecurity Tool Suite")
		s += "\n\nSelect a category:\n\n"

		for i, cat := range m.categories {
			cursor := " "
			if m.cursor == i {
				cursor = ">"
			}
			s += styles.MenuItem.Render(cursor + " " + cat.Name)
			s += "\n"
		}

		s += "\n(↑/↓) Navigate • (enter) Select • (q) Quit\n"
	} else {
		category := m.categories[m.activeCategory]
		s = styles.Title.Render("🛡️  " + category.Name)
		s += "\n\nSelect a tool:\n\n"

		for i, tool := range category.Tools {
			cursor := " "
			if m.cursor == i {
				cursor = ">"
			}
			s += styles.MenuItem.Render(cursor + " " + tool)
			s += "\n"
		}

		s += "\n(↑/↓) Navigate • (enter) Select • (esc) Back • (q) Quit\n"
	}

	return s
}
