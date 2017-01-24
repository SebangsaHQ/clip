package clip

import (
    "errors"
    "strings"
    "fmt"
    "github.com/gin-gonic/gin"
)

func PostHandler(c *gin.Context) {
    resp := NewResponse()
    var dataJson Grab
    url := c.PostForm("url")
    err := ValidateURL(&url)
    if err != nil {
        resp.ErrorResponse(400, err.Error())
        resp.JSON(c)
        return
    }

    g := NewGrab(url)
    err = g.GrabDefault()
    if err != nil {
        dataJson = g.GetGrab()
        if err.Error() == "blocked" {
            resp.ErrorResponse(400, fmt.Sprintf("This Domain / Site is Blocked by clip. Blocked URL : %s", dataJson.FinalURL))
            resp.JSON(c)
            return
        }
        if err.Error() != "redigo: nil returned" {
            Logr.Infof("ERROR GrabDefault is : %s", err.Error())
            Logr.Infof("Data JSON When error is : %#v", dataJson)
            g = NewGrab(dataJson.Url)
            err = g.GrabCURL()
            if err != nil {
                if g.GetGrab().Status == 500 {
                    resp.ErrorResponse(500, err.Error())
                    resp.JSON(c)
                    return
                }
                resp.ErrorResponse(400, err.Error())
                resp.JSON(c)
                return
            }
        }
    }

    //Coba ulang untuk grab final URL apabila 200 tapi ndak dapat apa-apa
    dataJson = g.GetGrab()
    if dataJson.Title == "" {
        if dataJson.FinalURL != url && dataJson.FinalURL != "" {
            Logr.Debug("Try To Grab Final URL")

            g = NewGrab(dataJson.FinalURL)
            err = g.GrabDefault()
            if err != nil {
                if g.GetGrab().Status == 500 {
                    resp.ErrorResponse(500, err.Error())
                    resp.JSON(c)
                    return
                }
                resp.ErrorResponse(400, err.Error())
                resp.JSON(c)
                return
            }
        }
    }

    //omit status,, *not needed by client but needed by clip
    dataJson = g.GetGrab()
    dataJson.Status = 0
    dataJson.Url = url
    if dataJson.ContentType == "" {
        dataJson.ContentType = "error"
    }

    //escape URL
    dataJson.FaviconUrl = EscapeURL(dataJson.FaviconUrl)
    dataJson.ImageThumbUrl = EscapeURL(dataJson.ImageThumbUrl)
    dataJson.MediaUrl = EscapeURL(dataJson.MediaUrl)

    //send response to client
    resp.SetData(dataJson)
    resp.JSON(c)
}

func ValidateURL(rawUrl *string) error {
    //not empty
    if *rawUrl == "" {
        return errors.New("parameter invalid")
    }

    //has a right protocol?
    if strings.Contains(*rawUrl, "://") {
        if !strings.HasPrefix(*rawUrl, "http://") && !strings.HasPrefix(*rawUrl, "https://") {
            return errors.New("protocol invalid! must http or https")
        }
    } else {
        //append http protocol if don't have protocol
        *rawUrl = fmt.Sprintf("http://%s", *rawUrl)
    }

    return nil
}

