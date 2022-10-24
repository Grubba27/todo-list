package main

import (
	"fmt"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"os"
	s "todo-list/src/spinner"
	t "todo-list/src/task"
)

var (
	columnStyle = lipgloss.NewStyle().
			Padding(1, 2)

	focusedStyle = lipgloss.NewStyle().
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			BorderBackground(lipgloss.Color("62"))

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241"))
)

type Model struct {
	focused   t.Status
	lists     []list.Model
	err       error
	isLoading bool
	spinner   s.Spinner
	quitting  bool
}

func New() *Model {
	return &Model{isLoading: true, quitting: false, spinner: s.New()}
}
func (m *Model) MoveToNext() tea.Msg {
	selectedItem := m.lists[m.focused].SelectedItem()
	if len(m.lists[m.focused].Items()) == 0 {
		return nil
	}
	selectedTask := selectedItem.(t.Task)
	m.lists[selectedTask.Status].
		RemoveItem(m.lists[m.focused].Index())
	selectedTask.Next()
	m.lists[selectedTask.Status].
		InsertItem(
			len(m.lists[selectedTask.Status].Items())-1,
			list.Item(selectedTask),
		)
	return nil
}

func (m *Model) Next() {
	if m.focused == t.Done {
		m.focused = t.Todo
	} else {
		m.focused++
	}
}

func (m *Model) Prev() {
	if m.focused == t.Todo {
		m.focused = t.Done
	} else {
		m.focused--
	}
}
func (m *Model) initList(width, height int) {
	defaultList := list.New(
		[]list.Item{},
		list.NewDefaultDelegate(),
		width/t.Padding,
		height/2,
	)
	m.lists = []list.Model{defaultList, defaultList, defaultList}

	m.lists[t.Todo].Title = "To Do"
	m.lists[t.Todo].SetItems([]list.Item{
		t.New(t.Todo, "Sleep", "dream some"),
		t.New(t.Todo, "Idk", "yay"),
	})

	m.lists[t.InProgress].Title = "In Progress"
	m.lists[t.InProgress].SetItems([]list.Item{
		t.New(t.InProgress, "Do something", "buy some groceries"),
		t.New(t.InProgress, "Do something2", "buy some groceries"),
	})

	m.lists[t.Done].Title = "Done"
	m.lists[t.Done].SetItems([]list.Item{
		t.New(t.Done, "Drink Coffee", "spicy coffee"),
	})

}

func (m *Model) Init() tea.Cmd {
	return nil
}
func (m *Model) View() string {
	if m.quitting {
		return "Bye!"
	}
	if m.isLoading {
		return m.spinner.View()
	}

	todoView := m.lists[t.Todo].View()
	inProgress := m.lists[t.InProgress].View()
	done := m.lists[t.Done].View()
	switch m.focused {
	case t.Done:
		return lipgloss.JoinHorizontal(
			lipgloss.Left,
			columnStyle.Render(todoView),
			columnStyle.Render(inProgress),
			focusedStyle.Render(done),
		)
	case t.InProgress:
		return lipgloss.JoinHorizontal(
			lipgloss.Left,
			columnStyle.Render(todoView),
			focusedStyle.Render(inProgress),
			columnStyle.Render(done),
		)
	default:
		return lipgloss.JoinHorizontal(
			lipgloss.Left,
			focusedStyle.Render(todoView),
			columnStyle.Render(inProgress),
			columnStyle.Render(done),
		)
	}

}
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		columnStyle.Width(msg.Width / t.Padding)
		focusedStyle.Width(msg.Width / t.Padding)
		columnStyle.Height(msg.Height - t.Padding)
		focusedStyle.Height(msg.Height - t.Padding)
		m.initList(msg.Width, msg.Height)
		m.isLoading = false
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		case "left", "h":
			m.Prev()
		case "right", "l":
			m.Next()
		case "enter":
			return m, m.MoveToNext
		default:
			return m, nil
		}

	}
	var cmd tea.Cmd
	m.lists[m.focused], cmd = m.lists[m.focused].Update(msg)
	return m, cmd
}

func main() {
	m := New()
	p := tea.NewProgram(m)

	err := p.Start()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
