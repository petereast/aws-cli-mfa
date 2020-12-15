# aws-cli-mfa

A simple tool to handle logging into AWS using their MFA command thing.

## Usage
- This application is self-contained, it doesn't need anything at runtime apart from the config file detailed below.
- You can build the binary with the following command: `cd src/ && go build main.go`, which will output an executable called `awsmfacli`. Put this binary somewhere in your `PATH` and you'll  be able to use it.

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
- [x] Split the different things into different files
- [ ] Add unit tests
