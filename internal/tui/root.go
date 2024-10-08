package tui

import (
	"fmt"
	"os"
	"strings"

	"github.com/9-Realms-Dev/gonvm/internal/tui/services"
	"github.com/9-Realms-Dev/gonvm/internal/tui/styles"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
)

type errMsg error

type model struct {
	spinner  spinner.Model
	current  services.NodeVersion
	versions []string
	quitting bool
	err      error
}

func initialModel() model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = styles.SpinnerStyle

	nodeVersion, err := services.GetCurrentDetails()
	if err != nil {
		os.Exit(1)
	}

	versions, err := services.GetVersions()
	if err != nil {
		os.Exit(1)
	}

	return model{
		spinner:  s,
		current:  *nodeVersion,
		versions: versions,
	}
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
	var s strings.Builder

	// Build the string to display at the top of the screen for the current version information
	current := fmt.Sprintf("Node: %sNPM: %s", m.current.NodeVersion, m.current.NpmVersion)
	s.WriteString(styles.HeaderStyle.Render(current))

	s.WriteString("\n\nLocal Versions:\n")
	for _, v := range m.versions {
		s.WriteString(fmt.Sprintf("\n%s", v))
	}

	// Help text
	s.WriteString("\n\nPress 'q' to quit\n")

	if m.quitting {
		return "Quitting..."
	}
	return styles.BaseStyle.Render(s.String())
}

func Dashboard() {
	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
