package honeybadger

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

type Server struct {
	APIKey *string
	URL    *string
}

func (s Server) Notify(feature Feature, payload Payload) error {
	client := &http.Client{
		Timeout: 3 * time.Second,
	}

	url, err := url.Parse(*s.URL)
	if err != nil {
		return err
	}
	url.Path = "v1/" + feature.Endpoint
	req, err := http.NewRequest("POST", url.String(), bytes.NewReader(payload.toJSON()))
	if err != nil {
		return err
	}

	req.Header.Set("X-API-Key", *s.APIKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf(
			"request failed status=%d expected=%d",
			resp.StatusCode,
			http.StatusCreated)
	}

	return nil
}
