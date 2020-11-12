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

func parse(groups []string, maxThreadsCount int) (err error) {
	openFile()
	defer closeFile()

	var wgGroups sync.WaitGroup
	wgGroups.Add(len(groups))

	for _, group := range groups {
		go func(group string) {
			defer wgGroups.Done()

			var wgGoods sync.WaitGroup

			chunks := chunks(group, maxThreadsCount)
			wgGoods.Add(len(chunks))

			for _, chunk := range chunks {
				go func(chunk []int) {
					defer wgGoods.Done()

					goods := getGoodsByPages(chunk, group)

					export(goods)
				}(chunk)
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

func getGoodsByPages(chunk []int, group string) (goods []GoodsOffer) {
	req := newRequest()
	req.Packet.Data.GoodsGroupID = group

	for i := 0; i < len(chunk); i++ {
		req.Packet.Data.Page = chunk[i]

		data, _ := getGoods(req)
		goods = append(goods, data...)

		logf("Group %v, Page %v is done", group, chunk[i])
	}

	return goods
}

func export(pageGoods []GoodsOffer) {
	for _, goods := range pageGoods {
		for _, offer := range goods.Offers {
			writeToFile(fmt.Sprintf("%v|%v|%v|%v|%v|%v|%v\n", contractorsMap[offer.ContractorID], offer.AddressShop, offer.MonitoringDate, goods.GoodsGroupName, goods.GoodsID, goods.GoodsName, offer.Price))
		}
	}
}

func chunks(group string, maxThreadsCount int) (chunks [][]int) {
	req := newRequest()
	req.Packet.Data.GoodsGroupID = group

	_, amountPages := getGoods(req)

	chankCapacity := amountPages / maxThreadsCount
	var threadsCount int
	if chankCapacity > 0 {
		threadsCount = maxThreadsCount
	} else {
		threadsCount = amountPages
		chankCapacity = 1
	}

	pages := make([]int, amountPages)
	for i := range pages {
		pages[i] = i + 1
	}

	remainsPages := pages[amountPages-(amountPages%threadsCount):]

	chunks = make([][]int, threadsCount)
	for i := 0; i < threadsCount; i++ {
		skip := i * chankCapacity
		take := skip + chankCapacity

		chunks[i] = make([]int, chankCapacity, chankCapacity+1)
		copy(chunks[i], pages[skip:take])

		if i < len(remainsPages) {
			chunks[i] = append(chunks[i], remainsPages[i:i+1]...)
			//chunks[i][len(chunks[i])-1] = remainsPages[i : i+1][i]
		}
	}

	return chunks
}
