package clip

import (
	"os"
	"strings"
	"rsc.io/pdf"
	"io/ioutil"
	"net/http"
	"fmt"
	"errors"
)

type PdfHandler struct {
	// Put field definition here
	Handler
	Reader *pdf.Reader
}

func NewPdfHandler() *PdfHandler {
	return &PdfHandler{}
}

func (h PdfHandler) Test(grab *Grab, contentType string) bool {
	return strings.HasPrefix(contentType, "pdf")
}

func (h *PdfHandler) Handle(me *Grab) (err error) {
	Logr.Info("Serving from PDF handler")
	
	// This copy paste from html_handler, maybe we will wrap this code
	// into new function
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

	// create temp file from bytes array
	tmpfile, err := ioutil.TempFile("", "grabber")
	if err != nil {
		return err
	}
	defer os.Remove(tmpfile.Name())

	// write bytes array to tmp file
	if _, err := tmpfile.Write(res); err != nil {
		return err
	}

	// read the file
	reader, err := pdf.Open(tmpfile.Name())
	if err != nil {
		return err
	}

	h.Reader = reader

	// This is not real extraction
	// The real extraction needs more work
	if h.Reader.NumPage() > 0 {
		page1 := h.Reader.Page(1)
		resource := page1.Content().Text
		
		var text = ""
		for idx, t := range resource {
			text += t.S
			if idx > 100 {
				break
			}
		}

		me.Description = text
		
	}

	return nil
}