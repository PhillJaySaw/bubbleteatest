package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

type task struct {
	id      string
	content string
}

type model struct {
	choices   []task
	cursor    string
	selected  map[string]struct{}
	inputMode bool
	newTask   string
}

func initializeModel() model {
	tasks := []task{
		{
			id:      "1",
			content: "Go running",
		},
		{
			id:      "2",
			content: "Groceries",
		},
		{
			id:      "3",
			content: "Do the dishes",
		},
	}

	return model{
		choices:   tasks,
		selected:  make(map[string]struct{}),
		cursor:    tasks[0].id,
		inputMode: false,
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		if m.inputMode {
			switch msg.String() {
			case "enter":
				if len(m.newTask) > 0 {
					id := getUnsafeId()
					m.choices = append(m.choices, task{id: id, content: m.newTask})
				}

				m.inputMode = false
				m.newTask = ""
			case "esc":
				m.inputMode = false
				m.newTask = ""
			case "backspace":
				if len(m.newTask) > 0 {
					m.newTask = m.newTask[:len(m.newTask)-1]
				}
			default:
				if len(msg.String()) == 1 && len(m.newTask) < 20 {
					m.newTask += msg.String()
				}
			}
		} else {
			switch msg.String() {

			case "ctrl+c", "q":
				return m, tea.Quit

			case "down", "j":
				m.Next()
			case "up", "k":
				m.Prev()

			case "enter", " ":
				_, ok := m.selected[m.cursor]
				if ok {
					delete(m.selected, m.cursor)
				} else {
					m.selected[m.cursor] = struct{}{}
				}
			case "a":
				m.inputMode = true

				if m.cursor == "" {
					m.cursor = m.choices[0].id
				}
			case "d":
				var indexToRemove int
				for i, choice := range m.choices {
					if choice.id == m.cursor {
						indexToRemove = i
						break
					}
				}
				copy(m.choices[indexToRemove:], m.choices[indexToRemove+1:])
				m.choices = m.choices[:len(m.choices)-1]

				if len(m.choices) == 0 {
					m.cursor = ""
				} else if indexToRemove == len(m.choices) {
					m.cursor = m.choices[len(m.choices)-1].id
				} else if indexToRemove == 0 {
					m.cursor = m.choices[0].id
				} else {
					m.cursor = m.choices[indexToRemove].id
				}
			}
		}

	}

	return m, nil
}

func (m model) View() string {
	s := "What do you need to do?\n\n"

	for _, choice := range m.choices {
		cursor := " "
		if m.cursor == choice.id {
			cursor = ">"
		}

		checked := " "
		if _, ok := m.selected[choice.id]; ok {
			checked = "x"
		}

		s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice.content)

	}

	if len(m.choices) == 0 {
		s += fmt.Sprintf("No tasks to dispaly\n")
	}

	if m.inputMode {
		s += fmt.Sprintf("\nNew task: %s_\n", m.newTask)
		s += "\nenter - add task | esc - cancel\n"
	} else {
		s += "\nj/k - up/down | a - add new task | space - check | d - delete\n"
	}

	return s
}

func (m *model) Next() {
	for i, task := range m.choices {
		if task.id == m.cursor {
			if i+1 < len(m.choices) {
				m.cursor = m.choices[i+1].id
				break
			}
		}
	}
}

func (m *model) Prev() {
	for i, task := range m.choices {
		if task.id == m.cursor {
			if i-1 >= 0 {
				m.cursor = m.choices[i-1].id
				break
			}
			break
		}
	}
}

func main() {
	p := tea.NewProgram(initializeModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error occured while trying to run the program: %v", err)
		os.Exit(1)
	}
}
