package views

import (
	"encoding/base64"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/mohammadv184/go-fido2/protocol/ctap2"
)

// CredentialListView is a view that displays a table of credentials.
type CredentialListView struct {
	t *table.Table
}

// NewCredentialListView creates a new CredentialListView.
func NewCredentialListView() *CredentialListView {
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Padding(0, 1).
		Faint(true).
		Border(lipgloss.NormalBorder(), false, false, true, false)

	cellStyle := lipgloss.NewStyle().
		Padding(0, 1)

	t := table.New().
		Border(lipgloss.HiddenBorder()).
		StyleFunc(func(row, _ int) lipgloss.Style {
			if row == table.HeaderRow {
				return headerStyle
			}
			return cellStyle
		}).
		Headers("RP", "USER", "DISPLAY NAME", "CREDENTIAL ID")

	return &CredentialListView{t: t}
}

// WithCredentials adds credentials to the view.
func (d *CredentialListView) WithCredentials(
	creds ...*ctap2.AuthenticatorCredentialManagementResponse,
) *CredentialListView {
	for _, cred := range creds {
		rpName := cred.RP.Name
		if rpName == "" {
			rpName = cred.RP.ID
		}

		userName := cred.User.Name
		displayName := cred.User.DisplayName

		credID := base64.RawURLEncoding.EncodeToString(cred.CredentialID.ID)

		d.t.Row(
			rpName,
			userName,
			displayName,
			credID,
		)
	}

	return d
}

// Render renders the view.
func (d *CredentialListView) Render() string {
	return d.t.Render()
}
