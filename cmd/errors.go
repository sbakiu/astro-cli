package cmd

import (
	"errors"
)

var (
	errInvalidBothAirflowAndRuntimeVersions = errors.New("You provided both a runtime version and an Airflow version. You have to provide only one of these to initialize your project.") //nolint

	errConfigProjectName = errors.New("project name is invalid")
	errProjectNameSpaces = errors.New("this project name is invalid, a project name cannot contain spaces. Try using '-' instead")

	errInvalidSetArgs    = errors.New("must specify exactly two arguments (key value) when setting a config")
	errInvalidConfigPath = errors.New("config does not exist, check your config key")
)
