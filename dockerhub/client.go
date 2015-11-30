package dockerhub

import (
	"fmt"
	"net/http/cookiejar"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/franela/goreq"
	"github.com/saulshanabrook/pypi-dockerhub/github"
)

type Auth struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Repo struct {
	Owner string
	Name  string
}

type webAuth struct {
	token string
	jar   *cookiejar.Jar
}

type Client struct {
	*Repo
	*webAuth
	github *github.Repo
}

// NewClient logs you into Docker Hub.
func NewClient(auth *Auth, repo *Repo, ghRepo *github.Repo) (c *Client, err error) {
	wa, err := createWebAuth(auth)
	if err != nil {
		return nil, wrapError(err, "logging in")
	}
	err = wa.verifyLoggedIn()
	return &Client{repo, wa, ghRepo}, wrapError(err, "verifying logged in")
}

func (wa *webAuth) callURL(url, method string, body interface{}, statusCode int, resJSON interface{}) (res *goreq.Response, err error) {
	req := goreq.Request{
		Method:      method,
		Uri:         url,
		Body:        body,
		CookieJar:   wa.jar,
		Accept:      "application/json",
		Host:        "hub.docker.com",
		ContentType: "application/json",
	}.WithHeader("Referer", "https://hub.docker.com/login/")
	if wa.token != "" {
		req = req.WithHeader("Authorization", fmt.Sprintf("JWT %v", wa.token))

	}
	res, err = req.Do()
	if err != nil {
		return res, wrapError(err, fmt.Sprintf("%v to %v", method, url))
	}
	if statusCode != 0 && res.StatusCode != statusCode {
		return res, wrongResponseError(res,
			fmt.Sprintf("%v to %v should have returned a %v", method, url, statusCode))
	}
	if resJSON == nil {
		return
	}
	err = res.Body.FromJsonTo(resJSON)
	if err != nil {
		return res, wrapError(err, fmt.Sprintf("extracting JSON from %v to %v", method, url))
	}
	return
}

func (wa *webAuth) callAPI(path, method string, body interface{}, statusCode int, resJSON interface{}) (res *goreq.Response, err error) {
	return wa.callURL(fmt.Sprintf("https://hub.docker.com/%v", path), method, body, statusCode, resJSON)
}

func (c *Client) callRepo(rel Release, path, method string, body interface{}, statusCode int, resJSON interface{}) (*goreq.Response, error) {
	return c.callAPI(
		fmt.Sprintf("v2/repositories/%v/%v/%v", c.Owner, c.Name, path),
		method,
		body,
		statusCode,
		resJSON)
}

// 1. POST JSON of auth to https://hub.docker.com/v2/users/login/, get back `{"token": "<whatever it is>"}`
// (? maybe not neccesary) 2. POST JSON of token to https://hub.docker.com/attempt-login/ as `{"jwt": "whatever it is"}` to get back JWT cookie
func createWebAuth(auth *Auth) (wa *webAuth, err error) {
	jar, err := cookiejar.New(nil)
	if err != nil {
		return nil, err
	}
	wa = &webAuth{jar: jar}
	goreq.SetConnectTimeout(10 * time.Second)
	log.WithFields(log.Fields{
		"auth": auth,
	}).Debug("Logging into DockerHub")
	var resJSON struct {
		Token string `json:"token"`
	}
	res, err := wa.callAPI("v2/users/login/", "POST", auth, 200, &resJSON)
	if err != nil {
		return nil, wrapError(err, "login")
	}
	if resJSON.Token == "" {
		return nil, fmt.Errorf("Didnt get a token back from the login")
	}
	wa.token = resJSON.Token
	if err = res.Body.Close(); err != nil {
		return nil, wrapError(err, "closing body of POST login")
	}

	log.WithFields(log.Fields{}).Debug("Posting login back in to get cookie")
	res, err = wa.callAPI("attempt-login/", "POST", struct {
		Jwt string `json:"jwt"`
	}{wa.token}, 200, nil)
	if err != nil {
		return nil, wrapError(err, "login")
	}
	if err = res.Body.Close(); err != nil {
		return nil, wrapError(err, "closing body of POST attempt-login")
	}
	return
}

func (wa *webAuth) verifyLoggedIn() error {
	log.WithFields(log.Fields{}).Debug("Verifying can get user")

	res, err := wa.callAPI("v2/user/", "GET", "", 200, nil)
	if err != nil {
		return wrapError(err, "verifyLoggedIn")
	}
	return wrapError(res.Body.Close(), "closing body on GET user")
}

func (c *Client) verifyAccessNamespace() error {
	log.WithFields(log.Fields{}).Debug("Verifying passed in namespace is within namespace")

	res, err := c.callAPI("v2/repositories/namespaces/", "GET", "", 200, nil)
	var rBody struct {
		Namespaces []string `json:"namespaces"`
	}
	if err = res.Body.FromJsonTo(&rBody); err != nil {
		return wrapError(err, "getting json from namespace return")
	}
	if err = res.Body.Close(); err != nil {
		return wrapError(err, "closing body on GET namespaces")
	}
	if !contains(rBody.Namespaces, c.Owner) {
		return fmt.Errorf(
			"The %s namespace is not in the ones in your account: %v",
			c.Owner,
			rBody.Namespaces,
		)
	}
	return nil
}
