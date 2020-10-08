package main

import (
    "fmt"
    "log"
    "net/http"
    "net/url"
    "io/ioutil"
    "encoding/json"
    "github.com/samarkanov/khanty-app/utils"
    "github.com/gorilla/mux"
    "github.com/rs/cors"
)

func db_config() (string, string){
    db_portno := utils.Portno("db")
    db_host := utils.Host("db")
    return db_portno, db_host
}

func handle_subscribe(w http.ResponseWriter, r * http.Request) {

    db_portno, db_host := db_config()
    topic := r.FormValue("topic")
    client_port := r.FormValue("portno")

    fmt.Println("topic" + topic)
    fmt.Println("client_port" + client_port)

    if len(topic) > 0 && len(client_port) > 0 {
        url_path := fmt.Sprintf("%s:%s/master", db_host, db_portno)

        form := url.Values{
            "topic": {topic},
            "value": {client_port},
        }

        _, err := http.PostForm(url_path, form)

        if err != nil {
            fmt.Fprintf(w, utils.ReplyError(err.Error()))
        }
    }
}

func handle_unsubscribe(w http.ResponseWriter, r * http.Request) {

    db_portno, db_host := db_config()
    vars := mux.Vars(r)
    client_port  := vars["client_id"]
    topic := vars["topic"]

    if len(topic) > 0 && len(client_port) > 0 {
        url_path := fmt.Sprintf("%s:%s/master/%s/%s",
                    db_host, db_portno, topic, client_port)

        client := &http.Client{}
        request, err := http.NewRequest("DELETE", url_path, nil)

        if err != nil {
            fmt.Fprintf(w, utils.ReplyError(err.Error()))
            return
        }

        client.Do(request)
    }
}

func get_client_host(client_portno string) (string) {
    var client_host string
    config_host := utils.Host("config")
    config_portno := utils.Portno("config")

    reply, _ := http.Get(config_host + ":" + config_portno)
    defer reply.Body.Close()
    body, _ := ioutil.ReadAll(reply.Body)

    var config map[string]struct{
        Host string `json:"host"`
        Port string `json:"port"`
    }

    if err := json.Unmarshal(body, &config);  err != nil {
        panic(err)
    }

    for _, config_ob := range config {
        if config_ob.Port == client_portno {
            client_host = config_ob.Host
        }
    }

    return client_host
}

func notify_client(w * http.ResponseWriter, client_portno string, topic string, value string) {

    // sending request towards client
    if len(client_portno) > 0 && len(topic) > 0 {
        client_host := get_client_host(client_portno)

        if len(client_host) > 0 {
            fmt.Println(client_host)
            url_path := fmt.Sprintf("%s:%s/update", client_host, client_portno)

            form := url.Values{
                "topic": {topic},
                "value": {value},
            }

            _, err := http.PostForm(url_path, form)

            if err != nil {
                fmt.Fprintf(*w, utils.ReplyError(err.Error()))
            }
        }
    }
}

func handle_notify(w http.ResponseWriter, r * http.Request) {

    topic := r.FormValue("topic")
    message := r.FormValue("value")
    var list_clients []string

    if len(topic) > 0 && len(message) > 0 {

        subscribers := get_subscribers_for_topic(topic)

        for _, client := range subscribers {
            notify_client(&w, client, topic, message)
            list_clients = append(list_clients, client)
        }

        res, _ := json.Marshal(struct{
            Success bool `json:"success"`
            Message string `json:"message"`
            Clients []string `json:"clients"`
        }{
            true,
            "clients notified OK",
            list_clients,
        })
        fmt.Fprintf(w,  string(res))
    }
}

func get_subscribers_for_topic(topic string) ([]string){
    db_portno, db_host := db_config()
    var res []string

    url_path := fmt.Sprintf("%s:%s/master/%s", db_host, db_portno, topic)
    client := &http.Client{}
    request, _ := http.NewRequest("GET", url_path, nil)

    resp, _ := client.Do(request)

    if resp.StatusCode == http.StatusOK {
        body, _ := ioutil.ReadAll(resp.Body)
        json.Unmarshal(body, &res)
    }

    return res
}

func handle_get_subscribers(w http.ResponseWriter, r * http.Request) {

    vars := mux.Vars(r)
    topic := vars["topic"]

    subscribers := get_subscribers_for_topic(topic)

    res, _ := json.Marshal(struct{
        Topic string `json:"topic"`
        Subscribers []string `json:"subscribers"`
    }{
        topic,
        subscribers,
    })
    fmt.Fprintf(w,  string(res))
}


func main() {
    r := mux.NewRouter()

    c := cors.New(cors.Options{
        AllowedOrigins: []string{"http://develop.valenoq.com:*"},
        AllowedMethods: []string{"POST", "GET", "DELETE"},
        AllowCredentials: true,
    })

    handler := c.Handler(r)

    // POST subscribe
    r.HandleFunc("/subscribe", handle_subscribe).Methods("POST")

    // POST notify
    r.HandleFunc("/notify", handle_notify).Methods("POST")

    // DELETE unsubscribe
    r.HandleFunc("/unsubscribe/{topic}/{client_id}", handle_unsubscribe).Methods("DELETE")

    // GET all subscribers for a given topic
    r.HandleFunc("/subscribers/{topic}", handle_get_subscribers).Methods("GET")

    http.Handle("/", r)
    log.Fatal(http.ListenAndServe(":"+utils.Portno("master"), handler))
}
