package main

import (
	"context"
	"fmt"

	"github.com/ethereum/go-ethereum/common"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	//"github.com/node-real/go-pkg/rlp"
)

const (
	user   = "bijieprd"
	passwd = "Csap2012"
	//DSN    = "bijieprd:Csap2012@tcp(tf-nodereal-prod-ethstates-db-ap-cluster-1.cluster-ro-cbwugcxdizut.ap-northeast-1.rds.amazonaws.com:3306)/eth_mainnet_states?parseTime=true&multiStatements=true"
	DSN   = "bijieprd:Csap2012@tcp(tf-nodereal-prod-ethdataset-global-db-cluster-2.cluster-ro-cbwugcxdizut.ap-northeast-1.rds.amazonaws.com:3306)/eth_mainnet_blocks?parseTime=true&multiStatements=true"
	DSN4X = "bijieprd:CtxfbMM2ov5git6m@tcp(test-qa-states-4xlarge.cluster-custom-cb6vaj1ctcqk.us-east-1.rds.amazonaws.com:3306)/statesdb?parseTime=true&multiStatements=true"
	DSN8X = "bijieprd:CtxfbMM2ov5git6m@tcp(states-8xlarge-single-test.cluster-custom-cb6vaj1ctcqk.us-east-1.rds.amazonaws.com:3306)/statesdb?parseTime=true&multiStatements=true"
)

type Storage struct {
	Address     common.Address `gorm:"column:address;type:BINARY(20);primaryKey"`
	Slot        common.Hash    `gorm:"column:slot;type:BINARY(32);primaryKey"`
	Incarnation uint64         `gorm:"column:incarnation;primaryKey"`
	Number      uint64         `gorm:"column:number;primaryKey"`
	Data        common.Hash    `gorm:"column:data;type:BINARY(32)"`
}

func storagesTable(pid uint64) string {
	return fmt.Sprintf("storages_part%v", pid)
}

type Code struct {
	Hash common.Hash `gorm:"column:hash;type:BINARY(32);primaryKey"`
	Code []byte      `gorm:"column:code;type:BLOB"`
}

func testDb(dsn string, iterNum int) {
	db, err := gorm.Open(mysql.Open(dsn),
		&gorm.Config{
			DisableForeignKeyConstraintWhenMigrating: true,
		},
	)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// storages
	// address 0x0000000000000000000000000000000000001000
	// slot 0x0000000000000000000000000000000000000000000000000000000000000003
	// incarnation 1
	// number 10030001

	fmt.Printf("%s connect successfully", dsn)
	hash := common.Hash{}
	hash.UnmarshalText([]byte("0x00004bbb305c6875f77ea6fa33724f09a0bebe74932b692339002befcdeae316"))
	for i := 0; i < iterNum; i++ {
		//var st Storage
		//if err := db.WithContext(context.Background()).Table(storagesTable(pid)).
		//	Where("address = ? AND slot = ? AND incarnation = ? AND number <= ?", address, slot, incarnation, number).
		//	Order("number DESC").
		//	Limit(1).
		//	Take(&st).Error; err != nil {
		//
		//	if errIsNotFound(err) {
		//		err = nil
		//	}
		//	return common.Hash{}, false, err
		//}
		var s Storage
		var c Code
		if err := db.WithContext(context.Background()).Table("storages_part501").Where("number = ?", 100301221).Take(&s).Error; err != nil {
			fmt.Printf(err.Error())
			return
		}

		fmt.Println(s)
		fmt.Println(c)
		return
		//if err := db.WithContext(context.Background()).Where("hash = ?", hash).Take(&c).Error; err != nil {
		//	fmt.Printf(err.Error())
		//	return
		//}
	}
}

func main() {
	testDb(DSN4X, 10)
}
