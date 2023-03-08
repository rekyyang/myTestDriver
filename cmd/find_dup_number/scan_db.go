package main

//
//import (
//	"context"
//	"database/sql/driver"
//	"fmt"
//	"time"
//
//	"errors"
//	"io"
//	"math/big"
//
//	"github.com/ethereum/go-ethereum/common"
//	"github.com/ethereum/go-ethereum/rlp"
//	"gorm.io/driver/mysql"
//	"gorm.io/gorm"
//
//	"github.com/ethereum/go-ethereum/core/types"
//	//"github.com/node-real/go-pkg/rlp"
//)
//
//const (
//	user   = "bijieprd"
//	passwd = "Csap2012"
//	//DSN    = "bijieprd:Csap2012@tcp(tf-nodereal-prod-ethstates-db-ap-cluster-1.cluster-ro-cbwugcxdizut.ap-northeast-1.rds.amazonaws.com:3306)/eth_mainnet_states?parseTime=true&multiStatements=true"
//	DSN = "bijieprd:Csap2012@tcp(tf-nodereal-prod-ethdataset-global-db-cluster-2.cluster-ro-cbwugcxdizut.ap-northeast-1.rds.amazonaws.com:3306)/eth_mainnet_blocks?parseTime=true&multiStatements=true"
//)
//
//type BlockMeta struct {
//	Hash            common.Hash
//	Number          uint64
//	Time            uint64
//	ParentHash      common.Hash
//	TotalDifficulty *big.Int
//}
//
//type Big big.Int
//
//func (i *Big) Scan(value interface{}) error {
//	bytes, ok := value.([]byte)
//	if !ok {
//		return errors.New(fmt.Sprint("Failed to unmarshal Big value:", value))
//	}
//
//	i.Raw().SetBytes(bytes)
//	return nil
//}
//
//func (i Big) Value() (driver.Value, error) {
//	return i.Raw().Bytes(), nil
//}
//
//func (i *Big) Raw() *big.Int {
//	return (*big.Int)(i)
//}
//
//func (i *Big) DecodeRLP(s *rlp.Stream) error {
//	return s.Decode((*big.Int)(i))
//}
//
//func (i *Big) EncodeRLP(w io.Writer) error {
//	return rlp.Encode(w, (*big.Int)(i))
//}
//
//type Hashes []common.Hash
//
//func (h *Hashes) Scan(value interface{}) error {
//	bytes, ok := value.([]byte)
//	if !ok {
//		return errors.New(fmt.Sprint("Failed to unmarshal Hashes:", value))
//	}
//
//	bytesSize := len(bytes)
//	if bytesSize == 0 {
//		return nil
//	}
//
//	*h = make([]common.Hash, bytesSize/common.HashLength)
//	if len(bytes) != h.Size()*common.HashLength {
//		return errors.New("invalid hashes")
//	}
//	for i := 0; i < h.Size(); i++ {
//		(*h)[i].SetBytes(bytes[i*common.HashLength : (i+1)*common.HashLength])
//	}
//	return nil
//}
//
//func (h Hashes) Value() (driver.Value, error) {
//	count := h.Size()
//	if count == 0 {
//		return make([]byte, 0), nil
//	}
//
//	bytes := make([]byte, count*common.HashLength)
//	for i := 0; i < count; i++ {
//		copy(bytes[i*common.HashLength:(i+1)*common.HashLength], h[i].Bytes())
//	}
//	return bytes, nil
//}
//
//func (h Hashes) Size() int {
//	return len(h)
//}
//
//type BlockKey struct {
//	Hash   common.Hash
//	Number uint64
//}
//
//type BlockKeys []BlockKey
//
//func (h *BlockKeys) Scan(value interface{}) error {
//	bytes, ok := value.([]byte)
//	if !ok {
//		return errors.New(fmt.Sprint("Failed to unmarshal BlockKeys:", value))
//	}
//
//	if len(bytes) == 0 {
//		return nil
//	}
//
//	return rlp.DecodeBytes(bytes, h)
//}
//
//func (h BlockKeys) Value() (driver.Value, error) {
//	if h.Size() == 0 {
//		return make([]byte, 0), nil
//	}
//
//	return rlp.EncodeToBytes(h)
//}
//
//func (h BlockKeys) Size() int {
//	return len(h)
//}
//
//type Bloom types.Bloom
//
//func (b *Bloom) Scan(value interface{}) error {
//	bytes, ok := value.([]byte)
//	if !ok {
//		return errors.New(fmt.Sprint("Failed to unmarshal Bloom:", value))
//	}
//
//	if len(bytes) == 0 {
//		return nil
//	}
//
//	copy(b[:], bytes)
//	return nil
//}
//
//func (b Bloom) Value() (driver.Value, error) {
//	return b[:], nil
//}
//
//type Header struct {
//	ID uint64 `gorm:"column:id;primaryKey" rlp:"-"`
//
//	Hash            common.Hash    `gorm:"column:hash;type:BINARY(32);uniqueIndex:idx_hash"`
//	Number          uint64         `gorm:"column:number;index:idx_number"`
//	ParentHash      common.Hash    `gorm:"column:parent_hash;type:BINARY(32)"`
//	Miner           common.Address `gorm:"column:miner;type:BINARY(20)"`
//	UnclesHash      common.Hash    `gorm:"column:uncles_hash;type:BINARY(32)"`
//	StateRoot       common.Hash    `gorm:"column:state_root;type:BINARY(32)"`
//	TxsRoot         common.Hash    `gorm:"column:txs_root;type:BINARY(32)"`
//	ReceiptsRoot    common.Hash    `gorm:"column:receipts_root;type:BINARY(32)"`
//	LogsBloom       Bloom          `gorm:"column:logs_bloom;type:VARBINARY(256)"`
//	Difficulty      *Big           `gorm:"column:difficulty;type:VARBINARY(32)"`
//	TotalDifficulty *Big           `gorm:"column:total_difficulty;type:VARBINARY(32)"`
//	GasLimit        uint64         `gorm:"column:gas_limit"`
//	GasUsed         uint64         `gorm:"column:gas_used"`
//	Timestamp       uint64         `gorm:"column:timestamp"`
//	MixDigest       common.Hash    `gorm:"column:mix_digest;type:BINARY(32)"`
//	Extra           []byte         `gorm:"column:extra;type:MEDIUMBLOB"`
//	Nonce           uint64         `gorm:"column:nonce"`
//	BaseFee         *Big           `gorm:"column:base_fee;type:VARBINARY(32)" rlp:"nil"`
//
//	UncleKeys BlockKeys `gorm:"column:uncle_keys;type:VARBINARY(512)"` // 512 bytes is enough, max uncle count is 2
//	TxHashes  Hashes    `gorm:"column:tx_hashes;type:MEDIUMBLOB"`
//
//	// rlp encode size (go-ethereum/core/types)
//	HeaderSize uint64 `gorm:"column:header_size"`
//	BlockSize  uint64 `gorm:"column:block_size"`
//}
//
//func (h *Header) ToTypes() *types.Header {
//	nh := &types.Header{
//		ParentHash:  h.ParentHash,
//		UncleHash:   h.UnclesHash,
//		Coinbase:    h.Miner,
//		Root:        h.StateRoot,
//		TxHash:      h.TxsRoot,
//		ReceiptHash: h.ReceiptsRoot,
//		Bloom:       types.Bloom(h.LogsBloom),
//		Difficulty:  h.Difficulty.Raw(),
//		Number:      big.NewInt(0).SetUint64(h.Number),
//		GasLimit:    h.GasLimit,
//		GasUsed:     h.GasUsed,
//		Time:        h.Timestamp,
//		Extra:       h.Extra,
//		MixDigest:   h.MixDigest,
//		Nonce:       types.EncodeNonce(h.Nonce),
//		BaseFee:     h.BaseFee.Raw(),
//	}
//	return nh
//}
//
//func (h *Header) Meta() *BlockMeta {
//	return &BlockMeta{
//		Hash:            h.Hash,
//		Number:          h.Number,
//		Time:            h.Timestamp,
//		ParentHash:      h.ParentHash,
//		TotalDifficulty: h.TotalDifficulty.Raw(),
//	}
//}
//
//func main() {
//	db, err := gorm.Open(mysql.Open(DSN),
//		&gorm.Config{
//			DisableForeignKeyConstraintWhenMigrating: true,
//		},
//	)
//	if err != nil {
//		fmt.Println(err.Error())
//		return
//	}
//
//	fmt.Println("connect successfully")
//
//	//start := 13858856
//	//start := 15000000 // x
//	//start := 15500000 // x
//	//start := 15750000 // y
//	//start := 15625000 // y
//	//start := 15575000 // y
//	//start := 15537500 // y
//	//start := 15518750 // x
//	//start := 15528750 // x
//	start := 15537381 // x
//	//end := 16774955
//	//start := 0
//	end := start + 100000
//
//	// last record dup bn 15537380
//	// Sep-15-2022 06:40:27 AM +UTC)
//	// Sep-15-2022 14:40:27 UTC+8)
//
//	// 北京时间2022年9月15日14时44分 eth2.0合并 15537394
//
//	//var count int64
//	//tableName := fmt.Sprintf("headers_part%v", start/5000000)
//	//_ = db.WithContext(context.Background()).
//	//	Table(tableName).Select("number").
//	//	Where("number >= ? and number <= ?", start, end).Count(&count)
//	//fmt.Println(count)
//	//return
//
//	for blockNumber := start; blockNumber <= end; blockNumber++ {
//		tableName := fmt.Sprintf("headers_part%v", blockNumber/5000000)
//
//		var count int64
//		db_ := db.WithContext(context.Background()).
//			Table(tableName).Select("number").
//			Where("number = ?", blockNumber).Count(&count)
//
//		if count != 1 {
//			fmt.Printf("bn: %d, count: %d\n", blockNumber, count)
//		}
//		_ = db_
//
//		//for i := 0; i < int(count); i++ {
//		//	var hdr Header
//		//	if err := db_.WithContext(context.Background()).
//		//		Table(tableName).
//		//		Where("number = ?", blockNumber).Offset(i).
//		//		Take(&hdr).Error; err != nil {
//		//		fmt.Printf("err : %s, bn : %d\n", err.Error(), blockNumber)
//		//		continue
//		//	}
//		//}
//		time.Sleep(10 * time.Millisecond)
//		if blockNumber%1000 == 0 {
//			fmt.Printf("%d\n", blockNumber)
//		}
//	}
//}
