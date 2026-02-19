package prompts

import (
	"errors"
	"fmt"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// PinPrompt is a prompt for entering a security key PIN.
type PinPrompt struct {
	textInput   textinput.Model
	quitting    bool
	submitted   bool
	retries     uint
	title       string
	placeholder string
	validate    func(string) error
	err         error
}

// NewPinPrompt creates a new PinPrompt.
func NewPinPrompt() *PinPrompt {
	ti := textinput.New()
	ti.Placeholder = "PIN"
	ti.Focus()
	ti.EchoMode = textinput.EchoPassword
	ti.EchoCharacter = '•'

	return &PinPrompt{
		textInput:   ti,
		title:       "Enter PIN",
		placeholder: "PIN",
	}
}

// WithTitle sets the title of the prompt.
func (p *PinPrompt) WithTitle(title string) *PinPrompt {
	p.title = title
	return p
}

// WithPlaceholder sets the placeholder text for the PIN input.
func (p *PinPrompt) WithPlaceholder(placeholder string) *PinPrompt {
	p.placeholder = placeholder
	p.textInput.Placeholder = placeholder
	return p
}

// WithRetries sets the number of remaining PIN retries to display a warning if necessary.
func (p *PinPrompt) WithRetries(retries uint) *PinPrompt {
	p.retries = retries
	return p
}

// WithValidation sets a validation function for the entered PIN.
func (p *PinPrompt) WithValidation(v func(string) error) *PinPrompt {
	p.validate = v
	return p
}

// Init initializes the bubbletea model.
func (p *PinPrompt) Init() tea.Cmd {
	return textinput.Blink
}

// Update handles message updates for the bubbletea model.
func (p *PinPrompt) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		p.err = nil // Reset error on any key press
		switch msg.String() {
		case "ctrl+c", "esc":
			p.quitting = true
			return p, tea.Quit
		case "enter":
			if p.validate != nil {
				if err := p.validate(p.textInput.Value()); err != nil {
					p.err = err
					return p, nil
				}
			}
			p.submitted = true
			return p, tea.Quit
		}
	}

	p.textInput, cmd = p.textInput.Update(msg)
	return p, cmd
}

// View renders the prompt view.
func (p *PinPrompt) View() string {
	if p.quitting {
		return ""
	}

	var warning string
	if p.retries == 1 {
		warning = lipgloss.NewStyle().
			Foreground(lipgloss.Color("201")).
			Render("Warning: This is your LAST attempt before the device is locked.") +
			"\n\n"
	}

	var errView string
	if p.err != nil {
		errView = lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Render(p.err.Error()) + "\n\n"
	}

	return fmt.Sprintf(
		"%s\n\n%s%s%s\n\n%s",
		p.title,
		warning,
		errView,
		p.textInput.View(),
		"(esc to quit)",
	) + "\n"
}

// Run executes the prompt and returns the entered PIN.
func (p *PinPrompt) Run() (string, error) {
	tm, err := tea.NewProgram(p).Run()
	if err != nil {
		return "", err
	}

	finalModel, ok := tm.(*PinPrompt)
	if !ok {
		return "", errors.New("failed to get result from prompt")
	}

	if finalModel.quitting || !finalModel.submitted {
		return "", errors.New("canceled")
	}

	return finalModel.textInput.Value(), nil
}
