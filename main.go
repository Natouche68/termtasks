package main

import (
	"fmt"
	"os"
	"os/user"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.org/x/term"
)

var styles = map[string]lipgloss.Style{
	"statusBar": lipgloss.NewStyle().
		Background(lipgloss.Color("#3e3e3e")).
		Foreground(lipgloss.Color("#dfdfdf")).
		Padding(0, 1),

	"title": lipgloss.NewStyle().
		Bold(true).
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

	"help": lipgloss.NewStyle().
		Italic(true).
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
	createTaskInput textinput.Model
	currentProject  int
	currentTask     int
	currentAction   string // Can be : tasks, help or add
	projects        []project
}

func initModel() model {
	ti := textinput.New()
	ti.CharLimit = 50
	ti.Width = 53

	return model{
		createTaskInput: ti,
		currentProject:  0,
		currentTask:     0,
		currentAction:   "tasks",
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
	return textinput.Blink
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit

		case "q":
			if m.currentAction != "add" {
				return m, tea.Quit
			}

		case "tab":
			m.currentProject++
			m.currentTask = 0
			if m.currentProject == len(m.projects) {
				m.currentProject = 0
			}

		case "shift+tab":
			m.currentProject--
			m.currentTask = 0
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
			if m.currentAction == "tasks" {
				m.projects[m.currentProject].tasks[m.currentTask].completed = !m.projects[m.currentProject].tasks[m.currentTask].completed
			} else if m.currentAction == "add" {
				if m.createTaskInput.Value() != "" {
					m.projects[m.currentProject].tasks = append(m.projects[m.currentProject].tasks, task{
						name:      m.createTaskInput.Value(),
						completed: false,
					})
					m.createTaskInput.Reset()
				}
				m.currentAction = "tasks"
				m.createTaskInput.Blur()
			}

		case "h":
			if m.currentAction == "help" {
				m.currentAction = "tasks"
			} else if m.currentAction == "tasks" {
				m.currentAction = "help"
			}

		case "a":
			if m.currentAction == "tasks" {
				m.currentAction = "add"
				m.createTaskInput.Focus()
				return m, nil
			}

		case "esc":
			if m.currentAction == "add" {
				m.createTaskInput.Blur()
				m.currentAction = "tasks"
			}

		case "d":
			if m.currentAction == "tasks" {
				m.projects[m.currentProject].tasks = append(
					m.projects[m.currentProject].tasks[:m.currentTask],
					m.projects[m.currentProject].tasks[m.currentTask+1:]...,
				)
				m.currentTask--
				if m.currentTask == -1 {
					m.currentTask = len(m.projects[m.currentProject].tasks) - 1
				}
			}
		}
	}

	var cmd tea.Cmd
	m.createTaskInput, cmd = m.createTaskInput.Update(msg)
	return m, cmd
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
		if len(m.projects[m.currentProject].tasks) == 0 {
			tasks = styles["help"].Render("\nThis project is empty...")
		}

		// Blank between tasks and status bar
		blankSpace := strings.Repeat(
			"\n",
			terminalHeight-lipgloss.Height(lipgloss.JoinVertical(lipgloss.Left, tabs, tasks))-3,
		)

		// Status Bar
		statusBarTitle := styles["title"].Render("TermTasks")
		statusBarHelp := styles["statusBar"].
			Copy().
			Foreground(lipgloss.Color("#919191")).
			Italic(true).
			Render("h : help  q : quit")
		statusBar := styles["statusBar"].
			Copy().
			Width(terminalWidth - (lipgloss.Width(statusBarTitle) + lipgloss.Width(statusBarHelp))).
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
				statusBarHelp,
			),
		)
	} else if m.currentAction == "add" {
		title := styles["title"].Render("Add a new task") + "\n"

		textInput := lipgloss.JoinVertical(
			lipgloss.Left,
			"What's the name of the new task ?",
			m.createTaskInput.View(),
			styles["help"].Render("(press esc to return to the home)"),
		)

		// Blank between input and status bar
		blankSpace := strings.Repeat(
			"\n",
			terminalHeight-lipgloss.Height(lipgloss.JoinVertical(lipgloss.Left, title, textInput))-3,
		)

		// Status Bar
		statusBarTitle := styles["title"].Render("TermTasks")
		statusBarHelp := styles["statusBar"].
			Copy().
			Foreground(lipgloss.Color("#919191")).
			Italic(true).
			Render("h : help  q : quit")
		statusBar := styles["statusBar"].
			Copy().
			Width(terminalWidth - (lipgloss.Width(statusBarTitle) + lipgloss.Width(statusBarHelp))).
			Render(currentUser.Username)

		return lipgloss.JoinVertical(
			lipgloss.Left,
			title,
			textInput,
			blankSpace,
			lipgloss.JoinHorizontal(
				lipgloss.Top,
				statusBarTitle,
				statusBar,
				statusBarHelp,
			),
		)
	} else if m.currentAction == "help" {
		title := styles["title"].Render("Help") + "\n"

		help := styles["help"].Render(fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s",
			"q or ctrl+c  : quit",
			"h            : help",
			"",
			"tab          : switch project (to the right)",
			"shift+tab    : switch project (to the left)",
			"up           : move the cursor up",
			"down         : move the cursor down",
			"",
			"a            : add a task",
			"d            : delete the current task",
		))

		// Blank between help and status bar
		blankSpace := strings.Repeat(
			"\n",
			terminalHeight-lipgloss.Height(lipgloss.JoinVertical(lipgloss.Left, title, help))-3,
		)

		// Status Bar
		statusBarTitle := styles["title"].Render("TermTasks")
		statusBarHelp := styles["statusBar"].
			Copy().
			Foreground(lipgloss.Color("#919191")).
			Italic(true).
			Render("h : help  q : quit")
		statusBar := styles["statusBar"].
			Copy().
			Width(terminalWidth - (lipgloss.Width(statusBarTitle) + lipgloss.Width(statusBarHelp))).
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
				statusBarHelp,
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
