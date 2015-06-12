package honeybadger

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

var (
	RateExceeded    = errors.New("Rate exceeded: slow down!")
	PaymentRequired = errors.New("Payment required: expired trial or credit card?")
	Unauthorized    = errors.New("Unauthorized: bad API key?")
)

func newServerBackend(config *Configuration) Server {
	return Server{
		URL:    &config.Endpoint,
		APIKey: &config.APIKey,
		Client: &http.Client{
			Transport: http.DefaultTransport,
			Timeout:   3 * time.Second,
		},
	}
}

type Server struct {
	APIKey *string
	URL    *string
	Client *http.Client
}

func (s Server) Notify(feature Feature, payload Payload) error {
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

	resp, err := s.Client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case 201:
		return nil
	case 429, 503:
		return RateExceeded
	case 402:
		return PaymentRequired
	case 403:
		return Unauthorized
	default:
		return fmt.Errorf(
			"request failed status=%d expected=%d",
			resp.StatusCode,
			http.StatusCreated)
	}
}
