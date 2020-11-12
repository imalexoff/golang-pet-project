package main

import (
	"fmt"
	"os"
	"time"
)

var file *os.File

func openFile() {
	now := time.Now()
	timestamp := now.UnixNano() / int64(time.Millisecond)
	nowLabel := now.Format("15-04-05")
	fileName := fmt.Sprintf("goods-%v-%v.csv", nowLabel, timestamp)

	logf("Create file - %v", fileName)

	f, err := os.Create(fileName)
	f.WriteString("ContractorName|AddressShop|MonitoringDate|GoodsGroupName|GoodsId|GoodsName|Price\n")

	if err != nil {
		panic(err)
	}

	file = f
}

func closeFile() {
	log("Close file...")

	file.Close()
}

func writeToFile(message string) {
	file.WriteString(message)
}
