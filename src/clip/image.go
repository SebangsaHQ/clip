package clip

import (
    "net/http"
    "image"
    _ "image/gif"
    _ "image/jpeg"
    _ "image/png"
    "fmt"
    "strconv"
    "os/exec"
    "os"
    "errors"
    "crypto/tls"
)

//Get size of image from specific url
func GetImageSize(url string) (height, width string) {
    // Ignore ssl
    tr := &http.Transport{
        TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
    }
    clientH := &http.Client{Transport: tr}

    req, err := http.NewRequest("GET", url, nil)
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
        Logr.Infof("[getImageSize] error on getting url, :: %s", err.Error())
        LogLine()
        return
    }

    defer resp.Body.Close()

    m, _, err := image.Decode(resp.Body)
    if err != nil {
        fmt.Println(url)
        Logr.Infof("[getImageSize] error on resp body, :: %s", err.Error())
        LogLine()
        return
    }
    g := m.Bounds()

    // Get height and width
    h := g.Dy()
    w := g.Dx()
    height = strconv.Itoa(h)
    width = strconv.Itoa(w)

    return
}

//Get size of image from specific path
func GetImageSizeLocal(path string) (height, width string) {
    f, err := os.Open(path)
    if err != nil {
        Logr.Infof("[GetImageSizeLocal] error when opening file, :: %s", err.Error())
        LogLine()
    }

    m, _, err := image.Decode(f)
    if err != nil {
        Logr.Infof("[GetImageSizeLocal] error on resp body, :: %s", err.Error())
        LogLine()
        return
    }
    g := m.Bounds()

    // Get height and width
    h := g.Dy()
    w := g.Dx()
    height = strconv.Itoa(h)
    width = strconv.Itoa(w)

    return
}

func ConvertToPNG(fileName, writeLoc, imgDirHost string) (newName string, erro error) {
    //convert using full path
    futureName := fmt.Sprintf("%s/%s%s", imgDirHost, fileName, ".png")
    source := fmt.Sprintf("%s[0]", writeLoc)
    cmd := exec.Command("convert", source, futureName)
    _, erro = cmd.Output()
    if erro != nil {
        Logr.Errorf("Convert ico to png error :: %s", erro.Error())
        LogLine()
        return
    }

    //check name.png
    if _, err := os.Stat(futureName); os.IsNotExist(err) {
        //not exists check name-0.png
        if _, err := os.Stat(fmt.Sprintf("%s/%s%s", imgDirHost, fileName, "-0.png")); os.IsNotExist(err) {
            return "", errors.New("ConvertToPNG fail to give name")
        }else {
            newName = fmt.Sprintf("%s%s", fileName, "-0.png")
            return newName, nil
        }
    }else {
        newName = fmt.Sprintf("%s%s", fileName, ".png")
        return newName, nil
    }

    return "", errors.New("ConvertToPNG fail to give name")
}
