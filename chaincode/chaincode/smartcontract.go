package chaincode

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

// SmartContract provides functions for managing an Asset
type SmartContract struct {
	contractapi.Contract
}

/* 해야하는 일
- asset 구조 정의
물품 거래 일자 및 시간 : timestamp string
물품 거래 가능 구역 : location string
판매자 이름, 구매자 이름, : seller, customer
판매자 인증서, 구매자 인증서 : seller_cert, customer_cert
중고 물품명, : name
물품 가격 : price
거래 상태(게시, 예약, 거래 완료) : state

- 성능 평가
- 내용 채워넣기
*/

type Asset struct {
	ID              string `json:"ID"`
	ProductName     string `json:"ProductName"`
	ProductPrice    int    `json:"ProductPrice"`
	TransactionDate string `json:"TransactionDate"`
	Seller          string `json:"Seller"`
	Buyer           string `json:"Buyer"`
	Location        string `json:"Location"`
	State           string `json:"State"`
}
type Participant struct {
	Name        string `json:"Name"`
	Certificate string `json:"Certificate"`
}
type CCEvent struct {
	Asset       Asset       `json:"Asset"`
	Participant Participant `json:"Participant"`
}

// InitLedger adds a base set of assets to the ledger
func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	txTimeAsPtr, errTx := s.GetTxTimestampChannel(ctx)
	if errTx != nil {
		return fmt.Errorf("Returning error in Transaction TimeStamp")
	}
	assets := []Asset{
		{ID: "asset1", ProductName: "BookA", ProductPrice: 1, TransactionDate: txTimeAsPtr, Seller: "seller1", Buyer: "buyer1", Location: "Busan", State: "on_sale"},
		{ID: "asset2", ProductName: "PencilB", ProductPrice: 2, TransactionDate: txTimeAsPtr, Seller: "seller2", Buyer: "buyer2", Location: "Busan", State: "on_sale"},
		{ID: "asset3", ProductName: "CellphoneC", ProductPrice: 5, TransactionDate: txTimeAsPtr, Seller: "seller3", Buyer: "buyer3", Location: "Busan", State: "booking"},
		{ID: "asset4", ProductName: "LaptopD", ProductPrice: 5, TransactionDate: txTimeAsPtr, Seller: "seller4", Buyer: "buyer4", Location: "Daegu", State: "booking"},
		{ID: "asset5", ProductName: "BookE", ProductPrice: 5, TransactionDate: txTimeAsPtr, Seller: "seller5", Buyer: "buyer5", Location: "Seoul", State: "sold_out"},
		{ID: "asset6", ProductName: "BookF", ProductPrice: 5, TransactionDate: txTimeAsPtr, Seller: "seller6", Buyer: "buyer6", Location: "Seoul", State: "sold_out"},
	}

	for _, asset := range assets {
		assetJSON, err := json.Marshal(asset)
		if err != nil {
			return err
		}

		err = ctx.GetStub().PutState(asset.ID, assetJSON)
		if err != nil {
			return fmt.Errorf("failed to put to world state. %v", err)
		}
	}

	return nil
}

// CreateAsset issues a new asset to the world state with given details.
func (s *SmartContract) CreateAsset(ctx contractapi.TransactionContextInterface, id string, name string, price int, timestamp string, seller string, buyer string, location string, state string) error {
	exists, err := s.AssetExists(ctx, id)
	if err != nil {
		return err
	}
	if exists {
		return fmt.Errorf("the asset %s already exists", id)
	}

	asset := Asset{
		ID:              id,
		ProductName:     name,
		ProductPrice:    price,
		TransactionDate: timestamp,
		Seller:          seller,
		Buyer:           buyer,
		Location:        location,
		State:           state,
	}
	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}
	creator, err := ctx.GetStub().GetCreator()
	if err != nil {
		return err
	}

	certificate := base64.StdEncoding.EncodeToString(creator)
	participant := Participant{
		Name:        seller,
		Certificate: certificate,
	}
	eventData := CCEvent{
		Asset:       asset,
		Participant: participant,
	}
	eventJSON, err := json.Marshal(eventData)
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	err = ctx.GetStub().PutState(id, assetJSON)
	log.Println("event:" + string(eventJSON) + " , asset: " + string(assetJSON))
	if err != nil {
		return fmt.Errorf("failed to store the asset: %v", err)
	}
	err = ctx.GetStub().SetEvent("CreateAsset", eventJSON)
	if err != nil {
		log.Println(string(err.Error()))
		return err
	}
	return err
}

// ReadAsset returns the asset stored in the world state with given id.
func (s *SmartContract) ReadAsset(ctx contractapi.TransactionContextInterface, id string) (*Asset, error) {
	assetJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return nil, fmt.Errorf("failed to read from world state: %v", err)
	}
	if assetJSON == nil {
		return nil, fmt.Errorf("the asset %s does not exist", id)
	}

	var asset Asset
	err = json.Unmarshal(assetJSON, &asset)
	if err != nil {
		return nil, err
	}
	return &asset, nil
}

// UpdateAsset updates an existing asset in the world state with provided parameters.
func (s *SmartContract) UpdateAsset(ctx contractapi.TransactionContextInterface, id string, name string, price int, timestamp string, seller string, buyer string, location string, state string) error {
	exists, err := s.AssetExists(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("the asset %s does not exist", id)
	}

	// overwriting original asset with new asset
	asset := Asset{
		ID:              id,
		ProductName:     name,
		ProductPrice:    price,
		TransactionDate: timestamp,
		Seller:          seller,
		Buyer:           buyer,
		Location:        location,
		State:           state,
	}
	creator, err := ctx.GetStub().GetCreator()
	if err != nil {
		return err
	}

	certificate := base64.StdEncoding.EncodeToString(creator)
	participant := Participant{
		Name:        seller,
		Certificate: certificate,
	}
	eventData := CCEvent{
		Asset:       asset,
		Participant: participant,
	}
	eventJSON, err := json.Marshal(eventData)
	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return err
	}
	err = ctx.GetStub().PutState(id, assetJSON)
	if err != nil {
		return fmt.Errorf("failed to store the asset: %v", err)
	}
	ctx.GetStub().SetEvent("UpdateAsset", eventJSON)
	return err
}

// DeleteAsset deletes an given asset from the world state.
func (s *SmartContract) DeleteAsset(ctx contractapi.TransactionContextInterface, id string) error {

	// exists, err := s.AssetExists(ctx, id)
	// if err != nil {
	// 	return err
	// }
	// if !exists {
	// 	return fmt.Errorf("the asset %s does not exist", id)
	// }
	assetJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return err
	}
	var asset Asset
	err = json.Unmarshal(assetJSON, &asset)
	if err != nil {
		return err
	}
	creator, err := ctx.GetStub().GetCreator()
	if err != nil {
		return err
	}

	certificate := base64.StdEncoding.EncodeToString(creator)
	participant := Participant{
		Name:        asset.Seller,
		Certificate: certificate,
	}
	eventData := CCEvent{
		Asset:       asset,
		Participant: participant,
	}
	eventJSON, err := json.Marshal(eventData)
	if err != nil {
		return err
	}
	err = ctx.GetStub().DelState(id)
	if err != nil {
		return fmt.Errorf("failed to store the asset: %v", err)
	}
	ctx.GetStub().SetEvent("DeleteAsset", eventJSON)

	return err
}

// AssetExists returns true when asset with given ID exists in world state
func (s *SmartContract) AssetExists(ctx contractapi.TransactionContextInterface, id string) (bool, error) {
	assetJSON, err := ctx.GetStub().GetState(id)
	if err != nil {
		return false, fmt.Errorf("failed to read from world state: %v", err)
	}

	return assetJSON != nil, nil
}

// TransferAsset updates the owner field of asset with given id in world state, and returns the old owner.
func (s *SmartContract) TransferAsset(ctx contractapi.TransactionContextInterface, id string, buyer string) (string, error) {
	asset, err := s.ReadAsset(ctx, id)
	if err != nil {
		return "", err
	}
	timeStamp, err := s.GetTxTimestampChannel(ctx)
	if err != nil {
		return "", err
	}
	seller := asset.Seller
	asset.Buyer = buyer
	asset.TransactionDate = timeStamp
	creator, err := ctx.GetStub().GetCreator()
	if err != nil {
		return "", err
	}
	certificate := base64.StdEncoding.EncodeToString(creator)
	participant := Participant{
		Name:        asset.Buyer,
		Certificate: certificate,
	}
	eventData := CCEvent{
		Asset:       *asset,
		Participant: participant,
	}
	eventJSON, err := json.Marshal(eventData)
	if err != nil {
		return "", err
	}

	assetJSON, err := json.Marshal(asset)
	if err != nil {
		return "", err
	}
	err = ctx.GetStub().PutState(id, assetJSON)
	if err != nil {
		return seller, fmt.Errorf("failed to store the asset: %v", err)
	}
	ctx.GetStub().SetEvent("TransferAsset", eventJSON)

	return seller, err
}

// GetAllAssets returns all assets found in world state
func (s *SmartContract) GetAllAssets(ctx contractapi.TransactionContextInterface) ([]*Asset, error) {
	// range query with empty string for startKey and endKey does an
	// open-ended query of all assets in the chaincode namespace.
	resultsIterator, err := ctx.GetStub().GetStateByRange("", "")
	if err != nil {
		return nil, err
	}
	defer resultsIterator.Close()

	var assets []*Asset
	var assetJSON []byte
	for resultsIterator.HasNext() {
		queryResponse, err := resultsIterator.Next()
		if err != nil {
			return nil, err
		}
		assetJSON = append(assetJSON, queryResponse.Value...)
		var asset Asset
		err = json.Unmarshal(queryResponse.Value, &asset)
		if err != nil {
			return nil, err
		}
		assets = append(assets, &asset)
	}
	ctx.GetStub().SetEvent("GetAllAssets", assetJSON)
	return assets, nil
}

// GetTxTimestampChannel Function gets the Transaction time when the chain code was executed it remains same on all the peers where chaincode executes
func (t *SmartContract) GetTxTimestampChannel(ctx contractapi.TransactionContextInterface) (string, error) {
	txTimeAsPtr, err := ctx.GetStub().GetTxTimestamp()
	if err != nil {
		fmt.Printf("Returning error in TimeStamp \n")
		return "Error", err
	}
	fmt.Printf("\t returned value from APIstub: %v\n", txTimeAsPtr)
	timeStr := time.Unix(txTimeAsPtr.Seconds, int64(txTimeAsPtr.Nanos)).String()

	return timeStr, nil
}
