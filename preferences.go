package main

import (
	"encoding/gob"
	"github.com/robfig/cron"
	"os"
)

type Preferences struct {
	RefreshRate        string   `json:"refresh_rate"`
	NewUserPermissions UserRole `json:"new_user_permissions"`
}

func getOrCreatePrefs() (*os.File, error) {
	f, err := os.OpenFile(ExePath + "/prefs.gob", os.O_RDWR, 0666)
	if err != nil && os.IsNotExist(err) {
		if _, err := os.Create(ExePath + "/prefs.gob"); err != nil {
			return nil, err
		} else {
			p := Preferences{"@every 30m", PublicRole}
			WritePreferences(&p)
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
	CRON.AddFunc(UserPreferences.RefreshRate, DownloadArticles)
	CRON.Start()
	return nil
}
