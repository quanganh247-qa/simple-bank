package db

import (
	"database/sql"
	"log"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

const (
	dbDriver = "postgres"
	dbSource = "postgres://root:fHWFyt98gPR51h3NxjcroWoIscjt7QOb@dpg-cp649mmn7f5s73a6r8ag-a.oregon-postgres.render.com/simple_bank_7qc2"
)

var testQueries *Queries
var db *sql.DB

func TestMain(m *testing.M) {
	// Establish a database connection
	var err error
	db, err = sql.Open(dbDriver, dbSource)
	if err != nil {
		log.Fatal("cannot connect to db:", err)
	}
	defer db.Close()

	// Initialize the testQueries object
	testQueries = New(db)
	if testQueries == nil {
		log.Fatalf("testQueries is nil")
	}

	// // Retrieve an account for testing
	// acc, err := testQueries.GetAccount(context.Background(), 1)
	// if err != nil {
	// 	log.Fatalf("error retrieving account: %v", err)
	// }
	// log.Printf("Account Details: Owner=%s, Balance=%d, Currency=%s", acc.Owner, acc.Balance, acc.Currency)

	// Run the tests
	code := m.Run()

	os.Exit(code)
}
