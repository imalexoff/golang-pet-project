package main

// Request ...
type Request struct {
	CRC    string `json:"CRC"`
	Packet Packet `json:"Packet"`
}

// Packet ...
type Packet struct {
	FromID    string `json:"FromId"`
	ServerKey string `json:"ServerKey"`
	Data      Data   `json:"Data"`
}

// Data ...
type Data struct {
	ContractorID       string `json:"ContractorId"`
	GoodsGroupID       string `json:"GoodsGroupId"`
	Page               int    `json:"Page"`
	Search             string `json:"Search"`
	OrderBy            int    `json:"OrderBy"`
	OrderByContractor  int    `json:"OrderByContractor"`
	CompareOntractorID int    `json:"CompareСontractorId"`
	CatalogType        int    `json:"CatalogType"`
}

// GoodsOffer ...
type GoodsOffer struct {
	GoodsGroupName string  `json:"GoodsGroupName"`
	GoodsID        int     `json:"GoodsId"`
	GoodsName      string  `json:"GoodsName"`
	Offers         []Offer `json:"Offers"`
}

// Offer ...
type Offer struct {
	ContractorID   int    `json:"ContractorId"`
	AddressShop    string `json:"AddressShop"`
	MonitoringDate string `json:"MonitoringDate"`
	Price          string `json:"Price"`
}

// Contractor ...
type Contractor struct {
	ContractorID   int    `json:"ContractorId"`
	ContractorName string `json:"ContractorName"`
}
