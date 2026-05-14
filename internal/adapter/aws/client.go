package aws

import (
	"context"
	"os"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
)

var (
	cfg  aws.Config
	err  error
	once sync.Once
)

// GetAWSConfig returns a shared AWS configuration.
// If AWS_ENDPOINT is set, it configures the BaseEndpoint (useful for LocalStack).
func GetAWSConfig(ctx context.Context) (aws.Config, error) {
	once.Do(func() {
		endpoint := os.Getenv("AWS_ENDPOINT")

		var opts []func(*config.LoadOptions) error
		if endpoint != "" {
			opts = append(opts, config.WithBaseEndpoint(endpoint))
		}

		cfg, err = config.LoadDefaultConfig(ctx, opts...)
	})

	return cfg, err
}
