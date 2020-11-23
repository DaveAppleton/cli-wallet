package client

import (
	"context"
	"github.com/spacemeshos/CLIWallet/localtypes"
	apitypes "github.com/spacemeshos/api/release/go/spacemesh/v1"
)

// AccountInfo returns basic account data such as balance and nonce from the global state
func (c *GRPCClient) AccountInfo(address []byte) (*localtypes.AccountState, error) {
	gsc := c.globalStateClient()

	resp, err := gsc.Account(context.Background(), &apitypes.AccountRequest{
		AccountId: &apitypes.AccountId{Address: address}})
	if err != nil {
		return nil, err
	}

	return &localtypes.AccountState{
		Balance: resp.AccountWrapper.StateCurrent.Balance.Value,
		Nonce:   resp.AccountWrapper.StateCurrent.Counter,
	}, nil

}
