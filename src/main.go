package main

import (
  "fmt"
  "awsmfacli/app"
  "awsmfacli/config"
  "awsmfacli/filepaths"
  "awsmfacli/credentials"
  "awsmfacli/awsclient"
  "awsmfacli/tokencode"
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
