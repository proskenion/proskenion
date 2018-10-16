package model

type ObjectFinder interface {
	// Query gets value from targetId
	Query(targetId string, value Unmarshaler) error
	// Append [targetId] = value
	Append(targetId string, value Marshaler) error
}

type TxFinder interface {
	// Query gets tx from txHash
	Query(txHash Hash) (Transaction, error)
	// Append tx
	Append(tx Transaction) error
}
