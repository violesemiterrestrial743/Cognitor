package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newVersionCommand(streams ioStreams) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print version",
		RunE: func(cmd *cobra.Command, args []string) error {
			_, _ = fmt.Fprintln(streams.stdout, Version)
			return nil
		},
	}
}
