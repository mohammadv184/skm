package config

import (
	"errors"
	"fmt"

	"github.com/mohammadv184/go-fido2"
	"github.com/mohammadv184/go-fido2/protocol/ctap2"
	"github.com/mohammadv184/skm/internal/skm/completion"
	"github.com/mohammadv184/skm/internal/ui/prompts"
	"github.com/spf13/cobra"
)

var enterpriseAttestationCMD = cobra.Command{
	Use:   "enterprise-attestation",
	Short: "Enable Enterprise Attestation",
	Long:  "Enable Enterprise Attestation on a security key, if supported. This is often required in enterprise environments for device-specific attestation.",
	Example: `  skm config enterprise-attestation
  skm config enterprise-attestation --device-path /dev/hidraw0 --pin 123456`,
	RunE: enterpriseAttestationHandler,
}

var (
	enterpriseAttestationDevicePath string
	enterpriseAttestationPin        string
)

func init() {
	enterpriseAttestationCMD.Flags().
		StringVarP(&enterpriseAttestationDevicePath, "device-path", "d", "", "Path to the security key device")
	_ = enterpriseAttestationCMD.RegisterFlagCompletionFunc("device-path", completion.CompleteDevicePath)
	enterpriseAttestationCMD.Flags().StringVarP(&enterpriseAttestationPin, "pin", "p", "", "PIN for the security key")
	rootCMD.AddCommand(&enterpriseAttestationCMD)
}

func enterpriseAttestationHandler(cmd *cobra.Command, _ []string) error {
	var selectedDev *fido2.DeviceDescriptor

	if enterpriseAttestationDevicePath != "" {
		devs, err := fido2.Enumerate()
		if err != nil {
			return err
		}
		for _, dev := range devs {
			if dev.Path == enterpriseAttestationDevicePath {
				selectedDev = &dev
				break
			}
		}
		if selectedDev == nil {
			return fmt.Errorf("device not found at path: %s", enterpriseAttestationDevicePath)
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

	pin := enterpriseAttestationPin
	if pin == "" {
		retries, _, _ := dev.GetPINRetries()
		pin, err = prompts.NewPinPrompt().WithRetries(retries).Run()
		if err != nil {
			return err
		}
	}

	token, err := dev.GetPinUvAuthTokenUsingPIN(pin, ctap2.PermissionAuthenticatorConfiguration, "")
	if err != nil {
		return err
	}

	err = dev.EnableEnterpriseAttestation(token)
	if err != nil {
		return err
	}

	cmd.Println("Enterprise Attestation enabled successfully.")
	return nil
}
