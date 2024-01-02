package dockerregistry

import "time"

type Tags struct {
	Count    int    `json:"count"`
	Next     string `json:"next"`
	Previous any    `json:"previous"`
	Results  []struct {
		Creator int `json:"creator"`
		ID      int `json:"id"`
		Images  []struct {
			Architecture string    `json:"architecture"`
			Features     string    `json:"features"`
			Variant      any       `json:"variant"`
			Digest       string    `json:"digest"`
			Os           string    `json:"os"`
			OsFeatures   string    `json:"os_features"`
			OsVersion    any       `json:"os_version"`
			Size         int       `json:"size"`
			Status       string    `json:"status"`
			LastPulled   time.Time `json:"last_pulled"`
			LastPushed   time.Time `json:"last_pushed"`
		} `json:"images"`
		LastUpdated         time.Time `json:"last_updated"`
		LastUpdater         int       `json:"last_updater"`
		LastUpdaterUsername string    `json:"last_updater_username"`
		Name                string    `json:"name"`
		Repository          int       `json:"repository"`
		FullSize            int       `json:"full_size"`
		V2                  bool      `json:"v2"`
		TagStatus           string    `json:"tag_status"`
		TagLastPulled       time.Time `json:"tag_last_pulled"`
		TagLastPushed       time.Time `json:"tag_last_pushed"`
		MediaType           string    `json:"media_type"`
		ContentType         string    `json:"content_type"`
		Digest              string    `json:"digest"`
	} `json:"results"`
}
