package model

type Hotelier struct {
	HotelierId int      `json:"HotelierId"`
	Name       string   `json:"Name"`
	Hotels     []*Hotel `json:"Hotels"`
}
