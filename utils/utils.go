package main

import (
    "net/http"
    "io/ioutil"
    // "fmt"
    "encoding/json"
)

var (
    CONFIG_PORT = "10002"
    CONFIG_HOST = "develop.valenoq.com"
)

type ErrorRreply struct {
    Result string `json:"result"`
    Error string `json:"error"`
    Message string `json:"message"`
}

func ReplyError(msg string) string {
    reply, _ := json.Marshal(
        &ErrorRreply{
            "NOK",
            "Illegal inputs",
            msg,
        })
    return string(reply)
}

func get_config(service string) map[string]string {
    reply, _ := http.Get("http://" + CONFIG_HOST + ":" + CONFIG_PORT + "/" + service)
    defer reply.Body.Close()
    body, _ := ioutil.ReadAll(reply.Body)

    var config_item map[string]string
    if err := json.Unmarshal(body, &config_item);  err != nil {
        panic(err)
    }

    if config_item["result"] == "NOK" {
        panic(config_item["message"])
    }

    return config_item
}

func Portno(service string) string {
    config := get_config(service)
    return config["port"]
}

func Host(service string) string {
    config := get_config(service)
    return config["host"]
}

func Get(service string, postfix string) []byte {
    config := get_config(service)
    reply, _ := http.Get(config["host"] + ":" + config["port"] + "/" + postfix)
    defer reply.Body.Close()
    body, _ := ioutil.ReadAll(reply.Body)
    return body
}
