package cmd

import (
	"context"
	clienttypes "github.com/sentinel-official/cli-client/types"
	"github.com/sentinel-official/cli-client/utils"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	hubtypes "github.com/sentinel-official/hub/types"
	providertypes "github.com/sentinel-official/hub/x/provider/types"
	"github.com/spf13/cobra"

	"github.com/sentinel-official/cli-client/x/provider/types"
)

var (
	header = []string{
		"Name",
		"Address",
		"Identity",
		"Website",
	}
)

func QueryProvider() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "provider [address]",
		Short: "Query a provider",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			outputFormat, err := cmd.Flags().GetString(clienttypes.FlagOutput)
			if err != nil {
				return err
			}

			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			address, err := hubtypes.ProvAddressFromBech32(args[0])
			if err != nil {
				return err
			}

			var (
				qsc = providertypes.NewQueryServiceClient(ctx)
			)

			result, err := qsc.QueryProvider(
				context.Background(),
				providertypes.NewQueryProviderRequest(address),
			)
			if err != nil {
				return err
			}

			var (
				item       = types.NewProviderFromRaw(&result.Provider)
				outputRows [][]string
			)

			outputRows = append(
				outputRows,
				[]string{
					item.Name,
					item.Address,
					item.Identity,
					item.Website,
				},
			)

			err = utils.WriteOutput(header, outputRows, outputFormat)
			if err != nil {
				return err
			}
			return nil
		},
	}

	flags.AddQueryFlagsToCmd(cmd)

	return cmd
}

func QueryProviders() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "providers",
		Short: "Query providers",
		RunE: func(cmd *cobra.Command, args []string) error {
			outputFormat, err := cmd.Flags().GetString(clienttypes.FlagOutput)
			if err != nil {
				return err
			}

			ctx, err := client.GetClientQueryContext(cmd)
			if err != nil {
				return err
			}

			pagination, err := client.ReadPageRequest(cmd.Flags())
			if err != nil {
				return err
			}

			var (
				qsc = providertypes.NewQueryServiceClient(ctx)
			)

			result, err := qsc.QueryProviders(
				context.Background(),
				providertypes.NewQueryProvidersRequest(pagination),
			)
			if err != nil {
				return err
			}

			var (
				items      = types.NewProvidersFromRaw(result.Providers)
				outputRows [][]string
			)

			for i := 0; i < len(items); i++ {
				outputRows = append(
					outputRows,
					[]string{
						items[i].Name,
						items[i].Address,
						items[i].Identity,
						items[i].Website,
					},
				)
			}

			err = utils.WriteOutput(header, outputRows, outputFormat)
			if err != nil {
				return err
			}
			return nil
		},
	}

	flags.AddQueryFlagsToCmd(cmd)
	flags.AddPaginationFlagsToCmd(cmd, "providers")

	return cmd
}
