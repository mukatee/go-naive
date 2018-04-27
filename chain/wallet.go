package chain

func balanceFor(address string) int {
	balance := 0
	for _, val := range unspentTxOuts {
		if val.Address == address {
			balance += val.Amount
		}
	}
	return balance
}

func findTxInsFor(address string, amount int) ([]TxIn, int) {
	balance := 0
	var unspents []TxIn
	for _, val := range unspentTxOuts {
		if val.Address == address {
			balance += val.Amount
			txIn := TxIn{val.TxId, val.TxIdx}
			unspents = append(unspents, txIn)
		}
		if balance >= amount {
			return unspents, balance
		}
	}
	return nil, -1
}

func splitTxIns(from string, to string, toSend int, total int) []TxOut {
	diff := total - toSend
	txOut := TxOut{to, toSend}
	var txOuts []TxOut
	txOuts = append(txOuts, txOut)
	if diff == 0 {
		return txOuts
	}
	txOut2 := TxOut{from, diff}
	txOuts = append(txOuts, txOut2)
	return txOuts
}