package main

import (
	"encoding/gob"
	"encoding/json"
	"github.com/robfig/cron"
	"net/http"
	"os"
	"time"
)

type Preferences struct {
	RefreshRate        string   `json:"refresh_rate"`
	NewUserPermissions UserRole `json:"new_user_permissions"`
}

func getOrCreatePrefs() (*os.File, error) {
	f, err := os.OpenFile(ExePath+"/prefs.gob", os.O_RDWR, 0666)
	if err != nil && os.IsNotExist(err) {
		if f, err := os.Create(ExePath + "/prefs.gob"); err != nil {
			return nil, err
		} else {
			p := Preferences{"30m", PublicRole}
			WritePreferences(&p)
			return f, nil
		}
	}
	return f, err
}

func WritePreferences(prefs *Preferences) error {
	prefFile, err := getOrCreatePrefs()
	if err != nil {
		return err
	}
	enc := gob.NewEncoder(prefFile)
	err = enc.Encode(prefs)
	return err
}

func ReadPreferences() (*Preferences, error) {
	f, err := getOrCreatePrefs()
	if err != nil {
		return nil, err
	}
	var prefs Preferences
	dec := gob.NewDecoder(f)
	err = dec.Decode(&prefs)
	if err != nil {
		return nil, err
	}
	return &prefs, nil
}

func ReloadPreferences() error {
	prefs, err := ReadPreferences()
	if err != nil {
		return err
	}
	UserPreferences = prefs
	if CRON != nil {
		CRON.Stop()
	}
	CRON = cron.New()
	CRON.AddFunc("@every "+UserPreferences.RefreshRate, UpdateArticles)
	CRON.Start()
	return nil
}

func GetPreferences(w http.ResponseWriter, req *http.Request) {
	sessionToken := req.Header.Get("x-session-token")
	if sessionToken == "" {
		WriteJSONError(w, http.StatusBadRequest, "Session token not provided")
		return
	}
	u, err, code := GetUserForSessionToken(sessionToken)
	if err != nil {
		WriteJSONError(w, code, err.Error())
	} else if u.Role != AdminRole {
		WriteJSONError(w, http.StatusUnauthorized, "You must be an administrator to read or write  server preferences")
		return
	}
	prefs, err := ReadPreferences()
	if err != nil {
		w.Header().Set("content-type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"error": "Couldn't load preferences"}`))
		return
	}
	enc := json.NewEncoder(w)
	enc.Encode(prefs)
}

func PutPreferences(w http.ResponseWriter, req *http.Request) {
	sessionToken := req.Header.Get("x-session-token")
	if sessionToken == "" {
		WriteJSONError(w, http.StatusBadRequest, "Session token not provided")
		return
	}
	u, err, code := GetUserForSessionToken(sessionToken)
	if err != nil {
		WriteJSONError(w, code, err.Error())
	} else if u.Role != AdminRole {
		WriteJSONError(w, http.StatusUnauthorized, "You must be an administrator to read or write  server preferences")
		return
	}
	dec := json.NewDecoder(req.Body)
	var prefs Preferences
	err = dec.Decode(&prefs)
	if err != nil {
		WriteJSONError(w, http.StatusBadRequest, "Malformed JSON or missing field")
		return
	}
	_, err = time.ParseDuration(prefs.RefreshRate)
	if err != nil {
		WriteJSONError(w, http.StatusBadRequest, "Invalid duration")
		return
	}
	err = WritePreferences(&prefs)
	if err != nil {
		WriteJSONError(w, http.StatusInternalServerError, "Couldn't write preferences")
		return
	}
	ReloadPreferences()
}
