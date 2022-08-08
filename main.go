package main

import (
	"fmt"
	"os"
	"os/user"
	"strings"

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

	"tab": lipgloss.NewStyle().
		Border(lipgloss.Border{
			Top:         "─",
			Bottom:      "─",
			Left:        "│",
			Right:       "│",
			TopLeft:     "╭",
			TopRight:    "╮",
			BottomLeft:  "┴",
			BottomRight: "┴",
		}, true).
		BorderForeground(lipgloss.Color("#00b202")).
		Padding(0, 1),

	"activeTab": lipgloss.NewStyle().
		Border(lipgloss.Border{
			Top:         "─",
			Bottom:      " ",
			Left:        "│",
			Right:       "│",
			TopLeft:     "╭",
			TopRight:    "╮",
			BottomLeft:  "┘",
			BottomRight: "└",
		}, true).
		BorderForeground(lipgloss.Color("#00b202")).
		Foreground(lipgloss.Color("#00b202")).
		Padding(0, 1),

	"tabGap": lipgloss.NewStyle().
		Border(lipgloss.Border{
			Bottom: "─",
		}, true).
		BorderTop(false).
		BorderLeft(false).
		BorderRight(false).
		BorderForeground(lipgloss.Color("#00b202")),

	"task": lipgloss.NewStyle().
		Border(lipgloss.Border{
			Left: "│",
		}, true).
		BorderTop(false).
		BorderBottom(false).
		BorderRight(false).
		BorderForeground(lipgloss.Color("#8a8a8a")).
		PaddingLeft(1).
		MarginBottom(1),

	"currentTask": lipgloss.NewStyle().
		Border(lipgloss.Border{
			Left: ">",
		}, true).
		BorderTop(false).
		BorderBottom(false).
		BorderRight(false).
		BorderForeground(lipgloss.Color("#00b202")).
		Foreground(lipgloss.Color("#00b202")).
		PaddingLeft(1).
		MarginBottom(1),

	"title": lipgloss.NewStyle().
		Bold(true).
		Italic(true).
		Padding(0, 1).
		MarginBottom(1).
		Background(lipgloss.Color("#00b202")),

	"help": lipgloss.NewStyle().
		Foreground(lipgloss.Color("#8a8a8a")),
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
	currentAction  string // Can be : tasks, help or add
	projects       []project
}

func initModel() model {
	return model{
		currentProject: 0,
		currentTask:    0,
		currentAction:  "tasks",
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

		case "tab":
			m.currentProject++
			if m.currentProject == len(m.projects) {
				m.currentProject = 0
			}

		case "shift+tab":
			m.currentProject--
			if m.currentProject == -1 {
				m.currentProject = len(m.projects) - 1
			}

		case "down":
			m.currentTask++
			if m.currentTask == len(m.projects[m.currentProject].tasks) {
				m.currentTask = 0
			}

		case "up":
			m.currentTask--
			if m.currentTask == -1 {
				m.currentTask = len(m.projects[m.currentProject].tasks) - 1
			}

		case "enter":
			m.projects[m.currentProject].tasks[m.currentTask].completed = !m.projects[m.currentProject].tasks[m.currentTask].completed

		case "h":
			if m.currentAction == "help" {
				m.currentAction = "tasks"
			} else {
				m.currentAction = "help"
			}
		}
	}

	return m, nil
}

func (m model) View() string {
	terminalWidth, terminalHeight, terminalSizeError := term.GetSize(int(os.Stdout.Fd()))
	if terminalSizeError != nil {
		return fmt.Sprintf("There was an error while getting the terminal's size : %v\n", terminalSizeError)
	}

	currentUser, userError := user.Current()
	if userError != nil {
		fmt.Printf("There was an error while getting the current user : %v\n", userError)
	}

	if m.currentAction == "tasks" {
		// Tabs
		var tabs string
		for i, openProject := range m.projects {
			var tab string
			if m.currentProject == i {
				tab = styles["activeTab"].Render(openProject.name)
			} else {
				tab = styles["tab"].Render(openProject.name)
			}
			tabs = lipgloss.JoinHorizontal(
				lipgloss.Top,
				tabs,
				tab,
			)
		}
		tabsGap := styles["tabGap"].Render(strings.Repeat(" ", terminalWidth-lipgloss.Width(tabs)))
		tabs = lipgloss.JoinHorizontal(lipgloss.Bottom, tabs, tabsGap)

		// Tasks
		var tasks string
		for i, task := range m.projects[m.currentProject].tasks {
			var renderedTask string
			if i == m.currentTask {
				if task.completed {
					renderedTask = styles["currentTask"].Copy().Strikethrough(true).Render(task.name)
				} else {
					renderedTask = styles["currentTask"].Render(task.name)
				}
			} else {
				if task.completed {
					renderedTask = styles["task"].Copy().Strikethrough(true).Render(task.name)
				} else {
					renderedTask = styles["task"].Render(task.name)
				}
			}

			tasks = lipgloss.JoinVertical(
				lipgloss.Left,
				tasks,
				renderedTask,
			)
		}

		// Blank between tasks and status bar
		blankSpace := strings.Repeat(
			"\n",
			terminalHeight-lipgloss.Height(lipgloss.JoinVertical(lipgloss.Left, tabs, tasks))-3,
		)

		// Status Bar
		statusBarTitle := styles["statusBarTitle"].Render("TermTasks")
		statusBar := styles["statusBar"].
			Copy().
			Width(terminalWidth - lipgloss.Width(statusBarTitle)).
			Render(currentUser.Username)

		return lipgloss.JoinVertical(
			lipgloss.Left,
			tabs,
			tasks,
			blankSpace,
			lipgloss.JoinHorizontal(
				lipgloss.Top,
				statusBarTitle,
				statusBar,
			),
		)
	} else if m.currentAction == "add" {
		return "You can't add a task now 😥..."
	} else if m.currentAction == "help" {
		title := styles["title"].Render("Help")

		help := styles["help"].Render(fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s",
			"q or ctrl+c  : quit",
			"h            : help",
			"",
			"tab          : switch project (to the right)",
			"shift+tab    : switch project (to the left)",
			"up           : move the cursor up",
			"down         : move the cursor down",
			"",
			"a            : add a task",
		))

		// Blank between help and status bar
		blankSpace := strings.Repeat(
			"\n",
			terminalHeight-lipgloss.Height(lipgloss.JoinVertical(lipgloss.Left, title, help))-3,
		)

		// Status Bar
		statusBarTitle := styles["statusBarTitle"].Render("TermTasks")
		statusBar := styles["statusBar"].
			Copy().
			Width(terminalWidth - lipgloss.Width(statusBarTitle)).
			Render(currentUser.Username)

		return lipgloss.JoinVertical(
			lipgloss.Left,
			title,
			help,
			blankSpace,
			lipgloss.JoinHorizontal(
				lipgloss.Top,
				statusBarTitle,
				statusBar,
			),
		)
	}

	return "There was an error, you seem to be trying to do an invalid action."
}

func main() {
	if err := tea.NewProgram(initModel(), tea.WithAltScreen()).Start(); err != nil {
		fmt.Printf("There was an error during the starting of the programm : %v\n", err)
		os.Exit(1)
	}
}
