package skm

import (
	"fmt"
	"strings"

	"github.com/mohammadv184/go-fido2"
	"github.com/mohammadv184/skm/internal/skm/completion"
	"github.com/mohammadv184/skm/internal/ui/prompts"
	"github.com/mohammadv184/skm/internal/ui/views"
	"github.com/spf13/cobra"
)

var infoCMD = cobra.Command{
	Use:   "info",
	Short: "Show information about a security key",
	Long:  `Display detailed technical information about a selected security key, including AAGUID, supported versions, extensions, and protocol options.`,
	Example: `  skm info
  skm info --all
  skm info --device-path /dev/hidraw0`,
	RunE: infoHandler,
}

var (
	infoDevicePath string
	infoAll        bool
)

func init() {
	infoCMD.Flags().StringVarP(&infoDevicePath, "device-path", "d", "", "Path to the security key device")
	_ = infoCMD.RegisterFlagCompletionFunc("device-path", completion.CompleteDevicePath)
	infoCMD.Flags().BoolVarP(&infoAll, "all", "a", false, "Show information for all connected security keys")
	rootCMD.AddCommand(&infoCMD)
}

func infoHandler(cmd *cobra.Command, _ []string) error {
	devs, err := fido2.Enumerate()
	if err != nil {
		return err
	}

	if len(devs) == 0 {
		cmd.Println("No security keys found.")
		return nil
	}

	var selectedDevs []fido2.DeviceDescriptor

	if infoAll {
		selectedDevs = devs
	} else if infoDevicePath != "" {
		for _, dev := range devs {
			if dev.Path == infoDevicePath {
				selectedDevs = append(selectedDevs, dev)
				break
			}
		}
		if len(selectedDevs) == 0 {
			return fmt.Errorf("device not found at path: %s", infoDevicePath)
		}
	} else {
		cmd.Println("Select a security key to view its information:")

		selected, err := prompts.NewDeviceSelectPrompt().WithDevices(devs...).Run()
		if err != nil {
			return err
		}
		selectedDevs = append(selectedDevs, *selected)
	}

	for i, sd := range selectedDevs {
		if i > 0 {
			cmd.Println("\n" + strings.Repeat("-", 40) + "\n")
		}

		dev, err := fido2.Open(sd)
		if err != nil {
			cmd.Printf("Error opening device %s: %v\n", sd.Path, err)
			continue
		}

		info := dev.Info()

		pinRetries, _, _ := dev.GetPINRetries()
		uvRetries, err := dev.GetUVRetries()
		hasUV := err == nil

		v := views.NewDeviceInfoView(&sd, info).WithRetries(pinRetries, uvRetries, hasUV)
		cmd.Println(v.Render())
		_ = dev.Close()
	}

	return nil
}
