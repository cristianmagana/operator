package release

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ssm"
)

type Release struct {
	// Name of the release
	Name string
	// Version of the release
	EksVersion string
	// Architecture of the release
	Architecture string
	// Value of the release
	Value string
}

type SSMGetParameterAPI interface {
	GetParameter(ctx context.Context,
		params *ssm.GetParameterInput,
		optFns ...func(*ssm.Options)) (*ssm.GetParameterOutput, error)
}

func findParameter(c context.Context, api SSMGetParameterAPI, input *ssm.GetParameterInput) (*ssm.GetParameterOutput, error) {
	return api.GetParameter(c, input)
}

func getArchitectureRelease(ver *string, arch *string) string {

	parameterName := ""

	switch *arch {
	case "arm64":
		parameterName = fmt.Sprintf("/aws/service/eks/optimized-ami/%s/amazon-linux-2-arm64/recommended/release_version", *ver)
	default:
		parameterName = fmt.Sprintf("/aws/service/eks/optimized-ami/%s/amazon-linux-2/recommended/release_version", *ver)
	}
	return parameterName
}

func (r *Release) GetRelease() (*Release, error) {

	currentEksVersion := os.Getenv("EKS_VERSION")
	slog.Debug("Current EKS version: " + currentEksVersion)

	nodeArchitecture := os.Getenv("NODE_ARCHITECTURE")
	slog.Debug("Current EKS version: " + currentEksVersion)

	parameterName := getArchitectureRelease(&currentEksVersion, &nodeArchitecture)

	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		panic("configuration error, " + err.Error())
	}

	client := ssm.NewFromConfig(cfg)

	input := &ssm.GetParameterInput{
		Name: &parameterName,
	}

	results, err := findParameter(context.TODO(), client, input)
	if err != nil {
		slog.Error(err.Error())
		return nil, err
	}

	release := Release{
		Name:         fmt.Sprintf("LATEST RELEASE AS OF %s", time.Now().Format(time.RFC850)),
		EksVersion:   currentEksVersion,
		Architecture: nodeArchitecture,
		Value:        *results.Parameter.Value,
	}

	slog.Info(*results.Parameter.Value)

	return &release, nil
}
