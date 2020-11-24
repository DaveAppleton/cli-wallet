package client

import (
	"context"
	"github.com/spacemeshos/CLIWallet/common"
	apitypes "github.com/spacemeshos/api/release/go/spacemesh/v1"
	gosmtypes "github.com/spacemeshos/go-spacemesh/common/types"
)

// AccountInfo returns basic account data such as balance and nonce from the global state
func (c *gRPCClient) AccountInfo(address gosmtypes.Address) (*common.AccountState, error) {
	gsc := c.globalStateClient()

	resp, err := gsc.Account(context.Background(), &apitypes.AccountRequest{
		AccountId: &apitypes.AccountId{Address: address.Bytes()}})
	if err != nil {
		return nil, err
	}

	return &common.AccountState{
		Balance: resp.AccountWrapper.StateCurrent.Balance.Value,
		Nonce:   resp.AccountWrapper.StateCurrent.Counter,
	}, nil

}
