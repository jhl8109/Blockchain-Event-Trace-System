package db

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"strings"
	"time"
)

var db *sql.DB

type TemporalData struct {
	ID                  int64
	UploadDate          sql.NullTime
	TransactionDate     sql.NullTime
	BalancePaymentDate  sql.NullTime
	VehicleDeliveryDate sql.NullTime
}

type SpatialData struct {
	ID                     int64
	VehicleDeliveryAddress sql.NullString
	AssignorAddress        sql.NullString
	AssigneeAddress        sql.NullString
}

type ParticipantData struct {
	ID                                 int64
	AssignorName                       sql.NullString
	AssignorResidentRegistrationNumber sql.NullString
	AssignorPhoneNumber                sql.NullString
	AssignorAddress                    sql.NullString
	AssigneeName                       sql.NullString
	AssigneeResidentRegistrationNumber sql.NullString
	AssigneePhoneNumber                sql.NullString
	AssigneeAddress                    sql.NullString
}
type TransactionData struct {
	ID                           int64
	TransactionState             sql.NullString
	VehicleRegistrationNumber    sql.NullString
	NewVehicleRegistrationNumber sql.NullString
	VehicleModelName             sql.NullString
	VehicleIdentificationNumber  sql.NullString
	TransactionDate              sql.NullTime
	TransactionAmount            sql.NullFloat64
	BalancePaymentDate           sql.NullTime
	VehicleDeliveryDate          sql.NullTime
	VehicleDeliveryAddress       sql.NullString
	Mileage                      sql.NullInt64
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
	dsn := "root:pwd@tcp(localhost:3306)/Offchain?charset=utf8&parseTime=True&loc=Local"

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
func InsertTemporalData(data TemporalData) {
	query := "INSERT INTO TemporalData (id, upload_date, transaction_date, balance_payment_date, vehicle_delivery_date) VALUES (?, ?, ?, ?, ?)"
	_, err := db.Exec(query, data.ID, data.UploadDate, data.TransactionDate, data.BalancePaymentDate, data.VehicleDeliveryDate)
	if err != nil {
		log.Fatal("Failed to insert data into temporal_data:", err)
	}
	log.Println(logQuery(query, data.ID, data.UploadDate, data.TransactionDate, data.BalancePaymentDate, data.VehicleDeliveryDate))
}

func InsertSpatialData(data SpatialData) {
	query := "INSERT INTO SpatialData (id, vehicle_delivery_address, assignor_address, assignee_address) VALUES (?, ?, ?, ?)"
	_, err := db.Exec(query, data.ID, data.VehicleDeliveryAddress, data.AssignorAddress, data.AssigneeAddress)
	if err != nil {
		log.Fatal("Failed to insert data into spatial_data:", err)
	}
	log.Println(logQuery(query, data.ID, data.VehicleDeliveryAddress, data.AssignorAddress, data.AssigneeAddress))
}

func InsertParticipantData(data ParticipantData) {
	query := "INSERT INTO ParticipantData (" +
		"id,assignor_name, assignor_resident_registration_number, assignor_phone_number, assignor_address, " +
		"assignee_name, assignee_resident_registration_number, assignee_phone_number, assignee_address) " +
		"VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)"
	_, err := db.Exec(query,
		data.ID, data.AssignorName, data.AssignorResidentRegistrationNumber, data.AssignorPhoneNumber, data.AssignorAddress,
		data.AssigneeName, data.AssigneeResidentRegistrationNumber, data.AssigneePhoneNumber, data.AssigneeAddress)
	if err != nil {
		log.Fatal("Failed to insert data into participant_data:", err)
	}
	log.Println(logQuery(query, data.ID, data.AssignorName, data.AssignorResidentRegistrationNumber, data.AssignorPhoneNumber, data.AssignorAddress,
		data.AssigneeName, data.AssigneeResidentRegistrationNumber, data.AssigneePhoneNumber, data.AssigneeAddress))
}

func InsertTransactionData(data TransactionData) {
	query := "INSERT INTO TransactionData (" +
		"id, transaction_state, vehicle_registration_number, new_vehicle_registration_number, " +
		"vehicle_model_name, vehicle_identification_number, transaction_date, transaction_amount, " +
		"balance_payment_date, vehicle_delivery_date, vehicle_delivery_address, mileage) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)"
	_, err := db.Exec(query,
		data.ID, data.TransactionState, data.VehicleRegistrationNumber, data.NewVehicleRegistrationNumber,
		data.VehicleModelName, data.VehicleIdentificationNumber, data.TransactionDate, data.TransactionAmount,
		data.BalancePaymentDate, data.VehicleDeliveryDate, data.VehicleDeliveryAddress, data.Mileage)
	if err != nil {
		log.Fatal("Failed to insert data into transaction_data:", err)
	}
	log.Println(logQuery(query, data.ID, data.TransactionState, data.VehicleRegistrationNumber, data.NewVehicleRegistrationNumber,
		data.VehicleModelName, data.VehicleIdentificationNumber, data.TransactionDate, data.TransactionAmount,
		data.BalancePaymentDate, data.VehicleDeliveryDate, data.VehicleDeliveryAddress, data.Mileage))

}
func UpdateTemporalData(data TemporalData) {
	query := "UPDATE TemporalData SET upload_date = ?, transaction_date = ?, balance_payment_date = ?, vehicle_delivery_date = ? WHERE id = ?"
	_, err := db.Exec(query, data.UploadDate, data.TransactionDate, data.BalancePaymentDate, data.VehicleDeliveryDate, data.ID)
	if err != nil {
		log.Fatal("Failed to update data in temporal_data:", err)
	}
	fmt.Println(logQuery(query, data.UploadDate, data.TransactionDate, data.BalancePaymentDate, data.VehicleDeliveryDate, data.ID))
}

func UpdateSpatialData(data SpatialData) {
	query := "UPDATE SpatialData SET vehicle_delivery_address = ?, assignor_address = ?, assignee_address = ? WHERE id = ?"
	_, err := db.Exec(query, data.VehicleDeliveryAddress, data.AssignorAddress, data.AssigneeAddress, data.ID)
	if err != nil {
		log.Fatal("Failed to update data in spatial_data:", err)
	}
	fmt.Println(logQuery(query, data.VehicleDeliveryAddress, data.AssignorAddress, data.AssigneeAddress, data.ID))
}

func UpdateParticipantData(data ParticipantData) {
	query := "UPDATE ParticipantData SET assignor_name = ?, assignor_resident_registration_number = ?, assignor_phone_number = ?, assignor_address = ?, " +
		"assignee_name = ?, assignee_resident_registration_number = ?, assignee_phone_number = ?, assignee_address = ?  WHERE id = ?"
	_, err := db.Exec(query, data.AssignorName, data.AssignorResidentRegistrationNumber, data.AssignorPhoneNumber, data.AssignorAddress,
		data.AssigneeName, data.AssigneeResidentRegistrationNumber, data.AssigneePhoneNumber, data.AssigneeAddress, data.ID)
	if err != nil {
		log.Fatal("Failed to update data in participant_data:", err)
	}
	fmt.Println(logQuery(query, data.AssignorName, data.AssignorResidentRegistrationNumber, data.AssignorPhoneNumber, data.AssignorAddress,
		data.AssigneeName, data.AssigneeResidentRegistrationNumber, data.AssigneePhoneNumber, data.AssigneeAddress, data.ID))
}

func UpdateTransactionData(data TransactionData) {
	query := "UPDATE TransactionData SET transaction_state = ?, vehicle_registration_number = ?, new_vehicle_registration_number = ?, vehicle_model_name = ?, " +
		"vehicle_identification_number = ?, transaction_date = ?, transaction_amount = ?, balance_payment_date = ?, " +
		"vehicle_delivery_date = ?, vehicle_delivery_address = ?, mileage = ? WHERE id = ?"
	_, err := db.Exec(query, data.TransactionState, data.VehicleRegistrationNumber, data.NewVehicleRegistrationNumber, data.VehicleModelName,
		data.VehicleIdentificationNumber, data.TransactionDate, data.TransactionAmount, data.BalancePaymentDate,
		data.VehicleDeliveryDate, data.VehicleDeliveryAddress, data.Mileage, data.ID)
	if err != nil {
		log.Fatal("Failed to update data in transaction_data:", err)
	}
	fmt.Println(logQuery(query, data.TransactionState, data.VehicleRegistrationNumber, data.NewVehicleRegistrationNumber, data.VehicleModelName,
		data.VehicleIdentificationNumber, data.TransactionDate, data.TransactionAmount, data.BalancePaymentDate,
		data.VehicleDeliveryDate, data.VehicleDeliveryAddress, data.Mileage, data.ID))
}

// Define a struct to hold the result set columns
type TemporalTransactionData struct {
	// Define fields based on your column names and types
	ID                           int
	UploadDate                   time.Time
	TransactionDate              time.Time
	BalancePaymentDate           time.Time
	VehicleDeliveryDate          time.Time
	TransactionState             string
	VehicleRegistrationNumber    string
	NewVehicleRegistrationNumber string
	VehicleModelName             string
	VehicleIdentificationNumber  string
	TransactionAmount            float64
	VehicleDeliveryAddress       string
	Mileage                      int
}

func SelectData() ([]TemporalTransactionData, error) {
	query := `SELECT * FROM Offchain.TemporalData t JOIN Offchain.TransactionData td ON t.id = td.id WHERE t.upload_date BETWEEN ? AND ? AND td.transaction_state = ? AND td.vehicle_model_name = ? AND td.transaction_amount BETWEEN ? AND ? AND td.mileage BETWEEN ? AND ?
	`
	startDate := time.Date(2022, 7, 13, 0, 0, 0, 0, time.UTC)
	endDate := time.Date(2023, 7, 13, 23, 59, 59, 0, time.UTC)
	transactionState := "Selling"
	vehicleModelName := "Tesla Model S"
	minTransactionAmount := 5000
	maxTransactionAmount := 30000
	minMileage := 0
	maxMileage := 200000

	rows, err := db.Query(query,
		startDate, endDate,
		transactionState, vehicleModelName,
		minTransactionAmount, maxTransactionAmount,
		minMileage, maxMileage,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to execute the query: %v", err)
	}
	defer rows.Close()

	var result []TemporalTransactionData

	for rows.Next() {
		var tempData TemporalTransactionData
		err = rows.Scan(
			&tempData.ID,
			&tempData.UploadDate,
			&tempData.TransactionDate,
			&tempData.BalancePaymentDate,
			&tempData.VehicleDeliveryDate,
			&tempData.TransactionState,
			&tempData.VehicleRegistrationNumber,
			&tempData.NewVehicleRegistrationNumber,
			&tempData.VehicleModelName,
			&tempData.VehicleIdentificationNumber,
			&tempData.TransactionAmount,
			&tempData.VehicleDeliveryAddress,
			&tempData.Mileage,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row values: %v", err)
		}
		result = append(result, tempData)
	}

	err = rows.Err()
	if err != nil {
		return nil, fmt.Errorf("error during result iteration: %v", err)
	}

	log.Println(logQuery(query, startDate, endDate, transactionState, vehicleModelName, minTransactionAmount, maxTransactionAmount, minMileage, maxMileage))
	return result, nil
}
