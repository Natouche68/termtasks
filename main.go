package main

import (
	"fmt"
	"os"
	"os/user"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

var styles = map[string]lipgloss.Style{
	"statusBar": lipgloss.NewStyle().
		Background(lipgloss.Color("#3e3e3e")).
		Foreground(lipgloss.Color("#dfdfdf")).
		PaddingLeft(1),

	"statusBarTitle": lipgloss.NewStyle().
		Background(lipgloss.Color("#00b202")).
		Foreground(lipgloss.Color("#b2ffb3")).
		PaddingLeft(1).
		PaddingRight(1),
}

type task struct {
	name      string
	completed bool
}

type project struct {
	name  string
	tasks []task
}

type model struct {
	currentProject int
	currentTask    int
	projects       []project
}

func initModel() model {
	return model{
		currentProject: 0,
		currentTask:    0,
		projects: []project{
			{
				name: "TermTasks",
				tasks: []task{
					{
						name:      "Create the view function",
						completed: true,
					},
					{
						name:      "Add COLORS ! 🎉",
						completed: false,
					},
				},
			},
			{
				name: "Other",
				tasks: []task{
					{
						name:      "Give carrots to Bruno",
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
	terminalWidth, _, terminalSizeError := term.GetSize(int(os.Stdout.Fd()))
	if terminalSizeError != nil {
		return fmt.Sprintf("There was an error while getting the terminal's size : %v\n", terminalSizeError)
	}

	var s string

	currentUser, userError := user.Current()
	if userError != nil {
		fmt.Printf("There was an error while getting the current user : %v\n", userError)
	}

	statusBarTitle := styles["statusBarTitle"].Render("TermTasks")
	statusBar := styles["statusBar"].Copy().Width(terminalWidth - lipgloss.Width(statusBarTitle)).Render(currentUser.Username)

	s = fmt.Sprintf("%s", lipgloss.JoinHorizontal(
		lipgloss.Top,
		statusBarTitle,
		statusBar,
	))
	return s
}

func main() {
	if err := tea.NewProgram(initModel(), tea.WithAltScreen()).Start(); err != nil {
		fmt.Printf("There was an error during the starting of the programm : %v\n", err)
		os.Exit(1)
	}
}
