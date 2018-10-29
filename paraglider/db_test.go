package paragliding

import (
	"testing"
)

func Test_Connect(t *testing.T) {
	db := &Database{URI: "mongodb://test:test12@ds141783.mlab.com:41783/a2-trackdb", Name: "a2-trackdb"}
	db.Connect()
	if db.conn == nil {
		t.Error("Failed to connect to database")
	}
}

func Test_Insert(t *testing.T) {

}
