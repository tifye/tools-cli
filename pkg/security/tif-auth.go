package security

import (
	"net/http"
)

type TifAuthTransport struct {
	wrappedTripper http.RoundTripper
	apiKey         string
	accessToken    string
}

func NewTifAuthTransport(wrappedTripper http.RoundTripper, apiKey, accessToken string) *TifAuthTransport {
	return &TifAuthTransport{
		wrappedTripper: wrappedTripper,
		apiKey:         apiKey,
		accessToken:    accessToken,
	}
}

func (t *TifAuthTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("x-api-key", "fruit-pie")
	req.Header.Set("token", t.accessToken)
	req.Header.Set("Ocp-Apim-Subscription-Key", t.apiKey)
	return t.wrappedTripper.RoundTrip(req)
}
