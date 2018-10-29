package paragliding

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/mongodb/mongo-go-driver/bson/objectid"
)

// WebHookMgr is the manager for webhooks
type WebHookMgr struct {
	DB     *Database
	Ticker *MgrTicker
}

// HandlerNewTrackWebHook is the handler for POST /api/webhook/new_track/.
// it registers a new webhook and reponds with the id assigned to it
func (whMgr *WebHookMgr) HandlerNewTrackWebHook(w http.ResponseWriter, r *http.Request) {
	var postData map[string]string
	err := json.NewDecoder(r.Body).Decode(&postData)
	if err == nil {
		triggerVal, err2 := strconv.Atoi(postData["minTriggerValue"])
		if err2 != nil {
			http.Error(w, "triggervalue is not a number", http.StatusBadRequest)
			return
		}
		minTriggerVal, _ := strconv.ParseInt(postData["minTriggerValue"], 10, 64) // guaranteed to be number cause regex checks in url
		wekbookInfo := WebhookInfo{ID: objectid.New(), WebhookURL: postData["webhookURL"], MinTriggerValue: int64(triggerVal), Counter: minTriggerVal, LatestTimestamp: (time.Now().UnixNano() / int64(time.Millisecond))}
		id, added := whMgr.DB.Insert("webhooks", wekbookInfo)
		if added {
			w.Header().Add("content-type", "application/json")
			json.NewEncoder(w).Encode(struct {
				ID string `json:"id"`
			}{id})
		} else {
			http.Error(w, "track already exists with id: "+id, http.StatusBadRequest)
		}
	} else if err == io.EOF {
		http.Error(w, "POST body is empty", http.StatusBadRequest)
	} else {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
	}
}

// HandlerGetWebhookHookByID is the handler for "GET /api/webhook/new_track/<webhook_id>"
// it responds with the webhoo url and the minimum trigger value
func (whMgr *WebHookMgr) HandlerGetWebhookHookByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "application/json")
	parts := strings.Split(r.URL.Path, "/")
	webhookInfo, found := whMgr.DB.GetWebhookByID(parts[len(parts)-1]) // guaranteed to be valid cause of regex in server.go
	if !found {
		http.Error(w, "the id does not exist", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(webhookInfo)
}

// HandlerDeleteWebhookHookByID is the handler for "DELETE /api/webhook/new_track/<webhook_id>"
// it deletes the webhook and reponds with the webhook's info
func (whMgr *WebHookMgr) HandlerDeleteWebhookHookByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "application/json")
	parts := strings.Split(r.URL.Path, "/")
	webhookInfo, found := whMgr.DB.GetWebhookByID(parts[len(parts)-1]) // guaranteed to be valid cause of regex in server.go
	if !found {
		http.Error(w, "the id does not exist", http.StatusNotFound)
		return
	}
	err := whMgr.DB.DeleteWebhookByID(parts[len(parts)-1])
	if err != nil {
		http.Error(w, "could not delete webhook", http.StatusInternalServerError)
	}
	json.NewEncoder(w).Encode(webhookInfo)
}

// InvokeNewWebHooks should be called when a new track is added. it will invoke the webohooks that should be invoked
func (whMgr *WebHookMgr) InvokeNewWebHooks() {
	webhooks, err := whMgr.DB.GetAllInvokeWebhooks()
	if err != nil {
		return
	}
	for _, v := range webhooks {
		startTime := time.Now()

		tickerResp, err := whMgr.Ticker.GetTickerByTimeStamp(v.LatestTimestamp)
		if err != nil {
			fmt.Println(err)
			continue
		}

		var trackIdsString string
		nOfNewTrack := len(tickerResp.TrackIDs)

		if nOfNewTrack <= 0 {
			continue
		}

		trackIdsString = tickerResp.TrackIDs[0].Hex()
		for i := 1; i < nOfNewTrack; i++ {
			trackIdsString += ", " + tickerResp.TrackIDs[i].Hex()
		}

		areOrIsString := "is: "
		if nOfNewTrack > 1 {
			areOrIsString = "are: "
		}

		v.LatestTimestamp = tickerResp.TLatest
		reponseString := "latest timestamp: " + strconv.FormatInt(tickerResp.TLatest, 10) +
			", " + strconv.Itoa(len(tickerResp.TrackIDs)) + " new tracks " + areOrIsString +
			trackIdsString + ". (processing: " + strconv.FormatFloat(float64(time.Since(startTime))/float64(time.Millisecond), 'f', 2, 64) + "ms)"

		var jsonStr = []byte(`{"content":"` + reponseString + `"}`)
		//req, err := http.NewRequest("POST", v.WebhookURL, bytes.NewBuffer(jsonStr))
		_, postErr := http.Post(v.WebhookURL, "application/json", bytes.NewBuffer(jsonStr))
		if postErr != nil {
			fmt.Println(postErr)
		}
		whMgr.Ticker.DB.ResetWebhookCounter(v)

	}

}
