package cmd

import (
	"github.com/cosmos/cosmos-sdk/client"
	"github.com/spf13/cobra"

	depositcmd "github.com/sentinel-official/cli-client/x/deposit/client/cmd"
	nodecmd "github.com/sentinel-official/cli-client/x/node/client/cmd"
	plancmd "github.com/sentinel-official/cli-client/x/plan/client/cmd"
	providercmd "github.com/sentinel-official/cli-client/x/provider/client/cmd"
	sessioncmd "github.com/sentinel-official/cli-client/x/session/client/cmd"
	subscriptioncmd "github.com/sentinel-official/cli-client/x/subscription/client/cmd"
)

func QueryCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:                        "query",
		Aliases:                    []string{"q"},
		Short:                      "Query subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	cmd.AddCommand(depositcmd.QueryDeposit())
	cmd.AddCommand(depositcmd.QueryDeposits())
	cmd.AddCommand(nodecmd.QueryNode())
	cmd.AddCommand(nodecmd.QueryNodes())
	cmd.AddCommand(providercmd.QueryProvider())
	cmd.AddCommand(providercmd.QueryProviders())
	cmd.AddCommand(plancmd.QueryPlan())
	cmd.AddCommand(plancmd.QueryPlans())
	cmd.AddCommand(subscriptioncmd.QuerySubscription())
	cmd.AddCommand(subscriptioncmd.QuerySubscriptions())
	cmd.AddCommand(subscriptioncmd.QueryQuota())
	cmd.AddCommand(subscriptioncmd.QueryQuotas())
	cmd.AddCommand(sessioncmd.QuerySession())
	cmd.AddCommand(sessioncmd.QuerySessions())

	return cmd
}
