package infra

import "net/url"

func GetExposedHost() string {
	u, err := url.Parse(GetEnv("EXPOSE_DOCKER_APP_URL"))
	if err != nil {
		panic(err)
	}

	return u.Host
}

