package filepaths

import (
	"fmt"
	"os/user"
)

type UnixFilePaths struct{}

func (_ UnixFilePaths) getHomedir() (homeDir string, err error) {
	usr, err := user.Current()

	homeDir = usr.HomeDir
	return
}

func (self UnixFilePaths) GetConfigLocation() (location string, err error) {
	homeDir, err := self.getHomedir()

	location = fmt.Sprintf("%s/.aws-mfa.json", homeDir)
	return
}

func (self UnixFilePaths) GetAwsCredentialsLocation() (location string, err error) {
	homeDir, err := self.getHomedir()

	location = fmt.Sprintf("%s/.aws/credentials", homeDir)
	return
}

type FilePaths interface {
	GetConfigLocation() (string, error)
	GetAwsCredentialsLocation() (string, error)
}
