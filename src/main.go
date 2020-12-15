package main

import (
  "fmt"
  "awsmfacli/src/app"
  "awsmfacli/src/config"
  "awsmfacli/src/filepaths"
  "awsmfacli/src/credentials"
  "awsmfacli/src/awsclient"
  "awsmfacli/src/tokencode"
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
