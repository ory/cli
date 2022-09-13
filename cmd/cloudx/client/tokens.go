package client

import (
	"os"

	"github.com/ory/x/stringsx"
)

// GetProjectAPIKeyFromEnvironment returns the project API key from the environment variable.
func GetProjectAPIKeyFromEnvironment() string {
	return stringsx.Coalesce(os.Getenv("ORY_API_KEY"), os.Getenv("ORY_PERSONAL_ACCESS_TOKEN"), os.Getenv("ORY_PAT"))
}
