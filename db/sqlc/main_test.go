package db

import (
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/zaid13/simplebank/db/util"
	"log"
	"os"
	"testing"
)

var testQueries *Queries
var testDB *sql.DB

func TestMain(m *testing.M) {

	config, err := util.LoadConfig("../..")

	if err != nil {
		log.Fatal("cannot start server:", err)
	}

	testDB, err = sql.Open(config.DBDriver, config.DBSource)

	if err != nil {
		log.Fatal("cannot connect to DB", err)
	}

	testQueries = New(testDB)

	os.Exit(m.Run())

}
