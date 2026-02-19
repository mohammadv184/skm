package pin

import (
	"errors"
	"fmt"

	"github.com/mohammadv184/go-fido2"
	"github.com/mohammadv184/skm/internal/skm/completion"
	"github.com/mohammadv184/skm/internal/ui/prompts"
	"github.com/spf13/cobra"
)

var setCMD = cobra.Command{
	Use:   "set",
	Short: "Set a new PIN on a security key",
	Long:  "Set a new PIN on a security key that currently doesn't have one.",
	Example: `  skm pin set
  skm pin set --device-path /dev/hidraw0 --pin 123456`,
	RunE: setHandler,
}

var (
	setDevicePath string
	setPin        string
)

func init() {
	setCMD.Flags().StringVarP(&setDevicePath, "device-path", "d", "", "Path to the security key device")
	_ = setCMD.RegisterFlagCompletionFunc("device-path", completion.CompleteDevicePath)
	setCMD.Flags().StringVarP(&setPin, "pin", "p", "", "New PIN for the security key")
	rootCMD.AddCommand(&setCMD)
}

func setHandler(cmd *cobra.Command, _ []string) error {
	var selectedDev *fido2.DeviceDescriptor

	if setDevicePath != "" {
		devs, err := fido2.Enumerate()
		if err != nil {
			return err
		}
		for _, dev := range devs {
			if dev.Path == setDevicePath {
				selectedDev = &dev
				break
			}
		}
		if selectedDev == nil {
			return fmt.Errorf("device not found at path: %s", setDevicePath)
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

	info := dev.Info()
	if isPINSet, ok := info.Options["clientPin"]; ok && isPINSet {
		return errors.New("PIN is already set, use 'skm pin change' to update it")
	}

	newPIN := setPin
	if newPIN == "" {
		newPIN, err = prompts.NewPinPrompt().
			WithTitle("Enter New PIN").
			WithValidation(func(s string) error {
				if len(s) < 4 {
					return errors.New("PIN must be at least 4 characters long")
				}
				return nil
			}).Run()
		if err != nil {
			return err
		}

		_, err = prompts.NewPinPrompt().
			WithTitle("Confirm New PIN").
			WithValidation(func(s string) error {
				if s != newPIN {
					return errors.New("PINs do not match")
				}
				return nil
			}).Run()
		if err != nil {
			return err
		}
	} else {
		if len(newPIN) < 4 {
			return errors.New("PIN must be at least 4 characters long")
		}
	}

	err = dev.SetPIN(newPIN)
	if err != nil {
		return err
	}

	cmd.Println("PIN set successfully.")
	return nil
}
