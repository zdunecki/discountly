package infra

import "net/url"

func GetHost() string {
	u, err := url.Parse(GetEnv("DOCKER_APP_URL"))
	if err != nil {
		panic(err)
	}

	return u.Host
}
