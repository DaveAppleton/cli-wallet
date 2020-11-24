package repl

const (
	initialTransferMsg         = "Transfer coins from local account to another account."
	destAddressMsg             = "Enter or paste destination address: "
	txIdMsg                    = "Enter or paste transaction id: "
	smesherIdMsg               = "Enter or paste a Smesher id: "
	amountToTransferMsg        = "Enter amount to transfer in Smidge: "
	confirmTransactionMsg      = "Confirm transaction (y/n): "
	confirmDeleteDataMsg       = "Delete smeshing smeshing data files (y/n)"
	createAccountMsg           = "Account alias (name): "
	useDefaultGasMsg           = "Use default transaction fee of 1 Smidge? (y/n) "
	enterGasPrice              = "Enter transaction fee (Smidge):"
	smeshingDatadirMsg         = "Enter data file directory: "
	smeshingSpaceAllocationMsg = "Enter space allocation (GB): "
	msgSignMsg                 = "Enter message to sign (in hex): "
	msgTextSignMsg             = "Enter text message to sign: "
	coinUnitName               = "Smidge"
)

const splash = `

                                    .++++++++++++++++++++++++++.
                                    %@@@@@@@@@@@@@@@@@@@@@@@@@@%
                                   -@@@@@@@##############@@@@@@@-
                                     +@@@@@*.          .*@@@@@+
                                      .+@@@@@*.      .*@@@@@+.
                                        .*@@@@@+.  .+@@@@@*.
                                          .*@@@@@++@@@@@*.
                                            .*@@@@@@@@*.
                                              *@@@@@@*
                                            =@@@@@@@@@@=
                                          =@@@@@#::#@@@@@=
                                        =%@@@@%:    :#@@@@%=
                                      -%@@@@%-        -%@@@@%-
                                    -%@@@@%-            -%@@@@%-
                                   *@@@@%-                -%@@@@*
                                   *@@@@#:                :#@@@@*
                                    =@@@@@#:            :#@@@@@=
                                      =@@@@@#:        :#@@@@@=
                                        =@@@@@#:    :#@@@@@=
                                          +@@@@@*..*@@@@@+
                                            +@@@@@@@@@@+
                                             .*@@@@@@*
                                            .*@@@@@@@@*.
                                          .+@@@@@**@@@@@+.
                                         +@@@@@*.  .*@@@@@+
                                       +@@@@@*.      .*@@@@@+
                                     +@@@@@*.          .*@@@@@+
                                   -@@@@@@@##############@@@@@@@-
                                    %@@@@@@@@@@@@@@@@@@@@@@@@@@%
                                    .++++++++++++++++++++++++++.

`
