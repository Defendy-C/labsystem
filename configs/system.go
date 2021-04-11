package configs

import (
	"os"
	"strings"
)

const (
	ProjectName = "labsystem"

	RSAPrivateKeyPEM = "/configs/private_key.pem"
	RSAPublicKeyPEM  = "/configs/public_key.pem"
)

// project path
var curProjectPath string

func CurProjectPath() string {
	if curProjectPath == "" {
		path, err := os.Getwd()
		if err != nil {
			panic("don't located current project position")
		}
		curProjectPath = path[:strings.Index(path, ProjectName)] + ProjectName
	}

	return curProjectPath
}

// system environment
type Environment string

const (
	Development Environment = "development"
	Production  Environment = "production"
)
