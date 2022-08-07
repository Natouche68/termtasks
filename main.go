package main

import (
	"fmt"
	"os"
	"os/user"

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
	currentUser, userError := user.Current()
	if userError != nil {
		fmt.Printf("There was an error while getting the current user : %v\n", userError)
	}

	titlebarStyle := lipgloss.NewStyle().
		Align(lipgloss.Center).
		Foreground(lipgloss.Color("#9de2ff")).
		Background(lipgloss.Color("#0b1f88")).
		Width(terminalWidth).
		PaddingTop(1).
		PaddingBottom(1)

	titlebar := fmt.Sprintf(
		"%s%s",
		titlebarStyle.Copy().Bold(true).PaddingBottom(0).Render("TermTasks"),
		titlebarStyle.Copy().Italic(true).Render(currentUser.Username),
	)

	return titlebar
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
