package core_test

import (
	"io/ioutil"
	"log"
	"net"
	"testing"

	"core"
	"pos"
)

func testServer(app *core.AppContext) func() {

	// Log statements muck up the tests output
	if !testing.Verbose() {
		log.SetOutput(ioutil.Discard)
	}

	// If short tests, use in memory database
	if testing.Short() {
		memoryDB := core.MemoryDB{}
		memoryDB.Init()

		*app = core.AppContext{
			DB:            &memoryDB,
			Kounta:        &pos.MockKounta{},
			SiteWhitelist: []int64{29716, 87654},
		}

		return func() {}
	} else {
		postgresDB, err := core.CleanPG("postgres://postgres@localhost:5432/rize_core_test?sslmode=disable")
		if err != nil {
			switch err.(type) {
			case *net.OpError:
				log.Fatal("Cannot connect to database, make sure PostgreSQL is running")
			default:
				log.Fatal(err.Error())
			}
		}

		*app = core.AppContext{
			DB:            postgresDB,
			Kounta:        &pos.MockKounta{},
			SiteWhitelist: []int64{29716, 87654},
		}

		return func() {
			postgresDB.Close()
		}
	}
}
