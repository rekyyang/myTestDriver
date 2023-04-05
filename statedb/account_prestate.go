package statedb

import (
	"bytes"
	"io"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/node-real/blocktree"
	"github.com/node-real/blocktree/models"
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
	accountPrestateMap       map[common.Address]AccountPrestate
	view                     blocktree.View
	accountPrestateMapChange bool
	allowUpload              bool
	rwMutex                  sync.RWMutex
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

func NewPreState(view blocktree.View) *Prestates {
	return &Prestates{
		accountPrestateMap:       make(map[common.Address]AccountPrestate),
		view:                     view,
		accountPrestateMapChange: false,
		allowUpload:              false,
		rwMutex:                  sync.RWMutex{},
	}
}

func (p *Prestates) ContainState() bool {
	return len(p.accountPrestateMap) != 0
}

func (p *Prestates) EncodeRLP(_w io.Writer) error {
	a1 := (*PrestatesRLP)(p)
	w := rlp.NewEncoderBuffer(_w)
	_tmp0 := w.List()
	mapLen := len(a1.accountPrestateMap)
	w.WriteUint64(uint64(mapLen))
	if mapLen == 0 {
		w.ListEnd(_tmp0)
		return w.Flush()
	}
	for addr, accountPrestate := range a1.accountPrestateMap {
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
	_tmp.accountPrestateMap = make(map[common.Address]AccountPrestate)
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
		_tmp.accountPrestateMap[addr] = accountPrestate
	}
	if err = s.ListEnd(); err != nil {
		return err
	}
	p.accountPrestateMap = _tmp.accountPrestateMap
	//*p = _tmp
	return nil
}

func (p *Prestates) Copy() *Prestates {
	tmp := NewPreState(p.view)
	for addr, accountPrestate := range p.accountPrestateMap {
		tmpAccountPrestate := AccountPrestate{
			Nonce:       accountPrestate.Nonce,
			Incarnation: accountPrestate.Incarnation,
			//Balance:     accountPrestate.Balance,
			Root:     accountPrestate.Root[:],
			CodeHash: accountPrestate.CodeHash[:],
			//AccountCode: accountPrestate.AccountCode[:],
			Storages: make(map[common.Hash]common.Hash),
		}
		if accountPrestate.Balance == nil {
			tmpAccountPrestate.Balance = new(big.Int)
		} else {
			//tmpAccountPrestate.Balance = big.NewInt(accountPrestate.Balance.Int64())
			//tmpBigInt := new(big.Int)
			//tmpBigInt.SetBytes(accountPrestate.Balance.Bytes())
			//tmpAccountPrestate.Balance = tmpBigInt
			tmpAccountPrestate.Balance = new(big.Int).Set(accountPrestate.Balance)
		}
		for k, v := range accountPrestate.Storages {
			tmpAccountPrestate.Storages[k] = v
		}
		tmp.accountPrestateMap[addr] = tmpAccountPrestate
	}
	tmp.accountPrestateMapChange = p.accountPrestateMapChange
	tmp.allowUpload = false
	return tmp
}

func (p *Prestates) IsArchive() bool {
	return p.view.IsArchive()
}

func (p *Prestates) Storage(addr common.Address, incarnation uint64, slot common.Hash) (common.Hash, error) {
	p.rwMutex.RLock()
	if accountPrestate, ok := p.accountPrestateMap[addr]; ok {
		if _, ok1 := accountPrestate.Storages[slot]; ok1 {
			p.rwMutex.RUnlock()
			return accountPrestate.Storages[slot], nil
		}
	}
	p.rwMutex.RUnlock()
	value, err := p.view.Storage(addr, incarnation, slot)
	if err != nil {
		return common.Hash{}, err
	}
	if p.allowUpload {
		p.accountPrestateMapChange = true
		p.rwMutex.Lock()
		//log.Infof("Prestates Storage, slot:%s, data:%s, blockNumber:%d", slot.Hex(), value.Hex(), p.view.CurrentBlockNumber())
		accountPrestateTmp, ok := p.accountPrestateMap[addr]
		if ok {
			accountPrestateTmp.Storages[slot] = value
		} else {
			accountPrestateTmp = AccountPrestate{}
			accountPrestateTmp.Storages = make(map[common.Hash]common.Hash)
			accountPrestateTmp.Storages[slot] = value
			p.accountPrestateMap[addr] = accountPrestateTmp
		}
		p.rwMutex.Unlock()
	}
	return value, err
}

func (p *Prestates) ContractCode(addr common.Address, hash common.Hash) ([]byte, error) {
	//p.rwMutex.RLock()
	//accountPrestate, ok := p.accountPrestateMap[addr]
	//p.rwMutex.RUnlock()
	//if ok && len(accountPrestate.AccountCode) != 0 {
	//	return accountPrestate.AccountCode, nil
	//}

	// get code and save the code only when code is not nil
	//code, err := p.view.ContractCode(hash)
	//if err == nil && len(code) != 0 && p.allowUpload {
	//	p.accountPrestateMapChange = true
	//	p.rwMutex.Lock()
	//	//log.Infof("Prestates ContractCode code(%d), blockNumber:%d", len(code), p.view.CurrentBlockNumber())
	//	accountPrestateTmp, ok1 := p.accountPrestateMap[addr]
	//	if ok1 {
	//		accountPrestateTmp.AccountCode = code
	//	} else {
	//		accountPrestateTmp = AccountPrestate{}
	//		accountPrestateTmp.Storages = make(map[common.Hash]common.Hash)
	//		accountPrestateTmp.AccountCode = code
	//		p.accountPrestateMap[addr] = accountPrestateTmp
	//	}
	//	p.rwMutex.Unlock()
	//}
	// code is too big, determine not save
	return p.view.ContractCode(hash)
}

func (p *Prestates) Account(addr common.Address) (*models.Account, error) {
	p.rwMutex.RLock()
	accountPrestateTmp, ok := p.accountPrestateMap[addr]
	p.rwMutex.RUnlock()
	if ok {
		acc := &models.Account{
			Nonce:       accountPrestateTmp.Nonce,
			Balance:     accountPrestateTmp.Balance,
			Root:        accountPrestateTmp.Root,
			CodeHash:    accountPrestateTmp.CodeHash,
			Incarnation: accountPrestateTmp.Incarnation,
		}
		//log.Infof("getPA suc, b s:%d", acc.Balance.Sign())
		return acc, nil
	}
	acc, err := p.view.Account(addr)
	if err == nil && acc != nil && p.allowUpload {
		p.accountPrestateMapChange = true
		p.rwMutex.Lock()
		//log.Infof("Prestates r(%d), i(%d), n(%d), b(%v), c(%d), number:%d", len(acc.Root),
		//	acc.Incarnation, acc.Nonce, acc.Balance, len(acc.CodeHash), p.view.CurrentBlockNumber())
		accountPrestate, ok := p.accountPrestateMap[addr]
		if ok {
			accountPrestate.Root, accountPrestate.Incarnation, accountPrestate.Nonce, accountPrestate.Balance,
				accountPrestate.CodeHash = acc.Root, acc.Incarnation, acc.Nonce, acc.Balance, acc.CodeHash
		} else {
			accountPrestate = AccountPrestate{}
			accountPrestate.Storages = make(map[common.Hash]common.Hash)
			accountPrestate.Root, accountPrestate.Incarnation, accountPrestate.Nonce, accountPrestate.Balance,
				accountPrestate.CodeHash = acc.Root, acc.Incarnation, acc.Nonce, acc.Balance, acc.CodeHash
			p.accountPrestateMap[addr] = accountPrestate
		}
		p.rwMutex.Unlock()
	}
	return acc, err
}
