package logger

import (
	"fmt"

	"git.cto.ai/sdk-go/pkg/sdk"
)

func LogSlackError(sdk *sdk.Sdk, err error) {
	slackPrint := fmt.Sprintf("%v", err)
	err = sdk.Print(slackPrint)
	if err != nil {
		fmt.Println(err)
	}
}

func LogSlack(sdk *sdk.Sdk, str string) {
	err := sdk.Print(str)
	if err != nil {
		fmt.Println(err)
	}
}
