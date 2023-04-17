package main

import (
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

const (
	DSN = "bijieprd:sFOUsoMtvrWIKp9ppjYR@tcp(tf-nodereal-qa-dataplatform-db-db.cluster-cb6vaj1ctcqk.us-east-1.rds.amazonaws.com:3306)/uniswapv3?parseTime=true&multiStatements=true&loc=Local"
)

func main() {
	db, err := gorm.Open(mysql.Open(DSN), &gorm.Config{})
	if err != nil {
		// 处理错误
	}
	sqlDB, err := db.DB()
	if err != nil {
		// 处理错误
	}

	// 设置最大打开连接数
	sqlDB.SetMaxOpenConns(100)

	// 设置最大空闲连接数
	sqlDB.SetMaxIdleConns(10)

	// 设置连接最大寿命
	sqlDB.SetConnMaxLifetime(time.Hour)

	var item struct {
		BlockNumber uint64 `gorm:"column:block_number"`
	}
	// select * from pancake_txs limit 1;
	if err := db.Table("v3_transaction").Limit(1).Scan(&item).Error; err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(item)

	defer sqlDB.Close()
}
