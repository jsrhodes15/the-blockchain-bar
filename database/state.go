package database

import (
	"bufio"
	"encoding/json"
	"os"
	"path/filepath"
)

type State struct {
	Balances  map[Account]uint
	txMempool []Tx

	dbFile *os.File
}

func NewStateFromDisk() (*State, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	// start with 'genesis' file
	genesisFilePath := filepath.Join(cwd, "database", "genesis.json")
	gen, err := loadGenesis(genesisFilePath)
	if err != nil {
		return nil, err
	}

	// build a map of balances for easy lookup
	balances := make(map[Account]uint)
	for account, balance := range gen.Balances {
		balances[account] = balance
	}

	txDbFilePath := filepath.Join(cwd, "database", "tx.db")
	f, err := os.OpenFile(txDbFilePath, os.O_APPEND|os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(f)
	state := &State{balances, make([]Tx, 0), f}

	// Iterate over each line in tx.db file (transaction)
	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return nil, err
		}

		var tx Tx
		json.Unmarshal(scanner.Bytes(), &tx)

		if err := state.apply(tx); err != nil {
			return nil, err
		}
	}

	return state, nil
}

func (s *State) Add(tx Tx) error {
	if err := s.apply(tx); err != nil {
		return err
	}

	s.txMempool = append(s.txMempool, tx)

	return nil
}

func (s *State) Persist() error {
	// make a copy of mempool because s.txMempool will be modified in the loop
	mempool := make([]Tx, len(s.txMempool))
	copy(mempool, s.txMempool)

	for i := 0; i < len(mempool); i++ {
		txJson, err := json.Marshal(s.txMempool[i])
		if err != nil {
			return err
		}

		if _, err = s.dbFile.Write(append(txJson, '\n')); err != nil {
			return err
		}

		// remove the TX written to file from mempool
		s.txMempool = append(s.txMempool[:i], s.txMempool[i+1:]...)
	}

	return nil
}

func (s *State) apply(tx Tx) error {
	if tx.IsReward {
		// do something
	}
}
