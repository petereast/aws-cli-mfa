package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"reflect"
  "os/exec"
	"strings"
)

type Config struct {
	SecretAccessKey string
	AccessKeyId string
	MfaDeviceArn string
}

type Credentials struct {
	SecretAccessKey string
	AccessKeyId string
	SessionToken string `json:"SessionToken,omitempty"`
}

type AwsResponse struct {
    Credentials Credentials
}

func main() {
	fmt.Println("Hello world!")
	// This will do a few things:
	// Read the config file from the root directory, to get the aws credentials and the ARN of the mfa device
	jsonBytes, err := ioutil.ReadFile("/Users/petereast/.aws-mfa.json")

	if err != nil {
		fmt.Print("Can't open file")
		return
	}

	var config Config
	err = json.Unmarshal(jsonBytes, &config)

	if err != nil {
		fmt.Print("Can't parse file")
		return
	}

  var initialCreds Credentials
  err = json.Unmarshal(jsonBytes, &initialCreds)
  WriteCredentials(initialCreds)

  fmt.Printf("Enter token code (device: %s):\n|>", config.MfaDeviceArn)
  var tokenCode string
  fmt.Scanf("%s", &tokenCode)

  creds, _ := StsCall(config.MfaDeviceArn, tokenCode)
	// Write this config to .aws/credentials
	// Call `aws sts get-access-code --mfa-device ARN --token-code
	// Write the new data to .aws/credentials

	WriteCredentials( creds)

}

func StsCall(deviceArn string, tokenCode string) (Credentials, *string) {
// # aws sts get-session-token --serial-number arn:aws:iam::627518313974:mfa/peter.east@cyted.ai --token-code
  output, err := exec.Command("aws", "sts", "get-session-token", "--serial-number", deviceArn, "--token-code", tokenCode).Output()

  if err != nil {
    panic("Can't authenticate OTP")
  }

  var awsResponse AwsResponse
  err = json.Unmarshal(output, &awsResponse)

  if err != nil {
    panic("Can't parse response from aws")
  }

  return awsResponse.Credentials, nil
}

func WriteCredentials(credentials Credentials) (bool, string) {
	// We need to encode the credentials into the right format
	/*
	  [default]
	  key = value
	*/

  ioutil.WriteFile("/Users/petereast/.aws/credentials", []byte(ConfigEncoder("default", credentials)), 0777)

	return false, string("Something went wrong")
}

func ConfigEncoder(title string, config interface{}) string {
	// This will read the fields of the interface and create a config file

	t := reflect.TypeOf(config)
	values := reflect.ValueOf(config)

	maxIndex := t.NumField()
	output := fmt.Sprintf("[%s]\n", title)

	for i := 0; i < maxIndex; i++ {
		f := t.Field(i)

		value := values.FieldByName(f.Name).String()
		name := ToSnakeCase(f.Name)

    if len(value) != 0 {
      output += fmt.Sprintf("aws%s = %s\n", name, value)
  }
	}

	return output
}

func ToSnakeCase(input string) string {
	output := ""
	for _, v := range input {
		val := string(v)
		lowerV := strings.ToLower(val)
		if lowerV != val {
			output += "_"
		}
		output += lowerV
	}

	return output
}
