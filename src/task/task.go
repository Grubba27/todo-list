package task

type Column int

const Padding = 3

const (
	Todo Column = iota
	InProgress
	Done
)

type Task struct {
	index       int
	Status      Column
	title       string
	description string
}

func New(status Column, title string, description string) Task {
	return Task{Status: status, title: title, description: description}
}

func NewWithIndex(status Column, title string, description string, index int) Task {
	return Task{Status: status, title: title, description: description, index: index}
}

func (t Task) FilterValue() string {
	return t.title
}
func (t Task) Title() string {
	return t.title
}

func (t Task) Description() string {
	return t.description
}
func (t *Task) Next() {
	if t.Status == Done {
		t.Status = Todo
	} else {
		t.Status++
	}
}
