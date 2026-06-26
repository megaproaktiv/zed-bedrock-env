package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"aws-zed/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	awsCfg, err := awsconfig.LoadDefaultConfig(ctx, awsconfig.WithRegion(cfg.Region))
	if err != nil {
		return fmt.Errorf("load AWS config: %w", err)
	}

	stsClient := sts.NewFromConfig(awsCfg)

	if cfg.AccountID == "" && cfg.RoleARN == "" {
		identity, err := stsClient.GetCallerIdentity(ctx, &sts.GetCallerIdentityInput{})
		if err != nil {
			return fmt.Errorf("get caller identity: %w", err)
		}
		cfg.AccountID = aws.ToString(identity.Account)
	}

	roleARN, err := cfg.ResolveRoleARN()
	if err != nil {
		return err
	}

	out, err := stsClient.AssumeRole(ctx, &sts.AssumeRoleInput{
		RoleArn:         aws.String(roleARN),
		RoleSessionName: aws.String(cfg.SessionName),
		DurationSeconds: aws.Int32(cfg.DurationSeconds),
	})
	if err != nil {
		return fmt.Errorf("assume role %q: %w", roleARN, err)
	}

	creds := out.Credentials
	if creds == nil || creds.AccessKeyId == nil || creds.SecretAccessKey == nil || creds.SessionToken == nil {
		return fmt.Errorf("assume role returned incomplete credentials")
	}

	fmt.Printf("export ZED_ACCESS_KEY_ID=\"%s\"\n", *creds.AccessKeyId)
	fmt.Printf("export ZED_SECRET_ACCESS_KEY=\"%s\"\n", *creds.SecretAccessKey)
	fmt.Printf("export ZED_SESSION_TOKEN=\"%s\"\n", *creds.SessionToken)
	fmt.Printf("export ZED_AWS_REGION=%s\n", cfg.Region)

	return nil
}
