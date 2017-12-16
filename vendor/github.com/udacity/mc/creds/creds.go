package creds

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sync"
	"time"

	hoth "github.com/udacity/go-hoth"
	"github.com/udacity/mc/config"
)

type credentials struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	IEmail   string `json:"impersonating_email"`
	IUID     string `json:"impersonating_uid"`
}

type credentialer interface {
	load() (credentials, error)
	save(credentials, string) error
	remove() error
	impersonate(string, string) (string, error)
	getJWT(string, string) (string, error)
}

type defaultCredentialer struct{}

// load local creds.
func (defaultCredentialer) load() (credentials, error) {
	b, err := ioutil.ReadFile(config.CredsFilepath())
	if err != nil {
		return credentials{}, err
	}

	creds := credentials{}
	err = json.Unmarshal(b, &creds)
	if err != nil {
		return credentials{}, err
	}

	return creds, nil
}

// save writes credentials to files with compatibility for other mentorship
// CLIs. jwt is treated separately because while mc requests a fresh JWT every
// time it runs, other CLIs look for a hardcoded JWT, which we save in a
// separate file.
func (defaultCredentialer) save(creds credentials, jwt string) error {
	c, err := json.Marshal(creds)
	if err != nil {
		return err
	}

	// save creds for `mc`
	err = ioutil.WriteFile(config.CredsFilepath(), c, 0644)
	if err != nil {
		return err
	}

	// save creds for other mentorship CLIs
	err = ioutil.WriteFile(config.MentorshipJWTFilepath(), []byte(jwt), 0644)
	if err != nil {
		return err
	}
	return nil
}

// delete the creds files
func (defaultCredentialer) remove() error {
	files := []string{config.CredsFilepath(), config.MentorshipJWTFilepath()}

	for _, f := range files {
		if _, err := os.Stat(f); !os.IsNotExist(err) {
			err = os.Remove(f)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// impersonate requests the JWT for another user.
func (defaultCredentialer) impersonate(jwt string, uid string) (string, error) {
	var httpClient = &http.Client{
		Timeout: time.Second * config.HTTPTimeoutInSeconds,
	}

	type payload struct {
		UserID string `json:"user_id"`
	}
	type response struct {
		JWT string `json:"jwt"`
	}

	p, err := json.Marshal(payload{uid})
	if err != nil {
		return "", err
	}
	br := bytes.NewReader(p)

	req, err := http.NewRequest("POST", "https://hoth.udacity.com/v2/impersonate", br)
	if err != nil {
		return "", err
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", jwt))

	res, err := httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer func() {
		err := res.Body.Close()
		if err != nil {
			fmt.Println(err)
		}
	}()

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	if res.StatusCode >= 400 {
		return "", fmt.Errorf("HTTP Status:%v Response: %v", res.StatusCode, string(resBody))
	}

	r := response{}
	err = json.Unmarshal(resBody, &r)

	// Return the JWT if everything is good.
	return r.JWT, nil
}

func (defaultCredentialer) getJWT(email string, password string) (string, error) {
	return hoth.GetJWT(email, password)
}

// Login as a new user. Also ends impersonations.
func Login(email string, password string) error {
	return login(defaultCredentialer{}, email, password)
}

func login(c credentialer, email string, password string) error {
	// test the credentials
	jwt, err := c.getJWT(email, password)
	if err != nil {
		return err
	}

	creds := credentials{
		email, password, "", "",
	}

	return c.save(creds, jwt)
}

// Load provides email, password.
func Load() (string, string, error) {
	return load(defaultCredentialer{})
}

func load(c credentialer) (string, string, error) {
	creds, err := c.load()
	if err != nil {
		return "", "", err
	}

	return creds.Email, creds.Password, nil
}

// LoadI provides the email address of the logged in user and whether that is
// the original user or impersonatee.
func LoadI() (email string, impersonating bool, err error) {
	return loadi(defaultCredentialer{})
}

func loadi(c credentialer) (string, bool, error) {
	creds, err := c.load()
	if err != nil {
		return "", false, err
	}

	if creds.IEmail != "" {
		return creds.IEmail, true, nil
	}

	return creds.Email, false, nil
}

type jwtCache struct {
	InUse        string
	Original     string
	Impersonatee string
	mux          sync.Mutex
}

func (c *jwtCache) Lock() {
	c.mux.Lock()
}

func (c *jwtCache) Unlock() {
	c.mux.Unlock()
}

func (c *jwtCache) Clear() {
	c.InUse = ""
	c.Original = ""
	c.Impersonatee = ""
}

var cache = jwtCache{}

// FetchJWT gets a JWT of the user or user being impersonated.
func FetchJWT() (string, error) {
	return fetchJWT(defaultCredentialer{})
}

func fetchJWT(c credentialer) (string, error) {
	cache.Lock()
	defer cache.Unlock()

	if cache.InUse != "" {
		return cache.InUse, nil
	}

	creds, err := c.load()
	if err != nil {
		return "", err
	}

	// jwt starts as the original user...
	var jwt string
	jwt, err = getOriginalJWT(c)
	if err != nil {
		return "", err
	}

	// ...but if there is an impersonatee, we overwrite it
	if creds.IUID != "" {
		jwt, err = getImpersonateeJWT(c, creds.IUID)
		if err != nil {
			return "", err
		}
	}

	cache.InUse = jwt
	return cache.InUse, nil
}

func getOriginalJWT(c credentialer) (string, error) {
	if cache.Original != "" {
		return cache.Original, nil
	}

	creds, err := c.load()
	if err != nil {
		return "", err
	}

	if creds.Email == "" || creds.Password == "" {
		return "", fmt.Errorf("error: no email or password")
	}

	jwt, err := c.getJWT(creds.Email, creds.Password)
	if err != nil {
		return "", err
	}

	cache.Original = jwt
	return cache.Original, nil
}

func getImpersonateeJWT(c credentialer, uid string) (string, error) {
	if cache.Impersonatee != "" {
		return cache.Impersonatee, nil
	}

	oldjwt, err := getOriginalJWT(c)
	if err != nil {
		return "", err
	}

	jwt, err := c.impersonate(oldjwt, uid)
	if err != nil {
		return "", err
	}

	cache.Impersonatee = jwt
	return cache.Impersonatee, nil
}

// Impersonate will ensure that future commands run with the JWT of a different
// user.
func Impersonate(uid string, email string) error {
	return impersonate(defaultCredentialer{}, uid, email)
}

func impersonate(c credentialer, uid string, email string) error {
	jwt, err := getImpersonateeJWT(c, uid)
	if err != nil {
		return err
	}

	creds, err := c.load()
	if err != nil {
		return err
	}

	creds.IEmail = email
	creds.IUID = uid

	return c.save(creds, jwt)
}

// StopImpersonating removes impersonatee credentials and sets JWT back to the
// original user.
func StopImpersonating() error {
	return stopImpersonating(defaultCredentialer{})
}

func stopImpersonating(c credentialer) error {
	creds, err := c.load()
	if err != nil {
		return err
	}

	jwt, err := getOriginalJWT(c)
	if err != nil {
		return err
	}

	creds.IEmail = ""
	creds.IUID = ""
	return c.save(creds, jwt)
}

// Logout removes credentials.
func Logout() error {
	return remove(defaultCredentialer{})
}

func remove(c credentialer) error {
	return c.remove()
}

// LoggedIn returns an error if not logged in.
func LoggedIn() bool {
	return loggedIn(defaultCredentialer{})
}

func loggedIn(c credentialer) bool {
	creds, err := c.load()
	if err != nil {
		fmt.Println(err)
		return false
	}

	if creds.Email == "" || creds.Password == "" {
		return false
	}

	if creds.IUID != "" {
		fmt.Printf("Impersonating %s...\n", creds.IEmail)
	}

	return true
}
