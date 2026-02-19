package prompts

import (
	"encoding/base64"
	"errors"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mohammadv184/go-fido2/protocol/ctap2"
)

// CredentialSelectPrompt is a prompt for selecting a credential from a list.
type CredentialSelectPrompt struct {
	table    table.Model
	creds    []*ctap2.AuthenticatorCredentialManagementResponse
	selected *ctap2.AuthenticatorCredentialManagementResponse
	quitting bool
}

// NewCredentialSelectPrompt creates a new CredentialSelectPrompt.
func NewCredentialSelectPrompt() *CredentialSelectPrompt {
	columns := []table.Column{
		{Title: "RP", Width: 20},
		{Title: "USER", Width: 20},
		{Title: "DISPLAY NAME", Width: 20},
		{Title: "CREDENTIAL ID", Width: 20},
	}

	s := table.DefaultStyles()
	s.Header = s.Header.
		Bold(true).
		Padding(0, 1).
		Faint(true).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.AdaptiveColor{Light: "#bbbbbb", Dark: "#555555"}).
		BorderBottom(true)
	s.Cell = s.Cell.
		Padding(0, 1)
	s.Selected = s.Selected.
		Foreground(lipgloss.AdaptiveColor{Light: "#11998e", Dark: "#4ecdc4"}).
		Bold(true)

	t := table.New(
		table.WithColumns(columns),
		table.WithFocused(true),
		table.WithHeight(10),
	)
	t.SetStyles(s)

	return &CredentialSelectPrompt{
		table: t,
	}
}

// WithCredentials sets the credentials to display in the prompt.
func (p *CredentialSelectPrompt) WithCredentials(
	creds ...*ctap2.AuthenticatorCredentialManagementResponse,
) *CredentialSelectPrompt {
	p.creds = creds
	rows := make([]table.Row, len(creds))
	for i, cred := range creds {
		rpName := cred.RP.Name
		if rpName == "" {
			rpName = cred.RP.ID
		}

		credID := base64.RawURLEncoding.EncodeToString(cred.CredentialID.ID)
		if len(credID) > 20 {
			credID = credID[:17] + "..."
		}

		rows[i] = table.Row{
			rpName,
			cred.User.Name,
			cred.User.DisplayName,
			credID,
		}
	}
	p.table.SetRows(rows)
	return p
}

// Init initializes the bubbletea model.
func (p *CredentialSelectPrompt) Init() tea.Cmd {
	return nil
}

// Update handles message updates for the bubbletea model.
func (p *CredentialSelectPrompt) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			p.quitting = true
			return p, tea.Quit
		case "enter":
			idx := p.table.Cursor()
			if idx >= 0 && idx < len(p.creds) {
				p.selected = p.creds[idx]
			}
			return p, tea.Quit
		}
	}

	p.table, cmd = p.table.Update(msg)
	return p, cmd
}

// View renders the prompt view.
func (p *CredentialSelectPrompt) View() string {
	if p.quitting {
		return ""
	}

	helpStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("241")).MarginTop(1)
	help := helpStyle.Render("↑/↓: move • enter: select • q/esc: quit")

	return lipgloss.NewStyle().Margin(1, 2).Render(
		"Select a credential:\n\n" +
			p.table.View() + "\n" +
			help,
	)
}

// Run executes the prompt and returns the selected credential.
func (p *CredentialSelectPrompt) Run() (*ctap2.AuthenticatorCredentialManagementResponse, error) {
	tm, err := tea.NewProgram(p).Run()
	if err != nil {
		return nil, err
	}

	finalModel, ok := tm.(*CredentialSelectPrompt)
	if !ok {
		return nil, errors.New("failed to get result from prompt")
	}

	if finalModel.selected == nil {
		return nil, errors.New("no credential selected")
	}

	return finalModel.selected, nil
}
