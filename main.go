package main

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"os"
	a "todo-list/src"
	s "todo-list/src/spinner"
	t "todo-list/src/task"
)

func main() {
	app := a.New()
	app.Spinner = s.New()
	form := a.NewForm(t.Todo)
	a.Models = []tea.Model{app, form}
	m := a.Models[a.Column]
	p := tea.NewProgram(m)
	err := p.Start()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
