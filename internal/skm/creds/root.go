package creds

import "github.com/spf13/cobra"

var rootCMD = cobra.Command{
	Use:     "creds",
	Aliases: []string{"c", "credential", "credentials"},
	Short:   "Manage credentials stored on security keys",
	Long:    `Provide a set of subcommands to manage resident credentials stored directly on your FIDO2 security keys. This includes listing all credentials and deleting specific ones.`,
	Example: `  skm creds list
  skm creds delete`,
}

// Init initializes the creds command and its subcommands.
func Init(skmRoot *cobra.Command) {
	skmRoot.AddCommand(&rootCMD)
}
