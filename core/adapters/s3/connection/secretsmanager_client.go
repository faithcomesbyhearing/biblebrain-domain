package connection

import (
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/secretsmanager"
)

func GetSecretsManagerClient() *secretsmanager.Client {

	// NOTE: typically, if IS_OFFLINE is true, we would configure a local endpoint for the service.
	// However, it does not appear that serverless_offline_ssm exposes an endpoint.
	// TODO: if/when we identify an endpoint for serverless_offline_ssm, configure it here
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithDefaultRegion("us-west-2"),
	)

	if err != nil {
		log.Panic(err)
	}
	return secretsmanager.NewFromConfig(cfg)
}
