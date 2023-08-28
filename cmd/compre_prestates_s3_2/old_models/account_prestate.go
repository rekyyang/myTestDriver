package old_models

import (
	"bytes"
	"io"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
)

type AccountPrestate struct {
	Nonce       uint64
	Incarnation uint64
	Balance     *big.Int
	Root        []byte
	CodeHash    []byte
	//AccountCode []byte
	Storages map[common.Hash]common.Hash
}

type Prestates struct {
	AccountPrestateMap map[common.Address]AccountPrestate
}

type PrestatesRLP Prestates

var (
	// DeletedRoot to mark a deleted account
	// 0xfec91ad4a0d3a97a5aa5d0b8b79f71ff5a63866a2be7950c6bac67bce785d708
	DeletedRoot = crypto.Keccak256Hash([]byte("deleted"))

	// EmptyRoot is the known root hash of an empty trie.
	EmptyRoot = common.HexToHash("56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421")

	// EmptyCodeHash is the known hash of the empty EVM bytecode.
	EmptyCodeHash = crypto.Keccak256Hash(nil)
)

func (p *Prestates) ContainState() bool {
	return len(p.AccountPrestateMap) != 0
}

func (p *Prestates) EncodeRLP(_w io.Writer) error {
	a1 := PrestatesRLP(*p)
	w := rlp.NewEncoderBuffer(_w)
	_tmp0 := w.List()
	mapLen := len(a1.AccountPrestateMap)
	w.WriteUint64(uint64(mapLen))
	if mapLen == 0 {
		w.ListEnd(_tmp0)
		return w.Flush()
	}
	for addr, accountPrestate := range a1.AccountPrestateMap {
		w.WriteBytes(addr[:])
		w.WriteUint64(accountPrestate.Nonce)
		w.WriteUint64(accountPrestate.Incarnation)
		storageMapLen := len(accountPrestate.Storages)
		w.WriteUint64(uint64(storageMapLen))
		if accountPrestate.Balance == nil {
			w.WriteBigInt(big.NewInt(0))
		} else {
			if accountPrestate.Balance.Sign() == -1 {
				return rlp.ErrNegativeBigInt
			}
			w.WriteBigInt(accountPrestate.Balance)
		}
		if bytes.Equal(accountPrestate.Root, EmptyRoot.Bytes()) {
			w.WriteBytes(rlp.EmptyString)
		} else {
			w.WriteBytes(accountPrestate.Root)
		}
		if bytes.Equal(accountPrestate.CodeHash, EmptyCodeHash.Bytes()) {
			w.WriteBytes(rlp.EmptyString)
		} else {
			w.WriteBytes(accountPrestate.CodeHash)
		}
		//if bytes.Equal(accountPrestate.AccountCode, EmptyCodeHash.Bytes()) {
		//	w.WriteBytes(rlp.EmptyString)
		//} else {
		//	w.WriteBytes(accountPrestate.AccountCode)
		//}
		for storageHash, storageData := range accountPrestate.Storages {
			w.WriteBytes(storageHash[:])
			w.WriteBytes(storageData[:])
		}
	}
	w.ListEnd(_tmp0)
	return w.Flush()
}

func (p *Prestates) DecodeRLP(s *rlp.Stream) error {
	var _tmp Prestates
	if _, err := s.List(); err != nil {
		return err
	}
	mapLen, err := s.Uint64()
	if err != nil {
		return err
	}
	_tmp.AccountPrestateMap = make(map[common.Address]AccountPrestate)
	if mapLen == 0 {
		//*p = _tmp
		if err = s.ListEnd(); err != nil {
			return err
		}
		return nil
	}
	for i := uint64(0); i < mapLen; i++ {
		var addr common.Address
		if err = s.ReadBytes(addr[:]); err != nil {
			return err
		}
		var accountPrestate AccountPrestate
		nonce, err1 := s.Uint64()
		if err1 != nil {
			return err1
		}
		accountPrestate.Nonce = nonce
		incarnation, err1 := s.Uint64()
		if err1 != nil {
			return err1
		}
		accountPrestate.Incarnation = incarnation
		storageMapLen, err1 := s.Uint64()
		if err1 != nil {
			return err1
		}
		balance, err1 := s.BigInt()
		if err1 != nil {
			return err1
		}
		accountPrestate.Balance = balance
		root, err1 := s.Bytes()
		if err1 != nil {
			return err1
		}
		if bytes.Equal(root, rlp.EmptyString) {
			accountPrestate.Root = EmptyRoot.Bytes()
		} else {
			accountPrestate.Root = root
		}
		codeHash, err1 := s.Bytes()
		if err1 != nil {
			return err1
		}
		if bytes.Equal(codeHash, rlp.EmptyString) {
			accountPrestate.CodeHash = EmptyCodeHash.Bytes()
		} else {
			accountPrestate.CodeHash = codeHash
		}
		//accountCode, err1 := s.Bytes()
		//if err1 != nil {
		//	return err1
		//}
		//if bytes.Equal(accountCode, rlp.EmptyString) {
		//	accountPrestate.AccountCode = nil
		//} else {
		//	accountPrestate.AccountCode = accountCode
		//}
		accountPrestate.Storages = make(map[common.Hash]common.Hash)
		if storageMapLen != 0 {
			for j := uint64(0); j < storageMapLen; j++ {
				var storageHash common.Hash
				if err1 = s.ReadBytes(storageHash[:]); err1 != nil {
					return err1
				}
				var storageData common.Hash
				if err1 = s.ReadBytes(storageData[:]); err1 != nil {
					return err1
				}
				accountPrestate.Storages[storageHash] = storageData
			}
		}
		_tmp.AccountPrestateMap[addr] = accountPrestate
	}
	if err = s.ListEnd(); err != nil {
		return err
	}
	*p = _tmp
	return nil
}
