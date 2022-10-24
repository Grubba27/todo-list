package task

type Status int

const Padding = 4

const (
	Todo Status = iota
	InProgress
	Done
)

type Task struct {
	Status      Status
	title       string
	description string
}

func New(status Status, title string, description string) Task {
	return Task{Status: status, title: title, description: description}
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
