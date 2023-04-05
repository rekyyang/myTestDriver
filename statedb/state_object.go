// Copyright 2014 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package statedb

import (
	"bytes"
	"fmt"
	"io"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/metrics"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/node-real/blocktree/models"
	"github.com/node-real/go-pkg/log"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	emptyCodeHash = crypto.Keccak256(nil)

	// emptyRoot is the known root hash of an empty trie.
	emptyRoot = []byte("56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421")

	stateCodeDBMissCounter = promauto.NewCounter(prometheus.CounterOpts{
		Namespace: "states",
		Subsystem: "state",
		Name:      "code_db_miss",
	})
)

type Code []byte

func (c Code) String() string {
	return string(c) //strings.Join(Disassemble(c), " ")
}

type Storage map[common.Hash]common.Hash

func (s Storage) String() (str string) {
	for key, value := range s {
		str += fmt.Sprintf("%X : %X\n", key, value)
	}

	return
}

func (s Storage) Copy() Storage {
	cpy := make(Storage, len(s))
	for key, value := range s {
		cpy[key] = value
	}

	return cpy
}

// StateObject represents an Ethereum account which is being modified.
//
// The usage pattern is as follows:
// First you need to obtain a state object.
// Account values can be accessed and modified through the object.
type StateObject struct {
	address  common.Address
	addrHash common.Hash // hash of ethereum address of the account
	data     models.Account
	db       *StateDB

	code Code // contract bytecode, which gets set when code is loaded

	originStorage  Storage // Storage cache of original entries to dedup rewrites, reset for every transaction
	pendingStorage Storage // Storage entries that need to be flushed to disk, at the end of an entire block
	dirtyStorage   Storage // Storage entries that have been modified in the current transaction execution

	// Cache flags.
	// When an object is marked suicided it will be delete from the trie
	// during the "update" phase of the state transition.
	dirtyCode bool // true if the code was updated
	suicided  bool
	deleted   bool
}

// empty returns whether the account is considered empty.
func (s *StateObject) empty() bool {
	return s.data.Nonce == 0 && s.data.Balance.Sign() == 0 && bytes.Equal(s.data.CodeHash, emptyCodeHash)
}

// newObject creates a state object.
func newObject(db *StateDB, address common.Address, data models.Account) *StateObject {
	if data.Balance == nil {
		data.Balance = new(big.Int)
	}
	if data.CodeHash == nil {
		data.CodeHash = emptyCodeHash
	}
	if len(data.Root) == 0 || bytes.Equal(data.Root, (common.Hash{}).Bytes()) {
		data.Root = emptyRoot
	}
	return &StateObject{
		db:             db,
		address:        address,
		addrHash:       crypto.Keccak256Hash(address[:]),
		data:           data,
		originStorage:  make(Storage),
		pendingStorage: make(Storage),
		dirtyStorage:   make(Storage),
	}
}

// EncodeRLP implements rlp.Encoder.
func (s *StateObject) EncodeRLP(w io.Writer) error {
	return rlp.Encode(w, s.data)
}

func (s *StateObject) markSuicided() {
	s.suicided = true
}

func (s *StateObject) touch() {
	s.db.journal.append(touchChange{
		account: &s.address,
	})
	if s.address == ripemd {
		// Explicitly put it in the dirty-cache, which is otherwise generated from
		// flattened journals.
		s.db.journal.dirty(s.address)
	}
}

// GetState retrieves a value from the account storage trie.
func (s *StateObject) GetState(key common.Hash) common.Hash {
	// If we have a dirty value for this state entry, return it
	value, dirty := s.dirtyStorage[key]
	if dirty {
		return value
	}
	// Otherwise return the entry's original value
	return s.GetCommittedState(key)
}

// GetCommittedState retrieves a value from the committed account storage trie.
func (s *StateObject) GetCommittedState(key common.Hash) common.Hash {
	// If we have a pending write or clean cached, return that
	if value, pending := s.pendingStorage[key]; pending {
		return value
	}
	if value, cached := s.originStorage[key]; cached {
		return value
	}

	// If the object was destructed in *this* block (and potentially resurrected),
	// the storage has been cleared out, and we should *not* consult the previous
	// database about any storage values. The only possible alternatives are:
	//   1) resurrect happened, and new slot values were set -- those should
	//      have been handles via pendingStorage above.
	//   2) we don't have new values, and can deliver empty response back
	if _, destructed := s.db.stateObjectsDestruct[s.address]; destructed {
		return common.Hash{}
	}

	// If no live objects are available, attempt to use snapshots
	var (
		value common.Hash
		err   error
		meter *time.Duration
	)
	readStart := time.Now()
	if metrics.EnabledExpensive {
		// If the snap is 'under construction', the first lookup may fail. If that
		// happens, we don't want to double-count the time elapsed. Thus this
		// dance with the metering.
		defer func() {
			if meter != nil {
				*meter += time.Since(readStart)
			}
		}()
	}

	if metrics.EnabledExpensive {
		meter = &s.db.StorageReads
	}

	//if value, err = s.db.view.Storage(s.address, s.data.Incarnation, key); err != nil {
	if value, err = s.db.prestates.Storage(s.address, s.data.Incarnation, key); err != nil {
		s.db.setError(err)
		return common.Hash{}
	}

	s.originStorage[key] = value
	return value
}

// SetState updates a value in account storage.
func (s *StateObject) SetState(key, value common.Hash) {
	// If the new value is the same as old, don't set
	prev := s.GetState(key)
	if prev == value {
		return
	}
	// New value is different, update and journal the change
	s.db.journal.append(storageChange{
		account:  &s.address,
		key:      key,
		prevalue: prev,
	})
	s.setState(key, value)
}

func (s *StateObject) setState(key, value common.Hash) {
	s.dirtyStorage[key] = value
}

// finalise moves all dirty storage slots into the pending area to be hashed or
// committed later. It is invoked at the end of every transaction.
func (s *StateObject) finalise(prefetch bool) {
	slotsToPrefetch := make([][]byte, 0, len(s.dirtyStorage))
	for key, value := range s.dirtyStorage {
		s.pendingStorage[key] = value
		if value != s.originStorage[key] {
			slotsToPrefetch = append(slotsToPrefetch, common.CopyBytes(key[:])) // Copy needed for closure
		}
	}
	// TODO: do we need prefetcher
	//if s.db.prefetcher != nil && prefetch && len(slotsToPrefetch) > 0 && s.data.Root != emptyRoot {
	//	s.db.prefetcher.prefetch(s.data.Root, slotsToPrefetch)
	//}
	if len(s.dirtyStorage) > 0 {
		s.dirtyStorage = make(Storage)
	}
	_ = slotsToPrefetch // lint
}

// AddBalance adds amount to s's balance.
// It is used to add funds to the destination account of a transfer.
func (s *StateObject) AddBalance(amount *big.Int) {
	// EIP161: We must check emptiness for the objects such that the account
	// clearing (0,0,0 objects) can take effect.
	if amount.Sign() == 0 {
		if s.empty() {
			s.touch()
		}
		return
	}
	s.SetBalance(new(big.Int).Add(s.Balance(), amount))
}

// SubBalance removes amount from s's balance.
// It is used to remove funds from the origin account of a transfer.
func (s *StateObject) SubBalance(amount *big.Int) {
	if amount.Sign() == 0 {
		return
	}
	s.SetBalance(new(big.Int).Sub(s.Balance(), amount))
}

func (s *StateObject) SetBalance(amount *big.Int) {
	s.db.journal.append(balanceChange{
		account: &s.address,
		prev:    new(big.Int).Set(s.data.Balance),
	})
	s.setBalance(amount)
}

func (s *StateObject) setBalance(amount *big.Int) {
	s.data.Balance = amount
}

func (s *StateObject) deepCopy(db *StateDB) *StateObject {
	so := newObject(db, s.address, s.data)
	so.code = s.code
	so.dirtyStorage = s.dirtyStorage.Copy()
	so.originStorage = s.originStorage.Copy()
	so.pendingStorage = s.pendingStorage.Copy()
	so.suicided = s.suicided
	so.dirtyCode = s.dirtyCode
	so.deleted = s.deleted
	return so
}

//
// Attribute accessors
//

// Returns the address of the contract/account
func (s *StateObject) Address() common.Address {
	return s.address
}

// Code returns the contract code associated with this object, if any.
func (s *StateObject) Code() []byte {
	if s.code != nil {
		return s.code
	}
	if bytes.Equal(s.CodeHash(), emptyCodeHash) {
		return nil
	}

	//code, err := s.db.view.ContractCode(common.BytesToHash(s.CodeHash()))
	code, err := s.db.prestates.ContractCode(s.address, common.BytesToHash(s.CodeHash()))
	if err != nil || len(code) == 0 {
		stateCodeDBMissCounter.Inc()
		s.db.setError(fmt.Errorf("can't load code, hash: %x , err: %v", s.CodeHash(), err))
		log.Warnw("can't load code", "hash", hexutil.Bytes(s.CodeHash()), "err", err)
		return nil
	}

	s.code = code
	return code
}

// CodeSize returns the size of the contract code associated with this object,
// or zero if none. This method is an almost mirror of Code, but uses a cache
// inside the database to avoid loading codes seen recently.
func (s *StateObject) CodeSize() int {
	return len(s.Code())
}

func (s *StateObject) SetCode(codeHash common.Hash, code []byte) {
	prevcode := s.Code()
	s.db.journal.append(codeChange{
		account:  &s.address,
		prevhash: s.CodeHash(),
		prevcode: prevcode,
	})
	s.setCode(codeHash, code)
}

func (s *StateObject) setCode(codeHash common.Hash, code []byte) {
	s.code = code
	s.data.CodeHash = codeHash[:]
	s.dirtyCode = true
}

func (s *StateObject) SetNonce(nonce uint64) {
	s.db.journal.append(nonceChange{
		account: &s.address,
		prev:    s.data.Nonce,
	})
	s.setNonce(nonce)
}

func (s *StateObject) setNonce(nonce uint64) {
	s.data.Nonce = nonce
}

func (s *StateObject) CodeHash() []byte {
	return s.data.CodeHash
}

func (s *StateObject) Balance() *big.Int {
	return s.data.Balance
}

func (s *StateObject) Nonce() uint64 {
	return s.data.Nonce
}
