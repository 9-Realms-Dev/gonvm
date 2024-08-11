package styles

import "github.com/charmbracelet/lipgloss"

var (
	PromptStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#7aa2f7"))

	InfoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#5dbffc")).
			Bold(true)

	SuccessStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#5ae184")).
			Bold(true)

	WarnStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#f69058")).
			Bold(true)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#ff6b6b")).
			Bold(true)

	SpinnerStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#5dbffc"))

	BaseStyle = lipgloss.NewStyle().
			Padding(1, 2)

	HeaderStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), false, false, true).
			BorderForeground(lipgloss.Color("#5dbffc"))
)
