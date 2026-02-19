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

var alwaysUVCMD = cobra.Command{
	Use:   "always-uv",
	Short: "Toggle Always User Verification (UV)",
	Long:  "Toggle the 'Always UV' option on a security key. If enabled, the device will always require user verification (e.g. PIN or Biometrics) for all operations.",
	Example: `  skm config always-uv
  skm config always-uv --device-path /dev/hidraw0 --pin 123456`,
	RunE: alwaysUVHandler,
}

var (
	alwaysUVDevicePath string
	alwaysUVPin        string
)

func init() {
	alwaysUVCMD.Flags().StringVarP(&alwaysUVDevicePath, "device-path", "d", "", "Path to the security key device")
	_ = alwaysUVCMD.RegisterFlagCompletionFunc("device-path", completion.CompleteDevicePath)
	alwaysUVCMD.Flags().StringVarP(&alwaysUVPin, "pin", "p", "", "PIN for the security key")
	rootCMD.AddCommand(&alwaysUVCMD)
}

func alwaysUVHandler(cmd *cobra.Command, _ []string) error {
	var selectedDev *fido2.DeviceDescriptor

	if alwaysUVDevicePath != "" {
		devs, err := fido2.Enumerate()
		if err != nil {
			return err
		}
		for _, dev := range devs {
			if dev.Path == alwaysUVDevicePath {
				selectedDev = &dev
				break
			}
		}
		if selectedDev == nil {
			return fmt.Errorf("device not found at path: %s", alwaysUVDevicePath)
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

	pin := alwaysUVPin
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

	err = dev.ToggleAlwaysUV(token)
	if err != nil {
		return err
	}

	cmd.Println("Always UV toggled successfully.")
	return nil
}
