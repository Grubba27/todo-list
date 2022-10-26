package app

import (
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	s "todo-list/src/spinner"
	t "todo-list/src/task"
)

var Models []tea.Model

const (
	Column t.Column = iota
	form
)
const InitWidth = 300
const InitHeight = 300

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

type App struct {
	focused   t.Column
	lists     []list.Model
	err       error
	isLoading bool
	Spinner   s.Spinner
	quitting  bool
	showRaw   bool
}

func New() *App {
	return &App{isLoading: true, quitting: false, Spinner: s.New()}
}
func (m *App) GetSelectedTask() t.Task {
	selectedItem := m.lists[m.focused].SelectedItem()
	selectedTask := selectedItem.(t.Task)
	return selectedTask
}
func (m *App) DeleteTask() t.Task {
	selectedTask := m.GetSelectedTask()
	m.lists[selectedTask.Status].
		RemoveItem(m.lists[m.focused].Index())
	return selectedTask
}
func (m *App) MoveDown() tea.Msg {
	if len(m.lists[m.focused].Items()) == 0 {
		return nil
	}
	selectedList := &m.lists[m.focused]
	selectedList.CursorDown()
	return nil
}
func (m *App) MoveUp() tea.Msg {
	if len(m.lists[m.focused].Items()) == 0 {
		return nil
	}
	selectedList := &m.lists[m.focused]
	selectedList.CursorUp()
	return nil
}
func (m *App) MoveToNext() tea.Msg {
	if len(m.lists[m.focused].Items()) == 0 {
		return nil
	}
	selectedTask := m.DeleteTask()
	selectedTask.Next()
	m.lists[selectedTask.Status].
		InsertItem(
			len(m.lists[selectedTask.Status].Items())-1,
			list.Item(selectedTask),
		)
	return nil
}

func (m *App) Next() {
	if m.focused == t.Done {
		m.focused = t.Todo
	} else {
		m.focused++
	}
}

func (m *App) Prev() {
	if m.focused == t.Todo {
		m.focused = t.Done
	} else {
		m.focused--
	}
}
func (m *App) initList(width, height int) {
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
		t.New(t.Done, "Drink Coffee", "just coffee"),
	})

}
func (m *App) Refresh(list []list.Model) {
	m.lists = list
}

func (m *App) Init() tea.Cmd {
	m.initList(InitWidth, InitHeight)
	return nil
}
func (m *App) View() string {
	if m.quitting {
		return "Bye!"
	}
	if m.isLoading {
		return m.Spinner.View()
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
func (m *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {

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
		case "k", "up":
			m.MoveUp()
		case "j", "down":
			m.MoveDown()
		case "e":
			if len(m.lists[m.focused].Items()) == 0 {
				return m, nil
			}
			Models[Column] = m
			selectedTask := m.GetSelectedTask()
			Models[form] = EditForm(
				m.focused,
				selectedTask.Title(),
				selectedTask.Description(),
			)
			m.DeleteTask()
			return Models[form].Update(nil)
		case "enter":
			return m, m.MoveToNext
		case "n":
			Models[Column] = m
			Models[form] = NewForm(m.focused)
			return Models[form].Update(nil)
		case "d":
			if len(m.lists[m.focused].Items()) == 0 {
				return m, nil
			}
			m.DeleteTask()
		case "s":
			Models[Column] = m
			m.showRaw = true
			d := NewDebug("salve teste", m.lists)
			return d.Update(nil)

		default:
			return m, nil
		}
	case t.Task:
		task := msg
		return m, m.lists[m.focused].InsertItem(len(m.lists[task.Status].Items()), task)

	}

	if m.isLoading {
		return m.Spinner.Update(msg)
	}
	var cmd tea.Cmd
	m.lists[m.focused], cmd = m.lists[m.focused].Update(msg)
	return m, cmd
}
