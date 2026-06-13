package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

func newRulesCommand(streams ioStreams) *cobra.Command {
	rules := []string{
		"added-access-check",
		"added-bounds-check",
		"added-lifetime-reference",
		"added-token-check",
		"changed-ioctl-validation",
		"added-rpc-auth-validation",
		"changed-com-permissions",
		"added-alpc-validation",
		"added-registry-hardening",
		"added-service-hardening",
	}
	return &cobra.Command{
		Use:   "rules",
		Short: "List semantic hardening rules",
		RunE: func(cmd *cobra.Command, args []string) error {
			for _, rule := range rules {
				_, _ = fmt.Fprintln(streams.stdout, rule)
			}
			return nil
		},
	}
}
