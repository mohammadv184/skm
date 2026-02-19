package skm

import (
	"context"
	"os"

	"github.com/mohammadv184/skm/internal/skm/config"
	"github.com/mohammadv184/skm/internal/skm/creds"
	"github.com/mohammadv184/skm/internal/skm/pin"
	"github.com/spf13/cobra"
)

var rootCMD = cobra.Command{
	Use:   "skm",
	Short: "SKM: Security Key Manager",
	Long: `SKM (Security Key Manager) is a powerful CLI tool designed for managing FIDO2 security keys.
It provides functionalities to list connected devices, retrieve detailed device information,
and manage resident credentials stored on the keys.`,
	Example: `  skm list
  skm info
  skm creds list
  skm pin set
  skm config always-uv`,
}

func init() {
	creds.Init(&rootCMD)
	pin.Init(&rootCMD)
	config.Init(&rootCMD)
}

// Main is the entry point of the SKM CLI.
func Main(args []string) {
	rootCMD.SetArgs(args[1:])
	rootCMD.SetOut(os.Stdout)
	rootCMD.SetErr(os.Stderr)

	rootCMD.CompletionOptions.HiddenDefaultCmd = true
	rootCMD.DisableAutoGenTag = true
	rootCMD.SetHelpCommand(&cobra.Command{Hidden: true})

	_ = rootCMD.ExecuteContext(context.Background())
}
