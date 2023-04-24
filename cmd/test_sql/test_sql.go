package main

import (
	"context"
	"fmt"
	"time"

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

func testDb(dsn string, iterNum int, label string) {
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
	// number 20010100
	// data 0x000000000000000000000000000000000000000000000070f3866ee30d9a3a38
	var address common.Address
	address.UnmarshalText([]byte("0x0000000000000000000000000000000000001000"))
	var number = 20013101

	// code
	fmt.Printf("%s connect successfully \n", dsn)
	hash := common.Hash{}
	hash.UnmarshalText([]byte("0x00004bbb305c6875f77ea6fa33724f09a0bebe74932b692339002befcdeae316"))
	startTime := time.Now()
	for i := 0; i < iterNum; i++ {
		var s Storage
		if err := db.WithContext(context.Background()).
			Table("storages_part1000").
			Where("address = ? AND number <= ?", address, number).
			Order("number DESC").
			Limit(1).Take(&s).Error; err != nil {
			fmt.Printf(err.Error())
			return
		}

		var c Code
		if err := db.WithContext(context.Background()).Where("hash = ?", hash).Take(&c).Error; err != nil {
			fmt.Printf(err.Error())
			return
		}

		if i == 0 {
			_ = s
			_ = c
			//fmt.Println(s.Number)
			//fmt.Println(c.Hash)
		}
	}
	endTime := time.Now()
	fmt.Println()
	fmt.Println()
	fmt.Printf("%s latency: %v", label, endTime.Sub(startTime))
	fmt.Println()
	fmt.Println()
}

func main() {
	//testDb(DSN4X, 1, "4x")
	//testDb(DSN8X, 1, "8x")

	fmt.Println(time.Until(time.Now().Add(time.Minute)))
	fmt.Println(time.Since(time.Now().Add(time.Minute)))
}
