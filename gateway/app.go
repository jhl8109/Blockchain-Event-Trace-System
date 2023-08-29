package gateway

import (
	"Blockchain-Event-Trace-System/db"
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hyperledger/fabric-gateway/pkg/client"
	"log"
	"os"
	"strconv"
	"time"
)

const (
	channelName   = "used-car"
	chaincodeName = "transaction"
)

var contract = &client.Contract{}
var network = &client.Network{}

type Transaction struct {
	ID                 int64              `json:"id"`
	UploadDate         string             `json:"uploadDate"`
	Assignor           Participant        `json:"assignor"`
	Assignee           Participant        `json:"assignee"`
	TransactionDetails TransactionDetails `json:"transactionDetails"`
}

type Participant struct {
	Name                       string `json:"name"`
	ResidentRegistrationNumber string `json:"residentRegistrationNumber"`
	PhoneNumber                string `json:"phoneNumber"`
	Address                    string `json:"address"`
}

type TransactionDetails struct {
	TransactionState             string `json:"transactionState"`
	VehicleRegistrationNumber    string `json:"vehicleRegistrationNumber"`
	NewVehicleRegistrationNumber string `json:"newVehicleRegistrationNumber"`
	VehicleModelName             string `json:"vehicleModelName"`
	VehicleIdentificationNumber  string `json:"vehicleIdentificationNumber"`
	TransactionDate              string `json:"transactionDate"`
	TransactionAmount            string `json:"transactionAmount"`
	BalancePaymentDate           string `json:"balancePaymentDate"`
	VehicleDeliveryDate          string `json:"vehicleDeliveryDate"`
	VehicleDeliveryAddress       string `json:"vehicleDeliveryAddress"`
	Mileage                      string `json:"mileage"`
}

type CCEvent2 struct {
	Transaction Transaction `json:"transaction"`
}

func Connect() {
	clientConnection := newGrpcConnection()

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

	network = gateway.GetNetwork(channelName)
	contract = network.GetContract(chaincodeName)
	fmt.Printf("*** first:%s\n", contract)

	//ctx, _ := context.WithCancel(context.Background())

	//startChaincodeEventListening(ctx, network)
}

func startChaincodeEventListening(ctx context.Context, network *client.Network) {

	//blockEvents, blockErr := network.BlockEvents(ctx, client.WithStartBlock(1))
	//if blockErr != nil {
	//	panic(fmt.Errorf("failed to start chaincode event listening: %w", blockErr))
	//}
	//fmt.Println("\n*** Start Block event listening")

	ccEvents, ccErr := network.ChaincodeEvents(ctx, chaincodeName)
	if ccErr != nil {
		panic(fmt.Errorf("failed to start block event listening: %w", ccErr))
	}
	fmt.Println("\n*** Start chaincode event listening")

	//go func() {
	//	for event := range blockEvents {
	//		hashBytes := event.GetHeader().GetDataHash()
	//		hashString := fmt.Sprintf("%x", hashBytes)
	//		blockNumber := event.GetHeader().GetNumber()
	//		fmt.Printf("\n<-- Block event received: \n   Received block number : %d \n   Received block hash - %s\n", blockNumber, hashString)
	//	}
	//}()
	go func() {
		outputFile := "process.txt"
		file, err := os.OpenFile(outputFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		}

		for event := range ccEvents {
			startTime := time.Now().UnixNano()
			startTimeString := fmt.Sprintf("%d", startTime)

			eventStr := formatJSON(event.Payload)
			var eventData Transaction
			err := json.Unmarshal(event.Payload, &eventData)
			if err != nil {
				log.Println(err.Error())
			}
			switch event.EventName {
			case "SellVehicle":
				db.InsertTemporalData(eventData.toTemporalData())
				db.InsertSpatialData(eventData.toSpatialData())
				db.InsertParticipantData(eventData.toParticipantData())
				db.InsertTransactionData(eventData.toTransactionData())
				break
			case "BuyVehicle":
				db.UpdateTemporalData(eventData.toTemporalData())
				db.UpdateSpatialData(eventData.toSpatialData())
				db.UpdateParticipantData(eventData.toParticipantData())
				db.UpdateTransactionData(eventData.toTransactionData())
				break
			case "CompromiseTransaction":
				db.UpdateTemporalData(eventData.toTemporalData())
				db.UpdateSpatialData(eventData.toSpatialData())
				db.UpdateParticipantData(eventData.toParticipantData())
				db.UpdateTransactionData(eventData.toTransactionData())
				break
			default:
				fmt.Printf(event.EventName)
			}
			fmt.Printf("\n<-- Chaincode event received: %s - %s\n", event.EventName, eventStr)

			endTime := time.Now().UnixNano()
			endTimeString := fmt.Sprintf("%d", endTime)
			if _, err := file.WriteString(startTimeString + " " + endTimeString + "\n"); err != nil {
				log.Println(err)
			}
		}

	}()
}
func timeParse(timeStr string) (time.Time, error) {
	const timeLayout = "2006-01-02 15:04:05"
	location, err := time.LoadLocation("Asia/Seoul")
	if err != nil {
		return time.Time{}, fmt.Errorf("Failed to load location: %s", err.Error())
	}
	if timeStr == "" {
		return time.Time{}, fmt.Errorf("Empty string provided, cannot parse as time.Time")
	}
	parsedTime, err := time.ParseInLocation(timeLayout, timeStr, location)
	if err != nil {
		return time.Time{}, fmt.Errorf("Failed to parse string '%s' to time.Time: %s", timeStr, err.Error())
	}
	return parsedTime, nil
}

func (tx Transaction) toTemporalData() db.TemporalData {
	uploadDate, _ := timeParse(tx.UploadDate)
	transactionDate, _ := timeParse(tx.TransactionDetails.TransactionDate)
	balancePaymentDate, _ := timeParse(tx.TransactionDetails.BalancePaymentDate)
	vehicleDeliveryDate, _ := timeParse(tx.TransactionDetails.VehicleDeliveryDate)

	temporalData := db.TemporalData{
		ID:                  tx.ID,
		UploadDate:          sql.NullTime{Time: uploadDate, Valid: uploadDate != time.Time{}},
		TransactionDate:     sql.NullTime{Time: transactionDate, Valid: transactionDate != time.Time{}},
		BalancePaymentDate:  sql.NullTime{Time: balancePaymentDate, Valid: balancePaymentDate != time.Time{}},
		VehicleDeliveryDate: sql.NullTime{Time: vehicleDeliveryDate, Valid: vehicleDeliveryDate != time.Time{}},
	}
	return temporalData
}

func (tx Transaction) toSpatialData() db.SpatialData {
	spatialData := db.SpatialData{
		ID:                     tx.ID,
		VehicleDeliveryAddress: sql.NullString{String: tx.TransactionDetails.VehicleDeliveryAddress, Valid: tx.TransactionDetails.VehicleDeliveryAddress != ""},
		AssignorAddress:        sql.NullString{String: tx.Assignor.Address, Valid: tx.Assignor.Address != ""},
		AssigneeAddress:        sql.NullString{String: tx.Assignee.Address, Valid: tx.Assignee.Address != ""},
	}
	return spatialData
}

func (tx Transaction) toParticipantData() db.ParticipantData {
	participantData := db.ParticipantData{
		ID:                                 tx.ID,
		AssignorName:                       sql.NullString{String: tx.Assignor.Name, Valid: tx.Assignor.Name != ""},
		AssignorResidentRegistrationNumber: sql.NullString{String: tx.Assignor.ResidentRegistrationNumber, Valid: tx.Assignor.ResidentRegistrationNumber != ""},
		AssignorPhoneNumber:                sql.NullString{String: tx.Assignor.PhoneNumber, Valid: tx.Assignor.PhoneNumber != ""},
		AssignorAddress:                    sql.NullString{String: tx.Assignor.Address, Valid: tx.Assignor.Address != ""},
		AssigneeName:                       sql.NullString{String: tx.Assignee.Name, Valid: tx.Assignee.Name != ""},
		AssigneeResidentRegistrationNumber: sql.NullString{String: tx.Assignee.ResidentRegistrationNumber, Valid: tx.Assignee.ResidentRegistrationNumber != ""},
		AssigneePhoneNumber:                sql.NullString{String: tx.Assignee.PhoneNumber, Valid: tx.Assignee.PhoneNumber != ""},
		AssigneeAddress:                    sql.NullString{String: tx.Assignee.Address, Valid: tx.Assignee.Address != ""},
	}
	return participantData
}
func (tx Transaction) toTransactionData() db.TransactionData {
	transactionDate, _ := timeParse(tx.TransactionDetails.TransactionDate)
	transactionAmount, _ := strconv.ParseFloat(tx.TransactionDetails.TransactionAmount, 64)
	balancePaymentDate, _ := timeParse(tx.TransactionDetails.BalancePaymentDate)
	vehicleDeliveryDate, _ := timeParse(tx.TransactionDetails.VehicleDeliveryDate)
	mileage, _ := strconv.ParseInt(tx.TransactionDetails.Mileage, 10, 64)

	transactionData := db.TransactionData{
		ID:                           tx.ID,
		TransactionState:             sql.NullString{String: tx.TransactionDetails.TransactionState, Valid: tx.TransactionDetails.TransactionState != ""},
		VehicleRegistrationNumber:    sql.NullString{String: tx.TransactionDetails.VehicleRegistrationNumber, Valid: tx.TransactionDetails.VehicleRegistrationNumber != ""},
		NewVehicleRegistrationNumber: sql.NullString{String: tx.TransactionDetails.NewVehicleRegistrationNumber, Valid: tx.TransactionDetails.NewVehicleRegistrationNumber != ""},
		VehicleModelName:             sql.NullString{String: tx.TransactionDetails.VehicleModelName, Valid: tx.TransactionDetails.VehicleModelName != ""},
		VehicleIdentificationNumber:  sql.NullString{String: tx.TransactionDetails.VehicleIdentificationNumber, Valid: tx.TransactionDetails.VehicleIdentificationNumber != ""},
		TransactionDate:              sql.NullTime{Time: transactionDate, Valid: transactionDate != time.Time{}},
		TransactionAmount:            sql.NullFloat64{Float64: transactionAmount, Valid: transactionAmount != 0},
		BalancePaymentDate:           sql.NullTime{Time: balancePaymentDate, Valid: balancePaymentDate != time.Time{}},
		VehicleDeliveryDate:          sql.NullTime{Time: vehicleDeliveryDate, Valid: vehicleDeliveryDate != time.Time{}},
		VehicleDeliveryAddress:       sql.NullString{String: tx.TransactionDetails.VehicleDeliveryAddress, Valid: tx.TransactionDetails.VehicleDeliveryAddress != ""},
		Mileage:                      sql.NullInt64{Int64: mileage, Valid: mileage != 0},
	}
	return transactionData
}

func formatJSON(data []byte) string {
	var result bytes.Buffer
	if err := json.Indent(&result, data, "", "  "); err != nil {
		panic(fmt.Errorf("failed to parse JSON: %w", err))
	}
	return result.String()
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

func GetAsset(c *gin.Context) {
	db.SelectData()
}
