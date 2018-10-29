package paragliding

import (
	"fmt"
	"net/http"
)

type AdminMgr struct {
	DB *Database
}

// HandlerTrackCount is the handler for GET /admin/api/tracks_count.
//it responds with the number of trascks in the database
func (aMgr *AdminMgr) HandlerTrackCount(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "text/plain")
	trackCount, err := aMgr.DB.GetTrackCount()
	if err != nil {
		http.Error(w, "error getting count from database", http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, trackCount)
}

// HandlerDeleteAllTracks is the handler for DELETE /admin/api/tracks.
//it deletes the tracks and responds of with the number of tracks deleted
func (aMgr *AdminMgr) HandlerDeleteAllTracks(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("content-type", "text/plain")
	trackCount, err := aMgr.DB.DeleteAllTracks()
	if err != nil {
		http.Error(w, "error getting count from database", http.StatusInternalServerError)
		return
	}
	fmt.Fprint(w, trackCount)
}
