package portscan

import (
	"fmt"
	"strconv"
	"strings"

	"suite/pkg/styles"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type state int

const (
	stateInput state = iota
	stateScanning
	stateDone
)

// Model represents the port scanner UI
type Model struct {
	targetInput textinput.Model
	startInput  textinput.Model
	endInput    textinput.Model
	progress    progress.Model
	state       state
	results     []Result
	err         error
	width       int
	height      int
}

// NewModel creates a new port scanner UI model
func NewModel() Model {
	target := textinput.New()
	target.Placeholder = "Enter target (e.g., localhost)"
	target.Focus()

	start := textinput.New()
	start.Placeholder = "Start port (e.g., 1)"

	end := textinput.New()
	end.Placeholder = "End port (e.g., 1024)"

	prog := progress.New(progress.WithDefaultGradient())

	return Model{
		targetInput: target,
		startInput:  start,
		endInput:    end,
		progress:    prog,
		state:       stateInput,
	}
}

// Init implements tea.Model
func (m Model) Init() tea.Cmd {
	return textinput.Blink
}

// Update implements tea.Model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "tab":
			if m.state == stateInput {
				if m.targetInput.Focused() {
					m.targetInput.Blur()
					m.startInput.Focus()
				} else if m.startInput.Focused() {
					m.startInput.Blur()
					m.endInput.Focus()
				} else if m.endInput.Focused() {
					m.endInput.Blur()
					m.targetInput.Focus()
				}
			}
		case "enter":
			if m.state == stateInput {
				if err := m.validate(); err != nil {
					m.err = err
					return m, nil
				}
				m.state = stateScanning
				return m, m.startScan
			}
		}
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.progress.Width = msg.Width - 4
	case progressMsg:
		if msg.err != nil {
			m.err = msg.err
			m.state = stateInput
			return m, nil
		}
		if msg.done {
			m.results = msg.results
			m.state = stateDone
			return m, nil
		}
		cmd := m.progress.SetPercent(msg.progress)
		return m, cmd
	}

	if m.state == stateInput {
		m.targetInput, cmd = m.targetInput.Update(msg)
		return m, cmd
	}

	return m, nil
}

// View implements tea.Model
func (m Model) View() string {
	var s strings.Builder

	s.WriteString(styles.Title.Render("🔍 Port Scanner"))
	s.WriteString("\n\n")

	switch m.state {
	case stateInput:
		s.WriteString(fmt.Sprintf("Target: %s\n", m.targetInput.View()))
		s.WriteString(fmt.Sprintf("Start Port: %s\n", m.startInput.View()))
		s.WriteString(fmt.Sprintf("End Port: %s\n\n", m.endInput.View()))
		s.WriteString("(tab) Switch fields • (enter) Start scan • (ctrl+c) Quit\n")

	case stateScanning:
		s.WriteString(fmt.Sprintf("Scanning %s...\n\n", m.targetInput.Value()))
		s.WriteString(m.progress.View() + "\n\n")
		s.WriteString("Press ctrl+c to cancel\n")

	case stateDone:
		s.WriteString(fmt.Sprintf("Scan complete for %s\n\n", m.targetInput.Value()))
		if len(m.results) == 0 {
			s.WriteString("No open ports found.\n")
		} else {
			s.WriteString("Open ports:\n\n")
			for _, r := range m.results {
				s.WriteString(fmt.Sprintf("%-6d %-10s %s\n", r.Port, r.State, r.Service))
			}
		}
		s.WriteString("\nPress q to quit\n")
	}

	if m.err != nil {
		s.WriteString("\n" + styles.Error.Render(m.err.Error()) + "\n")
	}

	return styles.Container.Render(s.String())
}

func (m Model) validate() error {
	if m.targetInput.Value() == "" {
		return fmt.Errorf("target is required")
	}

	start, err := strconv.Atoi(m.startInput.Value())
	if err != nil || start < 1 {
		return fmt.Errorf("invalid start port")
	}

	end, err := strconv.Atoi(m.endInput.Value())
	if err != nil || end < start || end > 65535 {
		return fmt.Errorf("invalid end port")
	}

	return nil
}

type progressMsg struct {
	progress float64
	results  []Result
	done     bool
	err      error
}

func (m Model) startScan() tea.Msg {
	start, _ := strconv.Atoi(m.startInput.Value())
	end, _ := strconv.Atoi(m.endInput.Value())

	scanner := NewScanner(m.targetInput.Value(), start, end)
	progress := make(chan float64)

	go func() {
		results, err := scanner.Scan(progress)
		if err != nil {
			tea.NewProgram(m).Send(progressMsg{err: err})
			return
		}
		tea.NewProgram(m).Send(progressMsg{progress: 100, results: results, done: true})
	}()

	return progressMsg{progress: 0}
}
