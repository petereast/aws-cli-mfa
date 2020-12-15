package main

import (
	"awsmfacli/app"
	"awsmfacli/awsclient"
	"awsmfacli/config"
	"awsmfacli/credentials"
	"awsmfacli/filepaths"
	"awsmfacli/tokencode"
	"fmt"
)

func main() {
	filePaths := filepaths.UnixFilePaths{}
	err := app.App{
		AwsClient:        awsclient.AwsCliLoginClient{},
		CredentialWriter: credentials.FileCredentialWriter{FilePaths: filePaths},
		ConfigReader:     config.ConfigGetter{FilePaths: filePaths},
		TokenCodeGetter:  tokencode.ConsoleTokenCodeGetter{},
	}.Run()
	if err != nil {
		fmt.Println("Error! ", err)
	} else {
		fmt.Println("Done!")
	}
}
