package main

import (
  "fmt"
  "awsmfacli/src/config"
  "awsmfacli/src/filepaths"
  "awsmfacli/src/credentials"
  "awsmfacli/src/awsclient"
  "awsmfacli/src/tokencode"
)

func ToCreds(c config.Config) credentials.Credentials {
	return credentials.Credentials{
		SecretAccessKey: c.SecretAccessKey,
		AccessKeyId:     c.AccessKeyId,
	}
}



type App struct {
	awsClient        awsclient.AwsLoginClient
	credentialWriter credentials.CredentialWriter
	configReader     config.ConfigReader
	tokenCodeGetter  tokencode.TokenCodeGetter
}

func (app App) run() error {
	cfg, err := app.configReader.GetConfig()
	if err != nil {
		return err
	}

	err = app.credentialWriter.WriteCredentials(ToCreds(cfg))
	if err != nil {
		return err
	}

	newCreds, err := app.awsClient.StsCall(cfg.MfaDeviceArn, app.tokenCodeGetter.GetTokenCode(cfg))
	if err != nil {
		return err
	}

	err = app.credentialWriter.WriteCredentials(newCreds)
	if err != nil {
		return err
	}
	return nil
}

func main() {
	filePaths := filepaths.UnixFilePaths{}
	err := App{
		awsClient:        awsclient.AwsCliLoginClient{},
		credentialWriter: credentials.FileCredentialWriter{FilePaths: filePaths},
		configReader:     config.ConfigGetter{FilePaths: filePaths},
		tokenCodeGetter:  tokencode.ConsoleTokenCodeGetter{},
	}.run()
	if err != nil {
		fmt.Println("Error! ", err)
	} else {
		fmt.Println("Done!")
	}
}
