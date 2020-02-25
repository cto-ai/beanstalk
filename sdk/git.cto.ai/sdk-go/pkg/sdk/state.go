package sdk

import (
	"fmt"

	"git.cto.ai/sdk-go/pkg/utils"
)

var (
	sdkStateDir = []string{
		"SDK_STATE_DIR",
	}

	sdkHomeDir = []string{
		"SDK_HOME_DIR",
	}

	sdkConfigDir = []string{
		"SDK_CONFIG_DIR",
	}
)

func (s *Sdk) StateDir() (string, error) {
	return getSdkEnv(sdkStateDir)
}

func (s *Sdk) HomeDir() (string, error) {
	return getSdkEnv(sdkHomeDir)
}

func (s *Sdk) ConfigDir() (string, error) {
	return getSdkEnv(sdkConfigDir)
}

// getSdkEnv gets an SDK daemon env variable from the container's
// env and returns it to the user.
//
// The input is a known SDK env variable name.
// These are set at the top of this file.
// These are set as a slice of strings for the possibility
// of env vars having multiple names.
//
// Example:
//
// sdk := sdk.New()
// msg := sdk.HomeDir()
// fmt.Println(msg)
//
// terminal output:
//
// /root
//
func getSdkEnv(env []string) (string, error) {
	var dst string

	utils.SetFromEnvValOrFallback(&dst, env, "")
	if dst == "" {
		return "", fmt.Errorf("Error %s not set in container", env[0])
	}

	return dst, nil
}
