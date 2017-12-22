package clip

import (
    "net/http"
    "io/ioutil"
    "bytes"
    "strings"
    "net/url"
    "fmt"
    "errors"
)

//metaKind
const (
    og int = 1
    twitter int = 2
    custom int = 3
)

type Grabber interface {
    GrabDefault() (err error)
    GrabCURL() (err error)
    GetGrab() Grab

    getHTML() (res []byte, err error)
    getProperty(map[string]MetaData)
    wrapURL(Url *string, image bool, wrapped *bool)
    wrapContentType()
    wrapXml(html []byte)
    wrapEncoding(from string)
}

type Grab struct {
                                                                              //Core
    Url                 string  `json:"url"`
                                                                              //required
    Title               string  `json:"title"`
    Description         string  `json:"description"`
    ContentType         string  `json:"content_type"`       //html, image, video, audio, pdf, txt
    SiteName            string  `json:"site_name"`

                                                                              //Optional
    ImageThumbUrl       string  `json:"image_thumb"`
    ImageWidth          string  `json:"image_width"`
    ImageHeight         string  `json:"image_height"`
    FaviconUrl          string  `json:"favicon_url"`
    MediaType           string  `json:"media_type"`
    MediaUrl            string  `json:"media_url"`
    MediaWidth          string  `json:"media_width"` //kalau tipe image, audio, video, dokumen dll
    MediaHeight         string  `json:"media_height"`
    PublishedDate       string  `json:"published_date"`

    Status              int     `json:"status,omitempty"`

    ImageThumbUrlStatus bool    `json:"-"`
    FaviconUrlStatus    bool    `json:"-"`
    FinalURL            string  `json:"-"`  //bila terjadi redirect, maka url akhir akan disimpan disini.
    Host                string  `json:"-"`
    Charset             string  `json:"-"`

    Client              *http.Client  //Client http
}

func NewGrab(url string) Grabber {
    grab := &Grab{
        Url : url,
    }
    return grab
}

func (me Grab) GetGrab() Grab {
    return me
}

func (me *Grab) GrabCURL() (err error)  {
    Logr.Infof("grab Using CURL : %s", me.Url)

    if err != nil {
        return
    }

    curlRes, err := Curl(me.Url)
    if err != nil {
        return
    }

    me.Status = curlRes.Status
    me.ContentType = curlRes.ContentType

    me.wrapContentType()

    if strings.HasPrefix(me.ContentType, "html") || strings.HasPrefix(me.ContentType, "xml") {
        if me.Status != 200 && me.Status != 404 {
            return errors.New(fmt.Sprintf("error. Response code is %d, (error must be 200 or 404)", me.Status))
        }


        if strings.HasPrefix(me.ContentType, "xml") {
            me.wrapXml(curlRes.Content)
        }

        meta := getMeta(curlRes.Content)
        me.getProperty(meta)

        //Specific host
        if StringInSlice(me.Host, HostSpecific) {
            MapMetaSpecific(bytes.NewBuffer(curlRes.Content), me)
        }

        //Wrap Encoding
        if me.Charset != "" && me.Charset != "utf-8" {
            me.wrapEncoding(me.Charset)
        }

    }

    return
}

// It will grab the properties of web to Grab struct
// it's flow? see http://internalgit.ojodowo.com/thoriq/clip/issues/3
// GrabDefault is using net http
func (me *Grab) GrabDefault() (err error) {
    Logr.Infof("grab URL : %s", me.Url)

    me.wrapSIDURL()

    clientH := &http.Client{}
    req, err := http.NewRequest("GET", me.Url, nil)
    if err != nil {
        Logr.Errorf("Error when build NewRequest, :: %s", err.Error())
        LogLine()
        return err
    }

    userAgent := "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/46.0.2490.80 Safari/537.36"
    req.Header.Set("User-Agent", userAgent)
    req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,*/*;q=0.8")
    resp, err := clientH.Do(req)
    if err != nil {
        me.Status = 500
        Logr.Errorf("Error when getting url, :: %s", err.Error())
        LogLine()
        return err
    }
    defer resp.Body.Close()

    Logr.Infof("http GET resp Status : %s", resp.Status)
    Logr.Infof("http final URL : %s", resp.Request.URL.String())
    me.FinalURL = resp.Request.URL.String()
    erv := ValidateURL(&me.FinalURL)
    if erv != nil {
        LogLine()
        return errors.New("blocked")
    }

    _, ok := resp.Header["Content-Type"]
    if !ok {
        me.Status = 500
        Logr.Error("Resp content type on header is not ok")
        LogLine()
        return err
    }

    me.ContentType = resp.Header["Content-Type"][0]

    me.wrapContentType()

    return SelectHandler(me, me.ContentType)
}

//Get All Meta concurrently
func getMeta(res []byte) map[string]MetaData {
    metaInspect := []string{
        "og:",
        "twitter:",
        "custom",
    }

    ch := make(chan MetaSpace)
    defer close(ch)
    mergedMaps := make(map[string]MetaData)

    for _, inspect := range metaInspect {
        go MapMeta(bytes.NewBuffer(res), inspect, ch)
    }

    //Merge maps
    for i := 0; i < len(metaInspect); i++ {
        ms := <-ch
        nSpace := ms.NameSpace
        mData := ms.MetaData
        mergedMaps[nSpace] = mData
    }

    return mergedMaps
}

// get property of meta
// informasi urut dari og > twitter > custom > file
// bila ada informasi yang kosong maka akan dilengapi, dan yang sudah ada tidak diganti
// misalnya bila pada og tidak ada title makan akan di lengkapi title yang ada di twitter
func (me *Grab) getProperty(metaData map[string]MetaData) {
    url, errParse := url.Parse(me.Url)
    if errParse != nil {
        return
    }

    me.Host          = url.Host
    Logr.Infof("url.host from url.Parse %s", me.Host)

    me.Title         = metaData["og:"]["title"]
    me.Description   = metaData["og:"]["description"]

    me.SiteName      = metaData["og:"]["site_name"]

    me.ImageThumbUrl = metaData["og:"]["image"]

    me.ImageWidth    = metaData["og:"]["image:width"]
    me.ImageHeight   = metaData["og:"]["image:height"]

    me.MediaType     = metaData["og:"]["type"]
    if strings.HasPrefix(me.MediaType, "video") || strings.HasSuffix(me.MediaType, "video") {
        me.MediaUrl      = metaData["og:"]["video:url"]
        me.MediaWidth    = metaData["og:"]["video:width"]
        me.MediaHeight   = metaData["og:"]["video:height"]
    }else {
        me.MediaUrl      = metaData["og:"][""]
    }

    if me.ImageThumbUrl != "" {
        me.wrapURL(&me.ImageThumbUrl, true, &me.ImageThumbUrlStatus)
    }

    if me.Title == "" {
        me.Title     = metaData["twitter:"]["title"]
    }
    if me.Description == "" {
        me.Description   = metaData["twitter:"]["description"]
    }
    if me.SiteName == "" {
        me.SiteName      = metaData["twitter:"]["site"]
    }
    if me.ImageThumbUrl == "" {
        me.ImageThumbUrl = metaData["twitter:"]["image"]
        if me.ImageThumbUrl != "" {
            me.wrapURL(&me.ImageThumbUrl, true, &me.ImageThumbUrlStatus)
        }
    }
    if me.MediaType == "" {
        me.MediaType     = metaData["twitter:"]["card"]
    }
    if me.MediaUrl == "" {
        me.MediaUrl      = metaData["twitter:"]["player"]
        me.MediaHeight   = metaData["twitter:"]["player:height"]
        me.MediaWidth    = metaData["twitter:"]["player:width"]
    }

    if me.Title == "" {
        me.Title         = metaData["custom"]["title"]
    }
    if me.Description == "" {
        me.Description   = metaData["custom"]["description"]
    }

    if me.SiteName == "" {
        if errParse == nil {
            me.SiteName  = url.Host
        }
    }

    if me.ImageThumbUrl == "" {
        me.ImageThumbUrl = metaData["custom"]["first_img"]
        if me.ImageThumbUrl != "" {
            me.wrapURL(&me.ImageThumbUrl, true, &me.ImageThumbUrlStatus)
        }
    }

    //Get favicon
    me.FaviconUrl    = metaData["custom"]["favicon"]
    if me.FaviconUrl != "" {
        me.wrapURL(&me.FaviconUrl, true, &me.FaviconUrlStatus)
    }

    //try to catch from basic url
    if me.FaviconUrl == "" {
        me.FaviconUrl = fmt.Sprintf("http://%s/favicon.ico", url.Host)
        me.wrapURL(&me.FaviconUrl, true, &me.FaviconUrlStatus)
    }

    //get published date
    me.PublishedDate = metaData["custom"]["published"]


    //get charset
    me.Charset = metaData["custom"]["charset"]

    if me.MediaUrl != "" {
        dum := false
        me.wrapURL(&me.MediaUrl, false, &dum)
    }

    //Get image size manually
    if me.ImageThumbUrl != "" {
        if me.ImageWidth == "" && me.ImageHeight == "" {
            width, height := GetImageSize(me.ImageThumbUrl)
            if width != "0" && height != "0" {
                me.ImageWidth = width
                me.ImageHeight = height
            }
        }
    }
}

//Get all HTML og grabbed URL
func (me *Grab) getHTML() (res []byte, err error) {
    resp, err := http.Get(me.Url)
    if err != nil {
        Logr.Error("Error when getting url")
        LogLine()
        return
    }
    defer resp.Body.Close()

    res, err = ioutil.ReadAll(resp.Body)
    if err != nil {
        Logr.Error("error pas read body response")
        LogLine()
        return
    }

    return
}


//escape URL
func EscapeURL(oldUrl string) string {
    Url, err := url.Parse(oldUrl)
    if err != nil {
        Logr.Errorf("Failed escape URL : %s", oldUrl)
        return oldUrl
    }
    return Url.String()
}
