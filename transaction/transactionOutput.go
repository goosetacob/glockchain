package transaction

// TXOutput defines a transaction on the glockchain
type TXOutput struct {
	Value        int    // value in transaction
	ScriptPubKey string // arbitrary lock
}

// CanBeUnlockedWith check unlockingData matches the ScriptPubKey
func (out *TXOutput) CanBeUnlockedWith(unlockingData string) bool {
	return out.ScriptPubKey == unlockingData
}
