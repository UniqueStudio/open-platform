package handler

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
)

type resultForCreation struct {
	url  string
	code string
}

var decToB62Map = map[int]string{0: "A", 1: "B", 2: "C", 3: "D", 4: "E", 5: "F", 6: "G", 7: "H", 8: "I", 9: "J", 10: "K", 11: "L", 12: "M", 13: "N", 14: "O", 15: "P", 16: "Q", 17: "R", 18: "S", 19: "T", 20: "U", 21: "V", 22: "W", 23: "X", 24: "Y", 25: "Z", 26: "a", 27: "b", 28: "c", 29: "d", 30: "e", 31: "f", 32: "g", 33: "h", 34: "i", 35: "j", 36: "k", 37: "l", 38: "m", 39: "n", 40: "o", 41: "p", 42: "q", 43: "r", 44: "s", 45: "t", 46: "u", 47: "v", 48: "w", 49: "x", 50: "y", 51: "z", 52: "0", 53: "1", 54: "2", 55: "3", 56: "4", 57: "5", 58: "6", 59: "7", 60: "8", 61: "9"}

type shorter struct {
	Url      string   `json:"url"`
	UrlList  []string `json:"urllist"`
	Number   int      `json:"number`
	ShortUrl []string `json:"shorturl"`
	HashCode []string `json:"hashcode"`
}

func CreateShortUrlHandler(c *gin.Context) {
	var data shorter
	c.BindJSON(&data)

	url := data.Url
	urlList := data.UrlList
	number := data.Number
	shortUrlPrefix := "http://uniqs.cc"

	if url == "" && len(urlList) == 0 || number == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Data missing parameter"})
	}

	urlList = append(urlList, url)

	resultCode := make(chan resultForCreation)
	for _, v := range urlList {
		go func(lUrl string) {
			shorturl, err := GenShortUrl(lUrl)
			if err != nil {
				logs.Error("gen shortUrl failed, error: " + err.Error())
				resultCode <- result{lUrl, err.Error()}
			} else {
				logs.Info("[create]: " + lUrl + " => " + shortUrl)
				resultCode <- result{lUrl, shortUrl}
			}

		}(v)
	}

	var results = make(map[string]interface{})
	var count = 1
	for {
		res := <-resultCode
		results[res.url] = shortUrlPrefix + res.code
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
	result := new(Short_Url)
	has, err := ORM_ShortUrl.Where("hashcode=?", urlmd5).Get(result)

	if err != nil {
		return "", err
	}

	if has == true {
		return result.shorturl, nil
	} else {

		newone := &Short_Url{
			Url:      Url,
			HashCode: urlmd5,
		}
		_, err := ORM_ShortUrl.Insert(newone)

		if err != nil {
			return "", err
		}

		id := newone.Id
		if id == 0 {
			return "", errors.New("get id failed")
		}

		shorturl := TransToCode(id)

		if shorturl == "" {
			return "", errors.New("gen code failed")
		}
		newone.Shorturl = shorturl
		_, err = ORM_ShortUrl.Id(id).Update(newone)

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
func TransToCode(id int) string {
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
