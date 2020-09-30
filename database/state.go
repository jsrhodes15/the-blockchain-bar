package database

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"sort"
)

type State struct {
	Balances  map[Account]uint

	dbFile *os.File

	latestBlock     Block
	latestBlockHash Hash
	hasGenesisBlock bool
}

func NewStateFromDisk(dataDir string) (*State, error) {
	err := initDataDirIfNotExists(dataDir)
	if err != nil {
		return nil, err
	}
	// start with a 'genesis' file
	gen, err := loadGenesis(getGenesisJsonFilePath(dataDir))
	if err != nil {
		return nil, err
	}
	// build a map of balances for easy lookup
	balances := make(map[Account]uint)
	for account, balance := range gen.Balances {
		balances[account] = balance
	}

	dbFilePath := getBlocksDbFilePath(dataDir)
	f, err := os.OpenFile(dbFilePath, os.O_APPEND|os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(f)

	state := &State{balances, f, Block{}, Hash{}, false}
	// Iterate over each line in tx.db file (transaction)
	for scanner.Scan() {
		if err := scanner.Err(); err != nil {
			return nil, err
		}

		blockFsJson := scanner.Bytes()

		if len(blockFsJson) == 0 {
			break
		}

		var blockFs BlockFS
		err = json.Unmarshal(blockFsJson, &blockFs)
		if err != nil {
			return nil, err
		}

		err = applyTXs(blockFs.Value.TXs, state)
		if err != nil {
			return nil, err
		}

		state.latestBlock = blockFs.Value
		state.latestBlockHash = blockFs.Key
		state.hasGenesisBlock = true
	}

	return state, nil
}

func (s *State) AddBlocks(blocks []Block) error {
	for _, b := range blocks {
		_, err := s.AddBlock(b)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *State) AddBlock(b Block) (Hash, error) {
	pendingState := s.copy()
	// Validate block meta + payload. Replays transactions to verify balances
	err := applyBlock(b, &pendingState)
	if err != nil {
		return Hash{}, err
	}

	blockHash, err := b.Hash()
	if err != nil {
		return Hash{}, err
	}

	blockFs := BlockFS{blockHash, b}

	blockFsJson, err := json.Marshal(blockFs)
	if err != nil {
		return Hash{}, err
	}

	fmt.Printf("\nPersisting new Block to disk:\n")
	fmt.Printf("\t%s\n", blockFsJson)
	// Write to disk
	_, err = s.dbFile.Write(append(blockFsJson, '\n'))
	if err != nil {
		return Hash{}, err
	}
	// All TXs are valid and no error writing to disk -> update main state
	s.Balances = pendingState.Balances
	s.latestBlockHash = blockHash
	s.latestBlock = b
	s.hasGenesisBlock = true

	return blockHash, nil
}

func (s *State) NextBlockNumber() uint64 {
	if !s.hasGenesisBlock {
		return uint64(0)
	}

	return s.LatestBlock().Header.Number + 1
}

func (s *State) LatestBlock() Block {
	return s.latestBlock
}

func (s *State) LatestBlockHash() Hash {
	return s.latestBlockHash
}

func (s *State) Close() error {
	return s.dbFile.Close()
}

func (s *State) copy() State {
	// For validation purposes, we want to make a copy of State, without any pointers to the original State{}
	c := State{}
	c.latestBlock = s.latestBlock
	c.latestBlockHash = s.latestBlockHash
	c.Balances = make(map[Account]uint)

	for acc, balance := range s.Balances {
		c.Balances[acc] = balance
	}

	return c
}

// verifies if a block can be added to the blockchain
// block metadata are verified as well as transactions within (sufficient balances, etc).
func applyBlock(b Block, s *State) error {
	nextExpectedBlockNumber := s.latestBlock.Header.Number + 1

	if s.hasGenesisBlock && b.Header.Number != nextExpectedBlockNumber {
		return fmt.Errorf(
			"next expected block must be '%d' not '%d'",
			nextExpectedBlockNumber,
			b.Header.Number,
		)
	}

	if s.hasGenesisBlock && s.latestBlock.Header.Number > 0 && !reflect.DeepEqual(b.Header.Parent, s.latestBlockHash) {
		return fmt.Errorf(
			"next block parent hash must be '%x' not '%x'",
			s.latestBlockHash,
			b.Header.Parent,
		)
	}

	hash, err := b.Hash()
	if err != nil {
		return err
	}

	if !IsBlockHashValid(hash) {
		return fmt.Errorf("invalid block %x", hash)
	}

	err = applyTXs(b.TXs, s)
	if err != nil {
		return err
	}

	return nil
}

func applyTXs(txs []Tx, s *State) error {
	sort.Slice(txs, func(i, j int) bool {
		return txs[i].Time < txs[j].Time
	})

	for _, tx := range txs {
		err := applyTx(tx, s)
		if err != nil {
			return err
		}
	}
	return nil
}

func applyTx(tx Tx, s *State) error {
	if tx.IsReward() {
		s.Balances[tx.To] += tx.Value
		return nil
	}

	if tx.Value > s.Balances[tx.From] {
		return fmt.Errorf("bad TX. Sender '%s' balance is %d TBB. Tx cost is %d",
			tx.From,
			s.Balances[tx.From],
			tx.Value,
		)
	}

	s.Balances[tx.From] -= tx.Value
	s.Balances[tx.To] += tx.Value

	return nil
}
