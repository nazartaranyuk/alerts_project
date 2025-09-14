package trippers

import "net/http"

type AuthRoundTripper struct {
	APIKey string
	Next   http.RoundTripper
}

func (art *AuthRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	request := req.Clone(req.Context())
	request.Header.Add("Authorization", art.APIKey)
	return art.Next.RoundTrip(request)
}
