package handler

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"net/http"
	"open-platform/db"
	"open-platform/utils"

	"github.com/gin-gonic/gin"
)

type resultForCreation struct {
	url  string
	code string
	err error
}


type shorter struct {

	UrlList  []string `json:"urllist"`
	Number   int      `json:"number"`
	ShortUrl []string `json:"shorturl"`
	HashCode []string `json:"hashcode"`
}

func CreateShortUrlHandler(c *gin.Context) {
	var data shorter
	c.BindJSON(&data)

	urlList := data.UrlList
	number := data.Number
	shortUrlPrefix := utils.AppConfig.ShortUrl  //"http://uniqs.cc"

	if len(urlList) == 0 || number == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Data missing parameter"})
		return
	}

	resultCode := make(chan resultForCreation)
	for _, v := range urlList {
		go func(lUrl string) {
			shortUrl, err := GenShortUrl(lUrl)
			if err != nil {
				//logs.Error("gen shortUrl failed, error: " + err.Error())
				resultCode <- resultForCreation{lUrl, "",err}
			} else {
				//logs.Info("[create]: " + lUrl + " => " + shortUrl)
				resultCode <- resultForCreation{lUrl, shortUrl,nil}
			}

		}(v)
	}

	var results = make(map[string]interface{})
	var count = 1
	for {
		res := <-resultCode

		if res.err !=nil{
			results[res.url] = res.err.Error()
		}else{
			results[res.url] = shortUrlPrefix + res.code
		}

		if count == len(urlList) {
			close(resultCode)
			c.JSON(http.StatusOK, gin.H{"urls": results})
			return
		}
		count++
	}
}

func GenShortUrl(Url string) (ShortUrl string, err error) {
	urlmd5 := MD5(Url)
	result := new(db.Short_Url)
	has, err := db.ORM.Where("hashcode=?", urlmd5).Get(result)

	if err != nil {
		return "", err
	}

	if has == true {
		return result.Shorturl, nil
	} else {

		newUrl := &db.Short_Url{
			Url:      Url,
			Hashcode: urlmd5,
		}
		_, err := db.ORM.Insert(newUrl)

		if err != nil {
			return "", err
		}

		id := newUrl.Id
		if id == 0 {
			return "", errors.New("get id failed")
		}

		shorturl := TransToCode(id)

		if shorturl == "" {
			return "", errors.New("gen code failed")
		}
		newUrl.Shorturl = shorturl
		_, err = db.ORM.Id(id).Update(newUrl)

		if err != nil {
			return "", err
		}

		return shorturl, nil
	}

}

func MD5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}
func TransToCode(id int64) string {
	bytes := []byte("0lv12NUJ3789qazwegbyhnujmipQAZWsxSXEDCR4kt56FVTGBYHMIodcrfKLOP")

	var code string
	for m := id; m > 0; m = m / 62 {
		n := m % 62
		code += string(bytes[n])
		if m < 62 {
			break
		}
	}
	return code
}
