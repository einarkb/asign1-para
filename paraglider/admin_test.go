package paragliding

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_HandlerTrackCount(t *testing.T) {
	// create a database connection
	db := &Database{URI: "mongodb://test:test12@ds141783.mlab.com:41783/a2-trackdb", Name: "a2-trackdb"}
	db.Connect()
	adminMgr := AdminMgr{DB: db}
	// creating request
	req, err := http.NewRequest("GET", "/admin/api/tracks_count", nil)
	if err != nil {
		t.Error(err)
	}
	res := httptest.NewRecorder()
	handler := http.HandlerFunc(adminMgr.HandlerTrackCount)

	handler.ServeHTTP(res, req)

	// server mock request
	if res.Code != http.StatusOK {
		t.Error("Bad status response: expected %i got %i", http.StatusOK, res.Code)
	}
}

func Test_HandlerDeleteAllTracks(t *testing.T) {
	// create a database connection
	db := &Database{URI: "mongodb://test:test12@ds143893.mlab.com:43893/a2-testddb", Name: "a2-testddb"}
	db.Connect()
	adminMgr := AdminMgr{DB: db}
	// creating request
	req, err := http.NewRequest("DELETE", "/admin/api/tracks", nil)
	if err != nil {
		t.Error(err)
	}
	res := httptest.NewRecorder()
	handler := http.HandlerFunc(adminMgr.HandlerDeleteAllTracks)

	handler.ServeHTTP(res, req)

	// server mock request
	if res.Code != http.StatusOK {
		t.Error("Bad status response: expected %i got %i", http.StatusOK, res.Code)
	}
}
