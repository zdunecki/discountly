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
//eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoiMTAzNjUyOTYzMTA4NDEwNjExMTE0IiwiZXhwIjoxNTQ5NzM1ODM1fQ.KMOT_al4nMvk02_F13hFnO6OlX_u-v528knuBEkIxJo