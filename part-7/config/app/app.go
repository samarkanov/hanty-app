package main

import (
    "fmt"
    "log"
    "strings"
    "encoding/json"
    "net/http"
    "khanty-config/config"
    "github.com/samarkanov/khanty-app/utils"
)

type success_reply struct {
    Result string `json:"result"`
    Host string `json:"host"`
    Port string `json:"port"`
}

func reply_success(host string, portno string) string {
    reply, _ := json.Marshal(
        &success_reply{
            "OK",
            host,
            portno,
        })
    return string(reply)
}

func handle_request(w http.ResponseWriter, r * http.Request) {
    var reply string
    config := *config.New() // getting from memory

    if r.URL.Path == "/" {
        res, _ := json.Marshal(config)
        fmt.Fprintf(w, string(res))
    } else {
        // return data for a given service
        service := strings.Split(r.URL.Path, "/")[1]
        service = strings.Split(service, "/")[0]
        portno, host, err := config.GetData(service)
        if err != nil {
            reply = utils.ReplyError(err.Error())
        } else {
            reply = reply_success(host, portno)
        }
        fmt.Fprintf(w, reply)
    }
}

func main() {

    config := config.New() // reading the file and saving to memory
    portno, _, err := config.GetData("config")

    if err != nil {
        panic(err)
    }

    http.HandleFunc("/", handle_request)
    log.Fatal(http.ListenAndServe(":" + portno, nil))
}
