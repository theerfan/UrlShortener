package main

import (
	"github.com/theerfan/urlshortener/util"
	"github.com/labstack/echo"
	// "github.com/theerfan/UrlShortener/Server/"
	"github.com/theerfan/urlshortener/Server/Database/MongoDB"
	"net/http"
	"strings"
	"time"
	"bytes"
	"encoding/base64"
	// "hash/fnv"
	"fmt"
	"crypto/md5"
	"encoding/json"
	"io/ioutil"
)

const base = "ggo.gl/"
type ClientRequest util.ClientRequest
// type URL util.URL

const servAddr = "localhost:3333"

func extractURL(url string) *util.URL {
	split := strings.Split(url, "//")
	protocolLength := len(split[0])
	var orig string
	var protocol string
	if len(split) == 1 {
		protocol = ""
		orig = split[0]
	} else {
		protocol = split[0][:protocolLength]
		orig = split[1]
	}
	expTime := time.Now().Local().Add(time.Hour * time.Duration(1))
	rawURL := util.URL{
		Protocol: protocol,
		Orig: orig,
		ExpTime: expTime }
	return &rawURL
}

func shortenNumber(url *util.URL) *util.URL {
	count := mongodb.GiveCount() + 1
	url.Short = base + string(count)
	return url
}

func shortenHash(url *util.URL) *util.URL {
	md5 := md5.Sum([]byte(url.Orig))
	hash := base64.StdEncoding.EncodeToString(md5[:])
	url.Short = hash[:6]
	return url
}

func main() {
	mongodb.Init()
	e := echo.New()
	e.POST("/shortener", func (c echo.Context) error {
		var req ClientRequest
		request := c.Request()
		if request.Body != nil {
			decoder := json.NewDecoder(request.Body)
			err := decoder.Decode(&req)
			if err != nil {
				fmt.Println(err)
				return err
			}
			url := extractURL(req.Url)
			if req.Method == "hash" {
				shortenHash(url)
			} else if req.Method == "counter" {
				shortenNumber(url)
			} else {
				return c.String(http.StatusNotAcceptable, "Invalid method!")
			}
			mongodb.PutIntoDatabase(*url)
			return c.String(http.StatusOK, url.Short)
		}
		return c.String(http.StatusNotFound, "Empty request!")
	})

	e.GET("/shortener", func (c echo.Context) error {
		var bodyBytes []byte
		var body string
		request := c.Request()
		if request.Body != nil {
			bodyBytes, _ = (ioutil.ReadAll(request.Body))
			body = string(bodyBytes)
			resultURL := mongodb.GetFromDatabase(body)
			ans := resultURL.Protocol + "://" + resultURL.Orig
			if resultURL != nil {
				return c.String(http.StatusOK, ans)
			}
			return c.String(http.StatusNotFound, "Haven't been shortened!")
		}
		c.Request().Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
		return c.String(http.StatusNotFound, "Haven't been shortened!")
	})
	e.Start(servAddr)
}