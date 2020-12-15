# aws-cli-mfa

A simple tool to handle logging into AWS using their MFA command thing.

## Usage
- This application is self-contained, build and install it with `go build`.

## Config
This will expect you to have a json config file at `~/.aws-mfa.json`. This should take the following format:
```json
{
  "SecretAccessKey": "AWS SECRET KEY",
  "AccessKeyId": "AWS ACCESS KEY ID",
  "MfaDeviceArn": "MFA device ARN"
}
```

## Next Steps
- [x] Don't use a hardcoded file path
- [ ] Split the different things into different files
- [ ] Add unit tests
