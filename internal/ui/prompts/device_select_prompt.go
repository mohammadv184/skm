package prompts

import (
	"errors"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mohammadv184/go-fido2"
)

// DeviceSelectPrompt is a prompt for selecting a security key from a list.
type DeviceSelectPrompt struct {
	table    table.Model
	devices  []fido2.DeviceDescriptor
	selected *fido2.DeviceDescriptor
	quitting bool
}

// NewDeviceSelectPrompt creates a new DeviceSelectPrompt.
func NewDeviceSelectPrompt() *DeviceSelectPrompt {
	columns := []table.Column{
		{Title: "PATH", Width: 20},
		{Title: "PRODUCT", Width: 30},
		{Title: "MANUFACTURER", Width: 40},
		{Title: "SERIAL", Width: 20},
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

	return &DeviceSelectPrompt{
		table: t,
	}
}

// WithDevices sets the devices to display in the prompt.
func (p *DeviceSelectPrompt) WithDevices(devs ...fido2.DeviceDescriptor) *DeviceSelectPrompt {
	p.devices = devs
	rows := make([]table.Row, len(devs))
	for i, dev := range devs {
		rows[i] = table.Row{
			dev.Path,
			dev.Product,
			dev.Manufacturer,
			dev.SerialNumber,
		}
	}
	p.table.SetRows(rows)
	return p
}

// Init initializes the bubbletea model.
func (p *DeviceSelectPrompt) Init() tea.Cmd {
	return nil
}

// Update handles message updates for the bubbletea model.
func (p *DeviceSelectPrompt) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			p.quitting = true
			return p, tea.Quit
		case "enter":
			idx := p.table.Cursor()
			if idx >= 0 && idx < len(p.devices) {
				p.selected = &p.devices[idx]
			}
			return p, tea.Quit
		}
	}

	p.table, cmd = p.table.Update(msg)
	return p, cmd
}

// View renders the prompt view.
func (p *DeviceSelectPrompt) View() string {
	if p.quitting {
		return ""
	}

	helpStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("241")).MarginTop(1)
	help := helpStyle.Render("↑/↓: move • enter: select • q/esc: quit")

	return lipgloss.NewStyle().Padding(1, 0, 1, 1).Render(
		"Select a security key:\n\n" +
			p.table.View() + "\n" +
			help,
	)
}

// Run executes the prompt and returns the selected device descriptor.
func (p *DeviceSelectPrompt) Run() (*fido2.DeviceDescriptor, error) {
	tm, err := tea.NewProgram(p).Run()
	if err != nil {
		return nil, err
	}

	finalModel, ok := tm.(*DeviceSelectPrompt)
	if !ok {
		return nil, errors.New("failed to get result from prompt")
	}

	if finalModel.selected == nil {
		return nil, errors.New("no device selected")
	}

	return finalModel.selected, nil
}
