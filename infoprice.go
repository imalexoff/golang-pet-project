package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/tidwall/gjson"
)

type request struct {
	CRC    string `json:"CRC"`
	Packet packet `json:"Packet"`
}

type packet struct {
	FromID    string `json:"FromId"`
	ServerKey string `json:"ServerKey"`
	Data      data   `json:"Data"`
}

type data struct {
	ContractorID       string `json:"ContractorId"`
	GoodsGroupID       string `json:"GoodsGroupId"`
	Page               int    `json:"Page"`
	Search             string `json:"Search"`
	OrderBy            int    `json:"OrderBy"`
	OrderByContractor  int    `json:"OrderByContractor"`
	CompareOntractorID int    `json:"CompareСontractorId"`
	CatalogType        int    `json:"CatalogType"`
}

type goodsOffer struct {
	GoodsID   int     `json:"GoodsId"`
	GoodsName string  `json:"GoodsName"`
	Offers    []offer `json:"Offers"`
}

type offer struct {
	ContractorID   int    `json:"ContractorId"`
	AddressShop    string `json:"AddressShop"`
	MonitoringDate string `json:"MonitoringDate"`
}

type contractor struct {
	ContractorID   int    `json:"ContractorId"`
	ContractorName string `json:"ContractorName"`
}

var contractors map[int]string

func getGoods(groups string) error {
	openFile()
	defer closeFile()

	request := initRequest(groups)
	initContractors(*request)

	var currPage int = 1

	for {
		request.Packet.Data.Page = currPage
		reqBody, _ := json.Marshal(request)

		resp, err := http.Post("https://api.infoprice.by/InfoPrice.Goods?v=2", "application/json", bytes.NewBuffer(reqBody))
		if err != nil {
			log(err.Error())
			return err
		}
		respData, _ := ioutil.ReadAll(resp.Body)
		resp.Body.Close()

		amountPages := gjson.GetBytes(respData, "Table.0.GeneralData.0.AmountPages").Num
		goodsJSON := gjson.GetBytes(respData, "Table.0.GoodsOffer").Raw

		log(fmt.Sprintf("page %v from %v", currPage, amountPages))

		var pageGoods []goodsOffer
		json.Unmarshal([]byte(goodsJSON), &pageGoods)

		export(pageGoods)

		if currPage == int(amountPages) {
			break
		}

		currPage++
	}

	return nil
}

func initRequest(groups string) *request {
	return &request{
		Packet: packet{
			FromID:    "10003001",
			ServerKey: "omt5W465fjwlrtxcEco97kew2dkdrorqqq",
			Data: data{
				GoodsGroupID: groups,
			},
		},
	}
}

func initContractors(request request) {
	reqBody, _ := json.Marshal(request)
	resp, err := http.Post("https://api.infoprice.by/InfoPrice.Contractors?v=3", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		log(err.Error())
		panic(err)
	}
	defer resp.Body.Close()

	respData, _ := ioutil.ReadAll(resp.Body)
	contractorsJSON := gjson.GetBytes(respData, "Table").Raw
	var _contractors []contractor
	json.Unmarshal([]byte(contractorsJSON), &_contractors)

	contractors = make(map[int]string)
	for _, v := range _contractors {
		contractors[v.ContractorID] = v.ContractorName
	}
}

func export(pageGoods []goodsOffer) {
	for _, goods := range pageGoods {
		for _, offer := range goods.Offers {
			writeToFile(fmt.Sprintf("[%v][%v][%v][%v][%v]\n", contractors[offer.ContractorID], offer.AddressShop, offer.MonitoringDate, goods.GoodsID, goods.GoodsName))
		}
	}
}
