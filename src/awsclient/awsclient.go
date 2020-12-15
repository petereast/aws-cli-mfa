package awsclient

import (
  "awsmfacli/src/credentials"
  "os/exec"
  "encoding/json"
)


type AwsResponse struct {
	Credentials credentials.Credentials
}

type AwsLoginClient interface {
	StsCall(string, string) (credentials.Credentials, error)
}

type AwsCliLoginClient struct{}

func (_ AwsCliLoginClient) StsCall(deviceArn string, tokenCode string) (creds credentials.Credentials, err error) {
	// # aws sts get-session-token --serial-number arn:aws:iam::627518313974:mfa/peter.east@cyted.ai --token-code
	output, err := exec.Command("aws", "sts", "get-session-token", "--serial-number", deviceArn, "--token-code", tokenCode).Output()

	var awsResponse AwsResponse
	err = json.Unmarshal(output, &awsResponse)

	creds = awsResponse.Credentials
	return
}

