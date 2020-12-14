package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os/exec"
	"reflect"
	"strings"
)

type Config struct {
	SecretAccessKey string
	AccessKeyId     string
	MfaDeviceArn    string
}

type Credentials struct {
	SecretAccessKey string
	AccessKeyId     string
	SessionToken    string `json:"SessionToken,omitempty"`
}

type AwsResponse struct {
	Credentials Credentials
}

func (credentials Credentials) WriteCredentials() {
	ioutil.WriteFile("/Users/petereast/.aws/credentials", []byte(ConfigEncoder("default", credentials)), 0777)
}

func main() {
	jsonBytes, err := ioutil.ReadFile("/Users/petereast/.aws-mfa.json")

	if err != nil {
		fmt.Print("Can't open file")
		return
	}

	var config Config
	err = json.Unmarshal(jsonBytes, &config)

	if err != nil {
		fmt.Print("Can't parse config file")
		return
	}

	var initialCreds Credentials
	err = json.Unmarshal(jsonBytes, &initialCreds)
	initialCreds.WriteCredentials()

	fmt.Printf("Enter token code (device: %s):\n|>", config.MfaDeviceArn)
	var tokenCode string
	fmt.Scanf("%s", &tokenCode)

	creds, err := StsCall(config.MfaDeviceArn, tokenCode)
	if err != nil {
		fmt.Print("Error!", err)
		return
	}

	creds.WriteCredentials()
	fmt.Println("Done!")
}

// TODO: DI the shit out of this
func StsCall(deviceArn string, tokenCode string) (creds Credentials, err error) {
	// # aws sts get-session-token --serial-number arn:aws:iam::627518313974:mfa/peter.east@cyted.ai --token-code
	output, err := exec.Command("aws", "sts", "get-session-token", "--serial-number", deviceArn, "--token-code", tokenCode).Output()

	var awsResponse AwsResponse
	err = json.Unmarshal(output, &awsResponse)

	creds = awsResponse.Credentials
	return
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
