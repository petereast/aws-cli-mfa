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

type AwsLoginClient interface {
  stsCall(string, string) (Credentials, error)
}

type CredentialWriter interface {
  writeCredentials(Credentials) error
}

type ConfigReader interface {
  getConfig(string) (Config, error)
}

type TokenCodeGetter interface {
  getTokenCode(config Config) string
}

type App struct {
  awsClient AwsLoginClient
  credentialWriter CredentialWriter
  configReader ConfigReader
  tokenCodeGetter TokenCodeGetter
}

func (app App) run() error {
  cfg, err := app.configReader.getConfig("/Users/petereast/.aws/credentials")
  if err != nil {
    return err
  }

  err = app.credentialWriter.writeCredentials(cfg.ToCreds())
  if err != nil {
    return err
  }

  newCreds, err := app.awsClient.stsCall(cfg.MfaDeviceArn, app.tokenCodeGetter.getTokenCode(cfg))
  if err != nil {
    return err
  }

  err = app.credentialWriter.writeCredentials(newCreds)
  if err != nil {
    return err
  }
  return nil
}

func (c Config) ToCreds() Credentials {
  return Credentials {
    SecretAccessKey: c.SecretAccessKey,
    AccessKeyId: c.AccessKeyId,
  }
}

// Cred writer
type FileCredentialWriter struct {}

func (_ FileCredentialWriter) writeCredentials(creds Credentials) (err error) {
  err = ioutil.WriteFile("/Users/petereast/.aws/credentials", []byte(ConfigEncoder("default", creds)), 0777)
  return
}

// Config Getter
type ConfigGetter struct {}

func (_ ConfigGetter) getConfig(string) (config Config, err error) {
	jsonBytes, err := ioutil.ReadFile("/Users/petereast/.aws-mfa.json")

	err = json.Unmarshal(jsonBytes, &config)

  return
}

// Done!
func (credentials Credentials) WriteCredentials() {
	ioutil.WriteFile("/Users/petereast/.aws/credentials", []byte(ConfigEncoder("default", credentials)), 0777)
}

func main() {
  err := App {
    awsClient: AwsCliLoginClient {},
    credentialWriter: FileCredentialWriter {},
    configReader: ConfigGetter {},
    tokenCodeGetter: ConsoleTokenCodeGetter{},
  }.run()
  if err != nil {
    fmt.Println("Error! ", err)
  } else {
  fmt.Println("Done aaaaaaaaaaa!")
}
}

type AwsCliLoginClient struct {}

func (_ AwsCliLoginClient) stsCall(deviceArn string, tokenCode string) (creds Credentials, err error) {
	// # aws sts get-session-token --serial-number arn:aws:iam::627518313974:mfa/peter.east@cyted.ai --token-code
	output, err := exec.Command("aws", "sts", "get-session-token", "--serial-number", deviceArn, "--token-code", tokenCode).Output()

	var awsResponse AwsResponse
	err = json.Unmarshal(output, &awsResponse)

	creds = awsResponse.Credentials
	return
}

type ConsoleTokenCodeGetter struct {}

func (_ ConsoleTokenCodeGetter) getTokenCode(config Config) (tokenCode string) {
	fmt.Printf("Enter token code (device: %s):\n|>", config.MfaDeviceArn)
	fmt.Scanf("%s", &tokenCode)

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
