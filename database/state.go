package database

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type State struct {
	Balances  map[Account]uint
	txMempool []Tx

	dbFile          *os.File
	latestBlockHash Hash
}

func NewStateFromDisk() (*State, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	// start with 'genesis' file
	gen, err := loadGenesis(filepath.Join(cwd, "database", "genesis.json"))
	if err != nil {
		return nil, err
	}

	// build a map of balances for easy lookup
	balances := make(map[Account]uint)
	for account, balance := range gen.Balances {
		balances[account] = balance
	}

	f, err := os.OpenFile(filepath.Join(cwd, "database", "block.db"), os.O_APPEND|os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(f)

	state := &State{balances, make([]Tx, 0), f, Hash{}}

	// Iterate over each line in tx.db file (transaction)
	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return nil, err
		}

		blockFsJson := scanner.Bytes()
		var blockFs BlockFS
		err = json.Unmarshal(blockFsJson, &blockFs)
		if err != nil {
			return nil, err
		}

		err = state.applyBlock(blockFs.Value)
		if err != nil {
			return nil, err
		}

		state.latestBlockHash = blockFs.Key
	}

	return state, nil
}

func (s *State) LatestBlockHash() Hash {
	return s.latestBlockHash
}

func (s *State) AddBlock(b Block) error {
	for _, tx := range b.TXs {
		err := s.AddTx(tx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *State) AddTx(tx Tx) error {
	err := s.apply(tx)
	if err != nil {
		return err
	}

	s.txMempool = append(s.txMempool, tx)

	return nil
}

func (s *State) Persist() (Hash, error) {
	// make a copy of mempool because s.txMempool will be modified in the loop
	block := NewBlock(
		s.latestBlockHash,
		uint64(time.Now().Unix()),
		s.txMempool,
	)

	blockHash, err := block.Hash()
	if err != nil {
		return Hash{}, err
	}

	blockFs := BlockFS{blockHash, block}

	// Encode it into a JSON string
	blockFsJson, err := json.Marshal(blockFs)
	if err != nil {
		return Hash{}, err
	}

	fmt.Printf("Persisting new block to disk:\n")
	fmt.Printf("\t%s\n", blockFsJson)

	// Write it to the DB file on a new line
	_, err = s.dbFile.Write(append(blockFsJson, '\n'))
	if err != nil {
		return Hash{}, err
	}

	s.latestBlockHash = blockHash

	// Reset the mempool
	s.txMempool = []Tx{}

	return s.latestBlockHash, nil
}

func (s *State) Close() error {
	return s.dbFile.Close()
}

func (s *State) applyBlock(b Block) error {
	for _, tx := range b.TXs {
		err := s.apply(tx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *State) apply(tx Tx) error {
	if tx.IsReward() {
		s.Balances[tx.To] += tx.Value
		return nil
	}

	if s.Balances[tx.From]-tx.Value < 0 {
		return fmt.Errorf("insufficient balance")
	}

	s.Balances[tx.From] -= tx.Value
	s.Balances[tx.To] += tx.Value

	return nil
}
