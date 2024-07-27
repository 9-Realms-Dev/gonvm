package tui

import (
	"fmt"
	"github.com/9-Realms-Dev/gonvm/internal/tui/styles"

	tea "github.com/charmbracelet/bubbletea"
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
			return styles.PromptStyle.Render("Confirmed!\n")
		}
		return styles.ErrorStyle.Render("Cancelled.\n")
	}
	return styles.PromptStyle.Render(fmt.Sprintf("%s (y/N) ", m.question))
}

func ConfirmPrompt(question string) (bool, error) {
	p := tea.NewProgram(confirmModel{question: question})
	m, err := p.Run()
	if err != nil {
		return false, err
	}
	return m.(confirmModel).confirm, nil
}
