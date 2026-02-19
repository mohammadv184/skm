package views

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/mohammadv184/go-fido2"
	"github.com/mohammadv184/go-fido2/protocol/ctap2"
)

// DeviceInfoView is a view that displays detailed information about a security key.
type DeviceInfoView struct {
	desc       *fido2.DeviceDescriptor
	info       *ctap2.AuthenticatorGetInfoResponse
	pinRetries uint
	uvRetries  uint
	hasUV      bool
}

// NewDeviceInfoView creates a new DeviceInfoView.
func NewDeviceInfoView(desc *fido2.DeviceDescriptor, info *ctap2.AuthenticatorGetInfoResponse) *DeviceInfoView {
	return &DeviceInfoView{
		desc: desc,
		info: info,
	}
}

// WithRetries sets the PIN and UV retries for the view.
func (v *DeviceInfoView) WithRetries(pinRetries uint, uvRetries uint, hasUV bool) *DeviceInfoView {
	v.pinRetries = pinRetries
	v.uvRetries = uvRetries
	v.hasUV = hasUV
	return v
}

// Render renders the view.
func (v *DeviceInfoView) Render() string {
	titleStyle := lipgloss.NewStyle().
		Bold(true).
		PaddingBottom(1)

	labelStyle := lipgloss.NewStyle().
		Bold(true).
		Width(20)

	valueStyle := lipgloss.NewStyle()

	var b strings.Builder

	b.WriteString(titleStyle.Render("Security Key Information"))
	b.WriteString("\n")

	renderRow := func(label, value string) {
		b.WriteString(labelStyle.Render(label))
		b.WriteString(valueStyle.Render(value))
		b.WriteString("\n")
	}

	renderRow("Product:", v.desc.Product)
	renderRow("Manufacturer:", v.desc.Manufacturer)
	renderRow("Serial:", v.desc.SerialNumber)
	renderRow("Path:", v.desc.Path)
	renderRow("AAGUID:", v.info.AAGUID.String())

	renderRow("PIN Retries:", strconv.FormatUint(uint64(v.pinRetries), 10))
	if v.hasUV {
		renderRow("UV Retries:", strconv.FormatUint(uint64(v.uvRetries), 10))
	}

	versions := make([]string, len(v.info.Versions))
	for i, ver := range v.info.Versions {
		versions[i] = string(ver)
	}
	renderRow("Versions:", strings.Join(versions, ", "))

	extensions := make([]string, len(v.info.Extensions))
	for i, ext := range v.info.Extensions {
		extensions[i] = string(ext)
	}
	renderRow("Extensions:", strings.Join(extensions, ", "))

	if v.info.MaxMsgSize > 0 {
		renderRow("Max Msg Size:", fmt.Sprintf("%d bytes", v.info.MaxMsgSize))
	}

	if len(v.info.PinUvAuthProtocols) > 0 {
		protocols := make([]string, len(v.info.PinUvAuthProtocols))
		for i, p := range v.info.PinUvAuthProtocols {
			protocols[i] = p.String()
		}
		renderRow("PIN/UV Protocols:", strings.Join(protocols, ", "))
	}

	if v.info.Options != nil {
		b.WriteString("\n")
		b.WriteString(lipgloss.NewStyle().Bold(true).Underline(true).Render("Options:"))
		b.WriteString("\n")

		optionLabelStyle := labelStyle.PaddingLeft(2).Width(25)
		enabledStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("42"))   // Green
		disabledStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("241")) // Gray

		for k, val := range v.info.Options {
			b.WriteString(optionLabelStyle.Render(string(k) + ":"))
			if val {
				b.WriteString(enabledStyle.Render("✔ enabled"))
			} else {
				b.WriteString(disabledStyle.Render("✘ disabled"))
			}
			b.WriteString("\n")
		}
	}

	return b.String()
}
