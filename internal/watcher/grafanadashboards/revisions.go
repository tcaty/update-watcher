package grafanadashboards

import "time"

type Revisions struct {
	Items []struct {
		ID            int       `json:"id"`
		DashboardID   int       `json:"dashboardId"`
		DashboardName string    `json:"dashboardName"`
		Revision      int       `json:"revision"`
		Description   string    `json:"description"`
		CreatedAt     time.Time `json:"createdAt"`
		UpdatedAt     time.Time `json:"updatedAt"`
		Downloads     int       `json:"downloads"`
		OrgID         int       `json:"orgId"`
		OrgName       string    `json:"orgName"`
		OrgSlug       string    `json:"orgSlug"`
		Links         []struct {
			Rel  string `json:"rel"`
			Href string `json:"href"`
		} `json:"links"`
	} `json:"items"`
	OrderBy   string `json:"orderBy"`
	Direction string `json:"direction"`
	Links     []struct {
		Rel  string `json:"rel"`
		Href string `json:"href"`
	} `json:"links"`
}
