package main

import (
    "fmt"
    "log"
    "encoding/json"
    "net/http"
)

type error_reply struct {
    Error string `json: "error"`
    Message string `json: "message"`
}

func print_error(w * http.ResponseWriter, msg string) {
    reply_ := &error_reply{
        Error: "Illegal inputs",
        Message: msg}

    reply, _ := json.Marshal(reply_)
    fmt.Fprintf((*w), string(reply))
}

func handle_request(w http.ResponseWriter, r * http.Request){
    service_name, ok := r.URL.Query()["service"]

    if !ok {
        print_error(&w, "Please provide service name as a GET parameter")
    } else {
        config := config.New() // getting the config from memory
        portno, host, err := config.GetData(service_name[0])

        if err != nil {
            print_error(&w, err.Error())
        } else {
            reply, _ := json.Marshal(
                struct{
                    Host string `json:"host"`
                    Port string `json:"port"`
                }{
                    host,
                    portno,
                },
            )
            fmt.Fprintf(w, string(reply))
        }
    }
}

func main() {
    //config := config.New() // reading the config file
    // and saving it to memory

    //portno, _, err := config.GetData("db")

    //if err != nil {
    //    panic(err)
    //}

    http.HandleFunc("/config", handle_request)
    log.Fatal(http.ListenAndServe(":" + utils.Portno("db"), nil))
}
