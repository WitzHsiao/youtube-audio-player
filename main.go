package main

import (
    "os/exec"
    "fmt"
    "encoding/json"
    "io/ioutil"
    "log"

    "./goytdl"
)

type Track struct {
    Url string
}

func main() {
    file, err := ioutil.ReadFile("list.json")
    if err != nil {
        log.Fatal(err)
    }
    
    var tracks []Track
    json.Unmarshal(file, &tracks)

    for _, track := range tracks {
        ytInfo := goytdl.GetYTInfo(track.Url)
        cmd := fmt.Sprintf("mplayer -novideo '%s'", ytInfo.Urls[0])
        fmt.Printf("%s %s\n", ytInfo.Author, ytInfo.Title)
        //fmt.Println("CMD: ", cmd)
        output, err := exec.Command("sh", "-c", cmd).Output()
        if err != nil {
            log.Println(err)
        }
        fmt.Println(string(output))
    }
}
