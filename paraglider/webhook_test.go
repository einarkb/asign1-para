package paragliding

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_HandlerGetWebhookHookByID(t *testing.T) {
	// create a database connection
	db := &Database{URI: "mongodb://test:test12@ds141783.mlab.com:41783/a2-trackdb", Name: "a2-trackdb"}
	db.Connect()
	whMgr := WebHookMgr{DB: db}
	// creating request
	req, err := http.NewRequest("GET", "/api/webhook/new_track/3248343893839498", nil)
	if err != nil {
		t.Error(err)
	}
	res := httptest.NewRecorder()
	handler := http.HandlerFunc(whMgr.HandlerGetWebhookHookByID)

	handler.ServeHTTP(res, req)

	// server mock request
	if res.Code == http.StatusOK {
		t.Error("Bad status response: expected %i got %i", http.StatusOK, res.Code)
	}
}
