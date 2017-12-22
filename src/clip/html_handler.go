package clip

import (
	"strings"
	"net/http"
	"io/ioutil"
	"errors"
	"bytes"
	"fmt"
)

type HtmlHandler struct {
	// Put field definition here
	Handler
}

func NewHtmlHandler() *HtmlHandler {
	return &HtmlHandler{}
}

func (h HtmlHandler) Test(grab *Grab, contentType string) bool {
	return strings.HasPrefix(contentType, "html") || strings.HasPrefix(contentType, "xml")
}

func (h HtmlHandler) Handle(me *Grab) (err error) {
	Logr.Info("Serving from html handler")
	client := &http.Client{}
	req, err := http.NewRequest("GET", me.Url, nil)
	if err != nil {
		Logr.Errorf("Error when build NewRequest, :: %s", err.Error())
		LogLine()
		return err
	}

	userAgent := "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/46.0.2490.80 Safari/537.36"
	req.Header.Set("User-Agent", userAgent)

	resp, err := client.Do(req)
	if err != nil {
		me.Status = 500
		Logr.Errorf("Error when getting url, :: %s", err.Error())
		LogLine()
		return err
	}
	defer resp.Body.Close()

	me.Status      = resp.StatusCode
	Logr.Infof("http GET resp Status : %d", me.Status)

	if me.Status != 200 && me.Status != 404 {
		return errors.New(fmt.Sprintf("error. Response code is %d, (error must be 200 or 404)", me.Status))
	}

	res, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		Logr.Error("error When Read resp.Body")
		LogLine()
		return err
	}

	if strings.HasPrefix(me.ContentType, "xml") {
		me.wrapXml(res)
	}

	meta := getMeta(res)
	me.getProperty(meta)

	//Specific host
	if StringInSlice(me.Host, HostSpecific) {
		MapMetaSpecific(bytes.NewBuffer(res), me)
	}

	Logr.Debugf("Charset is : %s", me.Charset)

	//Wrap Encoding
	if me.Charset != "" && me.Charset != "utf-8" {
		me.wrapEncoding(me.Charset)
	}

	return nil
}