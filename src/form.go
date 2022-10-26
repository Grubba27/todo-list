package app

import (
	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"todo-list/src/task"
)

type Form struct {
	column      task.Column
	title       textinput.Model
	description textarea.Model
}

func NewForm(column task.Column) *Form {
	f := &Form{column: column}
	f.title = textinput.New()
	f.title.Focus()
	f.title.Placeholder = "Place your title"
	f.description = textarea.New()
	f.description.Placeholder = " Place your description"
	return f
}
func EditForm(column task.Column, title, description string) *Form {
	f := &Form{column: column}
	f.title = textinput.New()
	f.title.Focus()
	f.title.SetValue(title)
	f.description = textarea.New()
	f.description.SetValue(description)
	return f
}
func (m Form) CreateTask() tea.Msg {
	t := task.New(m.column, m.title.Value(), m.description.Value())
	return t
}
func (m Form) Init() tea.Cmd {
	return nil
}
func (m Form) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "enter":
			if m.title.Focused() {
				m.title.Blur()
				m.description.Focus()
				return m, textarea.Blink
			} else {
				Models[form] = m
				return Models[Column], m.CreateTask
			}
		}
	}
	if m.title.Focused() {
		m.title, cmd = m.title.Update(msg)
		return m, cmd
	} else {
		m.description, cmd = m.description.Update(msg)
		return m, cmd
	}
}
func (m Form) View() string {
	return lipgloss.JoinVertical(lipgloss.Left, m.title.View(), m.description.View())
}
