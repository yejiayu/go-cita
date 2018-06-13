package errors

var (
	// ExecutionError
	NotEnoughBaseGas        = new("NotEnoughBaseGas")
	BlockGasLimitReached    = new("BlockGasLimitReached")
	AccountGasLimitReached  = new("AccountGasLimitReached")
	InvalidTransactionNonce = new("InvalidTransactionNonce")
	NotEnoughCash           = new("NotEnoughCash")
	NoTransactionPermission = new("NoTransactionPermission")
	NoContractPermission    = new("NoContractPermission")
	NoCallPermission        = new("NoCallPermission")
	ExecutionInternal       = new("ExecutionInternal")
	TransactionMalformed    = new("TransactionMalformed")

	// EvmError
	OutOfGas                   = new("OutOfGas")
	BadJumpDestination         = new("BadJumpDestination")
	BadInstruction             = new("BadInstruction")
	StackUnderflow             = new("StackUnderflow")
	OutOfStack                 = new("OutOfStack")
	Internal                   = new("Internal")
	MutableCallInStaticContext = new("MutableCallInStaticContext")
	OutOfBounds                = new("OutOfBounds")
	Reverted                   = new("Reverted")
)
