package db

import (
	"fmt"
)

// Dict represents a Dict
type Dict struct {
	Tx     *Tx
	Word   []byte
	Result []byte
}

func (p *Dict) bucket() []byte {
	return []byte("dicts")
}

func (d *Dict) get() ([]byte, error) {
	result := d.Tx.Bucket(d.bucket()).Get(d.Word)
	if result == nil {
		return nil, fmt.Errorf("找不到%s", string(d.Word))
	}

	return result, nil
}

func (d *Dict) Load() error {
	result, err := d.get()
	if err != nil {
		return err
	}

	d.Result = result

	return nil
}

func (d *Dict) Save() error {
	if len(d.Word) == 0 {
		return fmt.Errorf("word is empty")
	}

	return d.Tx.Bucket(d.bucket()).Put(d.Word, d.Result)
}
