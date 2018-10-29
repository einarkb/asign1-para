package paragliding

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mongodb/mongo-go-driver/bson/objectid"
)

func Test_HandlerLatestTick(t *testing.T) {
	// create a database connection
	db := &Database{URI: "mongodb://test:test12@ds141783.mlab.com:41783/a2-trackdb", Name: "a2-trackdb"}
	db.Connect()
	mgrTicker := MgrTicker{DB: db, PageCap: 5}
	// creating request
	req, err := http.NewRequest("GET", "/paragliding/api/ticker/lastest", nil)
	if err != nil {
		t.Error(err)
	}
	res := httptest.NewRecorder()
	handler := http.HandlerFunc(mgrTicker.HandlerLatestTick)

	handler.ServeHTTP(res, req)

	// server mock request
	if res.Code != http.StatusOK {
		t.Error("Bad status response: expected %i got %i", http.StatusOK, res.Code)
	}
}

func Test_HandlerTicker(t *testing.T) {
	// create a database connection
	db := &Database{URI: "mongodb://test:test12@ds141783.mlab.com:41783/a2-trackdb", Name: "a2-trackdb"}
	db.Connect()
	mgrTicker := MgrTicker{DB: db, PageCap: 5}
	// creating request
	req, err := http.NewRequest("GET", "/paragliding/api/ticker/", nil)
	if err != nil {
		t.Error(err)
	}
	res := httptest.NewRecorder()
	handler := http.HandlerFunc(mgrTicker.HandlerTicker)

	// server mock request
	handler.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Error("Bad status response: expected %i got %i", http.StatusOK, res.Code)
	}
}

func Test_HandlerTickerByTimestamp(t *testing.T) {
	// create a database connection
	db := &Database{URI: "mongodb://test:test12@ds141783.mlab.com:41783/a2-trackdb", Name: "a2-trackdb"}
	db.Connect()
	mgrTicker := MgrTicker{DB: db, PageCap: 5}
	// creating request
	req, err := http.NewRequest("GET", "/paragliding/api/ticker/1", nil)
	if err != nil {
		t.Error(err)
	}
	res := httptest.NewRecorder()
	handler := http.HandlerFunc(mgrTicker.HandlerTickerByTimestamp)

	// server mock request
	handler.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Error("Bad status response: expected %i got %i", http.StatusOK, res.Code)
	}
}

func Test_GetTickerByTimeStamp(t *testing.T) {
	db := &Database{URI: "mongodb://test:test12@ds143893.mlab.com:43893/a2-testddb", Name: "a2-testddb"}
	db.Connect()
	mgrTicker := MgrTicker{DB: db, PageCap: 5}

	// add tracks to the db
	tracks := [10]TrackInfo{}
	for i := 0; i < 10; i++ {
		tracks[i] = TrackInfo{ID: objectid.New(), HDate: "somedate", Pilot: "ole",
			Glider: "sometype", GliderID: "someID", TrackLength: "10",
			TrackURL: "www.trackurl.com", Timestamp: int64(i)}
		db.Insert("tracks", tracks[i])
	}

	// ask for a reponse and check if it is as expected
	res, err := mgrTicker.GetTickerByTimeStamp(2)
	if err != nil {
		t.Error("failed to get ticker by timestamp")
	}
	if res.TLatest != 9 {
		t.Error("latest timestamp is incorrect")
	}
	if len(res.TrackIDs) != 5 {
		t.Error("wrong number of track ids received")
	}
	if res.TStart != 3 {
		t.Error("TStarts is incorrect")
	}
	if res.TStop != 7 {
		t.Error("TStop is inccorrect")
	}

}
