package views

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
	"github.com/mohammadv184/go-fido2"
)

// DevicesListView is a view that displays a table of connected security keys.
type DevicesListView struct {
	t *table.Table
}

// NewDevicesListView creates a new DevicesListView.
func NewDevicesListView() *DevicesListView {
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Padding(0, 1).
		Faint(true)

	cellStyle := lipgloss.NewStyle().
		Padding(0, 1)

	t := table.New().
		BorderStyle(lipgloss.NewStyle().Faint(true)).
		BorderRight(false).BorderLeft(false).BorderBottom(false).BorderTop(false).
		BorderColumn(false).
		StyleFunc(func(row, _ int) lipgloss.Style {
			if row == table.HeaderRow {
				return headerStyle
			}
			return cellStyle
		}).
		Headers("PATH", "PRODUCT", "MANUFACTURER", "SERIAL")

	return &DevicesListView{t: t}
}

// WithDevices adds devices to the view.
func (d *DevicesListView) WithDevices(devs ...fido2.DeviceDescriptor) *DevicesListView {
	for _, dev := range devs {
		d.t.Row(
			dev.Path,
			dev.Product,
			dev.Manufacturer,
			dev.SerialNumber,
		)
	}

	return d
}

// Render renders the view.
func (d *DevicesListView) Render() string {
	return lipgloss.NewStyle().Padding(1, 0, 1, 0).Render(d.t.Render())
}
