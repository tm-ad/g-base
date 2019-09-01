package g_base

import "os"

const (
	// DevelopmentEnv is the name of the environment variable used to mark the running mode of the application
	DevelopmentEnv = "GO_ENV"
	// WorkingPathEnv is the name of the environment variable used to mark the directory where the application runs
	WorkingPathEnv = "GO_WORKING_PATH"

	// DEV is a development pattern markup constant
	DEV = "development"
	// TESTING is a testing pattern markup constant
	TESTING = "testing"
	// PROD is a production pattern markup constant
	PROD = "production"
)

// GOENV gets application running mode
func GOENV() string {
	return Env(DevelopmentEnv)
}

// Env gets the specified environment variable
func Env(name string) string {
	return os.Getenv(name)
}

// SetEnv sets environment variable
func SetEnv(name, value string) error {
	return os.Setenv(name, value)
}

// Development determines whether or not it is currently in development mode
// GO_ENV = development
func Development() bool {
	return Env(DevelopmentEnv) == DEV
}

// Testing determines whether or not it is currently in testing mode
// GO_ENV = testing
func Testing() bool {
	return Env(DevelopmentEnv) == TESTING
}

// Production determines whether or not it is currently in production mode
// GO_ENV != testing && GO_ENV != development
func Production() bool {
	return Env(DevelopmentEnv) != DEV && Env(DevelopmentEnv) != TESTING
}
