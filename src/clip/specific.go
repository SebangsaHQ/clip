//Specific? maksudnya adalah ketika suatu web yang di grab ga' bisa dengan cara umum. Dia perlu perlakuan khusus
//Nah kamu dapat menambahkan perlakuan khususnya di sini
//Apa yang perlu dilakukan untuk menambah perlakuan khusus?
//1. Tambahkan nama host nya pada HostSpecific ... (Nama host bisa didapatkan dari fungsi url.Parse() yaa)
//2. Tambahkan nama host pada list switch grab.Host
//3. Tambahkan fungsi-fungsi khususnya

package clip
import (
    "io"
    "golang.org/x/net/html"
    "strings"
    "fmt"
)

//1.
//List of specific host URL
var HostSpecific = []string{
    "krjogja.com",
    "en.wikipedia.org",
}

// MapMeta function extracts metadata from a html document.
// specific for one host only
func MapMetaSpecific(doc io.Reader, grab *Grab) {
    z := html.NewTokenizer(doc)
    complete := false
    for {
        tt := z.Next()
        if tt == html.ErrorToken {
            break
        }

        t := z.Token()

        if t.Type == html.EndTagToken && t.Data == "body" {
            break
        }

        //2.
        switch grab.Host {
        case "krjogja.com":
            complete = krjogja(t, z, grab)
        case "en.wikipedia.org":
            complete = enwikipediaorg(t, z, grab)
        }

        //jika element yang dicari sudah ketemu maka iterasi dihentikan
        if complete {
            break
        }

    }
}


//3. FUNGSI - FUNGSI KHUSUS

//Metani krjogja.com
//look for description
func krjogja(t html.Token, z *html.Tokenizer, grab *Grab) bool {
    if t.Data == "div" {
        var prop string
        for _, a := range t.Attr {
            switch a.Key {
            case "class":
                prop = a.Val
            }
        }
        if strings.Contains(prop, "detail-content") {
            for {
                tt := z.Next()
                if tt == html.ErrorToken {
                    break
                }
                t := z.Token()
                if t.Type == html.EndTagToken && t.Data == "div" {
                    break
                }

                if tt == html.TextToken {
                    data := t.Data
                    grab.Description = fmt.Sprintf("%s%s", grab.Description, data)
                }
            }
            return true
        }
    }
    return false
}


//Metani en.wikipedia.org
//look for description
func enwikipediaorg(t html.Token, z *html.Tokenizer, grab *Grab) bool {
    if t.Data == "p" {

        for {
            tt := z.Next()
            if tt == html.ErrorToken {
                break
            }
            t := z.Token()
            if t.Type == html.EndTagToken && t.Data == "p" {
                break
            }

            if tt == html.TextToken {
                data := t.Data
                grab.Description = fmt.Sprintf("%s%s", grab.Description, data)
            }
        }
        return true
    }
    return false
}
