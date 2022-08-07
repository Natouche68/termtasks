package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

type task struct {
	name      string
	completed bool
}

type project struct {
	name  string
	tasks []task
}

type model struct {
	projects []project
}

var terminalWidth int
var terminalHeight int
var terminalSizeError error

func initModel() model {
	return model{
		projects: []project{
			{
				name: "TermTasks",
				tasks: []task{
					{
						name:      "Create the view function",
						completed: false,
					},
					{
						name:      "Add COLORS ! 🎉",
						completed: false,
					},
				},
			},
		},
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		}
	}

	return m, nil
}

func (m model) View() string {
	style := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#37dcff")).
		Background(lipgloss.Color("#000f3b")).
		Padding(4, 12, 4, 12).
		Margin(2, 1, 2, 1).
		Width(terminalWidth - 2).
		Height(terminalHeight - 4).
		Align(lipgloss.Center)

	return style.Render("Hello World !")
}

func main() {
	terminalWidth, terminalHeight, terminalSizeError = term.GetSize(int(os.Stdout.Fd()))
	if terminalSizeError != nil {
		fmt.Printf("There was an error while getting the terminal's size : %v\n", terminalSizeError)
	}

	if err := tea.NewProgram(initModel(), tea.WithAltScreen()).Start(); err != nil {
		fmt.Printf("There was an error during the starting of the programm : %v\n", err)
		os.Exit(1)
	}
}
