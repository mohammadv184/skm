package skm

import (
	"errors"
	"fmt"

	"github.com/mohammadv184/go-fido2"
	"github.com/mohammadv184/skm/internal/skm/completion"
	"github.com/mohammadv184/skm/internal/ui/prompts"
	"github.com/spf13/cobra"
)

var resetCMD = cobra.Command{
	Use:   "reset",
	Short: "Factory reset a security key",
	Long:  "Completely wipe all credentials and reset the PIN of a security key. This action is IRREVERSIBLE. Most keys require physical touch within a short window after power-on to perform a reset.",
	Example: `  skm reset
  skm reset --device-path /dev/hidraw0 --yes`,
	RunE: resetHandler,
}

var (
	resetDevicePath string
	resetYes        bool
)

func init() {
	resetCMD.Flags().StringVarP(&resetDevicePath, "device-path", "d", "", "Path to the security key device")
	_ = resetCMD.RegisterFlagCompletionFunc("device-path", completion.CompleteDevicePath)
	resetCMD.Flags().BoolVarP(&resetYes, "yes", "y", false, "Confirm reset without prompting")
	rootCMD.AddCommand(&resetCMD)
}

func resetHandler(cmd *cobra.Command, _ []string) error {
	var selectedDev *fido2.DeviceDescriptor

	if resetDevicePath != "" {
		devs, err := fido2.Enumerate()
		if err != nil {
			return err
		}
		for _, dev := range devs {
			if dev.Path == resetDevicePath {
				selectedDev = &dev
				break
			}
		}
		if selectedDev == nil {
			return fmt.Errorf("device not found at path: %s", resetDevicePath)
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

	confirm := resetYes
	if !confirm {
		var err error
		confirm, err = prompts.NewConfirmPrompt(
			"Are you absolutely sure?",
			"This will PERMANENTLY delete all credentials and reset the PIN. Most keys require physical touch after this command is sent.",
		).Run()
		if err != nil {
			return err
		}
	}

	if !confirm {
		cmd.Println("Reset canceled.")
		return nil
	}

	cmd.Println("Performing reset. Please touch your security key if it starts blinking.")
	err = dev.Reset()
	if err != nil {
		return err
	}

	cmd.Println("Security key reset successfully.")
	return nil
}
