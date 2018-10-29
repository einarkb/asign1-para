package paragliding

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func (mgrTicker *MgrTicker) Test_HandlerLatestTick(t *testing.T) {
	req, err := http.NewRequest("GET", "/paragliding/api/ticker/", nil)
	if err != nil {
		t.Error(err)
	}
	res := httptest.NewRecorder()
	handler := http.HandlerFunc(mgrTicker.HandlerLatestTick)

	handler.ServeHTTP(res, req)

	if res.Code != http.StatusOK {
		t.Error("Bad status response: expected %i got %i", http.StatusOK, res.Code)
	}
}
