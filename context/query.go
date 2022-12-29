package context

import (
	"context"

	"github.com/cosmos/cosmos-sdk/client"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/query"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	hubtypes "github.com/sentinel-official/hub/types"
	deposittypes "github.com/sentinel-official/hub/x/deposit/types"
	nodetypes "github.com/sentinel-official/hub/x/node/types"
	plantypes "github.com/sentinel-official/hub/x/plan/types"
	providertypes "github.com/sentinel-official/hub/x/provider/types"
	sessiontypes "github.com/sentinel-official/hub/x/session/types"
	subscriptiontypes "github.com/sentinel-official/hub/x/subscription/types"
	"github.com/spf13/cobra"
	rpchttp "github.com/tendermint/tendermint/rpc/client/http"

	clitypes "github.com/sentinel-official/cli-client/types"
)

type QueryContext struct {
	client.Context
}

func NewQueryContextFromCmd(cmd *cobra.Command) (ctx QueryContext, err error) {
	ctx.Context = client.GetClientContextFromCmd(cmd)

	ctx.NodeURI, err = cmd.Flags().GetString(clitypes.FlagRPCAddress)
	if err != nil {
		return ctx, err
	}

	ctx.Client, err = rpchttp.New(ctx.NodeURI, "/websocket")
	if err != nil {
		return ctx, err
	}

	return ctx, nil
}

func (c QueryContext) WithContext(v client.Context) QueryContext {
	c.Context = v
	return c
}

func (c *QueryContext) QueryAccount(accAddr sdk.AccAddress) (authtypes.AccountI, error) {
	var (
		account   authtypes.AccountI
		qc        = authtypes.NewQueryClient(c)
		resp, err = qc.Account(
			context.Background(),
			&authtypes.QueryAccountRequest{
				Address: accAddr.String(),
			},
		)
	)

	if err != nil {
		return nil, err
	}
	if err := c.InterfaceRegistry.UnpackAny(resp.Account, &account); err != nil {
		return nil, err
	}

	return account, nil
}

func (c *QueryContext) QueryDeposit(accAddr sdk.AccAddress) (*deposittypes.Deposit, error) {
	var (
		qsc       = deposittypes.NewQueryServiceClient(c)
		resp, err = qsc.QueryDeposit(
			context.Background(),
			deposittypes.NewQueryDepositRequest(
				accAddr,
			),
		)
	)

	if err != nil {
		return nil, err
	}

	return &resp.Deposit, nil
}

func (c *QueryContext) QueryDeposits(pagination *query.PageRequest) (deposittypes.Deposits, error) {
	var (
		qsc       = deposittypes.NewQueryServiceClient(c)
		resp, err = qsc.QueryDeposits(
			context.Background(),
			deposittypes.NewQueryDepositsRequest(
				pagination,
			),
		)
	)

	if err != nil {
		return nil, err
	}

	return resp.Deposits, nil
}

func (c *QueryContext) QueryNode(nodeAddr hubtypes.NodeAddress) (*nodetypes.Node, error) {
	var (
		qsc       = nodetypes.NewQueryServiceClient(c)
		resp, err = qsc.QueryNode(
			context.Background(),
			nodetypes.NewQueryNodeRequest(
				nodeAddr,
			),
		)
	)

	if err != nil {
		return nil, err
	}

	return &resp.Node, nil
}

func (c *QueryContext) QueryNodes(status hubtypes.Status, pagination *query.PageRequest) (nodetypes.Nodes, error) {
	var (
		qsc       = nodetypes.NewQueryServiceClient(c)
		resp, err = qsc.QueryNodes(
			context.Background(),
			nodetypes.NewQueryNodesRequest(
				status,
				pagination,
			),
		)
	)

	if err != nil {
		return nil, err
	}

	return resp.Nodes, nil
}

func (c *QueryContext) QueryNodesForProvider(provAddr hubtypes.ProvAddress, status hubtypes.Status, pagination *query.PageRequest) (nodetypes.Nodes, error) {
	var (
		qsc       = nodetypes.NewQueryServiceClient(c)
		resp, err = qsc.QueryNodesForProvider(
			context.Background(),
			nodetypes.NewQueryNodesForProviderRequest(
				provAddr,
				status,
				pagination,
			),
		)
	)

	if err != nil {
		return nil, err
	}

	return resp.Nodes, nil
}

func (c *QueryContext) QueryPlan(id uint64) (*plantypes.Plan, error) {
	var (
		qsc       = plantypes.NewQueryServiceClient(c)
		resp, err = qsc.QueryPlan(
			context.Background(),
			plantypes.NewQueryPlanRequest(
				id,
			),
		)
	)
	if err != nil {
		return nil, err
	}

	return &resp.Plan, nil
}

func (c *QueryContext) QueryPlans(status hubtypes.Status, pagination *query.PageRequest) (plantypes.Plans, error) {
	var (
		qsc       = plantypes.NewQueryServiceClient(c)
		resp, err = qsc.QueryPlans(
			context.Background(),
			plantypes.NewQueryPlansRequest(
				status,
				pagination,
			),
		)
	)

	if err != nil {
		return nil, err
	}

	return resp.Plans, nil
}

func (c *QueryContext) QueryPlansForProvider(provAddr hubtypes.ProvAddress, status hubtypes.Status, pagination *query.PageRequest) (plantypes.Plans, error) {
	var (
		qsc       = plantypes.NewQueryServiceClient(c)
		resp, err = qsc.QueryPlansForProvider(
			context.Background(),
			plantypes.NewQueryPlansForProviderRequest(
				provAddr,
				status,
				pagination,
			),
		)
	)

	if err != nil {
		return nil, err
	}

	return resp.Plans, nil
}

func (c *QueryContext) QueryProvider(provAddr hubtypes.ProvAddress) (*providertypes.Provider, error) {
	var (
		qsc       = providertypes.NewQueryServiceClient(c)
		resp, err = qsc.QueryProvider(
			context.Background(),
			providertypes.NewQueryProviderRequest(
				provAddr,
			),
		)
	)

	if err != nil {
		return nil, err
	}

	return &resp.Provider, nil
}

func (c *QueryContext) QueryProviders(pagination *query.PageRequest) (providertypes.Providers, error) {
	var (
		qsc       = providertypes.NewQueryServiceClient(c)
		resp, err = qsc.QueryProviders(
			context.Background(),
			providertypes.NewQueryProvidersRequest(
				pagination,
			),
		)
	)

	if err != nil {
		return nil, err
	}

	return resp.Providers, nil
}

func (c *QueryContext) QuerySession(id uint64) (*sessiontypes.Session, error) {
	var (
		qsc       = sessiontypes.NewQueryServiceClient(c)
		resp, err = qsc.QuerySession(
			context.Background(),
			sessiontypes.NewQuerySessionRequest(
				id,
			),
		)
	)

	if err != nil {
		return nil, err
	}

	return &resp.Session, nil
}

func (c *QueryContext) QuerySessions(pagination *query.PageRequest) (sessiontypes.Sessions, error) {
	var (
		qsc       = sessiontypes.NewQueryServiceClient(c)
		resp, err = qsc.QuerySessions(
			context.Background(),
			sessiontypes.NewQuerySessionsRequest(
				pagination,
			),
		)
	)

	if err != nil {
		return nil, err
	}

	return resp.Sessions, nil
}

func (c *QueryContext) QuerySessionsForAddress(accAddr sdk.AccAddress, status hubtypes.Status, pagination *query.PageRequest) (sessiontypes.Sessions, error) {
	var (
		qsc       = sessiontypes.NewQueryServiceClient(c)
		resp, err = qsc.QuerySessionsForAddress(
			context.Background(),
			sessiontypes.NewQuerySessionsForAddressRequest(
				accAddr,
				status,
				pagination,
			),
		)
	)

	if err != nil {
		return nil, err
	}

	return resp.Sessions, nil
}

func (c *QueryContext) QuerySubscription(id uint64) (*subscriptiontypes.Subscription, error) {
	var (
		qsc       = subscriptiontypes.NewQueryServiceClient(c)
		resp, err = qsc.QuerySubscription(
			context.Background(),
			subscriptiontypes.NewQuerySubscriptionRequest(
				id,
			),
		)
	)

	if err != nil {
		return nil, err
	}

	return &resp.Subscription, nil
}

func (c *QueryContext) QuerySubscriptions(pagination *query.PageRequest) (subscriptiontypes.Subscriptions, error) {
	var (
		qsc       = subscriptiontypes.NewQueryServiceClient(c)
		resp, err = qsc.QuerySubscriptions(
			context.Background(),
			subscriptiontypes.NewQuerySubscriptionsRequest(
				pagination,
			),
		)
	)

	if err != nil {
		return nil, err
	}

	return resp.Subscriptions, nil
}

func (c *QueryContext) QuerySubscriptionsForAddress(accAddr sdk.AccAddress, status hubtypes.Status, pagination *query.PageRequest) (subscriptiontypes.Subscriptions, error) {
	var (
		qsc       = subscriptiontypes.NewQueryServiceClient(c)
		resp, err = qsc.QuerySubscriptionsForAddress(
			context.Background(),
			subscriptiontypes.NewQuerySubscriptionsForAddressRequest(
				accAddr,
				status,
				pagination,
			),
		)
	)

	if err != nil {
		return nil, err
	}

	return resp.Subscriptions, nil
}

func (c *QueryContext) QueryQuota(id uint64, accAddr sdk.AccAddress) (*subscriptiontypes.Quota, error) {
	var (
		qsc       = subscriptiontypes.NewQueryServiceClient(c)
		resp, err = qsc.QueryQuota(
			context.Background(),
			subscriptiontypes.NewQueryQuotaRequest(
				id,
				accAddr,
			),
		)
	)

	if err != nil {
		return nil, err
	}

	return &resp.Quota, nil
}

func (c *QueryContext) QueryQuotas(id uint64, pagination *query.PageRequest) (subscriptiontypes.Quotas, error) {
	var (
		qsc       = subscriptiontypes.NewQueryServiceClient(c)
		resp, err = qsc.QueryQuotas(
			context.Background(),
			subscriptiontypes.NewQueryQuotasRequest(
				id,
				pagination,
			),
		)
	)

	if err != nil {
		return nil, err
	}

	return resp.Quotas, nil
}

func (c *QueryContext) QueryActiveSession(accAddr sdk.AccAddress) (*sessiontypes.Session, error) {
	var (
		qsc       = sessiontypes.NewQueryServiceClient(c)
		resp, err = qsc.QuerySessionsForAddress(
			context.Background(),
			sessiontypes.NewQuerySessionsForAddressRequest(
				accAddr,
				hubtypes.StatusActive,
				&query.PageRequest{
					Limit: 1,
				},
			),
		)
	)

	if err != nil {
		return nil, err
	}
	if len(resp.Sessions) == 0 {
		return nil, nil
	}

	return &resp.Sessions[0], nil
}
