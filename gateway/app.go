package gateway

import (
	"Blockchain-Event-Trace-System/db"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hyperledger/fabric-gateway/pkg/client"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

const (
	channelName   = "mychannel"
	chaincodeName = "used_market3"
)

var now = time.Now()
var assetID = fmt.Sprintf("asset%d", now.Unix()*1e3+int64(now.Nanosecond())/1e6)
var contract = &client.Contract{}
var network = &client.Network{}
var startTime = time.Now()
var commitTime = time.Now()
var saveTime = time.Now()

type TransactionRequest struct {
	ID              string `json:"ID"`
	ProductName     string `json:"ProductName"`
	ProductPrice    int    `json:"ProductPrice"`
	TransactionDate string `json:"TransactionDate"`
	Seller          string `json:"Seller"`
	Buyer           string `json:"Buyer"`
	Location        string `json:"Location"`
	State           string `json:"State"`
}

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

func Connect() {
	clientConnection := newGrpcConnection()
	//defer clientConnection.Close()

	id := newIdentity()
	sign := newSign()

	gateway, err := client.Connect(
		id,
		client.WithSign(sign),
		client.WithClientConnection(clientConnection),
		client.WithEvaluateTimeout(5*time.Second),
		client.WithEndorseTimeout(15*time.Second),
		client.WithSubmitTimeout(5*time.Second),
		client.WithCommitStatusTimeout(1*time.Minute),
	)
	if err != nil {
		panic(err)
	}
	//defer gateway.Close()

	network = gateway.GetNetwork(channelName)
	contract = network.GetContract(chaincodeName)
	fmt.Printf("*** first:%s\n", contract)
	// Context used for event listening
	ctx, _ := context.WithCancel(context.Background())
	//defer cancel()

	// Listen for events emitted by subsequent transactions
	startChaincodeEventListening(ctx, network)

	//replayChaincodeEvents(ctx, network, 1)
}

func queryPerformanceTest() {
	startID := 1
	endID := 100
	db.QueryTxDataByID("asset1")
	file, err := os.OpenFile("sql_query_performance.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("failed to open file: %v", err)
	}
	defer file.Close()

	for id := startID; id <= endID; id++ {
		assetID := fmt.Sprintf("asset%d", id)
		dbStart := time.Now()

		txData, err := db.QueryTxDataByID(assetID)
		if err != nil {
			log.Printf("Error fetching data for ID %d: %v\n", assetID, err)
			continue
		}

		dbEnd := time.Now()
		elapsedTime := dbEnd.Sub(dbStart)
		elapsedTimeMs := float64(elapsedTime) / float64(time.Millisecond)
		fmt.Println(txData)
		fmt.Printf("ID: %d, Time: %.3f ms\n", id, elapsedTimeMs)
		_, err = fmt.Fprintf(file, "%.3f\n", elapsedTimeMs)
		if err != nil {
			log.Fatalf("failed to write to file: %v", err)
		}
	}
}

func CreateAsset(c *gin.Context) {
	startTime = time.Now()
	var transactionRequest TransactionRequest
	if err := c.BindJSON(&transactionRequest); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"err": err})
	}
	createAsset(contract, transactionRequest)
	commitTime = time.Now()
}
func UpdateAsset(c *gin.Context) {

	var transactionRequest TransactionRequest
	if err := c.BindJSON(&transactionRequest); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"err": err})
	}
	updateAsset(contract, transactionRequest)

}
func TransferAsset(c *gin.Context) {
	startTime = time.Now()
	var transactionRequest TransactionRequest
	if err := c.BindJSON(&transactionRequest); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"err": err})
	}
	transferAsset(contract, transactionRequest)
	commitTime = time.Now()
}
func DeleteAsset(c *gin.Context) {
	var transactionRequest TransactionRequest
	if err := c.BindJSON(&transactionRequest); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"err": err})
	}
	deleteAsset(contract, transactionRequest)
}
func GetAsset(c *gin.Context) {
	readStart := time.Now()
	var transactionRequest TransactionRequest
	if err := c.BindJSON(&transactionRequest); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"err": err})
	}
	readAsset(contract, transactionRequest)
	readEnd := time.Now()
	readTimeMs := readEnd.Sub(readStart).Seconds() * 1000
	fmt.Printf("%.3f\n", readTimeMs)
	file, err := os.OpenFile("ledger_query_performance.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0777)
	if err != nil {
		panic(fmt.Errorf("failed to open file: %w", err))
	}
	defer file.Close()
	_, err = fmt.Fprintf(file, "%.3f\n", readTimeMs)
	if err != nil {
		panic(fmt.Errorf("failed to write to file: %w", err))
	}
}
func GetAllAssets(c *gin.Context) {
	var transactionRequest TransactionRequest
	if err := c.BindJSON(&transactionRequest); err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"err": err})
	}
	getAllAssets(contract, transactionRequest)
}

func startChaincodeEventListening(ctx context.Context, network *client.Network) {

	blockEvents, blockErr := network.BlockEvents(ctx, client.WithStartBlock(1))
	if blockErr != nil {
		panic(fmt.Errorf("failed to start chaincode event listening: %w", blockErr))
	}
	fmt.Println("\n*** Start Block event listening")

	ccEvents, ccErr := network.ChaincodeEvents(ctx, chaincodeName)
	if ccErr != nil {
		panic(fmt.Errorf("failed to start block event listening: %w", ccErr))
	}
	fmt.Println("\n*** Start chaincode event listening")

	go func() {
		for event := range blockEvents {
			hashBytes := event.GetHeader().GetDataHash()
			hashString := fmt.Sprintf("%x", hashBytes)
			blockNumber := event.GetHeader().GetNumber()
			fmt.Printf("\n<-- Block event received: \n   Received block number : %d \n   Received block hash - %s\n", blockNumber, hashString)
		}
	}()
	go func() {
		for event := range ccEvents {
			asset := formatJSON(event.Payload)
			var ccEvent CCEvent
			err := json.Unmarshal(event.Payload, &ccEvent)
			if err != nil {
				log.Println(err.Error())
			}
			switch event.EventName {
			case "CreateAsset":
				db.InsertTemporalData(ccEvent.Asset.ID, ccEvent.Asset.TransactionDate)
				db.InsertParticipantData(ccEvent.Asset.ID, ccEvent.Participant.Name, ccEvent.Participant.Certificate)
				db.InsertTransactionData(ccEvent.Asset.ID, ccEvent.Asset.ProductName, ccEvent.Asset.ProductPrice, ccEvent.Asset.Seller, ccEvent.Asset.Buyer, ccEvent.Asset.State)
				db.InsertSpatialData(ccEvent.Asset.ID, ccEvent.Asset.Location)
			case "UpdateAsset":
				db.UpdateTransactionData1(ccEvent.Asset.ID, ccEvent.Asset.ProductName, ccEvent.Asset.ProductPrice, ccEvent.Asset.Seller, ccEvent.Asset.Buyer, ccEvent.Asset.State)
			case "DeleteAsset":
			case "TransferAsset":
				db.UpdateTemporalData(ccEvent.Asset.ID, ccEvent.Asset.TransactionDate)
				db.UpdateSpatialData(ccEvent.Asset.ID, ccEvent.Asset.Location)
				db.UpdateParticipantData(ccEvent.Asset.ID, ccEvent.Participant.Name, ccEvent.Participant.Certificate)
				db.UpdateTransactionData2(ccEvent.Asset.ID, ccEvent.Asset.Buyer, ccEvent.Asset.State)
			default:
				fmt.Printf(event.EventName)
			}
			saveTime = time.Now()
			commitTimeMs := commitTime.Sub(startTime).Seconds() * 1000
			saveTimeMS := saveTime.Sub(startTime).Seconds() * 1000
			fmt.Printf("%.3f %.3f\n", commitTimeMs, saveTimeMS)
			file, err := os.OpenFile("txquery.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0777)
			if err != nil {
				panic(fmt.Errorf("failed to open file: %w", err))
			}
			defer file.Close()
			_, err = fmt.Fprintf(file, "%.3f %.3f\n", commitTimeMs, saveTimeMS)
			if err != nil {
				panic(fmt.Errorf("failed to write to file: %w", err))
			}
			fmt.Printf("\n<-- Chaincode event received: %s - %s\n", event.EventName, asset)
		}
	}()
}

func formatJSON(data []byte) string {
	var result bytes.Buffer
	if err := json.Indent(&result, data, "", "  "); err != nil {
		panic(fmt.Errorf("failed to parse JSON: %w", err))
	}
	return result.String()
}
func getAllAssets(contract *client.Contract, request TransactionRequest) {
	fmt.Printf("\n--> Query Assets: GetAllAssets, getAllAssets\n")

	_, err := contract.EvaluateTransaction("GetAllAssets")
	if err != nil {
		panic(fmt.Errorf("failed to evaluate transaction: %w", err))
	}
	//result := formatJSON(evaluateResult)

	fmt.Println("\n*** Query get all assets successfully")
	queryPerformanceTest()
}

func readAsset(contract *client.Contract, request TransactionRequest) {
	fmt.Printf("\n--> Query Assets: ReadAsset, ReadAsset\n")
	_, err := contract.EvaluateTransaction("ReadAsset", request.ID)
	if err != nil {
		panic(fmt.Errorf("failed to evaluate transaction: %w", err))
	}
	//result := formatJSON(evaluateResult)

	fmt.Println("\n*** Query get all assets successfully")
}

func createAsset(contract *client.Contract, request TransactionRequest) uint64 {
	fmt.Printf("\n--> Submit transaction: CreateAsset, %s starts selling %s with price %d\n", request.Seller, request.ID, request.ProductPrice)
	price := strconv.Itoa(request.ProductPrice)
	_, commit, err := contract.SubmitAsync("CreateAsset",
		client.WithArguments(request.ID, request.ProductName, price, request.TransactionDate, request.Seller, request.Buyer, request.Location, "on_sale"))
	if err != nil {
		panic(fmt.Errorf("failed to submit transaction: %w", err))
	}

	status, err := commit.Status()
	if err != nil {
		panic(fmt.Errorf("failed to get transaction commit status: %w", err))
	}

	if !status.Successful {
		panic(fmt.Errorf("failed to commit transaction with status code %v", status.Code))
	}

	fmt.Println("\n*** CreateAsset committed successfully")

	return status.BlockNumber
}

func updateAsset(contract *client.Contract, request TransactionRequest) {
	fmt.Printf("\n--> Submit transaction: UpdateAsset, %s updates %s with price %d$\n", request.Seller, request.ID, request.ProductPrice)

	_, err := contract.SubmitTransaction("UpdateAsset",
		request.ID, request.ProductName, string(request.ProductPrice), request.TransactionDate, request.Seller, request.Buyer, request.Location, request.State)
	if err != nil {
		panic(fmt.Errorf("failed to submit transaction: %w", err))
	}

	fmt.Println("\n*** UpdateAsset committed successfully")
}

func transferAsset(contract *client.Contract, request TransactionRequest) {
	fmt.Printf("\n--> Submit transaction: TransferAsset, %s sells %s to %s with price %d$\n", request.Seller, request.ID, request.Buyer, request.ProductPrice)

	_, err := contract.SubmitTransaction("TransferAsset", request.ID, request.Buyer)
	if err != nil {
		panic(fmt.Errorf("failed to submit transaction: %w", err))
	}

	fmt.Println("\n*** TransferAsset committed successfully")
}

func deleteAsset(contract *client.Contract, request TransactionRequest) {
	fmt.Printf("\n--> Submit transaction: DeleteAsset, %s\n", assetID)

	_, err := contract.SubmitTransaction("DeleteAsset", request.ID)
	if err != nil {
		panic(fmt.Errorf("failed to submit transaction: %w", err))
	}

	fmt.Println("\n*** DeleteAsset committed successfully")
}

func replayChaincodeEvents(ctx context.Context, network *client.Network, startBlock uint64) {
	fmt.Println("\n*** Start chaincode event replay")

	events, err := network.ChaincodeEvents(ctx, chaincodeName, client.WithStartBlock(startBlock))
	if err != nil {
		panic(fmt.Errorf("failed to start chaincode event listening: %w", err))
	}

	for {
		select {
		case <-time.After(10 * time.Second):
			panic(errors.New("timeout waiting for event replay"))

		case event := <-events:
			asset := formatJSON(event.Payload)
			fmt.Printf("\n<-- Chaincode event replayed: %s - %s\n", event.EventName, asset)

			if event.EventName == "DeleteAsset" {
				// Reached the last submitted transaction so return to stop listening for events
				return
			}
		}
	}
}
