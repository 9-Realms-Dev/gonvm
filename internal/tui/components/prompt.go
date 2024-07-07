package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var (
	PromptStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4"))

	InfoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FFAA")).
			Bold(true)

	SuccessStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FFAA")).
			Bold(true)

	WarnStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFAA00")).
			Bold(true)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000")).
			Bold(true)
)

type confirmModel struct {
	question string
	confirm  bool
	quitting bool
}

func (m confirmModel) Init() tea.Cmd {
	return nil
}

func (m confirmModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "y", "Y":
			m.confirm = true
			m.quitting = true
			return m, tea.Quit
		case "n", "N", "q", "Q", "ctrl+c":
			m.confirm = false
			m.quitting = true
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m confirmModel) View() string {
	if m.quitting {
		if m.confirm {
			return PromptStyle.Render("Confirmed!\n")
		}
		return ErrorStyle.Render("Cancelled.\n")
	}
	return PromptStyle.Render(fmt.Sprintf("%s (y/N) ", m.question))
}

func ConfirmPrompt(question string) (bool, error) {
	p := tea.NewProgram(confirmModel{question: question})
	m, err := p.Run()
	if err != nil {
		return false, err
	}
	return m.(confirmModel).confirm, nil
}
