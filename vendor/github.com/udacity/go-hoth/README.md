# Go Hoth API

Go client for the Udacity [Hoth authentication microservice API][1].

Originally created by Brad Erickson (brad.erickson@udacity.com) for the
[BizDev b2b-dashboard][2] project.

[1]: https://github.com/udacity/hoth
[2]: https://github.com/udacity/b2b-dashboard

## Features

* Simplifies generation of Hoth JWTs.
* A backend HTTP client:
  * Exponential backoff request retry.
  * Handles unmarshalling JSON to Structs.
  * Used by Udacity Hoth-auth API clients:
    * go-students: https://github.com/udacity/go-students
    * go-grading: https://github.com/udacity/go-grading
    * go-classroom-content: https://github.com/udacity/go-classroom-content

## Usage

```
$ go get github.com/udacity/go-hoth
```

See `hoth_test.go` or full implementations in the above Udacity API clients for
further details.

## Tests

Run tests with:
```
go test -v
```
