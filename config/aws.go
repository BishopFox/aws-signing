package config

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/BishopFox/aws-signing/signing"
	"github.com/aws/aws-sdk-go-v2/aws/signer/v4"
	awscfg "github.com/aws/aws-sdk-go-v2/config"
)

var (
	DefaultsAws = Aws{
		Use:      false,
		Endpoint: "http://localhost:9200",
		Service:  "es",
	}
	EnvAws = Aws{
		Use:      hasEnvVar("AWS_SIGNING"),
		Endpoint: os.Getenv("AWS_ENDPOINT"),
		Service:  os.Getenv("AWS_SERVICE"),
	}
)

type Aws struct {
	Use      bool   `hcl:"enabled"`
	Service  string `hcl:"service"`
	Endpoint string `hcl:"endpoint"`
}

func (a Aws) EndpointUrl() *url.URL {
	return parseUrl(a.Endpoint, nil)
}

func MergeAws(cfgs ...Aws) Aws {
	rv := Aws{}
	for _, cur := range cfgs {
		rv.Use = rv.Use || cur.Use
		if cur.Endpoint != "" {
			rv.Endpoint = cur.Endpoint
		}
		if cur.Service != "" {
			rv.Service = cur.Service
		}
	}
	return rv
}

func (a Aws) Transport() (http.RoundTripper, error) {
	cfg, err := awscfg.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, fmt.Errorf("error loading aws config: %s", err)
	}
	if region := os.Getenv("AWS_REGION"); region != "" {
		cfg.Region = region
	}
	signer := v4.NewSigner()
	creds, err := cfg.Credentials.Retrieve(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve credentials: %w", err)
	}
	return signing.NewTransport(signer, creds, a.Service, cfg.Region), nil
}
