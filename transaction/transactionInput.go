package transaction

// TXInput defines am input to a transaction on the glockchain
type TXInput struct {
	Txid      []byte // transaction id
	Vout      int    // index of an output in transaction
	ScriptSig string // user defined wallet address
}

// CanUnlockOutputWith check unlockingData matches the ScriptSig
func (in *TXInput) CanUnlockOutputWith(unlockingData string) bool {
	return in.ScriptSig == unlockingData
}
