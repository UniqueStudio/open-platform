package handler

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"net/http"
	"open-platform/db"
	"open-platform/utils"
	"sync"

	"github.com/gin-gonic/gin"
)

type shorter struct {
	UrlList  []string `json:"urllist"`
	Number   int      `json:"number"`
}

//CreateShortUrlHandler is a func to handle shorturl generation
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

	var results = make(map[string]interface{})
	var lock sync.Mutex
	var wg sync.WaitGroup
	wg.Add(len(urlList))

	for _, v := range urlList {
		go func(lUrl string) {
			lock.Lock()
			defer lock.Unlock()

			shortUrl, err := GenShortUrl(lUrl)
			if err != nil {
				results[lUrl] = err.Error()
			}else{
				results[lUrl] = shortUrlPrefix + shortUrl
			}

			wg.Done()
		}(v)
	}

	wg.Wait()
	c.JSON(http.StatusOK, gin.H{"urls": results})

}

//GenShortUrl is a func to generate shorturl and insert the item to database
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
