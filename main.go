package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os/exec"
	"os/user"
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
	getConfig() (Config, error)
}

type TokenCodeGetter interface {
	getTokenCode(config Config) string
}

type FilePaths interface {
	getConfigLocation() (string, error)
	getAwsCredentialsLocation() (string, error)
}

type App struct {
	awsClient        AwsLoginClient
	credentialWriter CredentialWriter
	configReader     ConfigReader
	tokenCodeGetter  TokenCodeGetter
}

func (app App) run() error {
	cfg, err := app.configReader.getConfig()
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
	return Credentials{
		SecretAccessKey: c.SecretAccessKey,
		AccessKeyId:     c.AccessKeyId,
	}
}

// Cred writer
type FileCredentialWriter struct {
	filePaths FilePaths
}

func (self FileCredentialWriter) writeCredentials(creds Credentials) (err error) {
	path, err := self.filePaths.getAwsCredentialsLocation()
	err = ioutil.WriteFile(path, []byte(ConfigEncoder("default", creds)), 0777)
	return
}

// Config Getter
type ConfigGetter struct {
	filePaths FilePaths
}

func (self ConfigGetter) getConfig() (config Config, err error) {
	configPath, err := self.filePaths.getConfigLocation()
	jsonBytes, err := ioutil.ReadFile(configPath)

	err = json.Unmarshal(jsonBytes, &config)

	return
}

type AwsCliLoginClient struct{}

func (_ AwsCliLoginClient) stsCall(deviceArn string, tokenCode string) (creds Credentials, err error) {
	// # aws sts get-session-token --serial-number arn:aws:iam::627518313974:mfa/peter.east@cyted.ai --token-code
	output, err := exec.Command("aws", "sts", "get-session-token", "--serial-number", deviceArn, "--token-code", tokenCode).Output()

	var awsResponse AwsResponse
	err = json.Unmarshal(output, &awsResponse)

	creds = awsResponse.Credentials
	return
}

type ConsoleTokenCodeGetter struct{}

func (_ ConsoleTokenCodeGetter) getTokenCode(config Config) (tokenCode string) {
	fmt.Printf("Enter token code (device: %s):\n|>", config.MfaDeviceArn)
	fmt.Scanf("%s", &tokenCode)

	return
}

type UnixFilePaths struct{}

func (_ UnixFilePaths) getHomedir() (homeDir string, err error) {
	usr, err := user.Current()

	homeDir = usr.HomeDir
	return
}

func (self UnixFilePaths) getConfigLocation() (location string, err error) {
	homeDir, err := self.getHomedir()

	location = fmt.Sprintf("%s/.aws-mfa.json", homeDir)
	return
}

func (self UnixFilePaths) getAwsCredentialsLocation() (location string, err error) {
	homeDir, err := self.getHomedir()

	location = fmt.Sprintf("%s/.aws/credentials", homeDir)
	return
}

// Utilities

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

func main() {
	filePaths := UnixFilePaths{}
	err := App{
		awsClient:        AwsCliLoginClient{},
		credentialWriter: FileCredentialWriter{filePaths: filePaths},
		configReader:     ConfigGetter{filePaths: filePaths},
		tokenCodeGetter:  ConsoleTokenCodeGetter{},
	}.run()
	if err != nil {
		fmt.Println("Error! ", err)
	} else {
		fmt.Println("Done aaaaaaaaaaa!")
	}
}
