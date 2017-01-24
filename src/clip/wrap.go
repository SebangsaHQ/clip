package clip
import (
    "strings"
    "fmt"
    "net/url"
    "net/http"
    "gopkg.in/iconv.v1"
)


//Ngetes URL, kalo ga bisa di get, url dikosongin kemudian bila bisa diget akan melakukan download
//download the url and save it to local dir, Hasilnya url akan direplace menjadi url local
//param url : url gambar, image : true apabila file adalah gambar (imageThumb, Icon), wrapped : apabila url telah berhasil divalidasi
//maka diset true, agar tidak diganti gambar lain
func (me *Grab) wrapURL(Url *string, image bool, wrapped *bool) {
    if !*wrapped {
        originalURL := *Url
        //setup url with protocol first if haven't protocol
        if !strings.Contains(*Url, "://") && *Url != "" {
            if strings.Contains(*Url, "//") {
                *Url = fmt.Sprintf("http:%s", *Url)
            }else {
                rootUrl, errParse := url.Parse(me.Url)
                if errParse == nil {
                    *Url = fmt.Sprintf("http://%s%s", rootUrl.Host, *Url)
                    _, err := http.Head(*Url)
                    if err != nil {
                        urlPaths := strings.Split(rootUrl.Path, "/")
                        urlPaths = urlPaths[:len(urlPaths)-1]
                        urlPath := strings.Join(urlPaths, "/")
                        *Url = fmt.Sprintf("http://%s%s/%s", rootUrl.Host, urlPath, originalURL)
                    }
                }else {
                    *Url = ""
                }
            }
        }

        Logr.Infof("wrap URL : %s", *Url)

        clientH := &http.Client{}
        req, err := http.NewRequest("GET", *Url, nil)
        if err != nil {
            Logr.Errorf("Error when build NewRequest, :: %s", err.Error())
            LogLine()
            return
        }

        userAgent := "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/46.0.2490.80 Safari/537.36"
        req.Header.Set("User-Agent", userAgent)
        req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
        resp, err := clientH.Do(req)
        if err != nil {
            Logr.Info("[wrapURL] error on getting url, then url omitted :: %s", err.Error())
            *Url = ""
            return
        }

        defer resp.Body.Close()

        if resp.StatusCode != 200 {
            LogLine()
            *Url = ""
            return
        }
        cLength := resp.ContentLength
        if err != nil {
            *Url = ""
            LogLine()
            return
        }

        if cLength == 0 {
            *Url = ""
            LogLine()
            return
        }

        var tipe string
        if len(resp.Header["Content-Type"]) > 0 {
            tipe = resp.Header["Content-Type"][0]
            if image {
                if !strings.Contains(strings.ToLower(tipe), "image") {
                    LogLine()
                    *Url = ""
                    return
                }
            }
        }

        *wrapped = true
    }
}

//konversi tipe konten sesuai dengan tipe-tipe yang ada di embed.ly
//see http://internalgit.ojodowo.com/thoriq/clip/issues/5
func (me *Grab) wrapContentType() {
    types := []string{
        "html",
        "xml",
        "text",
        "image",
        "video",
        "audio",
        "json",
        "powerpoint",
    }

    for _, t := range types {
        if strings.Contains(me.ContentType, t) {
            me.ContentType = t
            if t == "powerpoint" {
                me.ContentType = "ppt"
            }
            return
        }
    }
    me.ContentType = "link"
}


//konversi xml apakah rss atau atom
func (me *Grab) wrapXml(html []byte) {
    if isAtom(html) {
        me.ContentType = "atom"
        return
    }else if isRSS(html) {
        me.ContentType = "rss"
        return
    }
}

//Wrap Any Encoding to UTF-8
func (me *Grab) wrapEncoding(from string) {
    to := "utf8"

    cd, err := iconv.Open(to, from)
    if err != nil {
        Logr.Errorf("[wrapEncoding] iconv open fail. Type = %s", from)
        LogLine()
        return
    }
    defer cd.Close()

    me.Title = cd.ConvString(me.Title)
    me.Description = cd.ConvString(me.Description)
}
