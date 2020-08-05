package get

import (
	"context"
	"fmt"

	"github.com/depscloud/api/v1alpha/tracker"
	"github.com/depscloud/cli/internal/writer"

	"github.com/spf13/cobra"
)

func DependentsCommand(
	dependencyClient tracker.DependencyServiceClient,
	writer writer.Writer,
) *cobra.Command {
	req := &tracker.DependencyRequest{}

	cmd := &cobra.Command{
		Use:     "dependents",
		Aliases: []string{"dependent"},
		Short:   "Get the list of modules that depend on the given module",
		Example: "deps get dependents -l go -o github.com -m depscloud/api",
		RunE: func(cmd *cobra.Command, args []string) error {
			if req.Language == "" || req.Organization == "" || req.Module == "" {
				return fmt.Errorf("language, organization, and module must be provided")
			}

			ctx := cmd.Context()

			response, err := dependencyClient.ListDependents(ctx, req)
			if err != nil {
				return err
			}

			for _, dependent := range response.Dependents {
				_ = writer.Write(dependent)
			}

			return nil
		},
	}

	topologyCmd := topologyCommand(writer,
		func(req *tracker.DependencyRequest, ctx context.Context) ([]*tracker.Dependency, error) {
			resp, err := dependencyClient.ListDependents(ctx, req)
			if err != nil {
				return nil, err
			}
			return resp.Dependents, nil
		})

	cmd.AddCommand(topologyCmd)
	addDependencyRequestFlags(cmd, req)

	return cmd
}
