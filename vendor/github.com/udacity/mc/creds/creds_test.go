package creds

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

type testCredentialer struct {
	loggedin      bool
	impersonating bool
	loaderr       bool
	saveerr       bool
	removeerr     bool
	getjwterr     bool
}

var (
	originalEmail     = "someone@udacity.com"
	originalJWT       = "originaljwtasdf"
	originalPassword  = "hunter2"
	impersonateeEmail = "someoneelse@udacity.com"
	impersonateeJWT   = "impersonateejwtfdsa"
	impersonateeUID   = "u123456789"
)

func before() {
	cache.Clear()
	callsToGetJWT = 0
	callsToLoad = 0
	callsToSave = 0
	callsToRemove = 0
	callsToImpersonate = 0
}

func sleepRandom() {
	r := rand.Int() % 10
	time.Sleep(time.Millisecond * time.Duration(r))
}

var callsToLoad = 0

func (t testCredentialer) load() (credentials, error) {
	callsToLoad++
	if t.loaderr {
		return credentials{}, fmt.Errorf("load error: triggered")
	}

	creds := credentials{}
	if !t.loggedin {
		return creds, nil
	}

	creds.Email = originalEmail
	creds.Password = originalPassword

	if !t.impersonating {
		return creds, nil
	}

	creds.IEmail = impersonateeEmail
	creds.IUID = impersonateeUID
	return creds, nil
}

var callsToSave = 0

func (t testCredentialer) save(credentials, string) error {
	sleepRandom()
	callsToSave++
	if t.saveerr {
		return fmt.Errorf("save error: triggered")
	}
	return nil
}

var callsToRemove = 0

func (t testCredentialer) remove() error {
	sleepRandom()
	callsToRemove++
	if t.removeerr {
		return fmt.Errorf("remove error: triggered")
	}
	return nil
}

var callsToGetJWT = 0

func (t testCredentialer) getJWT(email string, password string) (string, error) {
	sleepRandom()
	callsToGetJWT++

	if t.getjwterr {
		return "", fmt.Errorf("getJWT error: on purpose")
	}

	if email == impersonateeEmail {
		return "", fmt.Errorf("getJWT error: fetching jwt of impersonatee")
	}

	if email == originalEmail && password == originalPassword {
		return originalJWT, nil
	}

	return "", fmt.Errorf("getJWT error: bad email/password '%s'/'%s'", email, password)
}

var callsToImpersonate = 0

func (t testCredentialer) impersonate(oldjwt string, uid string) (string, error) {
	sleepRandom()
	callsToImpersonate++
	if oldjwt != originalJWT {
		return "", fmt.Errorf("impersonate error: wrong original JWT")
	}

	if uid == impersonateeUID {
		return impersonateeJWT, nil
	}

	return "", fmt.Errorf("impersonate error: bad uid '%s'", uid)
}

func TestLogin(t *testing.T) {
	t.Run("when first time logging in with valid email/password", func(t *testing.T) {
		before()
		c := testCredentialer{}

		err := login(c, originalEmail, originalPassword)
		if err != nil {
			t.Error(err)
		}

		if callsToSave != 1 {
			t.Errorf("expected 1 call to save, found %d", callsToSave)
		}
	})

	t.Run("when using invalid email/password", func(t *testing.T) {
		before()
		c := testCredentialer{}

		err := login(c, "bademail", "badpassword")
		if err == nil {
			t.Errorf("expected error but found none")
		}

		if callsToSave != 0 {
			t.Errorf("expected 0 calls to save, found %d", callsToSave)
		}
	})

	t.Run("when already logged in", func(t *testing.T) {
		before()
		c := testCredentialer{}
		c.loggedin = true

		err := login(c, originalEmail, originalPassword)
		if err != nil {
			t.Error(err)
		}

		if callsToSave != 1 {
			t.Errorf("expected 1 calls to save, found %d", callsToSave)
		}
	})
}

func TestLogout(t *testing.T) {
	t.Run("when logged in", func(t *testing.T) {
		before()
		c := testCredentialer{}
		c.loggedin = true

		err := remove(c)
		if err != nil {
			t.Error(err)
		}

		if callsToRemove != 1 {
			t.Errorf("expected 1 remove call, found %d", callsToRemove)
		}
	})

	t.Run("when not logged in", func(t *testing.T) {
		before()
		c := testCredentialer{}

		err := remove(c)
		if err != nil {
			t.Error(err)
		}

		if callsToRemove != 1 {
			t.Errorf("expected 1 remove call, found %d", callsToRemove)
		}
	})
}

func TestLoggedIn(t *testing.T) {
	t.Run("when logged in", func(t *testing.T) {
		before()
		c := testCredentialer{}
		c.loggedin = true

		if !loggedIn(c) {
			t.Errorf("expected to be logged in")
		}
	})

	t.Run("when not logged in", func(t *testing.T) {
		before()
		c := testCredentialer{}

		if loggedIn(c) {
			t.Errorf("expected to not be logged in")
		}
	})
}

func TestImpersonate(t *testing.T) {
	t.Run("when logged in", func(t *testing.T) {
		c := testCredentialer{}
		c.loggedin = true
		t.Run("with a valid impersonatee", func(t *testing.T) {
			before()
			err := impersonate(c, impersonateeUID, impersonateeEmail)
			if err != nil {
				t.Error(err)
			}

			if callsToImpersonate != 1 {
				t.Errorf("expected 1 call to impersonate, found %d", callsToImpersonate)
			}
			if callsToSave != 1 {
				t.Errorf("expected 1 call to save, found %d", callsToSave)
			}
		})

		t.Run("with an invalid impersonatee", func(t *testing.T) {
			before()
			err := impersonate(c, "nottherightuid", "nottherightemail")
			if err == nil {
				t.Errorf("expected an error, got none")
			}

			if callsToImpersonate != 1 {
				t.Errorf("expected 1 call to impersonate, found %d", callsToImpersonate)
			}
			if callsToSave != 0 {
				t.Errorf("expected 0 calls to save, found %d", callsToSave)
			}
		})
	})

	t.Run("when not logged in", func(t *testing.T) {
		c := testCredentialer{}
		t.Run("with a valid impersonatee", func(t *testing.T) {
			before()
			err := impersonate(c, impersonateeUID, impersonateeEmail)
			if err == nil {
				t.Errorf("expected an error, got none")
			}

			if callsToGetJWT != 0 {
				t.Errorf("expected 0 calls to getJWT, found %d", callsToGetJWT)
			}
			if callsToImpersonate != 0 {
				t.Errorf("expected 0 call to impersonate, found %d", callsToImpersonate)
			}
			if callsToSave != 0 {
				t.Errorf("expected 0 calls to save, found %d", callsToSave)
			}
		})
	})
}

func TestStopImpersonating(t *testing.T) {
	t.Run("when not logged in", func(t *testing.T) {
		before()
		c := testCredentialer{}

		err := stopImpersonating(c)
		if err == nil {
			t.Errorf("expected an error, found none")
		}
	})

	t.Run("when not actually impersonating", func(t *testing.T) {
		before()
		c := testCredentialer{}
		c.loggedin = true
		c.impersonating = false

		err := stopImpersonating(c)
		if err != nil {
			t.Error(err)
		}

		if callsToGetJWT != 1 {
			t.Errorf("expected 1 calls to getJWT, found %d", callsToGetJWT)
		}
		if callsToSave != 1 {
			t.Errorf("expected 1 calls to save, found %d", callsToSave)
		}
	})

	t.Run("when impersonating", func(t *testing.T) {
		before()
		c := testCredentialer{}
		c.loggedin = true
		c.impersonating = true

		err := stopImpersonating(c)
		if err != nil {
			t.Error(err)
		}

		if callsToGetJWT != 1 {
			t.Errorf("expected 1 calls to getJWT, found %d", callsToGetJWT)
		}
		if callsToSave != 1 {
			t.Errorf("expected 1 calls to save, found %d", callsToSave)
		}
	})
}

func TestWhoAmI(t *testing.T) {
	t.Run("when not logged in", func(t *testing.T) {
		before()
		c := testCredentialer{}
		c.loaderr = true

		_, _, err := loadi(c)
		if err == nil {
			t.Errorf("expected an error but found none")
		}
	})

	t.Run("when logged in", func(t *testing.T) {
		before()
		c := testCredentialer{}
		c.loggedin = true

		email, impersonating, err := loadi(c)
		if err != nil {
			t.Error(err)
		}

		if email != originalEmail {
			t.Errorf("expected email %s, got '%s'", originalEmail, email)
		}
		if impersonating {
			t.Errorf("expected to not be impersonating")
		}
	})

	t.Run("when logged in", func(t *testing.T) {
		before()
		c := testCredentialer{}
		c.loggedin = true
		c.impersonating = true

		email, impersonating, err := loadi(c)
		if err != nil {
			t.Error(err)
		}

		if email != impersonateeEmail {
			t.Errorf("expected email %s, got '%s'", impersonateeEmail, email)
		}
		if !impersonating {
			t.Errorf("expected to be impersonating")
		}
	})
}

func TestFetchJWT(t *testing.T) {
	t.Run("when not logged in (but the creds file exists)", func(t *testing.T) {
		before()
		c := testCredentialer{}
		c.loggedin = false

		jwt, err := fetchJWT(c)
		if err == nil {
			t.Errorf("expected error but did not find one")
		}
		if jwt != "" {
			t.Errorf("expected to find an empty string jwt, got %s", jwt)
		}
	})

	t.Run("when not logged in (and the creds file hasn't been created)", func(t *testing.T) {
		before()
		c := testCredentialer{}
		c.loggedin = false
		c.loaderr = true

		jwt, err := fetchJWT(c)
		if err == nil {
			t.Errorf("expected an error but did not get one")
		}
		if jwt != "" {
			t.Errorf("expected to find an empty string jwt, got %s", jwt)
		}
	})

	t.Run("when logged in", func(t *testing.T) {
		before()
		c := testCredentialer{}
		c.loggedin = true

		jwt, err := fetchJWT(c)
		if err != nil {
			t.Error(err)
		}
		if jwt != originalJWT {
			t.Errorf("expected jwt %s, got '%s'", originalJWT, jwt)
		}
	})

	t.Run("when impersonating", func(t *testing.T) {
		before()
		c := testCredentialer{}
		c.loggedin = true
		c.impersonating = true

		jwt, err := fetchJWT(c)
		if err != nil {
			t.Error(err)
		}
		if jwt != impersonateeJWT {
			t.Errorf("expected jwt %s, got '%s'", impersonateeEmail, jwt)
		}
	})

	t.Run("caches when run more than once", func(t *testing.T) {
		before()
		c := testCredentialer{}
		c.loggedin = true

		jwt1, err := fetchJWT(c)
		if err != nil {
			t.Error(err)
		}
		jwt2, err := fetchJWT(c)
		if err != nil {
			t.Error(err)
		}

		if jwt1 != jwt2 {
			t.Errorf("expected both jwts to match, got '%s' and '%s'", jwt1, jwt2)
		}

		if callsToGetJWT != 1 {
			t.Errorf("expected 1 call to getJWT, found %d", callsToGetJWT)
		}
	})

	t.Run("uses the cache when run on more than one goroutine", func(t *testing.T) {
		before()
		c := testCredentialer{}
		c.loggedin = true
		done := make(chan bool)

		makeClient := func(d chan bool) {
			_, err := fetchJWT(c)
			if err != nil {
				t.Error(err)
			} else {
				d <- true
			}
		}

		go makeClient(done)
		go makeClient(done)

		for i := 0; i < 2; i++ {
			select {
			case <-done:
				break
			}
		}

		// each runthrough of fetchJWT incurs 2 loads (if no impersonation)
		if callsToLoad != 2 {
			t.Errorf("expected 2 calls to load, found %d", callsToLoad)
		}

		if callsToGetJWT != 1 {
			t.Errorf("expected 1 call to getJWT, found %d", callsToGetJWT)
		}
	})
}
