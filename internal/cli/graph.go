package cli

import (
	"encoding/json"
	"fmt"

	"github.com/kernelstub/cognitor/internal/graph"
	"github.com/kernelstub/cognitor/internal/store"
	"github.com/spf13/cobra"
)

func newGraphCommand(streams ioStreams, configPath *string) *cobra.Command {
	var dbPath, query string
	cmd := &cobra.Command{
		Use:   "graph",
		Short: "Run graph-oriented triage queries",
		RunE: func(cmd *cobra.Command, args []string) error {
			if _, err := loadConfig(*configPath); err != nil {
				return err
			}
			db, err := store.Open(dbPath)
			if err != nil {
				return err
			}
			defer db.Close()
			findings, err := db.LoadFindings(cmd.Context())
			if err != nil {
				return err
			}
			var value any
			switch query {
			case "newly-protected":
				value = graph.FunctionsNewlyProtected(findings)
			case "validation-additions":
				value = graph.BinariesWithValidationAdditions(findings)
			case "sibling-potential":
				value = graph.SiblingPotential(findings)
			default:
				return fmt.Errorf("unknown graph query %q", query)
			}
			data, err := json.MarshalIndent(value, "", "  ")
			if err != nil {
				return err
			}
			_, err = fmt.Fprintln(streams.stdout, string(data))
			return err
		},
	}
	cmd.Flags().StringVar(&dbPath, "db", "", "findings database")
	cmd.Flags().StringVar(&query, "query", "newly-protected", "query: newly-protected, validation-additions, sibling-potential")
	_ = cmd.MarkFlagRequired("db")
	return cmd
}
