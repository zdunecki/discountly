package search

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/zdunecki/discountly/features/search/models"
	"io/ioutil"
	"net/http"
	"sync"
)

func FindBestDiscounts(c *gin.Context) {
	var s search.Search
	err := c.Bind(&s)
	if err != nil {
		e, _ := json.Marshal(err)
		c.JSON(http.StatusInternalServerError, e)
		return
	}

	wg := &sync.WaitGroup{}
	response, err := findBestDiscounts(wg, s)
	wg.Wait()

	if err != nil {
		e, _ := json.Marshal(err)
		c.JSON(http.StatusInternalServerError, e)
		return
	}

	c.JSON(http.StatusOK, response)
}

func ReceiveNearbyHook(c *gin.Context) {
	b, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	var hook search.Hook

	if err := json.Unmarshal(b, &hook); err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	if err := receiveHook(hook); err != nil {
		c.JSON(http.StatusInternalServerError, err)
		return
	}

	c.JSON(http.StatusOK, "ok")
}
