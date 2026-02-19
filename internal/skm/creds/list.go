package creds

import (
	"errors"
	"fmt"

	"github.com/mohammadv184/go-fido2"
	"github.com/mohammadv184/go-fido2/protocol/ctap2"
	"github.com/mohammadv184/skm/internal/skm/completion"
	"github.com/mohammadv184/skm/internal/ui/prompts"
	"github.com/mohammadv184/skm/internal/ui/views"
	"github.com/spf13/cobra"
)

var listCMD = cobra.Command{
	Use:     "list",
	Aliases: []string{"ls"},
	Short:   "List credentials stored on a security key",
	Long:    `Retrieve and display a list of all discoverable (resident) credentials stored on a selected security key. You will be prompted to select a device and enter its PIN.`,
	Example: `  skm creds list
  skm creds list --device-path /dev/hidraw0 --pin 123456`,
	RunE: listHandler,
}

var (
	listDevicePath string
	listPin        string
)

func init() {
	listCMD.Flags().StringVarP(&listDevicePath, "device-path", "d", "", "Path to the security key device")
	_ = listCMD.RegisterFlagCompletionFunc("device-path", completion.CompleteDevicePath)
	listCMD.Flags().StringVarP(&listPin, "pin", "p", "", "PIN for the security key")
	rootCMD.AddCommand(&listCMD)
}

func listHandler(cmd *cobra.Command, _ []string) error {
	var selectedDev *fido2.DeviceDescriptor

	if listDevicePath != "" {
		devs, err := fido2.Enumerate()
		if err != nil {
			return err
		}
		for _, dev := range devs {
			if dev.Path == listDevicePath {
				selectedDev = &dev
				break
			}
		}
		if selectedDev == nil {
			return fmt.Errorf("device not found at path: %s", listDevicePath)
		}
	} else {
		devs, err := fido2.Enumerate()
		if err != nil {
			return err
		}

		if len(devs) == 0 {
			return errors.New("no security keys found")
		}

		selectedDev, err = prompts.NewDeviceSelectPrompt().WithDevices(devs...).Run()
		if err != nil {
			return err
		}
	}

	dev, err := fido2.Open(*selectedDev)
	if err != nil {
		return err
	}
	defer func() {
		_ = dev.Close()
	}()

	pin := listPin
	if pin == "" {
		retries, _, _ := dev.GetPINRetries()
		pin, err = prompts.NewPinPrompt().WithRetries(retries).Run()
		if err != nil {
			return err
		}
	}

	token, err := dev.GetPinUvAuthTokenUsingPIN(pin, ctap2.PermissionCredentialManagement, "")
	if err != nil {
		return err
	}

	allCreds := make([]*ctap2.AuthenticatorCredentialManagementResponse, 0)

	for rp, err := range dev.EnumerateRPs(token) {
		if err != nil {
			return err
		}

		allCreds = append(allCreds, rp)
	}

	for i, cred := range allCreds {
		for c, err := range dev.EnumerateCredentials(token, cred.RPIDHash) {
			if err != nil {
				return err
			}

			c.RP = cred.RP
			allCreds[i] = c
		}
	}

	if len(allCreds) == 0 {
		cmd.Println("No credentials found on this device.")
		return nil
	}

	v := views.NewCredentialListView().WithCredentials(allCreds...)
	cmd.Println(v.Render())

	return nil
}
