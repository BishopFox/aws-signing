package signing

import (
	"net/http"
	"testing"

	"github.com/aws/aws-sdk-go-v2/aws"
	v4 "github.com/aws/aws-sdk-go-v2/aws/signer/v4"
)

type mockTransport struct {
	RoundTripFunc func(req *http.Request) (*http.Response, error)
}

func (m mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return m.RoundTripFunc(req)
}

func TestTransport_All(t *testing.T) {
	creds := aws.Credentials{
		AccessKeyID:     "a",
		SecretAccessKey: "b",
		SessionToken:    "c",
	}
	signer := v4.NewSigner()
	tests := []struct {
		signer  Signer
		service string
		region  string
		wantErr string
	}{
		{nil, "", "", MissingSigner.Error()},
		{signer, "", "", MissingService.Error()},
		{signer, "es", "", MissingRegion.Error()},
		{signer, "es", "us-east-1", ""},
	}

	for i, test := range tests {
		transport := NewTransport(test.signer, creds, test.service, test.region)
		gotAuthorization := ""
		transport.BaseTransport = mockTransport{
			RoundTripFunc: func(req *http.Request) (*http.Response, error) {
				gotAuthorization = req.Header.Get("Authorization")
				if req.Header.Get("X-Forwarded-For") != "" {
					t.Error("X-Forwarded-For should be removed from signed request")
				}
				return nil, nil
			},
		}
		req, _ := http.NewRequest("GET", "/", nil)
		_, err := transport.RoundTrip(req)
		if err == nil {
			if test.wantErr != "" {
				t.Errorf("[%d] got no error, expected %q", i, test.wantErr)
			}
			if gotAuthorization == "" {
				t.Error("expected Authorization header to be set, was empty")
			}
		} else {
			if test.wantErr == "" {
				t.Errorf("[%d] expected no error, got %q", i, err)
			} else if err.Error() != test.wantErr {
				t.Errorf("[%d] expected %q, got %q", i, test.wantErr, err)
			}
		}
	}
}
