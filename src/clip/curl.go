package clip

import (
    "os/exec"
    "strings"
)

type CURLInfo struct {
    Status int
    ContentType string
    Content []byte
}

func Curl(site string) (res CURLInfo, err error){
    //gain header info
    cmd := exec.Command("curl", "-s", "-I", "-L", site)
    headerInfo, err := cmd.Output()
    if err != nil {
        Logr.Infof("site : %s", site)
        Logr.Infof("headerInfo : %s", string(headerInfo))
        Logr.Errorf("curl exec error :: %s", err.Error())
        LogLine()
        return
    }

    perLines := strings.Split(string(headerInfo),"\n")
    for _, perLine := range perLines {
        if strings.HasPrefix(perLine, "HTTP") {
            perSpaces := strings.Split(perLine, " ")
            for _, perSpace := range perSpaces {
                if perSpace == "200" {
                    res.Status = 200
                }
            }

            if res.Status != 0 && res.ContentType != "" {
                break
            }
        }
        if strings.HasPrefix(perLine, "Content-Type") {
            perSpaces := strings.Split(perLine, " ")
            for _, perSpace := range perSpaces {
                if !strings.Contains(perSpace, "=") && !strings.Contains(perSpace, ":") {
                    res.ContentType =  strings.Replace(perSpace, ";", "", -1)
                }
            }

            if res.Status != 0 && res.ContentType != ""{
                break
            }
        }
    }

    //gain content
    cmd = exec.Command("curl","-L","-A", "\"Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/49.0.2623.87 Safari/537.36\"", "-c", "cookie.mw", site)
    res.Content, err = cmd.Output()
    if err != nil {
        Logr.Errorf("curl exec error :: %s", err.Error())
        LogLine()
        return
    }
    Logr.Infof("Status Code : %d", res.Status)
    Logr.Infof("Content Type : %s", res.ContentType)

    return
}
