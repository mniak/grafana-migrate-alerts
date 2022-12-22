package main

import (
	"fmt"
	"net/http"

	"gopkg.in/resty.v1"
)

func GetGrafanaRules(baseURL, grafanaSession string) (Folders, error) {
	client := resty.New().
		SetHostURL(baseURL).
		SetCookie(&http.Cookie{
			Name:   "grafana_session",
			Value:  grafanaSession,
			Path:   "/",
			Domain: baseURL,
		})

	resp, err := client.R().
		SetResult(Folders{}).
		Get("api/ruler/grafana/api/v1/rules")
	if err != nil {
		return Folders{}, err
	}

	if !resp.IsSuccess() {
		return Folders{}, fmt.Errorf("failed with status code: %d", resp.StatusCode())
	}
	return *resp.Result().(*Folders), nil
}
