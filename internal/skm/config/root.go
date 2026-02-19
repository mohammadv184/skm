package config

import (
	"github.com/spf13/cobra"
)

var rootCMD = cobra.Command{
	Use:     "config",
	Aliases: []string{"cfg"},
	Short:   "Manage security key configuration",
	Long:    "Commands for managing security key features such as Always UV and Enterprise Attestation.",
	Example: `  skm config always-uv
  skm config enterprise-attestation`,
}

// Init initializes the config command and its subcommands.
func Init(skmRoot *cobra.Command) {
	skmRoot.AddCommand(&rootCMD)
}
