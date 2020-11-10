package main

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/tidwall/gjson"
)

var contractorsMap map[int]string

func init() {
	contractors := getContractors()
	contractorsMap = make(map[int]string)
	for _, value := range contractors {
		contractorsMap[value.ContractorID] = value.ContractorName
	}
}

func parse(groups []string, threadsCount int) (err error) {
	openFile()
	defer closeFile()

	var wgGroups sync.WaitGroup
	wgGroups.Add(len(groups))

	for _, group := range groups {
		go func(group string) {
			defer wgGroups.Done()

			var wgGoods sync.WaitGroup
			wgGoods.Add(threadsCount)

			packetsCount := calcPackets(group, threadsCount)

			for i := 0; i < threadsCount; i++ {
				go func(i int) {
					defer wgGoods.Done()

					startPage := (i * packetsCount) + 1
					endPage := (i * packetsCount) + packetsCount
					goods := getGoodsByPages(startPage, endPage, group)

					export(goods)
				}(i)
			}

			wgGoods.Wait()
		}(group)
	}

	wgGroups.Wait()

	return nil
}

func getContractors() (contractors []Contractor) {
	req := newRequest()
	data := executeRequest("https://api.infoprice.by/InfoPrice.Contractors?v=3", req)

	contractorsJSON := gjson.GetBytes(data, "Table").Raw

	json.Unmarshal([]byte(contractorsJSON), &contractors)

	return contractors
}

func getGoods(req *Request) (goods []GoodsOffer, amountPages int) {
	data := executeRequest("https://api.infoprice.by/InfoPrice.Goods?v=2", req)

	amountPages = int(gjson.GetBytes(data, "Table.0.GeneralData.0.AmountPages").Num)
	goodsJSON := gjson.GetBytes(data, "Table.0.GoodsOffer").Raw

	json.Unmarshal([]byte(goodsJSON), &goods)

	return goods, amountPages
}

func getGoodsByPages(startPage, endPage int, group string) (goods []GoodsOffer) {
	req := newRequest()
	req.Packet.Data.GoodsGroupID = group

	for i := startPage; i <= endPage; i++ {
		req.Packet.Data.Page = i

		data, _ := getGoods(req)
		goods = append(goods, data...)

		logf("Page %v is done", i)
	}

	return goods
}

func export(pageGoods []GoodsOffer) {
	for _, goods := range pageGoods {
		for _, offer := range goods.Offers {
			writeToFile(fmt.Sprintf("[%v][%v][%v][%v][%v]\n", contractorsMap[offer.ContractorID], offer.AddressShop, offer.MonitoringDate, goods.GoodsID, goods.GoodsName))
		}
	}
}

func calcPackets(group string, threadsCount int) (packetsCount int) {
	req := newRequest()
	req.Packet.Data.GoodsGroupID = group

	_, amountPages := getGoods(req)

	if amountPages%threadsCount == 0 {
		packetsCount = amountPages / threadsCount
	} else {
		packetsCount = (amountPages / threadsCount) + 1
	}

	return packetsCount
}
