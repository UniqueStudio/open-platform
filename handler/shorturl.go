package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"open-platform/db"
	"strings"
)

//MapShortUrlHandler is a function to redirecte a request with shorturl to the origin url
func MapShortUrlHandler(c *gin.Context)  {
	var result db.Short_Url
	shorturl := c.Param("Shorturl")
	has,err := db.ORM.Where("Shorturl=?",shorturl).Get(&result)

	if err != nil {
		c.JSON(http.StatusBadGateway,gin.H{"message":"something wrong with db"})
		return
	}
	if has != true {
		c.JSON(http.StatusNotFound,gin.H{"message":"404 not found"})
		return
	}

	if !strings.Contains(result.Url,"http://") && !strings.Contains(result.Url,"https://"){
		result.Url = "http://"+result.Url
	}

	c.Header("Location",result.Url)
	c.JSON(http.StatusFound,gin.H{"url":result.Url})

	return
}
