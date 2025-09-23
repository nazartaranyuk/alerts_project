package trippers

import "net/http"

type ApiKeyRoundTripper struct {
	APIKey string
	Next   http.RoundTripper
}

func (art *ApiKeyRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	request := req.Clone(req.Context())
	request.Header.Add("Authorization", art.APIKey)
	return art.Next.RoundTrip(request)
}
