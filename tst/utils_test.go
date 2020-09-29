package main

import (
    "testing"
    // "fmt"
    "encoding/json"
)

func TestPortnoPanicIllegalService(t * testing.T) {
    defer func() { recover() }()
    Portno("dbi")
    t.Errorf("did not panic")
}

func TestPortnoDB(t * testing.T) {
    portno := Portno("db")
    if portno != "10003" {
        t.Errorf("Wrong port for db")
    }
}

func TestPortnoCONFIG(t * testing.T) {
    portno := Portno("config")
    if portno != "10002" {
        t.Errorf("Wrong port for config")
    }
}

func TestPortnoTOPIC(t * testing.T) {
    portno := Portno("topic")
    if portno != "10001" {
        t.Errorf("Wrong port for topic")
    }
}

func TestHost(t * testing.T) {
    portno := Host("db")
    if portno != "http://develop.valenoq.com" {
        t.Errorf("Wrong host")
    }
}

func TestGetConfig(t * testing.T) {
    res := Get("config", "")

    type config struct {
        Host string
        Port string
    }

    var data map[string]config
    json.Unmarshal(res, &data)

    if data["config"].Host != "http://develop.valenoq.com" {
        t.Errorf("Wrong host in TestGet")
    }

    if data["config"].Port != "10002" {
        t.Errorf("Wrong port for config in TestGetConfig")
    }

    if data["db"].Host != "http://develop.valenoq.com" {
        t.Errorf("Wrong host in TestGet")
    }

    if data["db"].Port != "10003" {
        t.Errorf("Wrong port for config in TestGetConfig")
    }

    if data["topic"].Host != "http://develop.valenoq.com" {
        t.Errorf("Wrong host in TestGet")
    }

    if data["topic"].Port != "10001" {
        t.Errorf("Wrong port for config in TestGetConfig")
    }
}

func TestGetConfigWithParam(t * testing.T) {
    res := Get("config", "/db")

    type config struct {
        Result string
        Host string
        Port string
    }

    var data config
    json.Unmarshal(res, &data)

    if data.Result != "OK" {
        t.Errorf("Wrong result in TestGetConfigWithParam")
    }
    if data.Host != "http://develop.valenoq.com" {
        t.Errorf("Wrong host in TestGetConfigWithParam")
    }
    if data.Port != "10003" {
        t.Errorf("Wrong port in TestGetConfigWithParam")
    }

}
