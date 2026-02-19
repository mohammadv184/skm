package pin

import (
	"github.com/spf13/cobra"
)

var rootCMD = cobra.Command{
	Use:     "pin",
	Aliases: []string{"p"},
	Short:   "Manage security key PIN",
	Long:    "Commands for setting, changing, and checking PIN retries on security keys.",
	Example: `  skm pin set
  skm pin change
  skm pin retries`,
}

// Init initializes the pin command and its subcommands.
func Init(skmRoot *cobra.Command) {
	skmRoot.AddCommand(&rootCMD)
}
