## Quickstart

### Install the package

`> go get github.com/udacity/uy-go`

### Simple Example

```go
package main

import (
    "log"
    uy "github.com/udacity/uy-go"
    "io/ioutil"
)

func main() {
	// set up the api client
	client := uy.NewAPIClient()

	// authorize the api client
	err := client.Auth(os.Getenv("UY_USERNAME"), os.Getenv("UY_PASSWORD"))
	if err != nil {
		log.Printf("Error Authorizing. Error: %v", err)
		return
	}

	// example: get users/me
	resp, err := client.Get("users/me")
	if err != nil {
		panic(err)
	}
	// the api library doesn't (yet) offer any helpers for response handling
	// here's an example of reading the response body into a slice of bytes
    // you probably want to look at the `encoding/json` package
	defer resp.Body.Close()
	jb, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	log.Printf("ME: %s", string(jb))
}
```

### Contributing

Project Maintainer: @artgillespie

Found a bug? Have a feature request? Please use the [github issues](https://github.com/udacity/uy-go/issues) for this repo.

Have a patch? Please submit a pull request against the `master` branch.
