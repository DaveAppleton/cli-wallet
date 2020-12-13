package client

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
	"syscall"

	smWallet "github.com/DaveAppleton/smWallet"
	xdr "github.com/davecgh/go-xdr/xdr2"
	"github.com/spacemeshos/CLIWallet/common"
	"github.com/spacemeshos/CLIWallet/log"
	pb "github.com/spacemeshos/api/release/go/spacemesh/v1"
	"github.com/spacemeshos/ed25519"
	gosmtypes "github.com/spacemeshos/go-spacemesh/common/types"
	"golang.org/x/crypto/ssh/terminal"
)

const accountsFileName = "accounts.json"

// WalletBackend wallet holder
type WalletBackend struct {
	*gRPCClient // Embedded interface
	//common.Store
	//accountsFilePath string
	wallet *smWallet.Wallet
	//currentAccount   *common.LocalAccount
}

func getPassword() (string, error) {
	fmt.Print("Enter password: ")
	bytePassword, err := terminal.ReadPassword(int(syscall.Stdin)) // no history
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(bytePassword)), nil

}

// OpenWalletBackend = open an existing wallet
func OpenWalletBackend(wallet string, grpcServer string, secureConnection bool) (wbx *WalletBackend, err error) {
	var wbe WalletBackend
	wbx = nil
	if wbe.wallet, err = smWallet.LoadWallet(wallet); err != nil {
		return
	}
	password, err := getPassword()
	if err != nil {
		return
	}
	fmt.Println("loading...")
	if err = wbe.wallet.Unlock(password); err != nil {
		return
	}
	ne, err := wbe.wallet.GetNumberOfAccounts()
	if err != nil {
		return nil, err
	}
	fmt.Println(wbe.wallet.Meta.DisplayName, "successfully opened with ", ne, "accounts")
	wbe.gRPCClient = newGRPCClient(grpcServer, secureConnection)
	if err = wbe.gRPCClient.Connect(); err != nil {
		// failed to connect to grpc server
		log.Error("failed to connect to the grpc server: %s", err)
		return
	}
	return &wbe, nil
}

// NewWalletBackend set up a wallet -
func NewWalletBackend(walletName string, grpcServer string, secureConnection bool) (wbx *WalletBackend, err error) {
	var wbe WalletBackend
	wbx = nil
	password, err := getPassword()
	if err != nil {
		return
	}
	fmt.Println("ok")
	if wbe.wallet, err = smWallet.NewWallet(walletName, password); err != nil {
		fmt.Println("failur to create wallet", err)
		return
	}
	if err = wbe.wallet.SaveWalletAs("myWallet_"); err != nil {

	}

	fmt.Println(wbe.wallet.Meta.DisplayName, "successfully created")
	wbe.gRPCClient = newGRPCClient(grpcServer, secureConnection)
	if err = wbe.gRPCClient.Connect(); err != nil {
		// failed to connect to grpc server
		log.Error("failed to connect to the grpc server: %s", err)
		return
	}
	return &wbe, nil
}

// CurrentAccount - get the latest account into cli-wallet format
func (w *WalletBackend) CurrentAccount() (*common.LocalAccount, error) {

	ca, err := w.wallet.CurrentAccount()
	if err != nil {
		return nil, err
	}
	pk, err := ca.PrivateKey()
	if err != nil {
		return nil, err
	}
	return &common.LocalAccount{Name: ca.DisplayName, PrivKey: pk, PubKey: smWallet.PublicKey(pk)}, nil
}

func (w *WalletBackend) CreateAccount(displayName string) (la *common.LocalAccount, err error) {
	pos, err := w.wallet.GenerateNewPair(displayName)
	if err != nil {
		return nil, err
	}
	if err = w.wallet.SetCurrent(pos); err != nil {
		return
	}
	return w.CurrentAccount()
}

func (w *WalletBackend) SetCurrentAccount(accountNumber int) error {
	return w.wallet.SetCurrent(accountNumber)
}

func interfaceToBytes(i interface{}) ([]byte, error) {
	var w bytes.Buffer
	if _, err := xdr.Marshal(&w, &i); err != nil {
		return nil, err
	}
	return w.Bytes(), nil
}

func (w *WalletBackend) StoreAccounts() error {
	return w.wallet.SaveWallet()
}

// Transfer creates a sign coin transaction and submits it
func (w *WalletBackend) Transfer(recipient gosmtypes.Address, nonce, amount, gasPrice, gasLimit uint64, key ed25519.PrivateKey) (*pb.TransactionState, error) {
	tx := common.SerializableSignedTransaction{}
	tx.AccountNonce = nonce
	tx.Amount = amount
	tx.Recipient = recipient
	tx.GasLimit = gasLimit
	tx.Price = gasPrice

	buf, _ := interfaceToBytes(&tx.InnerSerializableSignedTransaction)
	copy(tx.Signature[:], ed25519.Sign2(key, buf))
	b, err := interfaceToBytes(&tx)
	if err != nil {
		return nil, err
	}
	return w.SubmitCoinTransaction(b)
}

func (w *WalletBackend) GetAccount(accountName string) (*common.LocalAccount, error) {
	numberOfAccounts, err := w.wallet.GetNumberOfAccounts()
	if err != nil {
		log.Error("failed to retrieve number of accounts", err)
		return nil, err
	}
	for j := 0; j < numberOfAccounts; j++ {
		dn, err := w.wallet.GetAccountDisplayName(j)
		if err != nil {
			log.Error("failed to retrieve display names", err)
			return nil, err
		}
		if dn == accountName {
			pk, err := w.wallet.GetPrivateKey(j)
			if err != nil {
				log.Error("failed to retrieve private key", err)
				return nil, err
			}
			return &common.LocalAccount{Name: accountName, PrivKey: pk, PubKey: smWallet.PublicKey(pk)}, nil
		}
	}
	err = errors.New("failed to find :" + accountName)
	log.Error(err.Error())
	return nil, err
}

func (w *WalletBackend) ListAccounts() (res []string, err error) {
	numberOfAccounts, err := w.wallet.GetNumberOfAccounts()
	if err != nil {
		log.Error("failed to retrieve number of accounts", err)
		return []string{}, err
	}
	for j := 0; j < numberOfAccounts; j++ {
		dn, err := w.wallet.GetAccountDisplayName(j)
		if err != nil {
			log.Error("failed to retrieve display names", err)
			return []string{}, err
		}
		res = append(res, dn)
	}

	return res, nil

}
