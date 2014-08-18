package goytdl

import (
    "net/url"
    "fmt"
    "log"
    "net/http"
    "io/ioutil"
    //"strings"
)

const YT_URL = "http://www.youtube.com/get_video_info?hl=en_US&el=detailpage&video_id="

type YTInfo struct {
    Title string
    Author string
    Urls []string
}

func GetVideoId(urlString string) string {
    u, err := url.Parse(urlString)
    if err != nil {
        log.Fatal(err)
    }
    values := u.Query()
    return values.Get("v")
}

func GetYTInfo(urlString string) YTInfo {
    u := fmt.Sprintf("%s%s", YT_URL, GetVideoId(urlString))
    resp, err := http.Get(u)
    if err != nil {
        log.Fatal(err)
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        log.Fatal(err)
    }

    info, err := url.ParseQuery(string(body));
    if err != nil {
        log.Fatal(err)
    }
    ytInfo := YTInfo{}
    ytInfo.Title = info.Get("title")
    ytInfo.Author = info.Get("author")
    info, err = url.ParseQuery(info.Get("url_encoded_fmt_stream_map"))
    if err != nil {
        log.Fatal(err)
    }
    ytInfo.Urls = info["url"]
    
    return ytInfo
}
