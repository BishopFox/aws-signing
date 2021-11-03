package cli

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"

	awscfg "github.com/aws/aws-sdk-go-v2/config"

	"github.com/BishopFox/aws-signing/config"
)

type AwsArgs struct {
	use       bool
	endpoint  string
	service   string
	dumpCreds bool
}

func (a *AwsArgs) AddFlags(flags *flag.FlagSet) {
	flags.BoolVar(&a.use, "aws", false, "use aws request signing")
	flags.BoolVar(&a.use, "a", false, "use aws request signing (shorthand)")

	flags.StringVar(&a.endpoint, "aws-endpoint", "", "aws endpoint url")
	flags.StringVar(&a.endpoint, "e", "", "aws endpoint url (shorthand)")

	flags.StringVar(&a.service, "aws-service", "", "aws service")
	flags.StringVar(&a.service, "s", "", "aws service (shorthand)")

	flags.BoolVar(&a.dumpCreds, "creds", false, "emit creds")
}

func (AwsArgs) Options() string {
	return `
 -a, --aws                    Use AWS Request Signing
                              Default: false
                              Env Var: AWS_SIGNING

 -e, --aws-endpoint <url>     AWS Endpoint URL.
                              Default: http://localhost:9200
                              Env Var: AWS_ENDPOINT

 -s, --aws-service <service>  AWS Service.
                              Default: es
                              Env Var: AWS_SERVICE
`
}

func (a AwsArgs) Config() (config.Aws, error) {
	if a.endpoint != "" {
		if _, err := url.Parse(a.endpoint); err != nil {
			return config.Aws{}, fmt.Errorf("error parsing endpoint url %q: %s", a.endpoint, err)
		}
	}

	return config.Aws{
		Use:      a.use,
		Endpoint: a.endpoint,
		Service:  a.service,
	}, nil
}

func (a AwsArgs) Dump() {
	if !a.dumpCreds {
		return
	}
	if _, ok := os.LookupEnv("AWS_SIGNING"); a.use || ok {
		ctx := context.Background()
		cfg, err := awscfg.LoadDefaultConfig(ctx)
		if err != nil {
			return
		}
		log.Println(cfg.Credentials.Retrieve(ctx))
	}
}
