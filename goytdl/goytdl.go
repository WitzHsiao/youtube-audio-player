package goytdl

import (
    "net/url"
    "fmt"
    "log"
    "net/http"
    "io/ioutil"
    "strings"
    "net/http/cookiejar"
    "compress/gzip"
    "regexp"
)

const YT_URL = "http://www.youtube.com/get_video_info?hl=en_US&el=detailpage&video_id="

type YTInfo struct {
    Title string
    Author string
    UseSign bool
    Urls []string
}

func getVideoId(urlString string) string {
    u, err := url.Parse(urlString)
    if err != nil {
        log.Fatal(err)
    }
    values := u.Query()
    return values.Get("v")
}

func GetYTInfo(urlString string) YTInfo {
    u := fmt.Sprintf("%s%s", YT_URL, getVideoId(urlString))
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
    //fmt.Println(info)
    var ytInfo YTInfo
    ytInfo.Title = info.Get("title")
    ytInfo.Author = info.Get("author")
    if info.Get("use_cipher_signature") == "True" {
        ytInfo.UseSign = true
    } else {
        ytInfo.UseSign = false
    }
    if ytInfo.UseSign {
        videoUrl := getVideoUrlFromAnotherSite(urlString)
        ytInfo.Urls = append(ytInfo.Urls, videoUrl)
        return ytInfo
    }

    streams_raw := strings.Split(info.Get("url_encoded_fmt_stream_map"), ",")
    for _, streams := range streams_raw {
        info, err = url.ParseQuery(streams)
        if err != nil {
            log.Fatal(err)
        }
        for _, url := range info["url"] {
            ytInfo.Urls = append(ytInfo.Urls, url)
        }
    }
    return ytInfo
}


func getVideoUrlFromAnotherSite(urlString string) string {
    downloaderUrl := "http://9xbuddy.com/youtube?url="
    cookieJar, _ := cookiejar.New(nil)
    client := &http.Client {
        Jar: cookieJar,
    }
    
    resp, err := client.Get(fmt.Sprintf("%s%s", downloaderUrl, urlString))
    if err != nil {
        log.Fatal(err)
    }
    defer resp.Body.Close()

    req, _ := http.NewRequest("POST", 
        "http://9xbuddy.com/includes/main-post.php", 
        strings.NewReader(fmt.Sprintf("url=%s", url.QueryEscape(urlString))))
    req.Header.Add("Accept", "*/*")
    req.Header.Add("X-Requested-With", "XMLHttpRequest")
    req.Header.Add("Accept-Encoding", "gzip,deflate,sdch")
    req.Header.Add("Accept-Language", "en-US,en;q=0.8,zh-TW;q=0.6,zh;q=0.4,zh-CN;q=0.2")
    req.Header.Add("Content-Type", "application/x-www-form-urlencoded; charset=UTF-8")
    req.Header.Add("Host", "9xbuddy.com")
    req.Header.Add("Origin", "http://9xbuddy.com")
    req.Header.Add("Referer", fmt.Sprintf("%s%s", downloaderUrl, urlString))
    req.Header.Add("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_9_3) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/36.0.1985.143 Safari/537.36")

    
    resp, err = client.Do(req)
    if err != nil {
        log.Fatal(err)
    }
    defer resp.Body.Close()
    reader, _ := gzip.NewReader(resp.Body)
    body, err := ioutil.ReadAll(reader)
    if err != nil {
        log.Fatal(err)
    }
    string_body := string(body)
    re := regexp.MustCompile("\"downtube.php.+?\"")
    urls := re.FindAllString(string_body, -1)
    return "http://9xbuddy.com/" + strings.Trim(urls[1], "\"")
}

