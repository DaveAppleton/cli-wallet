package repl

import (
	"encoding/hex"
	"fmt"
	"github.com/spacemeshos/CLIWallet/localtypes"
	"github.com/spacemeshos/CLIWallet/log"
	apitypes "github.com/spacemeshos/api/release/go/spacemesh/v1"
	"github.com/spacemeshos/ed25519"
	gosmtypes "github.com/spacemeshos/go-spacemesh/common/types"
	"google.golang.org/genproto/googleapis/rpc/status"
	"os"
	"strconv"
	"strings"

	"github.com/c-bata/go-prompt"
)

const (
	prefix      = "$ "
	printPrefix = ">"
)


// TestMode variable used for check if unit test is running
var TestMode = false

type command struct {
	text        string
	description string
	fn          func()
}

type repl struct {
	commands []command
	client   Client
	input    string
}

// Client interface to REPL clients.
type Client interface {
	CreateAccount(alias string) *localtypes.LocalAccount
	CurrentAccount() *localtypes.LocalAccount
	SetCurrentAccount(a *localtypes.LocalAccount)
	AccountInfo(address gosmtypes.Address) (*localtypes.AccountState, error)
	NodeStatus() (*apitypes.NodeStatus, error)
	NodeInfo() (*localtypes.NodeInfo, error)
	Sanity() error
	Transfer(recipient gosmtypes.Address, nonce, amount, gasPrice, gasLimit uint64, key ed25519.PrivateKey) (*apitypes.TransactionState, error)
	ListAccounts() []string
	GetAccount(name string) (*localtypes.LocalAccount, error)
	StoreAccounts() error
	ServerUrl() string
	Smesh(datadir string, space uint, coinbase string) error
	GetMeshTransactions(address gosmtypes.Address, offset uint32, maxResults uint32) ([]*apitypes.Transaction, uint32, error)
	GetMeshActivations(address gosmtypes.Address, offset uint32, maxResults uint32) ([]*apitypes.Activation, uint32, error)
	SetCoinbase(coinbase gosmtypes.Address) (*status.Status, error)
	DebugAllAccounts() ([]*apitypes.Account, error)

	//Unlock(passphrase string) error
	//IsAccountUnLock(id string) bool
	//Lock(passphrase string) error
	//SetVariables(params, flags []string) error
	//GetVariable(key string) string
	//Restart(params, flags []string) error
	//NeedRestartNode(params, flags []string) bool
	//Setup(allocation string) error
}

// Start starts REPL.
func Start(c Client) {
	if !TestMode {
		r := &repl{client: c}
		r.initializeCommands()

		runPrompt(r.executor, r.completer, r.firstTime, uint16(len(r.commands)))
	} else {
		// holds for unit test purposes
		hold := make(chan bool)
		<-hold
	}
}

func (r *repl) initializeCommands() {
	r.commands = []command{
		{"new", "Create a new account (key pair) and set as current", r.createAccount},
		{"set", "Set one of the previously created accounts as current", r.chooseAccount},
		{"info", "Display the current account info", r.accountInfo},
		{"all-txs", "List all transactions (outgoing and incoming) for the current account", r.getMeshTransactions},
		{"send-coin", "Transfer coins from current account to another account", r.submitCoinTransaction},
		{"sign", "Sign a hex message with the current account private key", r.sign},
		{"textsign", "Sign a text message with the current account private key", r.textsign},
		{"rewards", "Set current account as rewards account in the node", r.setCoinbase},
		//{"smesh", "Start smeshing", r.smesh},
		{"node", "Get current p2p node info", r.nodeInfo},
		{"all", "Display all mesh accounts (debug)", r.debugAllAccounts},
		{"quit", "Quit the CLI", r.quit},

		//{"unlock accountInfo", "Unlock accountInfo.", r.unlockAccount},
		//{"lock accountInfo", "Lock LocalAccount.", r.lockAccount},
		//{"setup", "Setup POST.", r.setup},
		//{"restart node", "Restart node.", r.restartNode},
		//{"set", "change CLI flag or param. E.g. set param a=5 flag c=5 or E.g. set param a=5", r.setCLIFlagOrParam},
		//{"echo", "Echo runtime variable.", r.echoVariable},
	}
}

func (r *repl) executor(text string) {
	for _, c := range r.commands {
		if len(text) >= len(c.text) && text[:len(c.text)] == c.text {
			r.input = text
			//log.Debug(userExecutingCommandMsg, c.text)
			c.fn()
			return
		}
	}

	fmt.Println(printPrefix, "invalid command.")
}

func (r *repl) completer(in prompt.Document) []prompt.Suggest {
	suggets := make([]prompt.Suggest, 0)
	for _, command := range r.commands {
		s := prompt.Suggest{
			Text:        command.text,
			Description: command.description,
		}

		suggets = append(suggets, s)
	}

	return prompt.FilterHasPrefix(suggets, in.GetWordBeforeCursor(), true)
}

func (r *repl) firstTime() {
	fmt.Println(printPrefix, splash)

	if err := r.client.Sanity(); err != nil {
		log.Error("Failed to connect to node at %v: %v", r.client.ServerUrl(), err)
		r.quit()
	}

	fmt.Println("Welcome to Spacemesh. Connected to node at", r.client.ServerUrl())
}

func (r *repl) chooseAccount() {
	accs := r.client.ListAccounts()
	if len(accs) == 0 {
		r.createAccount()
		return
	}

	fmt.Println(printPrefix, "Choose an account to load:")
	accName := multipleChoice(accs)
	account, err := r.client.GetAccount(accName)
	if err != nil {
		panic("wtf")
	}
	fmt.Printf("%s Loaded account alias: `%s`, address: %s \n", printPrefix, account.Name, account.Address().String())

	r.client.SetCurrentAccount(account)
}

func (r *repl) createAccount() {
	fmt.Println(printPrefix, "Create a new account")
	alias := inputNotBlank(createAccountMsg)

	ac := r.client.CreateAccount(alias)
	err := r.client.StoreAccounts()
	if err != nil {
		log.Error("failed to create account: %v", err)
		return
	}

	fmt.Printf("%s Created account alias: `%s`, address: %s \n", printPrefix, ac.Name, ac.Address().String())
	r.client.SetCurrentAccount(ac)
}

func (r *repl) commandLineParams(idx int, input string) string {
	c := r.commands[idx]
	params := strings.Replace(input, c.text, "", -1)

	return strings.TrimSpace(params)
}

func (r *repl) accountInfo() {
	acc := r.client.CurrentAccount()
	if acc == nil {
		r.chooseAccount()
		acc = r.client.CurrentAccount()
	}

	address := gosmtypes.BytesToAddress(acc.PubKey)

	info, err := r.client.AccountInfo(address)
	if err != nil {
		log.Error("failed to get account info: %v", err)
		info = &localtypes.AccountState{}
	}

	fmt.Println(printPrefix, "Local alias: ", acc.Name)
	fmt.Println(printPrefix, "Address: ", address.String())
	fmt.Println(printPrefix, "Balance: ", info.Balance, coinUnitName)
	fmt.Println(printPrefix, "Nonce: ", info.Nonce)
	fmt.Println(printPrefix, fmt.Sprintf("Public key: 0x%s", hex.EncodeToString(acc.PubKey)))
	fmt.Println(printPrefix, fmt.Sprintf("Private key: 0x%s", hex.EncodeToString(acc.PrivKey)))
}

func (r *repl) nodeInfo() {

	info, err := r.client.NodeInfo()
	if err != nil {
		log.Error("failed to get node info: %v", err)
		return
	}

	fmt.Println(printPrefix, "Version:", info.Version)
	fmt.Println(printPrefix, "Build:", info.Build)
	fmt.Println(printPrefix, "API server:", r.client.ServerUrl())

	status, err := r.client.NodeStatus()
	if err != nil {
		log.Error("failed to get node status: %v", err)
		return
	}

	fmt.Println(printPrefix, "Synced:", status.IsSynced)
	fmt.Println(printPrefix, "Synced layer:", status.SyncedLayer.Number)
	fmt.Println(printPrefix, "Current layer:", status.TopLayer.Number)
	fmt.Println(printPrefix, "Verified layer:", status.VerifiedLayer.Number)
	fmt.Println(printPrefix, "Peers:", status.ConnectedPeers)

	/*
		fmt.Println(printPrefix, "Smeshing data directory:", info.SmeshingDatadir)
		fmt.Println(printPrefix, "Smeshing status:", info.SmeshingStatus)
		fmt.Println(printPrefix, "Smeshing coinbase:", info.SmeshingCoinbase)
		fmt.Println(printPrefix, "Smeshing remaining bytes:", info.SmeshingRemainingBytes)
	*/
}

func (r *repl) debugAllAccounts() {

	accounts, err := r.client.DebugAllAccounts()
	if err != nil {
		log.Error("failed to get debug all accounts: %v", err)
		return
	}

	for _, a := range accounts {
		fmt.Println(printPrefix, "Address:", gosmtypes.BytesToAddress(a.AccountId.Address).String())
		fmt.Println(printPrefix, "Balance:", a.StateCurrent.Balance.Value , coinUnitName)
		fmt.Println(printPrefix, "Nonce:", a.StateCurrent.Counter)
		fmt.Println(printPrefix, "-----")
	}
}

func (r *repl) submitCoinTransaction() {
	fmt.Println(printPrefix, initialTransferMsg)
	acc := r.client.CurrentAccount()
	if acc == nil {
		r.chooseAccount()
		acc = r.client.CurrentAccount()
	}

	srcAddress := gosmtypes.BytesToAddress(acc.PubKey)
	info, err := r.client.AccountInfo(srcAddress)
	if err != nil {
		log.Error("failed to get account info: %v", err)
		return
	}

	destAddressStr := inputNotBlank(destAddressMsg)
	destAddress := gosmtypes.HexToAddress(destAddressStr)

	amountStr := inputNotBlank(amountToTransferMsg) + coinUnitName

	gas := uint64(1)
	if yesOrNoQuestion(useDefaultGasMsg) == "n" {
		gasStr := inputNotBlank(enterGasPrice)
		gas, err = strconv.ParseUint(gasStr, 10, 64)
		if err != nil {
			log.Error("invalid transaction fee", err)
			return
		}
	}

	fmt.Println(printPrefix, "Transaction summary:")
	fmt.Println(printPrefix, "From:  ", srcAddress.String())
	fmt.Println(printPrefix, "To:    ", destAddress.String())
	fmt.Println(printPrefix, "Amount:", amountStr)
	fmt.Println(printPrefix, "Fee:   ", gas , coinUnitName)
	fmt.Println(printPrefix, "Nonce: ", info.Nonce)

	amount, err := strconv.ParseUint(amountStr, 10, 64)

	if yesOrNoQuestion(confirmTransactionMsg) == "y" {
		txState, err := r.client.Transfer(destAddress, info.Nonce, amount, gas, 100, acc.PrivKey)
		if err != nil {
			log.Error(err.Error())
			return
		}

		fmt.Println(printPrefix, "Transaction submitted.")
		fmt.Println(printPrefix, fmt.Sprintf("Transaction id: 0x%v", hex.EncodeToString(txState.Id.Id)))
		fmt.Println(printPrefix, fmt.Sprintf("Transaction state: 0x%v", txState.State.String()))
	}
}

func (r *repl) smesh() {
	acc := r.client.CurrentAccount()
	if acc == nil {
		r.chooseAccount()
		acc = r.client.CurrentAccount()
	}

	datadir := inputNotBlank(smeshingDatadirMsg)

	spaceStr := inputNotBlank(smeshingSpaceAllocationMsg)
	space, err := strconv.ParseUint(spaceStr, 10, 32)
	if err != nil {
		log.Error("failed to parse: %v", err)
		return
	}

	if err := r.client.Smesh(datadir, uint(space)<<30, acc.Address().String()); err != nil {
		log.Error("failed to start smeshing: %v", err)
		return
	}
}

func (r *repl) getMeshTransactions() {
	acc := r.client.CurrentAccount()
	if acc == nil {
		r.chooseAccount()
		acc = r.client.CurrentAccount()
	}

	// todo: request offset and total from user
	txs, total, err := r.client.GetMeshTransactions(acc.Address(), 0, 100)
	if err != nil {
		log.Error("failed to list transactions: %v", err)
		return
	}

	fmt.Println(printPrefix, fmt.Sprintf("Total mesh transactions: %d", total))
	for _, tx := range txs {
		printTransaction(tx)
	}
}

// helper method - prints tx info
func printTransaction(transaction *apitypes.Transaction) {
	// todo: implement me
}

func (r *repl) quit() {
	os.Exit(0)
}

func (r *repl) setCoinbase() {
	acc := r.client.CurrentAccount()
	if acc == nil {
		r.chooseAccount()
		acc = r.client.CurrentAccount()
	}

	status, err := r.client.SetCoinbase(acc.Address() 	)

	if err != nil {
		log.Error("failed to set rewards address: %v", err)
		return
	}

	if status.Code == 0 {
		fmt.Println(printPrefix, "Rewards address set to:", acc.Address().String())
	} else {
		// todo: what are possible non-zero status codes here?
		fmt.Println(printPrefix, fmt.Sprintf("Response status code: %d", status.Code))
	}
}

func (r *repl) sign() {
	acc := r.client.CurrentAccount()
	if acc == nil {
		r.chooseAccount()
		acc = r.client.CurrentAccount()
	}

	msgStr := inputNotBlank(msgSignMsg)
	msg, err := hex.DecodeString(msgStr)
	if err != nil {
		log.Error("failed to decode msg hex string: %v", err)
		return
	}

	signature := ed25519.Sign2(acc.PrivKey, msg)

	fmt.Println(printPrefix, fmt.Sprintf("signature (in hex): %x", signature))
}

func (r *repl) textsign() {
	acc := r.client.CurrentAccount()
	if acc == nil {
		r.chooseAccount()
		acc = r.client.CurrentAccount()
	}

	msg := inputNotBlank(msgTextSignMsg)
	signature := ed25519.Sign2(acc.PrivKey, []byte(msg))

	fmt.Println(printPrefix, fmt.Sprintf("signature (in hex): %x", signature))
}

/*
func (r *repl) unlockAccount() {
	passphrase := r.commandLineParams(1, r.input)
	err := r.client.Unlock(passphrase)
	if err != nil {
		log.Debug(err.Error())
		return
	}

	acctCmd := r.commands[3]
	r.executor(fmt.Sprintf("%s %s", acctCmd.text, passphrase))
}

func (r *repl) lockAccount() {
	passphrase := r.commandLineParams(2, r.input)
	err := r.client.Lock(passphrase)
	if err != nil {
		log.Debug(err.Error())
		return
	}

	acctCmd := r.commands[3]
	r.executor(fmt.Sprintf("%s %s", acctCmd.text, passphrase))
}*/
