package paragliding

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	igc "github.com/marni/goigc"
	"github.com/mongodb/mongo-go-driver/bson/objectid"
)

type TrackMgr struct {
	DB    *Database
	WHMgr *WebHookMgr
}

// HandlerPostTrack is the handler for POST /api/track. it registers the track and replies with the id
func (tMgr *TrackMgr) HandlerPostTrack(w http.ResponseWriter, r *http.Request) {
	var postData map[string]string
	err := json.NewDecoder(r.Body).Decode(&postData)
	if err == nil {
		track, err2 := igc.ParseLocation(postData["url"])
		if err2 != nil {
			http.Error(w, "could not get a track from url: "+postData["url"], http.StatusNotFound)
			return
		}
		trackInfo := TrackInfo{ID: objectid.New(), HDate: track.Date.String(), Pilot: track.Pilot,
			Glider: track.GliderType, GliderID: track.GliderID, TrackLength: CalculatedistanceFromPoints(track.Points),
			TrackURL: postData["url"], Timestamp: (time.Now().UnixNano() / int64(time.Millisecond))}
		id, added := tMgr.DB.Insert("tracks", trackInfo)
		if added {
			w.Header().Add("content-type", "application/json")
			json.NewEncoder(w).Encode(struct {
				ID string `json:"id"`
			}{id})
			tMgr.WHMgr.InvokeNewWebHooks() // invoke webhooks cause new track is added
		} else {
			http.Error(w, "track already exists with id: "+id, http.StatusBadRequest)
		}
	} else if err == io.EOF {
		http.Error(w, "POST body is empty", http.StatusBadRequest)
	} else {
		http.Error(w, "Something went wrong", http.StatusInternalServerError)
	}
}

// HandlerGetAllTracks is the handler for GET /api/track. it replies with an array of all track ids
func (tMgr *TrackMgr) HandlerGetAllTracks(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "application/json")
	ids, err := tMgr.DB.GetAllTrackIDs()
	if err != nil {
		http.Error(w, "Could not receive track list", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(ids)

}

// HandlerGetTrackByID is the handler for GET /api/track/<id>. it responds with info about the track
func (tMgr *TrackMgr) HandlerGetTrackByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "application/json")
	parts := strings.Split(r.URL.Path, "/")
	trackInfo, found := tMgr.DB.GetTrackByID(parts[len(parts)-1]) // guaranteed to be valid cause of regex in server.go
	if !found {
		http.Error(w, "the id does not exist", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(trackInfo)
}

// HandlerGetTrackFieldByID is the handler for GET /api/track/<id><field>. is reponds with the single informationc ontained in that field
func (tMgr *TrackMgr) HandlerGetTrackFieldByID(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "text/plain")
	parts := strings.Split(r.URL.Path, "/")
	trackInfo, found := tMgr.DB.GetTrackByID(parts[len(parts)-2]) // guaranteed to be valid cause of regex in server.go

	if !found {
		http.Error(w, "the id does not exist", http.StatusNotFound)
		return
	}
	field := parts[len(parts)-1]
	switch field {
	case "pilot":
		fmt.Fprintf(w, "pilot: %s", trackInfo.Pilot)
	case "glider":
		fmt.Fprintf(w, "glider: %s", trackInfo.Glider)
	case "glider_id":
		fmt.Fprintf(w, "glider_id: %s", trackInfo.GliderID)
	case "H_date":
		fmt.Fprintf(w, "H_date: %s", trackInfo.HDate)
	case "track_length":
		fmt.Fprintf(w, "track_length: %s", trackInfo.TrackLength)
	case "track_src_url":
		fmt.Fprintf(w, "track_src_url: %s", trackInfo.TrackLength)
	default:
		http.Error(w, "invalid field specified", http.StatusNotFound)
	}
}

// CalculatedistanceFromPoints take a set of points and retunr the total distance
func CalculatedistanceFromPoints(points []igc.Point) string {
	d := 0.0
	for i := 0; i < len(points)-1; i++ {
		d += points[i].Distance(points[i+1])
	}
	return strconv.FormatFloat(d, 'f', 2, 64)
}
