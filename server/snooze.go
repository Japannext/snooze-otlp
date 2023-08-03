package server

import (
	"bytes"
	"encoding/json"
  "fmt"
	"net/http"

  log "github.com/sirupsen/logrus"
)

type SnoozeClient struct {
	client *http.Client
  url string
  cacert string
}

var Snooze SnoozeClient

// Send a list of alerts to snooze
func (s *SnoozeClient) send(alerts []SnoozeAlertV1) error {
	buf := new(bytes.Buffer)

	err = json.NewEncoder(buf).Encode(alerts)
	if err != nil {
		return err
	}

  url := fmt.Sprintf("%s/api/alert", s.url)

  log.Debugf("Attempting POST %s", url)
  _, err := s.client.Post(url, "application/json", buf)
  if err != nil {
    log.Errorf("Error while writing to snooze: %s", err)
    return err
  }

	return nil
}

func initSnooze() {
  Snooze = SnoozeClient{
    client: &http.Client{},
    url: Config.SnoozeUrl,
    cacert: Config.SnoozeCaPath,
  }
}
