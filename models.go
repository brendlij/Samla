package main

type StorageLocation struct {
	ID           int64  `json:"id"`
	FriendlyName string `json:"friendlyName"`
	Room         string `json:"room"`
	Shelf        string `json:"shelf"`
	Compartment  string `json:"compartment"`
	Note         string `json:"note"`
}

type Box struct {
	ID         int64  `json:"id"`
	LocationID int64  `json:"locationId"`
	Code       string `json:"code"`
	Name       string `json:"name"`
}

type Bag struct {
	ID       int64  `json:"id"`
	BoxID    int64  `json:"boxId"`
	SerialNo string `json:"serialNo"`
}

type Manufacturer struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type Type struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
}

type Product struct {
	ID    int64  `json:"id"`
	SetID int64  `json:"setId"`
	Name  string `json:"name"`
	Kind  string `json:"kind"`
}

type BagInfo struct {
	ID                  int64  `json:"id"`
	SerialNo            string `json:"serialNo"`
	BoxID               int64  `json:"boxId"`
	BoxCode             string `json:"boxCode"`
	BoxName             string `json:"boxName"`
	LocationID          int64  `json:"locationId"`
	LocationName        string `json:"locationName"`
	LocationRoom        string `json:"locationRoom"`
	LocationShelf       string `json:"locationShelf"`
	LocationCompartment string `json:"locationCompartment"`
	LocationNote        string `json:"locationNote"`
}

type SetDetails struct {
	ID               int64     `json:"id"`
	Name             string    `json:"name"`
	ManufacturerID   *int64    `json:"manufacturerId"`
	ManufacturerName string    `json:"manufacturerName"`
	TypeID           *int64    `json:"typeId"`
	TypeName         string    `json:"typeName"`
	Bag              BagInfo   `json:"bag"`
	PhotoPath        string    `json:"photoPath"`
	PhotoSource      string    `json:"photoSource"`
	Tags             []string  `json:"tags"`
	Products         []Product `json:"products"`
}

type SetSearchResult struct {
	SetID            int64    `json:"setId"`
	SetName          string   `json:"setName"`
	ManufacturerName string   `json:"manufacturerName"`
	BoxCode          string   `json:"boxCode"`
	BoxName          string   `json:"boxName"`
	BagSerial        string   `json:"bagSerial"`
	LocationName     string   `json:"locationName"`
	Tags             []string `json:"tags"`
	ThumbnailPath    string   `json:"thumbnailPath"`
}
