package model

type ObjectFinder interface {
	// Query gets value from targetId
	Query(targetId Address, value Unmarshaler) error
	// Query All gets value from fromId
	QueryAll(fromId Address, value UnmarshalerFactory) ([]Unmarshaler, error)
	// Append [targetId] = value
	Append(targetId Address, value Marshaler) error
}

type TxFinder interface {
	// Query gets tx from txHash
	Query(txHash Hash) (Transaction, error)
	// Append tx
	Append(tx Transaction) error
}
