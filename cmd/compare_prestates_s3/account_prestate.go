package main

import (
	"io"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
)

var (
	// EmptyRoot is the known root hash of an empty trie.
	EmptyRoot = common.HexToHash("56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421")

	// EmptyCodeHash is the known hash of the empty EVM bytecode.
	EmptyCodeHash = crypto.Keccak256Hash(nil)
)

type AccountPrestate struct {
	Nonce             uint64
	Incarnation       uint64
	Balance           *big.Int
	Root              []byte
	CodeHash          []byte
	Storages          map[common.Hash]common.Hash `rlp:"-"`
	StoragesKeys      []common.Hash
	StoragesValues    []common.Hash
	Address           common.Address
	DeletedOrNotExist bool
}

func (acc AccountPrestate) IsDeletedOrNotExist() bool {
	return acc.DeletedOrNotExist
}

// nolock
func (acc *AccountPrestate) SetStorages(key, value common.Hash) {
	if acc.Storages == nil {
		acc.Storages = make(map[common.Hash]common.Hash)
	}
	acc.Storages[key] = value
}

// nolock
func (acc *AccountPrestate) GetStorages(key common.Hash) (common.Hash, bool) {
	if acc.Storages == nil {
		return common.Hash{}, false
	}
	v, ok := acc.Storages[key]
	return v, ok
}

func (acc *AccountPrestate) Copy() *AccountPrestate {
	// field `StorageKeys` and `StorageValues` are skipped
	newAcc := AccountPrestate{
		Nonce:             acc.Nonce,
		Incarnation:       acc.Incarnation,
		Root:              acc.Root,
		CodeHash:          acc.CodeHash,
		Storages:          map[common.Hash]common.Hash{},
		Address:           acc.Address,
		DeletedOrNotExist: acc.DeletedOrNotExist,
	}

	if acc.Balance == nil {
		newAcc.Balance = nil
	} else {
		newAcc.Balance = new(big.Int).Set(acc.Balance)
	}

	for k, v := range acc.Storages {
		newAcc.SetStorages(k, v)
	}

	return &newAcc
}

type AccountPrestateRLP AccountPrestate

func (acc *AccountPrestate) EncodeRLP(w io.Writer) error {
	lenOfKeys := len(acc.Storages)
	acc.StoragesKeys = make([]common.Hash, lenOfKeys)
	acc.StoragesValues = make([]common.Hash, lenOfKeys)
	idx := 0
	for k, v := range acc.Storages {
		acc.StoragesKeys[idx] = k
		acc.StoragesValues[idx] = v
		idx++
	}

	return rlp.Encode(w, (*AccountPrestateRLP)(acc))
}

func (acc *AccountPrestate) DecodeRLP(stream *rlp.Stream) error {
	accRLP := (*AccountPrestateRLP)(acc)
	if err := stream.Decode(accRLP); err != nil {
		return err
	}

	acc.Storages = make(map[common.Hash]common.Hash, len(acc.StoragesKeys))
	for idx := range acc.StoragesKeys {
		acc.Storages[acc.StoragesKeys[idx]] = acc.StoragesValues[idx]
	}

	return nil
}
