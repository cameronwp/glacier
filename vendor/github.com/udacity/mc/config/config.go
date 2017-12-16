package config

import (
	"log"
	"os/user"
	"path/filepath"
)

const (
	// HTTPTimeoutInSeconds denotes the number of seconds before signaling a
	// timeout for an external HTTP request.
	HTTPTimeoutInSeconds = 30
	// Dirname is where mc config and log files are stored.
	Dirname = ".mc"
	// UserAgent is the reported UA for HTTP requests.
	UserAgent = "mc"
)

// CredsFilepath returns the path of the creds file.
func CredsFilepath() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	credsFilename := "creds"
	return filepath.Join(usr.HomeDir, Dirname, credsFilename)
}

// LogsFilepath returns the path of the logs file.
func LogsFilepath() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	logFilename := "errs.log"
	return filepath.Join(usr.HomeDir, Dirname, logFilename)
}

// MentorshipJWTFilepath returns the location where other mentorship CLIs save
// their JWTs. Allows `mc` to be used to log in to other CLIs.
func MentorshipJWTFilepath() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}
	jwtFilename := ".hoth_jwt"
	return filepath.Join(usr.HomeDir, jwtFilename)
}
