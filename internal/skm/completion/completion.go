// Package completion provides shell completion functions for skm.
package completion

import (
	"encoding/base64"
	"fmt"

	"github.com/mohammadv184/go-fido2"
	"github.com/mohammadv184/go-fido2/protocol/ctap2"
	"github.com/spf13/cobra"
)

// CompleteDevicePath provides shell completion for FIDO2 device paths.
func CompleteDevicePath(_ *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
	devs, err := fido2.Enumerate()
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	paths := make([]string, 0, len(devs))
	for _, dev := range devs {
		paths = append(paths, dev.Path)
	}

	return paths, cobra.ShellCompDirectiveNoFileComp
}

// CompleteCredentialID provides shell completion for credential IDs stored on a device.
func CompleteCredentialID(cmd *cobra.Command, _ []string, _ string) ([]string, cobra.ShellCompDirective) {
	devicePath, _ := cmd.Flags().GetString("device-path")
	pin, _ := cmd.Flags().GetString("pin")

	if devicePath == "" || pin == "" {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	dev, err := fido2.OpenPath(devicePath)
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}
	defer func() {
		_ = dev.Close()
	}()

	token, err := dev.GetPinUvAuthTokenUsingPIN(pin, ctap2.PermissionCredentialManagement, "")
	if err != nil {
		return nil, cobra.ShellCompDirectiveError
	}

	var completions []string
	for rp, err := range dev.EnumerateRPs(token) {
		if err != nil {
			continue
		}
		for c, err := range dev.EnumerateCredentials(token, rp.RPIDHash) {
			if err != nil {
				continue
			}
			id := base64.RawURLEncoding.EncodeToString(c.CredentialID.ID)

			rpName := rp.RP.Name
			if rpName == "" {
				rpName = rp.RP.ID
			}

			completions = append(completions, fmt.Sprintf("%s\t%s (%s)", id, c.User.Name, rpName))
		}
	}

	return completions, cobra.ShellCompDirectiveNoFileComp
}
