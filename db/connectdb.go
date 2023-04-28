package db

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"strings"
)

var db *sql.DB

type TxData struct {
	ID           string
	ProductName  string
	ProductPrice int
	Seller       string
	Buyer        string
	State        string
}

func logQuery(query string, args ...interface{}) string {
	for _, arg := range args {
		var value string
		switch v := arg.(type) {
		case int, int64:
			value = fmt.Sprintf("%d", v)
		case float64:
			value = fmt.Sprintf("%f", v)
		case string:
			value = fmt.Sprintf("'%s'", v)
		default:
			value = fmt.Sprintf("'%v'", v)
		}
		query = strings.Replace(query, "?", value, 1)
	}
	return query
}

func ConnectionDB() *sql.DB {
	dsn := "root:pwd@tcp(localhost:3306)/trace?charset=utf8&parseTime=True&loc=Local"

	new_db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Failed to open database connection:", err)
	}
	//defer db.Close()

	// Test the connection to the database.
	err = new_db.Ping()
	if err != nil {
		log.Fatal("Failed to ping database:", err)
	}
	fmt.Println("Successfully connected to the database")
	db = new_db
	return new_db
}
func InsertTemporalData(id string, transactionDate string) {
	query := "INSERT INTO temporal_data (id, tx_date) VALUES (?, ?)"
	_, err := db.Exec(query, id, transactionDate)
	if err != nil {
		log.Fatal("Failed to insert data into temporal_data:", err)
	}
	log.Println(logQuery(query, id, transactionDate))

}

func InsertSpatialData(id string, location string) {
	query := "INSERT INTO spatial_data (id, location) VALUES (?, ?)"
	_, err := db.Exec(query, id, location)
	if err != nil {
		log.Fatal("Failed to insert data into spatial_data:", err)
	}
	log.Println(logQuery(query, id, location))
}

func InsertParticipantData(id string, userID string, certificate string) {
	query := "INSERT INTO participant_data (id,user_name, certificate) VALUES (?, ?, ?)"
	_, err := db.Exec(query, id, userID, certificate)
	if err != nil {
		log.Fatal("Failed to insert data into participant_data:", err)
	}
	log.Println(logQuery(query, id, userID, certificate))
}

//func InsertTransactionData(id string, productName string, productPrice float64, seller, buyer, state string, blockNumber int, blockHash string) {
func InsertTransactionData(id string, productName string, productPrice int, seller, buyer, state string) {
	query := "INSERT INTO tx_data (id, product_name, product_price, seller, buyer, state) VALUES (?, ?, ?, ?, ?, ?)"
	_, err := db.Exec(query, id, productName, productPrice, seller, buyer, state)
	if err != nil {
		log.Fatal("Failed to insert data into transaction_data:", err)
	}
	log.Println(logQuery(query, id, productName, productPrice, seller, buyer, state))

}
func UpdateTemporalData(id string, transactionDate string) {
	query := "UPDATE temporal_data SET tx_date = ? WHERE id = ?"
	_, err := db.Exec(query, transactionDate, id)
	if err != nil {
		log.Fatal("Failed to update data in temporal_data:", err)
	}
	fmt.Println(logQuery(query, transactionDate, id))
}

func UpdateSpatialData(id string, location string) {
	query := "UPDATE spatial_data SET location = ? WHERE id = ?"
	_, err := db.Exec(query, location, id)
	if err != nil {
		log.Fatal("Failed to update data in spatial_data:", err)
	}
	fmt.Println(logQuery(query, location, id))
}

func UpdateParticipantData(id, userID, certificate string) {
	query := "UPDATE participant_data SET user_name = ?,certificate = ?  WHERE id = ?"
	_, err := db.Exec(query, userID, certificate, id)
	if err != nil {
		log.Fatal("Failed to update data in participant_data:", err)
	}
	fmt.Println(logQuery(query, userID, certificate, id))
}

func UpdateTransactionData1(id string, productName string, productPrice int, seller, buyer, state string) {
	query := "UPDATE tx_data SET product_name = ?, product_price = ?, seller = ?, buyer = ?, state = ?, block_number = ?, block_hash = ? WHERE id = ?"
	_, err := db.Exec(query, productName, productPrice, seller, buyer, state, id)
	if err != nil {
		log.Fatal("Failed to update data in transaction_data:", err)
	}
	fmt.Println(logQuery(query, productName, productPrice, seller, buyer, state, id))
}
func UpdateTransactionData2(id string, buyer, state string) {
	query := "UPDATE tx_data SET  buyer = ?, state = ?  WHERE id = ?"
	_, err := db.Exec(query, buyer, state, id)
	if err != nil {
		log.Fatal("Failed to update data in transaction_data:", err)
	}
	fmt.Println(logQuery(query, buyer, state, id))

}
func UpdateTransactionData3(id string, buyer, state string) {
	query := "UPDATE tx_data SET  buyer = ?, state = ?  WHERE id = ?"
	_, err := db.Exec(query, buyer, state, id)
	if err != nil {
		log.Fatal("Failed to update data in transaction_data:", err)
	}
	log.Println(logQuery(query, buyer, state, id))
}
func QueryTxDataByID(id string) ([]TxData, error) {
	rows, err := db.Query("SELECT id, product_name, product_price, seller, buyer, state FROM tx_data")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var txDataList []TxData
	for rows.Next() {
		var txData TxData
		err = rows.Scan(&txData.ID, &txData.ProductName, &txData.ProductPrice, &txData.Seller, &txData.Buyer, &txData.State)
		if err != nil {
			return nil, err
		}
		txDataList = append(txDataList, txData)
	}

	return txDataList, nil
}
