package dockerhub

import (
	"fmt"

	"github.com/saulshanabrook/pypi-dockerhub/Godeps/_workspace/src/github.com/franela/goreq"
)

func wrapError(err error, message string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%v: %v", message, err.Error())
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func wrongResponseError(res *goreq.Response, mes string) error {
	body, err := res.Body.ToString()

	if err != nil {
		return wrapError(err, "trying to get string of body in WrongReponse ")
	}

	return fmt.Errorf("%v: %v: %v\n\n%v", res.Request, res.Status, res.Response.Header, body)
}
