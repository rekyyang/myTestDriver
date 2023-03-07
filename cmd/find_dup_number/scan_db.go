package main

import (
	"context"
	"fmt"
	"log"

	"github.com/node-real/blocktree/models"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	user   = "bijieprd"
	passwd = "Csap2012"
	DSN    = "bijieprd:Csap2012@tcp(tf-nodereal-prod-ethstates-db-ap-cluster-1.cluster-ro-cbwugcxdizut.ap-northeast-1.rds.amazonaws.com:3306)/eth_mainnet_states?parseTime=true&multiStatements=true"
)

func main() {
	db, err := gorm.Open(mysql.Open(DSN),
		&gorm.Config{
			DisableForeignKeyConstraintWhenMigrating: true,
		},
	)
	if err != nil {
		log.Panicf(err.Error())
	}
	for blockNumber := 0; blockNumber < 16774955; blockNumber++ {
		var hdr models.Header
		tableName := fmt.Sprintf("headers_part%v", blockNumber/5000000)
		if err := db.WithContext(context.Background()).
			Table(tableName).
			Where("number = ?", blockNumber).
			Take(&hdr).Error; err != nil {
			fmt.Printf("err : %s, bn : %d\n", err.Error(), blockNumber)
			continue
		}
	}
}
