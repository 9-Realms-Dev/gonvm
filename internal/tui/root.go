package tui

import (
	"fmt"
	"os"

	"github.com/9-Realms-Dev/gonvm/internal/tui/styles"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type errMsg error

// TODO: Add details from the curreent version of nvm in the direcorty
// TODO: Handle no versions in the directory
// TODO: Show the available versions in the directory
// TODO: Find out how to show the global imports in the current directory

type model struct {
	spinner  spinner.Model
	quitting bool
	err      error
}

func initialModel() model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = styles.SpinnerStyle
	return model{spinner: s}
}

func (m model) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "esc", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		default:
			return m, nil
		}

	case errMsg:
		m.err = msg
		return m, nil

	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
}

func (m model) View() string {
	if m.err != nil {
		return m.err.Error()
	}
	str := fmt.Sprintf("\n\n   %s coming soon please press q to quit\n\n", m.spinner.View())
	if m.quitting {
		return str + "\n"
	}
	return str
}

func Dashboard() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
