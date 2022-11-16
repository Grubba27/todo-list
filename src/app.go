package app

import (
	"fmt"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"todo-list/src/db"
	t "todo-list/src/task"
)

var Models []tea.Model

const (
	Column t.Column = iota
	form
)
const InitWidth = 300
const InitHeight = 300

type KeyMap struct {
	Quit     key.Binding
	Up       key.Binding
	Down     key.Binding
	Left     key.Binding
	Right    key.Binding
	Create   key.Binding
	Filter   key.Binding
	Edit     key.Binding
	Delete   key.Binding
	MoveNext key.Binding
}

var DefaultKeyMap = KeyMap{
	Quit: key.NewBinding(
		key.WithKeys("q", "esc", "ctrl+c"),
		key.WithHelp("q", "quit"),
	),
	Up: key.NewBinding(
		key.WithKeys("k", "up"),
		key.WithHelp("↑/k", "up"),
	),
	Down: key.NewBinding(
		key.WithKeys("j", "down"),
		key.WithHelp("↓/j", "down"),
	),
	Right: key.NewBinding(
		key.WithKeys("l", "right"),
		key.WithHelp("→/l", "right"),
	),
	Left: key.NewBinding(
		key.WithKeys("h", "left"),
		key.WithHelp("←/h", "left"),
	),
	Create: key.NewBinding(
		key.WithKeys("n"),
		key.WithHelp("n", "new"),
	),
	Filter: key.NewBinding(
		key.WithKeys("/"),
		key.WithHelp("/", "filter"),
	),
	Edit: key.NewBinding(
		key.WithKeys("e"),
		key.WithHelp("e", "edit"),
	),
	Delete: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "delete"),
	),
	MoveNext: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("↩", "move next"),
	),
}

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
	Spinner   spinner.Model
	quitting  bool
	showRaw   bool
}

func New() *App {
	spin := spinner.New()
	spin.Spinner = spinner.Dot
	spin.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	return &App{isLoading: true, quitting: false, Spinner: spin}
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
func InitAppUI(width, height int) list.Model {
	defaultList := list.New(
		[]list.Item{},
		list.NewDefaultDelegate(),
		width/t.Padding,
		height/2,
	)
	defaultList.KeyMap = list.KeyMap{
		CursorUp: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "up"),
		),
		CursorDown: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "down"),
		),
		PrevPage: key.NewBinding(
			key.WithKeys("left", "h", "pgup", "b", "u"),
			key.WithHelp("←/h/pgup", "prev page"),
		),
		NextPage: key.NewBinding(
			key.WithKeys("right", "l", "pgdown", "f", "d"),
			key.WithHelp("→/l/pgdn", "next page"),
		),
		GoToStart: key.NewBinding(
			key.WithKeys("home", "g"),
			key.WithHelp("g/home", "go to start"),
		),
		GoToEnd: key.NewBinding(
			key.WithKeys("end", "G"),
			key.WithHelp("G/end", "go to end"),
		),
		Filter: key.NewBinding(
			key.WithKeys("/"),
			key.WithHelp("/", "filter"),
		),
		ClearFilter: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "clear filter"),
		),

		// Filtering.
		CancelWhileFiltering: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "cancel"),
		),
		AcceptWhileFiltering: key.NewBinding(
			key.WithKeys("enter", "tab", "shift+tab", "ctrl+k", "up", "ctrl+j", "down"),
			key.WithHelp("enter", "apply filter"),
		),
		ShowFullHelp: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "more"),
		),
		CloseFullHelp: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "close help"),
		),

		// Quitting.
		Quit: key.NewBinding(
			key.WithKeys("q", "esc"),
			key.WithHelp("q", "quit"),
		),
		ForceQuit: key.NewBinding(key.WithKeys("ctrl+c")),
	}
	extraKeys := []key.Binding{
		DefaultKeyMap.Quit,
		DefaultKeyMap.Up,
		DefaultKeyMap.Down,
		DefaultKeyMap.Left,
		DefaultKeyMap.Right,
		DefaultKeyMap.Create,
		DefaultKeyMap.Filter,
		DefaultKeyMap.Edit,
		DefaultKeyMap.Delete,
		DefaultKeyMap.MoveNext,
	}
	defaultList.AdditionalShortHelpKeys = func() []key.Binding {
		return extraKeys
	}
	defaultList.SetShowHelp(false)
	return defaultList
}
func (m *App) initList(width, height int, finished func()) tea.Cmd {
	defaultList := InitAppUI(width, height)
	m.lists = []list.Model{defaultList, defaultList, defaultList}
	db.Connect()

	m.lists[t.Todo].Title = "To Do"
	todoList := db.GetTasksByColumn(t.Todo)
	m.lists[t.Todo].SetItems(todoList)

	m.lists[t.InProgress].Title = "In Progress"
	inProgressList := db.GetTasksByColumn(t.InProgress)
	m.lists[t.InProgress].SetItems(inProgressList)

	m.lists[t.Done].Title = "Done"
	doneList := db.GetTasksByColumn(t.Done)
	m.lists[t.Done].SetItems(doneList)

	defer finished()
	f := func() tea.Cmd {
		return nil
	}
	return f()
}
func (m *App) Refresh(list []list.Model) {
	m.lists = list
}

func (m *App) Init() tea.Cmd {
	return tea.Batch(m.Spinner.Tick, m.initList(InitWidth, InitHeight, func() {}))
}
func (m *App) View() string {
	if m.quitting {
		return "Bye!"
	}

	if m.isLoading {
		return fmt.Sprintf("\n\n  %s  loading... \n\n", m.Spinner.View())
	}

	todoView := m.lists[t.Todo]
	inProgress := m.lists[t.InProgress]
	done := m.lists[t.Done]

	todoView.SetShowHelp(false)
	inProgress.SetShowHelp(false)
	done.SetShowHelp(false)
	switch m.focused {
	case t.Done:
		done.SetShowHelp(true)
		return lipgloss.JoinHorizontal(
			lipgloss.Left,
			columnStyle.Render(todoView.View()),
			columnStyle.Render(inProgress.View()),
			focusedStyle.Render(done.View()),
		)
	case t.InProgress:
		inProgress.SetShowHelp(true)
		return lipgloss.JoinHorizontal(
			lipgloss.Left,
			columnStyle.Render(todoView.View()),
			focusedStyle.Render(inProgress.View()),
			columnStyle.Render(done.View()),
		)
	default:
		todoView.SetShowHelp(true)
		return lipgloss.JoinHorizontal(
			lipgloss.Left,
			focusedStyle.Render(todoView.View()),
			columnStyle.Render(inProgress.View()),
			columnStyle.Render(done.View()),
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
		go m.initList(msg.Width, msg.Height, func() {
			m.isLoading = false
		})
		if m.isLoading {
			var cmd tea.Cmd
			m.Spinner, cmd = m.Spinner.Update(msg)
			return m, cmd
		}
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, DefaultKeyMap.Quit):
			m.quitting = true
			return m, tea.Quit
		case key.Matches(msg, DefaultKeyMap.Up):
			m.MoveUp()
		case key.Matches(msg, DefaultKeyMap.Down):
			m.MoveDown()
		case key.Matches(msg, DefaultKeyMap.Left):
			m.Prev()
		case key.Matches(msg, DefaultKeyMap.Right):
			m.Next()
		case key.Matches(msg, DefaultKeyMap.Create):
			Models[Column] = m
			Models[form] = NewForm(m.focused)
			return Models[form].Update(nil)
		case key.Matches(msg, DefaultKeyMap.Edit):
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
		case key.Matches(msg, DefaultKeyMap.Delete):
			if len(m.lists[m.focused].Items()) == 0 {
				return m, nil
			}
			m.DeleteTask()
		case key.Matches(msg, DefaultKeyMap.MoveNext):
			return m, m.MoveToNext
		default:
			return m, nil
		}
	case t.Task:
		task := msg
		return m, m.lists[m.focused].InsertItem(len(m.lists[task.Status].Items()), task)
	}

	if m.isLoading {
		var cmd tea.Cmd
		m.Spinner, cmd = m.Spinner.Update(msg)
		return m, cmd
	}

	var cmd tea.Cmd
	m.lists[m.focused], cmd = m.lists[m.focused].Update(msg)
	return m, cmd
}
