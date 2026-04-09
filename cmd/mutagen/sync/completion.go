package sync

import (
	"context"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/mutagen-io/mutagen/cmd/mutagen/daemon"

	"github.com/mutagen-io/mutagen/pkg/selection"
	synchronizationsvc "github.com/mutagen-io/mutagen/pkg/service/synchronization"
)

// completeSessions provides tab-completion candidates for sync session
// identifiers and names by querying the daemon's session list.
func completeSessions(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	// Connect to the daemon without autostarting or enforcing a version match
	// to keep completion fast and side-effect-free.
	daemonConnection, err := daemon.Connect(false, false)
	if err != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}
	defer daemonConnection.Close()

	// List all sessions.
	synchronizationService := synchronizationsvc.NewSynchronizationClient(daemonConnection)
	response, err := synchronizationService.List(
		context.Background(),
		&synchronizationsvc.ListRequest{
			Selection: &selection.Selection{All: true},
		},
	)
	if err != nil || response.EnsureValid() != nil {
		return nil, cobra.ShellCompDirectiveNoFileComp
	}

	// Build completion candidates from each session's identifier and name.
	var completions []string
	for _, state := range response.SessionStates {
		session := state.Session

		description := fmt.Sprintf("%s <-> %s",
			session.Alpha.Format(""),
			session.Beta.Format(""),
		)

		completions = append(completions,
			fmt.Sprintf("%s\t%s", session.Identifier, description),
		)

		if session.Name != "" {
			completions = append(completions,
				fmt.Sprintf("%s\t%s (%s)", session.Name, description, session.Identifier),
			)
		}
	}

	return completions, cobra.ShellCompDirectiveNoFileComp
}
