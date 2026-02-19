package prompts

import (
	"errors"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// ConfirmPrompt is a simple yes/no confirmation prompt.
type ConfirmPrompt struct {
	title       string
	description string
	value       bool
	quitting    bool
	submitted   bool
}

// NewConfirmPrompt creates a new ConfirmPrompt.
func NewConfirmPrompt(title, description string) *ConfirmPrompt {
	return &ConfirmPrompt{
		title:       title,
		description: description,
	}
}

// Init initializes the bubbletea model.
func (p *ConfirmPrompt) Init() tea.Cmd {
	return nil
}

// Update handles message updates for the bubbletea model.
func (p *ConfirmPrompt) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc", "q":
			p.value = false
			p.quitting = true
			return p, tea.Quit
		case "y", "Y":
			p.value = true
			p.submitted = true
			return p, tea.Quit
		case "enter":
			p.submitted = true
			return p, tea.Quit
		}
	}
	return p, nil
}

// View renders the prompt view.
func (p *ConfirmPrompt) View() string {
	if p.quitting {
		return ""
	}

	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.AdaptiveColor{Light: "#11998e", Dark: "#4ecdc4"})
	descStyle := lipgloss.NewStyle().Faint(true)

	s := titleStyle.Render(p.title) + "\n"
	if p.description != "" {
		s += descStyle.Render(p.description) + "\n"
	}

	choices := " [y/N]"
	if p.value {
		choices = " [Y/n]"
	}

	return s + "\n" + "Confirm?" + choices + "\n"
}

// Run executes the prompt and returns the result.
func (p *ConfirmPrompt) Run() (bool, error) {
	tm, err := tea.NewProgram(p).Run()
	if err != nil {
		return false, err
	}

	finalModel, ok := tm.(*ConfirmPrompt)
	if !ok {
		return false, errors.New("failed to get result from prompt")
	}

	if finalModel.quitting {
		return false, errors.New("canceled")
	}

	return finalModel.value, nil
}
