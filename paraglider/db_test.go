package paragliding

import (
	"testing"

	"github.com/mongodb/mongo-go-driver/bson/objectid"
)

const testDBURI = "mongodb://test:test12@ds143893.mlab.com:43893/a2-testddb"
const testDBNAME = "a2-testddb"

// dont rly know how to handle user name and password in tests
// as environment varibles will fail here if teachers try to run it on their pc
// os.Getenv("DB_URI") and os.Getenv("DB_NAME") are the environment variables i use locally an don heroku

func Test_Connect(t *testing.T) {
	db := &Database{URI: "mongodb://test:test12@ds141783.mlab.com:41783/a2-trackdb", Name: "a2-trackdb"}
	db.Connect()
	if db.conn == nil {
		t.Error("Failed to connect to database")
	}
}

func Test_Insert(t *testing.T) {
	db := &Database{URI: testDBURI, Name: testDBNAME}
	db.Connect()
	db.DeleteAllTracksAndWebhooks()

	// inserting webhook
	wekbookInfo := WebhookInfo{ID: objectid.New(), WebhookURL: "www.testurl.com", MinTriggerValue: 2, Counter: 2, LatestTimestamp: 2}
	_, whAdded := db.Insert("webhooks", wekbookInfo)
	if !whAdded {
		t.Error("Failed to insert webhook")
	}

	// inserting track
	trackInfo := TrackInfo{ID: objectid.New(), HDate: "somedate", Pilot: "ole",
		Glider: "sometype", GliderID: "someID", TrackLength: "10",
		TrackURL: "www.trackurl.com", Timestamp: 120}
	_, trackAdded := db.Insert("tracks", trackInfo)
	if !trackAdded {
		t.Error("Failed to insert track")
	}
}

func Test_GetAllTrackIDs(t *testing.T) {
	db := &Database{URI: testDBURI, Name: testDBNAME}
	db.Connect()
	db.DeleteAllTracksAndWebhooks()

	// add 3 tracks
	var newIDs [3]objectid.ObjectID
	for i := 0; i < 3; i++ {
		newIDs[i] = objectid.New()
		db.Insert("tracks", TrackInfo{ID: newIDs[i], HDate: "somedate", Pilot: "ole",
			Glider: "sometype", GliderID: "someID", TrackLength: "10",
			TrackURL: "www.trackurl.com", Timestamp: 120})
	}

	// get all track ids from db and see if they match
	resIDs, err := db.GetAllTrackIDs()
	if err != nil || len(resIDs) != 3 {
		t.Error("coult not get list of ids")
	} else {
		for i := 0; i < 3; i++ {
			if resIDs[i] != newIDs[i] {
				t.Error("ids do not match")
			}
		}
	}

}

func Test_GetTrackByID(t *testing.T) {
	db := &Database{URI: testDBURI, Name: testDBNAME}
	db.Connect()
	db.DeleteAllTracksAndWebhooks()

	newID := objectid.New()
	db.Insert("tracks", TrackInfo{ID: newID, HDate: "somedate", Pilot: "ole",
		Glider: "sometype", GliderID: "someID", TrackLength: "10",
		TrackURL: "www.trackurl.com", Timestamp: 120})

	trackInfo, exists := db.GetTrackByID(newID.Hex())
	if !exists {
		t.Error("couldnt find the newly inserted track")
	} else if trackInfo.ID != newID {
		t.Error("ids dont match")
	}
}

func Test_GetTrackCount(t *testing.T) {
	db := &Database{URI: testDBURI, Name: testDBNAME}
	db.Connect()
	db.DeleteAllTracksAndWebhooks()

	db.Insert("tracks", TrackInfo{ID: objectid.New(), HDate: "somedate", Pilot: "ole",
		Glider: "sometype", GliderID: "someID", TrackLength: "10",
		TrackURL: "www.trackurl.com", Timestamp: 120})

	count, err := db.GetTrackCount()
	if err != nil {
		t.Error("couldnt get track count")
	} else if count != 1 {
		t.Error("got wrong track count")
	}
}

func Test_DeleteAllTracks(t *testing.T) {
	db := &Database{URI: testDBURI, Name: testDBNAME}
	db.Connect()
	db.DeleteAllTracksAndWebhooks()

	// insert a track
	db.Insert("tracks", TrackInfo{ID: objectid.New(), HDate: "somedate", Pilot: "ole",
		Glider: "sometype", GliderID: "someID", TrackLength: "10",
		TrackURL: "www.trackurl.com", Timestamp: 120})

	// delete tarcks and see if number of tracks deleted is 1
	count, err := db.DeleteAllTracks()
	if err != nil {
		t.Error("couldnt delete tracks")
	} else if count != 1 {
		t.Error("wrong number of tracks left in database")
	}
}

func Test_GetAllTracks(t *testing.T) {
	db := &Database{URI: testDBURI, Name: testDBNAME}
	db.Connect()
	db.DeleteAllTracksAndWebhooks()

	// add 3 tracks
	var newTracks [3]TrackInfo
	for i := 0; i < 3; i++ {
		track := TrackInfo{ID: objectid.New(), HDate: "somedate", Pilot: "ole",
			Glider: "sometype", GliderID: "someID", TrackLength: "10",
			TrackURL: "www.trackurl.com", Timestamp: 120}
		newTracks[i] = track
		db.Insert("tracks", track)
	}

	// get all track from db and see if they match with the ones added
	resTracks, err := db.GetAllTracks()
	if err != nil || len(resTracks) != 3 {
		t.Error("coult not get list of tracks")
	} else {
		for i := 0; i < 3; i++ {
			if resTracks[i] != newTracks[i] {
				t.Error("ids do not match")
			}
		}
	}
}

func Test_GetWebhookByID(t *testing.T) {
	db := &Database{URI: testDBURI, Name: testDBNAME}
	db.Connect()
	db.DeleteAllTracksAndWebhooks()

	// inserts a webhook into db
	newID := objectid.New()
	db.Insert("webhooks", WebhookInfo{ID: newID, WebhookURL: "www.testurl.com",
		MinTriggerValue: 2, Counter: 2, LatestTimestamp: 2})

	// attemps to get the newly inserted webhook by id
	whInfo, exists := db.GetWebhookByID(newID.Hex())
	if !exists {
		t.Error("couldnt find the newly inserted webhook")
	} else if whInfo.ID != newID {
		t.Error("ids dont match")
	}
}

func Test_DeleteWebhookByID(t *testing.T) {
	db := &Database{URI: testDBURI, Name: testDBNAME}
	db.Connect()
	db.DeleteAllTracksAndWebhooks()

	// inserts a webhook into db
	newID := objectid.New()
	db.Insert("webhooks", WebhookInfo{ID: newID, WebhookURL: "www.testurl.com",
		MinTriggerValue: 2, Counter: 2, LatestTimestamp: 2})

	// attempts to delete it
	err := db.DeleteWebhookByID(newID.Hex())
	if err != nil {
		t.Error("failed to delete webhook by id")
	}
}

func Test_GetAllInvokeWebhooks(t *testing.T) {
	db := &Database{URI: testDBURI, Name: testDBNAME}
	db.Connect()
	db.DeleteAllTracksAndWebhooks()

	// add two webhooks to the db. one that shoudl trigger in 1 call and one that should not.
	// then check if we received the correct one

	webhookNoInvoke := WebhookInfo{ID: objectid.New(), WebhookURL: "www.testurl.com",
		MinTriggerValue: 2, Counter: 2, LatestTimestamp: 2}
	webhookInvoke := WebhookInfo{ID: objectid.New(), WebhookURL: "www.testurl2.com",
		MinTriggerValue: 2, Counter: 1, LatestTimestamp: 2}

	db.Insert("webhooks", webhookNoInvoke)
	db.Insert("webhooks", webhookInvoke)

	webhooks, err := db.GetAllInvokeWebhooks()
	if err != nil {
		t.Error("could not get webhooks")
	} else if len(webhooks) != 1 {
		t.Error("received wrong number of webhooks")
	} else if webhooks[0] == webhookNoInvoke {
		t.Error("received the wrong webhook")
	}
}

func Test_ResetWebhookCounter(t *testing.T) {
	db := &Database{URI: testDBURI, Name: testDBNAME}
	db.Connect()
	db.DeleteAllTracksAndWebhooks()

	// create and add a webhook to the database
	webhook := WebhookInfo{ID: objectid.New(), WebhookURL: "www.testurl.com",
		MinTriggerValue: 2, Counter: 0, LatestTimestamp: 2}

	db.Insert("webhooks", webhook)
	webhook.LatestTimestamp = 10

	db.ResetWebhookCounter(webhook)
	webhook.Counter = 2 // counter of the updated webhook should be 2

	// get the webhook and see if it matches
	whInfo, exists := db.GetWebhookByID(webhook.ID.Hex())
	if !exists {
		t.Error("didnt receive the webhook from database")
	} else if whInfo != webhook {
		t.Error("webhook didnt update correctly")
	}
}
