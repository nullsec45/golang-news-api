package config

import (
	"context"
	"github.com/rs/zerolog/log"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
)

func (cfg Config) LoadAWSConfig() aws.Config{
	conf, err := awsConfig.LoadDefaultConfig(context.TODO(),
					awsConfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
						cfg.R2.ApiKey, cfg.R2.ApiSecret, "",
					)), awsConfig.WithRegion("auto"),
                )

	if err != nil {
		log.Fatal().Msgf("unable to load AWS Config, %v", err)
	}

	log.Info().Msgf("Success to loaded AWS Config.")

	return conf
}