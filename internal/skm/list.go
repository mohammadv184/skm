package skm

import (
	"github.com/mohammadv184/go-fido2"
	"github.com/mohammadv184/skm/internal/ui/views"
	"github.com/spf13/cobra"
)

var listCMD = cobra.Command{
	Use:     "list",
	Short:   "List connected security keys",
	Long:    `Enumerate and display all FIDO2 security keys currently connected to the system. It shows the device path, product name, manufacturer, and serial number.`,
	Example: `  skm list`,
	RunE:    listHandler,
}

func init() {
	rootCMD.AddCommand(&listCMD)
}

func listHandler(cmd *cobra.Command, _ []string) error {
	devs, err := fido2.Enumerate()
	if err != nil {
		return err
	}

	if len(devs) == 0 {
		cmd.Println("No security keys found.")
		return nil
	}

	t := views.NewDevicesListView().WithDevices(devs...)

	cmd.Println(t.Render())
	return nil
}
