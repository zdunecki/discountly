package oauth

import (
	"fmt"
	"github.com/zdunecki/discountly/infra"
	"io/ioutil"
	"math/rand"
	"net/http"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

func randomInt(min, max int) int {
	return min + rand.Intn(max-min)
}

func randomString(len int) string {
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		bytes[i] = byte(randomInt(65, 90))
	}
	return string(bytes)
}

var (
	googleOauthConfig *oauth2.Config
	randomState       string
	oauthStateString  = "pseudo-random"
)

func init() {
	clientId := infra.GetEnv("OAUTH_CLIENT_ID")
	clientSecret := infra.GetEnv("OAUTH_CLIENT_SECRET")
	appURL := infra.GetEnv("APP_URL")

	googleOauthConfig = &oauth2.Config{
		RedirectURL:  appURL + "/oauth/callback",
		ClientID:     clientId,
		ClientSecret: clientSecret,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint:     google.Endpoint,
	}

	randomState = randomString(10)
}

func getAuthUrl() string {
	return googleOauthConfig.AuthCodeURL(oauthStateString)
}

func getUserInfo(state string, code string, mobile bool) ([]byte, string, error) {
	if state != oauthStateString && !mobile {
		return nil, "", fmt.Errorf("invalid oauth state")
	}

	token, err := googleOauthConfig.Exchange(oauth2.NoContext, code)

	if err != nil {
		return nil, "", fmt.Errorf("code exchange failed: %s", err.Error())
	}

	response, err := getUserInfoFromToken(token.AccessToken)
	defer response.Body.Close()

	if err != nil {
		return nil, "", fmt.Errorf("failed getting user info: %s", err.Error())
	}

	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, "", fmt.Errorf("failed reading response body: %s", err.Error())
	}

	return contents, token.AccessToken, nil
}

func userInfoFromToken(token string) ([]byte, string, error) {
	response, err := getUserInfoFromToken(token)
	defer response.Body.Close()

	if err != nil {
		return nil, "", fmt.Errorf("failed getting user info: %s", err.Error())
	}

	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, "", fmt.Errorf("failed reading response body: %s", err.Error())
	}

	return contents, token, nil
}

func getUserInfoFromToken(accessToken string) (*http.Response, error) {
	return http.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + accessToken)
}
