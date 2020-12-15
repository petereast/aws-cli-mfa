package app

import (
  "awsmfacli/config"
  "awsmfacli/credentials"
  "awsmfacli/awsclient"
  "awsmfacli/tokencode"
)

func ToCreds(c config.Config) credentials.Credentials {
	return credentials.Credentials{
		SecretAccessKey: c.SecretAccessKey,
		AccessKeyId:     c.AccessKeyId,
	}
}



type App struct {
	AwsClient        awsclient.AwsLoginClient
	CredentialWriter credentials.CredentialWriter
	ConfigReader     config.ConfigReader
	TokenCodeGetter  tokencode.TokenCodeGetter
}

func (app App) Run() error {
	cfg, err := app.ConfigReader.GetConfig()
	if err != nil {
		return err
	}

	err = app.CredentialWriter.WriteCredentials(ToCreds(cfg))
	if err != nil {
		return err
	}

	newCreds, err := app.AwsClient.StsCall(cfg.MfaDeviceArn, app.TokenCodeGetter.GetTokenCode(cfg))
	if err != nil {
		return err
	}

	err = app.CredentialWriter.WriteCredentials(newCreds)
	if err != nil {
		return err
	}
	return nil
}

