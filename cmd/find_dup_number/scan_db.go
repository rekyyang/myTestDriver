package main

import (
	"context"
	"fmt"
	"time"

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
		fmt.Println(err.Error())
		return
	}

	fmt.Println("connect successfully")

	//end := 16774955
	end := 20000

	for blockNumber := 0; blockNumber < end; blockNumber++ {
		var hdr models.Header
		tableName := fmt.Sprintf("headers_part%v", blockNumber/5000000)
		if err := db.WithContext(context.Background()).
			Table(tableName).
			Where("number = ?", blockNumber).
			Take(&hdr).Error; err != nil {
			fmt.Printf("err : %s, bn : %d\n", err.Error(), blockNumber)
			continue
		}
		time.Sleep(10 * time.Millisecond)
		if blockNumber%1000 == 0 {
			fmt.Printf("%d", blockNumber)
		}
	}
}
