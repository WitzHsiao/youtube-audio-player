package main

import (
    //"os/exec"
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
        //song := "/Users/witzhsiao/Downloads/Italobrothers - This Is Nightlife (Video Edit)Mp3World.Lt.mp3"
        //cmd := fmt.Sprintf("mplayer '%s'", song)
        fmt.Printf("%s %s\n", ytInfo.Author, ytInfo.Title)
        fmt.Println("CMD: ", cmd)
        //output, err := exec.Command("sh", "-c", cmd).Output()
        //if err != nil {
            //log.Println(err)
        //}
        //fmt.Println(string(output))
    }
}
