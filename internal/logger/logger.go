package logger

import (
	"fmt"

	ctoai "github.com/cto-ai/sdk-go"
)

func LogSlackError(ux *ctoai.Ux, err error) {
	slackPrint := fmt.Sprintf("%v", err)
	err = ux.Print(slackPrint)
	if err != nil {
		fmt.Println(err)
	}
}

func LogSlack(ux *ctoai.Ux, str string) {
	err := ux.Print(str)
	if err != nil {
		fmt.Println(err)
	}
}
