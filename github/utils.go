package github

import "fmt"

func wrapError(err error, message string) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("Github(%v): %v", message, err.Error())
}
