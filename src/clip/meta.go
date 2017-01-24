package clip

import (
    "io"

    "golang.org/x/net/html"
    "fmt"
    "strings"
)

//data of meta key = meta name, value = content
type MetaData map[string]string

//meta and namespace used for channel
type MetaSpace struct {
    MetaData  MetaData
    NameSpace string
}

// MapMeta function extracts metadata from a html document.
// If no relevant metadata is found the result will be empty.
// The input is assumed to be UTF-8 encoded.
// About Namespace? insert exp : "og:" , "twitter:", "custom" >> unstructured
// More info about og at http://ogp.me and twitter at https://dev.twitter.com/cards/overview
func MapMeta(doc io.Reader, NameSpace string, ch chan MetaSpace) {
    md := make(MetaData)
    tags := MetaSpace{
        MetaData: md,
        NameSpace : NameSpace,
    }
    z := html.NewTokenizer(doc)

    for {
        tt := z.Next()
        if tt == html.ErrorToken {
            break
        }

        t := z.Token()

        if NameSpace == "custom" {
            if t.Type == html.EndTagToken && t.Data == "body" {
                break
            }
            //title
            if t.Data == "title" {
                tipe := z.Next()
                if tipe == html.TextToken {
                    title := string(z.Text())
                    if v, ok := tags.MetaData["title"]; ok {
                        tags.MetaData["title"] = fmt.Sprint(v, title)
                    }else {
                        tags.MetaData["title"] = title
                    }
                }
            }

            //published date
            if t.Data == "meta" {
                var prop, cont string
                for _, a := range t.Attr {
                    switch a.Key {
                    case "name", "property", "itemprop":
                        prop = a.Val
                    case "content":
                        cont = a.Val

                    //Charset
                    case "http-equiv":
                        if strings.HasPrefix(a.Val, "Content-Type"){
                            var charset, cType string
                            for _, va := range t.Attr {
                                if va.Key == "content"{
                                    cType = va.Val
                                    break
                                }
                            }
                            cType = strings.Replace(cType, " ", "", -1)
                            trimmedCT := strings.Split(cType, ";")

                            for _, value := range trimmedCT {
                                if strings.HasPrefix(value, "charset"){
                                    charset = strings.Replace(value, "charset=", "", -1)
                                }
                            }
                            tags.MetaData["charset"] = strings.ToLower(charset)
                            Logr.Debugf("METANI :: %s", tags.MetaData["charset"])
                        }

                    }
                }

                if strings.HasPrefix(prop, "description") {
                    tags.MetaData[prop] = cont
                }

                //publishedDate
                if !strings.Contains(strings.ToLower(prop), "publisher") {
                    if strings.Contains(strings.ToLower(prop), "publish") && cont != "" {
                        tags.MetaData["published"] = cont
                    }
                }
            }

            //favicon
            //<link rel="shortcut icon" href="favicon.ico">
            if t.Data == "link" {
                var prop, cont string
                for _, a := range t.Attr {
                    switch a.Key {
                    case "rel":
                        prop = a.Val
                    case "href":
                        cont = a.Val
                    }
                }

                if strings.HasSuffix(strings.ToLower(prop), "icon") && cont != "" {
                    if tags.MetaData["favicon"] != "" {
                        if !strings.HasSuffix(strings.ToLower(prop), "mask-icon") {
                            tags.MetaData["favicon"] = cont
                        }
                    }else {
                        tags.MetaData["favicon"] = cont
                    }
                }
            }

            //get first img
            //jika belum pernah dapat gambar
            if _, ok := tags.MetaData["first_img"]; !ok {
                if t.Data == "img" {
                    var src string
                    for _, a := range t.Attr {
                        switch a.Key {
                        case "src":
                            src = a.Val
                        }
                        if IsWebImage(src) {
                            tags.MetaData["first_img"] = src
                            tags.MetaData["first_img"] = src
                        }
                    }
                }
            }
        }else {
            if t.Type == html.EndTagToken && t.Data == "head" {
                break
            }
            if t.Data == "meta" {
                var prop, cont string
                for _, a := range t.Attr {
                    switch a.Key {
                    case "name", "property":
                        prop = a.Val
                    case "content", "value":
                        cont = a.Val
                    }
                }

                if strings.HasPrefix(prop, NameSpace) && cont != "" {
                    TrimmedProp := prop[len(NameSpace):]
                    tags.MetaData[TrimmedProp] = cont
                }
            }
        }

    }
    ch <- tags
}

func isRSS(html []byte) bool {
    if strings.Contains(string(html), "<rss ") {
        return true
    }
    return false
}

func isAtom(html []byte) bool {
    doc := string(html)
    if strings.Contains(doc, "xmlns:atom") && strings.Contains(strings.ToLower(doc), "/atom\">") {
        return true
    }
    return false
}

//if string has at least one of args substring, it will be true
func containsOf(src string, args ...string) bool {
    for _, arg := range args {
        if strings.Contains(src, arg) {
            return true
        }
    }
    return false
}

//is it containsOf(src, ".jpg", ".png", ".jpeg", ".gif") ?
func IsWebImage(src string) bool {
    return containsOf(src, ".jpg", ".png", ".jpeg", ".gif")
}
