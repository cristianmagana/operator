package release

import (
	"context"
	"encoding/json"
	"os"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/service/ssm"
	"github.com/aws/aws-sdk-go-v2/service/ssm/types"
	"github.com/aws/aws-sdk-go/aws"
)

type SSMGetParameterImpl struct{}

func (dt SSMGetParameterImpl) GetParameter(ctx context.Context, params *ssm.GetParameterInput,
	optFns ...func(*ssm.Options)) (*ssm.GetParameterOutput, error) {

	parameter := &types.Parameter{Value: aws.String("aws-docs-example-parameter-value")}

	output := &ssm.GetParameterOutput{
		Parameter: parameter,
	}

	return output, nil
}

func TestGetArchitectureRelease(t *testing.T) {

	amd64 := "amd64"
	arm64 := "arm64"

	ver := "1.20"

	amd64Release := getArchitectureRelease(&ver, &amd64)
	arm64Release := getArchitectureRelease(&ver, &arm64)

	if amd64Release != "/aws/service/eks/optimized-ami/1.20/amazon-linux-2/recommended/release_version" {
		t.Errorf("amd64Release = %s; want /aws/service/eks/optimized-ami/1.20/amazon-linux-2/recommended/release_version", amd64Release)
	}

	if arm64Release != "/aws/service/eks/optimized-ami/1.20/amazon-linux-2-arm64/recommended/release_version" {
		t.Errorf("amd64Release = %s; want /aws/service/eks/optimized-ami/1.20/amazon-linux-2-arm64/recommended/release_version", arm64Release)
	}
}

type Config struct {
	ParameterName string `json:"ParameterName"`
	Value         string `json:"Value"`
}

var globalConfig []Config

var configFileName = "config.json"

func populateConfig(t *testing.T) error {

	content, err := os.ReadFile(configFileName)
	if err != nil {
		return err
	}

	text := string(content)

	err = json.Unmarshal([]byte(text), &globalConfig)
	if err != nil {
		return err
	}

	return nil
}

func TestFindParameter(t *testing.T) {

	thisTime := time.Now()
	nowString := thisTime.Format("2006-01-02 15:04:05 Monday")
	t.Log("Starting unit test at " + nowString)

	err := populateConfig(t)
	if err != nil {
		t.Errorf("populateConfig failed: %s", err)
	}

	api := &SSMGetParameterImpl{}

	input := &ssm.GetParameterInput{
		Name: &globalConfig[0].ParameterName,
	}

	resp, err := findParameter(context.Background(), api, input)
	if err != nil {
		t.Errorf("findParameter failed: %s", err)
		return
	}
	t.Log("Parameter value: " + *resp.Parameter.Value)
}

func TestGetRelease(t *testing.T) {

	thisTime := time.Now()
	nowString := thisTime.Format("2006-01-02 15:04:05 Monday")
	t.Log("Starting unit test at " + nowString)

	var release *Release

	release, err := release.GetRelease()
	if err != nil {
		t.Errorf("GetRelease failed: %s", err)
		return
	}

	t.Log("Release value: " + release.ReleaseValue)
}
