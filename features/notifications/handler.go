package notifications

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
)

type PusherAuth struct {
	Auth string `json:"auth"`
}

func Auth(c *gin.Context) {
	params, _ := ioutil.ReadAll(c.Request.Body)

	response, err := Client.AuthenticatePrivateChannel(params)

	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{"errors": err.Error()})
		return
	}

	var auth PusherAuth

	if err := json.Unmarshal(response, &auth); err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{"errors": err.Error()})
		return
	}

	c.JSON(http.StatusOK, auth)
}

func ReceiveHook(c *gin.Context) {
	body, _ := ioutil.ReadAll(c.Request.Body)

	webHook, err := Client.Webhook(c.Request.Header, body)

	if err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{"errors": err.Error()})
		return
	}

	go receiveHook(webHook)

	c.JSON(http.StatusOK, http.StatusOK)
}
