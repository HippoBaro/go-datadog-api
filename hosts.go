package datadog

import "strconv"

type HostActionResp struct {
	Action   string `json:"action"`
	Hostname string `json:"hostname"`
	Message  string `json:"message,omitempty"`
}

type HostActionMute struct {
	Message  *string `json:"message,omitempty"`
	EndTime  *string `json:"end,omitempty"`
	Override *bool   `json:"override,omitempty"`
}

type HostDef struct {
	Name             string                 `json:"name"`
	Up               bool                   `json:"up"`
	IsMuted          bool                   `json:"is_muted"`
	LastReportedTime int                    `json:"last_reported_time"`
	Apps             []string               `json:"apps"`
	TagsBySource     map[string][]string    `json:"tags_by_source"`
	AwsName          string                 `json:"aws_name"`
	Metrics          map[string]interface{} `json:"metrics"`
	Sources          []string               `json:"sources"`
	Meta             map[string]interface{} `json:"meta"`
	HostName         string                 `json:"host_name"`
	ID               int                    `json:"id"`
	Aliases          []string               `json:"aliases"`
}

// SearchHosts searches through the hosts facet, returning matching hostnames.
func (client *Client) FilterHosts(search string) ([]HostDef, error) {
	type SearchHostsResult struct {
		TotalReturned int       `json:"total_returned"`
		HostList      []HostDef `json:"host_list"`
		TotalMatching int       `json:"total_matching"`
	}

	var start int
	api := func(res interface{}) (string, bool) {
		if res == nil {
			start = 0
		} else if res.(*SearchHostsResult).TotalReturned == 100 {
			start += 80
		} else {
			return "", false
		}
		return "/v1/hosts?filter=" + search + "&start=" + strconv.Itoa(start), true
	}

	var out SearchHostsResult
	resCombined := make([]HostDef, 0)
	combine := func(res interface{}) error {
		resCombined = append(resCombined, res.(*SearchHostsResult).HostList...)
		return nil
	}

	if err := client.doJsonRequestPaginated("GET", api, "", &out, combine); err != nil {
		return nil, err
	}
	return resCombined, nil
}

// MuteHost mutes all monitors for the given host
func (client *Client) MuteHost(host string, action *HostActionMute) (*HostActionResp, error) {
	var out HostActionResp
	uri := "/v1/host/" + host + "/mute"
	if err := client.doJsonRequest("POST", uri, action, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// UnmuteHost unmutes all monitors for the given host
func (client *Client) UnmuteHost(host string) (*HostActionResp, error) {
	var out HostActionResp
	uri := "/v1/host/" + host + "/unmute"
	if err := client.doJsonRequest("POST", uri, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}

// HostTotalsResp defines response to GET /v1/hosts/totals.
type HostTotalsResp struct {
	TotalUp     *int `json:"total_up"`
	TotalActive *int `json:"total_active"`
}

// GetHostTotals returns number of total active hosts and total up hosts.
// Active means the host has reported in the past hour, and up means it has reported in the past two hours.
func (client *Client) GetHostTotals() (*HostTotalsResp, error) {
	var out HostTotalsResp
	uri := "/v1/hosts/totals"
	if err := client.doJsonRequest("GET", uri, nil, &out); err != nil {
		return nil, err
	}
	return &out, nil
}
