package pin

import (
	"errors"
	"fmt"

	"github.com/mohammadv184/go-fido2"
	"github.com/mohammadv184/skm/internal/skm/completion"
	"github.com/mohammadv184/skm/internal/ui/prompts"
	"github.com/spf13/cobra"
)

var changeCMD = cobra.Command{
	Use:   "change",
	Short: "Change an existing PIN on a security key",
	Long:  "Change an existing PIN on a security key. You will be prompted for your current PIN and then your new PIN.",
	Example: `  skm pin change
  skm pin change --device-path /dev/hidraw0 --pin 123456 --new-pin 654321`,
	RunE: changeHandler,
}

var (
	changeDevicePath string
	changePin        string
	changeNewPin     string
)

func init() {
	changeCMD.Flags().StringVarP(&changeDevicePath, "device-path", "d", "", "Path to the security key device")
	_ = changeCMD.RegisterFlagCompletionFunc("device-path", completion.CompleteDevicePath)
	changeCMD.Flags().StringVarP(&changePin, "pin", "p", "", "Current PIN for the security key")
	changeCMD.Flags().StringVarP(&changeNewPin, "new-pin", "n", "", "New PIN for the security key")
	rootCMD.AddCommand(&changeCMD)
}

func changeHandler(cmd *cobra.Command, _ []string) error {
	var selectedDev *fido2.DeviceDescriptor

	if changeDevicePath != "" {
		devs, err := fido2.Enumerate()
		if err != nil {
			return err
		}
		for _, dev := range devs {
			if dev.Path == changeDevicePath {
				selectedDev = &dev
				break
			}
		}
		if selectedDev == nil {
			return fmt.Errorf("device not found at path: %s", changeDevicePath)
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
	if isPINSet, ok := info.Options["clientPin"]; !ok || !isPINSet {
		return errors.New("PIN is not set, use 'skm pin set' to set it")
	}

	currentPIN := changePin
	if currentPIN == "" {
		retries, _, _ := dev.GetPINRetries()
		currentPIN, err = prompts.NewPinPrompt().
			WithTitle("Enter Current PIN").
			WithRetries(retries).
			Run()
		if err != nil {
			return err
		}
	}

	newPIN := changeNewPin
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

	err = dev.ChangePIN(currentPIN, newPIN)
	if err != nil {
		return err
	}

	cmd.Println("PIN changed successfully.")
	return nil
}
