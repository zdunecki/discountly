package search

import (
	"encoding/json"
	"github.com/zdunecki/discountly/features/search/models"
	"github.com/gin-gonic/gin"
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
