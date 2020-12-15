package tokencode

import (
	"awsmfacli/config"
	"fmt"
)

type TokenCodeGetter interface {
	GetTokenCode(config config.Config) string
}

// Config Getter

type ConsoleTokenCodeGetter struct{}

func (_ ConsoleTokenCodeGetter) GetTokenCode(config config.Config) (tokenCode string) {
	fmt.Printf("Enter token code (device: %s):\n|>", config.MfaDeviceArn)
	fmt.Scanf("%s", &tokenCode)

	return
}
