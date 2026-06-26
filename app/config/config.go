package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

const (
	DefaultRoleName = "bedrock-direct-call-sonnet-4-6"
	DefaultRegion   = "eu-central-1"
)

// Settings holds AWS role assumption and region configuration.
type Settings struct {
	Role           string
	RoleARN        string
	AccountID      string
	Region         string
	SessionName    string
	DurationSeconds int32
}

// Load reads config.yaml (optional) and environment overrides.
func Load() (*Settings, error) {
	v := viper.New()

	v.SetDefault("role", DefaultRoleName)
	v.SetDefault("aws.region", DefaultRegion)
	v.SetDefault("aws.session_name", "zed-cli")
	v.SetDefault("aws.duration_seconds", 3600)

	v.SetConfigName("config-zbe")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")

	if home, err := os.UserHomeDir(); err == nil {
		v.AddConfigPath(filepath.Join(home, ".config"))
	}

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, fmt.Errorf("read config: %w", err)
		}
	}

	s := &Settings{
		Role:            v.GetString("role"),
		RoleARN:         strings.TrimSpace(v.GetString("role_arn")),
		Region:          strings.TrimSpace(v.GetString("aws.region")),
		SessionName:     v.GetString("aws.session_name"),
		DurationSeconds: int32(v.GetInt("aws.duration_seconds")),
	}

	if s.Region == "" {
		return nil, fmt.Errorf("aws.region is required")
	}

	return s, nil
}

// ResolveRoleARN returns a full IAM role ARN from config.
func (s *Settings) ResolveRoleARN() (string, error) {
	if s.RoleARN != "" {
		return s.RoleARN, nil
	}
	role := strings.TrimSpace(s.Role)
	if strings.HasPrefix(role, "arn:") {
		return role, nil
	}
	if s.AccountID == "" {
		return "", fmt.Errorf("account ID is required to build role ARN; it should be resolved via STS GetCallerIdentity before calling ResolveRoleARN")
	}
	return fmt.Sprintf("arn:aws:iam::%s:role/%s", s.AccountID, role), nil
}
