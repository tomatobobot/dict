package db

import "github.com/boltdb/bolt"

// Tx represents a BoltDB transaction
type Tx struct {
	*bolt.Tx
}

func (tx *Tx) Dict(word []byte) (*Dict, error) {
	d := &Dict{
		Tx:   tx,
		Word: word,
	}

	return d, d.Load()
}
