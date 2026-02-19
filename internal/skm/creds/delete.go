package creds

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"

	"github.com/mohammadv184/go-fido2"
	"github.com/mohammadv184/go-fido2/protocol/ctap2"
	"github.com/mohammadv184/skm/internal/skm/completion"
	"github.com/mohammadv184/skm/internal/ui/prompts"
	"github.com/spf13/cobra"
)

var deleteCMD = cobra.Command{
	Use:     "delete",
	Aliases: []string{"rm", "del", "remove"},
	Short:   "Delete a credential stored on a security key",
	Long:    `Permanently remove a specific discoverable (resident) credential from a selected security key. This action is irreversible.`,
	Example: `  skm creds delete
  skm creds delete --device-path /dev/hidraw0 --pin 123456 --credential-id base64-id`,
	RunE: deleteHandler,
}

var (
	deleteDevicePath   string
	deletePin          string
	deleteCredentialID string
)

func init() {
	deleteCMD.Flags().StringVarP(&deleteDevicePath, "device-path", "d", "", "Path to the security key device")
	_ = deleteCMD.RegisterFlagCompletionFunc("device-path", completion.CompleteDevicePath)
	deleteCMD.Flags().StringVarP(&deletePin, "pin", "p", "", "PIN for the security key")
	deleteCMD.Flags().
		StringVarP(&deleteCredentialID, "credential-id", "i", "", "ID of the credential to delete (base64 encoded)")
	_ = deleteCMD.RegisterFlagCompletionFunc("credential-id", completion.CompleteCredentialID)
	rootCMD.AddCommand(&deleteCMD)
}

func deleteHandler(cmd *cobra.Command, _ []string) error {
	var selectedDev *fido2.DeviceDescriptor

	if deleteDevicePath != "" {
		devs, err := fido2.Enumerate()
		if err != nil {
			return err
		}
		for _, dev := range devs {
			if dev.Path == deleteDevicePath {
				selectedDev = &dev
				break
			}
		}
		if selectedDev == nil {
			return fmt.Errorf("device not found at path: %s", deleteDevicePath)
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

	pin := deletePin
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

	var selectedCred *ctap2.AuthenticatorCredentialManagementResponse
	if deleteCredentialID != "" {
		decodedID, err := base64.RawURLEncoding.DecodeString(deleteCredentialID)
		if err != nil {
			return fmt.Errorf("failed to decode credential ID: %w", err)
		}
		for _, c := range allCreds {
			if bytes.Equal(c.CredentialID.ID, decodedID) {
				selectedCred = c
				break
			}
		}
		if selectedCred == nil {
			return fmt.Errorf("credential not found with ID: %s", deleteCredentialID)
		}
	} else {
		var err error
		selectedCred, err = prompts.NewCredentialSelectPrompt().WithCredentials(allCreds...).Run()
		if err != nil {
			return err
		}
	}

	err = dev.DeleteCredential(token, selectedCred.CredentialID)
	if err != nil {
		return err
	}

	cmd.Println("Credential deleted successfully.")
	return nil
}
