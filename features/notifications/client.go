package notifications

import (
	"github.com/pusher/pusher-http-go"
	"os"
)

var Client pusher.Client

func init() {
	Client = pusher.Client{
		AppId:   os.Getenv("PUSHER_APP_ID"),
		Key:     os.Getenv("PUSHER_APP_KEY"),
		Secret:  os.Getenv("PUSHER_APP_SECRET"),
		Cluster: "eu",
		Secure:  true,
	}
}
