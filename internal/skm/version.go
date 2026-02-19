package skm

import (
	"github.com/mohammadv184/skm/configs"
	"github.com/spf13/cobra"
)

var versionCMD = cobra.Command{
	Use:   "version",
	Short: "Show SKM version information",
	Long:  `Display the current version of SKM, along with the build date and commit hash.`,
	RunE:  versionHandler,
}

func init() {
	rootCMD.AddCommand(&versionCMD)
}

func versionHandler(cmd *cobra.Command, _ []string) error {
	cmd.Printf(
		"SKM version: %s %s (%s/%s - %s) built at %s\n",
		configs.Version,
		configs.Commit,
		configs.OS,
		configs.Arch,
		configs.Distribution,
		configs.BuildDate,
	)
	return nil
}
