package main

import (
    "fmt"
    "log"
    "encoding/json"
    "net/http"
    "appconfig/config"
    "github.com/khanty-app/utils"
)

func handle_request(w http.ResponseWriter, r * http.Request) {
    config := *config.New() // getting from memory
    res, _ := json.Marshal(config)
    fmt.Fprintf(w, string(res))
}

func main() {

    res := say_hi()

    config := config.New() // reading the file
                           // and saving to memory

    portno, _, err := config.GetData("config")

    if err != nil {
        panic(err)
    }

    http.HandleFunc("/", handle_request)
    log.Fatal(http.ListenAndServe(":" + portno, nil))
}
