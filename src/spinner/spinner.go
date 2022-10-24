package spinner

import (
	"fmt"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type errMsg error

type Spinner struct {
	spinner  spinner.Model
	quitting bool
	err      error
}

func (m Spinner) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m Spinner) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		default:
			return m, nil
		}

	case errMsg:
		m.err = msg
		return m, nil

	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
}

func (m Spinner) View() string {
	if m.err != nil {
		return m.err.Error()
	}
	str := fmt.Sprintf("  %s Loading...", m.spinner.View())
	if m.quitting {
		return str + "\n"
	}
	return str
}
func (m Spinner) Tick() {
	m.spinner.Tick()
}
func New() Spinner {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return Spinner{spinner: s}
}