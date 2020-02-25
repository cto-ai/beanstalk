package sdk

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"git.cto.ai/sdk-go/pkg/request"
)

const promptServiceRequest = "PromptServiceRequest"

// PromptServiceRequest private method generates a request to the SDK daemon from
// a PromptRequest type struct and fills in the appropriate request
// details for a prompt request.
// The request is subsequently executed and the response body is handled through unmarshalling. The response is known to contain a file path which is then
// opened and the contents of the file is parsed.
// Both the PromptOutput type and error are returned.
//
// The request method and struct are exported so that an advanced user
// may choose to use the direct request as opposed to using the
// constructor/activator methods defined below.
//
// Example:
//
// See the public methods:
// sdk.PromptInput()
// sdk.PromptNumber()
// sdk.PromptSecret()
// sdk.PromptPassword()
// sdk.PromptConfirm()
// sdk.PromptList()
// sdk.PromptAutocomplete()
// sdk.PromptCheckbox()
// sdk.PromptEditor()
// sdk.PromptDatetime()
type PromptRequest struct {
	Name       string   `json:"name"`
	PromptType string   `json:"type"`
	PromptMsg  string   `json:"message"`
	Choices    []string `json:"choices"`
	AllowEmpty bool     `json:"allowEmpty"`

	// Flag field is the alias to the Name field
	// when the op is run from the command line
	// the name or flag can be used
	// as supplying the fields beforehand
	// to skip the prompting
	//
	// e.g. if Name is set to "test" and Flag is set to "t"
	// ops run goOp --test promptresponse
	// is equivalent to
	// ops run goOp -t promptresponse
	// both of which will supply the answer to this prompt
	//
	Flag    string      `json:"flag"`
	Confirm bool        `json:"confirm"`
	Min     interface{} `json:"minimum"`
	Max     interface{} `json:"maximum"`
	Variant string      `json:"variant"`
	Default interface{} `json:"default"`
}

type PromptResponse struct {
	ReplyFilename string  `name:"replyFilename"`
}

type PromptOutput map[string]interface{}

func (s *Sdk) PromptRequest(prompt PromptRequest) (PromptOutput, error) {
	responseOutput := &PromptResponse{}
	output := PromptOutput{}

	bytes, err := json.Marshal(prompt)
	if err != nil {
		return output, fmt.Errorf("Failed to marshal request body in service %s: %s", promptServiceRequest, err)
	}

	op := &request.Operation{
		Name:       promptServiceRequest,
		HttpMethod: "POST",
		HttpPath:   "/prompt",
		Payload:    bytes,
	}

	req := s.NewRequest(op, responseOutput)
	err = req.Exec()
	if err != nil {
		return nil, fmt.Errorf("Error in sending request to sdk daemon: %s", err)
	}

	err = output.UnmarshalFile(responseOutput.ReplyFilename)
	if err != nil {
		return nil, fmt.Errorf("Error in deserializing response from sdk daemon file: %s", err)
	}

	return output, nil
}

const promptInput = "PromptInput"

// PromptInput is a constructor method that calls the PromptRequest method.
// This is the public facing API to enable developers to present an input
// prompt to the interface (i.e. terminal or slack).
//
// The name is the name of the prompt, which will be accessible in the future
// as metrics to the developer. It will also be used for cli based prompts.
//
// The msg is the prompt presented to the user.
//
// The defaultValue is the response output of the prompt if the user of the op
// returns with an empty response.
//
// The allowEmpty switch is used to indicate if the user can provide an empty response
// and thus the default value being used.
//
// The method returns the user's response as string.
//
// Example:
//
// s := sdk.New()
// resp, err := s.PromptInput("opinion" ,"What do you think of Go?", "good", "i", true) // user responds with "good"
// if err != nil {
//     panic(err)
// }
//
// fmt.Println(resp)
//
// Output:
// good
func (s *Sdk) PromptInput(name, msg, defaultValue, flag string, allowEmpty bool) (string, error) {
	r := PromptRequest{
		Name:       name,
		PromptType: "input",
		PromptMsg:  msg,
		Flag:       flag,
		Default:    defaultValue,
		AllowEmpty: allowEmpty,
	}

	output, err := s.PromptRequest(r)
	if err != nil {
		return "", fmt.Errorf("Error in %s: %s", promptInput, err)
	}

	return output[name].(string), nil
}

const promptNumber = "PromptNumber"

// PromptNumber is a constructor method that calls the PromptRequest method.
// This is the public facing API to enable developers to present an input
// prompt for numbers specificaly to the interface (i.e. terminal or slack).
//
// The name is the name of the prompt, which will be accessible in the future
// as metrics to the developer. It will also be used for cli based prompts.
//
// The msg is the prompt presented to the user.
//
// The defaultValue is the response output of the prompt if the user of the op
// returns with an empty response.
//
// The min is the minimum number that a user can input.
//
// The max is the maximum number that a user can input.
//
// The method returns the user's response as float64.
//
// Example:
//
// s := sdk.New()
// resp, err := s.PromptNumber("count", "How many hoorays do you want?", "n", 10, 0, 10) // user responds with 7
// if err != nil {
//     panic(err)
// }
//
// fmt.Println(resp)
//
// Output:
// 7
func (s *Sdk) PromptNumber(name, msg, flag string, defaultValue, min, max int) (float64, error) {
	r := PromptRequest{
		Name:       name,
		PromptType: "number",
		PromptMsg:  msg,
		Flag:       flag,
		Default:    defaultValue,
		Min:        min,
		Max:        max,
	}

	output, err := s.PromptRequest(r)
	if err != nil {
		return 0, fmt.Errorf("Error in %s: %s", promptNumber, err)
	}

	return output[name].(float64), nil
}

const promptSecret = "PromptSecret"

// PromptSecret is a constructor method that calls the PromptRequest method.
// This is the public facing API to enable developers to present an input
// prompt for secrets in the interface (i.e. terminal or slack).
//
// The secret can be entered by the user as clear plain text or if the user
// already has a team secrets store configured, the value from their store
// can also be entered. The SDK daemon will be able to get the secret from
// the store and can be used by the developer. The secret store's value name
// will also be in clear plain text on the interface.
//
// The name is the name of the prompt, which will be accessible in the future
// as metrics to the developer. It will also be used for cli based prompts.
//
// The msg is the prompt presented to the user.
//
// The method returns the user's response as string.
//
// Example:
//
// s := sdk.New()
// resp, err := s.PromptSecret("sshKey", "What SSH private key do you want to use?", "s") // user responds with 1234asdf
// if err != nil {
//     panic(err)
// }
//
// fmt.Println(resp)
//
// Output:
// 1234asdf
func (s *Sdk) PromptSecret(name, msg, flag string) (string, error) {
	r := PromptRequest{
		Name:       name,
		PromptType: "secret",
		PromptMsg:  msg,
		Flag:       flag,
	}

	output, err := s.PromptRequest(r)
	if err != nil {
		return "", fmt.Errorf("Error in %s: %s", promptSecret, err)
	}

	return output[name].(string), nil
}

const promptPassword = "PromptPassword"

// PromptPassword is a constructor method that calls the PromptRequest method.
// This is the public facing API to enable developers to present an input
// prompt for passwords in the interface (i.e. terminal or slack).
//
// The password can be entered by the user as hidden text or if the user
// already has a team secrets store configured, the value from their store
// can also be entered. The SDK daemon will be able to get the password from
// the store and can be used by the developer. The secret store's value name
// will also be in hidden text on the interface.
//
// The name is the name of the prompt, which will be accessible in the future
// as metrics to the developer. It will also be used for cli based prompts.
//
// The msg is the prompt presented to the user.
//
// The confirmation is a switch for asking the user a confirmation prompt
// for confirming the password.
//
// The method returns the user's response as string.
//
// Example:
//
// s := sdk.New()
// resp, err := s.PromptPassword("password", "What new password would you like to use?", "p") // user responds with 1234asdf
// if err != nil {
//     panic(err)
// }
//
// fmt.Println(resp)
//
// Output:
// 1234asdf
func (s *Sdk) PromptPassword(name, msg, flag string, confirmation bool) (string, error) {
	r := PromptRequest{
		Name:       name,
		PromptType: "password",
		PromptMsg:  msg,
		Confirm:    confirmation,
		Flag:       flag,
	}

	output, err := s.PromptRequest(r)
	if err != nil {
		return "", fmt.Errorf("Error in %s: %s", promptPassword, err)
	}

	return output[name].(string), nil
}

const promptConfirm = "PromptConfirm"

// PromptConfirm is a constructor method that calls the PromptRequest method.
// This is the public facing API to enable developers to present a yes/no answer
// to the user in the interface (i.e. terminal or slack).
//
// The name is the name of the prompt, which will be accessible in the future
// as metrics to the developer. It will also be used for cli based prompts.
//
// The msg is the prompt presented to the user.
//
// The defaultValue is the value that will be returned if the user presses enter.
//
// The method returns the user's response as bool type.
//
// Example:
//
// s := sdk.New()
// resp, err := s.PromptConfirm("verbose", "Do you want to run in verbose mode?", "c", false) // user responds with y
// if err != nil {
//     panic(err)
// }
//
// fmt.Println(resp)
//
// Output:
// true
func (s *Sdk) PromptConfirm(name, msg, flag string, defaultValue bool) (bool, error) {
	r := PromptRequest{
		Name:       name,
		PromptType: "confirm",
		Flag:       flag,
		PromptMsg:  msg,
		Default:    defaultValue,
	}

	output, err := s.PromptRequest(r)
	if err != nil {
		return false, fmt.Errorf("Error in %s: %s", promptConfirm, err)
	}

	return output[name].(bool), nil
}

const promptList = "PromptList"

// PromptList is a constructor method that calls the PromptRequest method.
// This is the public facing API to enable developers to present a list of
// options to the user select one item from in the interface
// (i.e. terminal or slack).
//
// The choices is the list of string options that the user can select from
// using up and down arrows.
//
// The name is the name of the prompt, which will be accessible in the future
// as metrics to the developer. It will also be used for cli based prompts.
//
// The msg is the prompt presented to the user.
//
// The defaultValue is the value that will be returned if the user presses enter.
//
// The method returns the user's response as a string.
//
// Example:
//
// s := sdk.New()
// resp, err := s.PromptList([]string{"AWS", "Google Cloud", "Azure"}, "platform", "What cloud platform would you like to deploy to?", "AWS", "L") // user selects Azure
// if err != nil {
//     panic(err)
// }
//
// fmt.Println(resp)
//
// Output:
// Azure
func (s *Sdk) PromptList(choices []string, name, msg, defaultValue, flag string) (string, error) {
	r := PromptRequest{
		Name:       name,
		PromptType: "list",
		PromptMsg:  msg,
		Choices:    choices,
		Default:    defaultValue,
		Flag:       flag,
	}

	output, err := s.PromptRequest(r)
	if err != nil {
		return "", fmt.Errorf("Error in %s: %s", promptList, err)
	}

	return output[name].(string), nil
}

const promptAutocomplete = "PromptAutocomplete"

// PromptAutocomplete is a constructor method that calls the PromptRequest method.
// This is the public facing API to enable developers to present a list of
// options to the user select one item from in the interface
// (i.e. terminal or slack).
//
// The choices is the list of string options that the user can select from
// using up and down arrows. The user can also start typing to filter the list
// of options.
//
// The name is the name of the prompt, which will be accessible in the future
// as metrics to the developer. It will also be used for cli based prompts.
//
// The msg is the prompt presented to the user.
//
// The defaultValue is the value that will be returned if the user presses enter.
//
// The method returns the user's response as a string.
//
// Example:
//
// s := sdk.New()
// resp, err := s.PromptAutocomplete([]string{"Alberta", "British Columbia", "Manitoba", "New Brunswick","Newfoundland and Labrador", "Nova Scotia", "Ontario", "Prince Edward Island", "Quebec", "Saskatchewan"}, "province", "To which province do you want to send your letter?", "Alberta", "a") // user selects Ontario
// if err != nil {
//     panic(err)
// }
//
// fmt.Println(resp)
//
// Output:
// Ontario
func (s *Sdk) PromptAutocomplete(choices []string, name, msg, defaultValue, flag string) (string, error) {
	r := PromptRequest{
		Name:       name,
		PromptType: "autocomplete",
		PromptMsg:  msg,
		Choices:    choices,
		Default:    defaultValue,
		Flag:       flag,
	}

	output, err := s.PromptRequest(r)
	if err != nil {
		return "", fmt.Errorf("Error in %s: %s", promptAutocomplete, err)
	}

	return output[name].(string), nil
}

const promptCheckbox = "PromptCheckbox"

// PromptCheckbox is a constructor method that calls the PromptRequest method.
// This is the public facing API to enable developers to present a list of
// options to the user select multiple items from in the interface
// (i.e. terminal or slack).
//
// The choices is the list of string options that the user can select from
// using up and down arrows.
//
// The name is the name of the prompt, which will be accessible in the future
// as metrics to the developer. It will also be used for cli based prompts.
//
// The msg is the prompt presented to the user.
//
// The method returns the user's response as a slice of strings.
//
// Example:
//
// s := sdk.New()
// resp, err := s.PromptCheckbox([]string{"Lua", "Node.js", "Perl", "Python 2", "Python 3", "Raku", "Ruby"}, "tools", "Which interpreters would you like to have included in your OS image?", "Lua", "C") // user selects Lua
// if err != nil {
//     panic(err)
// }
//
// fmt.Println(resp)
//
// Output:
// Lua
func (s *Sdk) PromptCheckbox(choices []string, name, msg, flag string) ([]string, error) {
	r := PromptRequest{
		Name:       name,
		PromptType: "checkbox",
		PromptMsg:  msg,
		Choices:    choices,
		Flag:       flag,
	}

	output, err := s.PromptRequest(r)
	if err != nil {
		return []string{}, fmt.Errorf("Error in %s: %s", promptCheckbox, err)
	}

	outputInferface := output[name].([]interface{})
	out := make([]string, 0, len(outputInferface))

	for _, v := range outputInferface {
		out = append(out, v.(string))
	}

	return out, nil
}

const promptEditor = "PromptEditor"

// PromptEditor is a constructor method that calls the PromptRequest method.
// This is the public facing API to enable developers to present
// a multi-line response from the user. If used in a terminal interface,
// the nano editor will be presented.
//
// The name is the name of the prompt, which will be accessible in the future
// as metrics to the developer. It will also be used for cli based prompts.
//
// The msg is the prompt presented to the user.
//
// The defaultValue is a template that is shown in the nano editor.
//
// The method returns the user's response as a string.
//
// Example:
//
// s := sdk.New()
//
// template := `Features:
//
// Fixes:
//
// Chores:
//`
//
// resp, err := s.PromptEditor("notes", "Please enter your release notes", template, "e")
// if err != nil {
//     panic(err)
// }
//
// Output:
// [Nano will be brought up with the template in the editor]
func (s *Sdk) PromptEditor(name, msg, defaultValue, flag string) (string, error) {
	r := PromptRequest{
		Name:       name,
		PromptType: "editor",
		PromptMsg:  msg,
		Default:    defaultValue,
		Flag:       flag,
	}

	output, err := s.PromptRequest(r)
	if err != nil {
		return "", fmt.Errorf("Error in %s: %s", promptEditor, err)
	}

	return output[name].(string), nil
}

const promptDatetime = "PromptDatetime"

// PromptDatetime is a constructor method that calls the PromptRequest method.
// This is the public facing API to enable developers to present a date in
// 2006-01-02 15:04:05 format to the interface (i.e. terminal or slack).
//
// The user can move between the YYYY, mm, dd, HH, MM, SS using the left
// and right arrows. And increment or decrement the values using the up
// and down arrows.
//
// The name is the name of the prompt, which will be accessible in the future
// as metrics to the developer. It will also be used for cli based prompts.
//
// The msg is the prompt presented to the user.
//
// The variant is an option between "date", "time", or "datetime".
//
// The minDate is the earliest date that the user can select.
//
// The maxDate is the latest date that the user can select.
//
// The defaultValue is the value that is first shown and
// will be returned if the user presses enter.
//
// Note that the dates will be converted to RFC3339 format.
//
// The method returns the user's response as a time.Time type.
//
// Example:
// import "time"
//
// s := sdk.New()
// resp, err := s.PromptDatetime("nextRun", "When do you want to run the code next?", "datetime", "T", time.Now(), time.Now().Add(time.Hour * 2), time.Now()) // user selects default
// if err != nil {
//     panic(err)
// }
//
// fmt.Println(resp)
//
// Output:
// [the output will equal time.Now() in 2006-01-02 15:04:05 format]
func (s *Sdk) PromptDatetime(name, msg, variant, flag string, minDate, maxDate, defaultValue time.Time) (time.Time, error) {
	r := PromptRequest{
		Name:       name,
		PromptType: "datetime",
		PromptMsg:  msg,
		Variant:    variant,
		Flag:       flag,
		Default:    defaultValue.Format(time.RFC3339),
		Min:        minDate.Format(time.RFC3339),
		Max:        maxDate.Format(time.RFC3339),
	}

	output, err := s.PromptRequest(r)
	if err != nil {
		return time.Time{}, fmt.Errorf("Error in %s: %s", promptDatetime, err)
	}

	const format = "2006-01-02 15:04:05"

	outputFmt, err := time.Parse(format, output[name].(string))
	if err != nil {
		return time.Time{}, fmt.Errorf("Error in %s: %s", promptDatetime, err)
	}

	return outputFmt, nil
}

// unmarshalFile reads a file from a specified path, parses file into memory,
// and deserializes into the promptOutput type (map[string]interface{}).
//
// This is needed for the prompt type specifically because the SDK daemon's
// response body to prompt requests is synchronous when requested. However,
// the user response to prompts are asynchronous. Thus, the response body reveals
// the file path in which the user response will be written. The file is a blocking
// file (FIFO special file, named pipe) and thus the go program will be blocked
// until user responds.
func (p *PromptOutput) UnmarshalFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}

	defer file.Close()

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	err = json.Unmarshal(bytes, p)
	if err != nil {
		return err
	}

	return nil
}
