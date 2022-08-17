package internal

import (
	"fmt"
	"os"
)

type AcConfig struct {
	Dir        string
	CookiePath string
	Endpoint   string
}

func NewAcConfig() *AcConfig {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	acConfigDir := home + "/.cpcli/ac/"

	return &AcConfig{
		Dir:        acConfigDir,
		CookiePath: acConfigDir + "cookie",
		Endpoint:   "https://atcoder.jp",
	}

}
