package credentials

import (
  "awsmfacli/filepaths"
  "awsmfacli/encoder"
  "io/ioutil"
)

type Credentials struct {
	SecretAccessKey string
	AccessKeyId     string
	SessionToken    string `json:"SessionToken,omitempty"`
}

// Cred writer
type FileCredentialWriter struct {
	FilePaths filepaths.FilePaths
}

func (self FileCredentialWriter) WriteCredentials(creds Credentials) (err error) {
	path, err := self.FilePaths.GetAwsCredentialsLocation()
	err = ioutil.WriteFile(path, []byte(encoder.ConfigEncoder("default", creds)), 0777)
	return
}

type CredentialWriter interface {
	WriteCredentials(Credentials) error
}
