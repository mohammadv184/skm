package skm

import (
	"fmt"
	"os"

	mcobra "github.com/muesli/mango-cobra"
	"github.com/muesli/roff"
	"github.com/spf13/cobra"
)

var mandocCMD = &cobra.Command{
	Use:    "mandoc",
	Short:  "Generate man pages for SKM",
	Hidden: true,
	RunE: func(cmd *cobra.Command, _ []string) error {
		manPage, err := mcobra.NewManPage(1, cmd.Root())
		if err != nil {
			return err
		}

		_, err = fmt.Fprint(os.Stdout, manPage.Build(roff.NewDocument()))
		return err
	},
}

func init() {
	rootCMD.AddCommand(mandocCMD)
}
