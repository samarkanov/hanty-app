package main

import (
    "fmt"
    "log"
    "os"
    "net/http"
    "net/url"
    "io/ioutil"
    "encoding/json"
    "github.com/samarkanov/khanty-app/utils"
    "github.com/gorilla/mux"
)

type Client struct {
    name string
}

func db_config() (string, string){
    db_portno := utils.Portno("db")
    db_host := utils.Host("db")
    return db_portno, db_host
}

func (client * Client) get_state() (* http.Response) {
    db_portno, db_host := db_config()
    url_path := fmt.Sprintf("%s:%s/%s", db_host, db_portno, client.name)

    http_client := &http.Client{}
    request, _ := http.NewRequest("GET", url_path, nil)

    resp, _ := http_client.Do(request)
    return resp
}

func (client * Client) handle_get(w http.ResponseWriter, r * http.Request) {
    resp := client.get_state()

    if resp.StatusCode == http.StatusOK {
        body, _ := ioutil.ReadAll(resp.Body)
        fmt.Fprintf(w, string(body))
    }
}

func (client * Client) handle_get_topic(w http.ResponseWriter, r * http.Request) {
    resp := client.get_state()
    vars := mux.Vars(r)
    topic  := vars["topic"]

    if resp.StatusCode == http.StatusOK {
        body, _ := ioutil.ReadAll(resp.Body)
        var reply []struct{
            Topic string
            Data  []string
        }
        json.Unmarshal(body, &reply)

        for _, dbitem := range reply {
            if dbitem.Topic == topic {
                fmt.Fprintf(w, dbitem.Data[len(dbitem.Data)-1])
            }
        }
    }
}

func (client * Client) handle_update(w http.ResponseWriter, r * http.Request) {
    db_portno, db_host := db_config()
    topic := r.FormValue("topic")
    value := r.FormValue("value")

    if len(topic) > 0 && len(value) > 0 {
        url_path := fmt.Sprintf("%s:%s/%s", db_host, db_portno, client.name)

        form := url.Values{
            "topic": {topic},
            "value": {value},
        }

        _, err := http.PostForm(url_path, form)

        if err != nil {
            fmt.Fprintf(w, utils.ReplyError(err.Error()))
        }
    }
}

func store_state(client_name string, property string, value string) {
    db_portno, db_host := db_config()
    url_path := fmt.Sprintf("%s:%s/%s", db_host, db_portno, client_name)

    form := url.Values{
        "topic": {property},
        "value": {value},
    }

    http.PostForm(url_path, form)
}

func (client * Client) init() {
    store_state(client.name, "ChangeColor", "#5e72e4")
    store_state(client.name, "SendMessage", "#SayHi")
}


func main() {
    r := mux.NewRouter()

    // client name:
    if len(os.Args) < 2 {
        fmt.Println("usage: fo run client.go <clientName>")
        return
    }
    client_name := os.Args[1]
    fmt.Println("starting on port " + utils.Portno(client_name))

    // initializing client
    client := &Client{name: client_name}
    client.init()

    // GET state
    r.HandleFunc("/", client.handle_get).Methods("GET")

    // GET current topic
    r.HandleFunc("/{topic}", client.handle_get_topic).Methods("GET")

    // POST update
    r.HandleFunc("/update", client.handle_update).Methods("POST")

    http.Handle("/", r)
    log.Fatal(http.ListenAndServe(":" + utils.Portno(client_name), nil))
}
