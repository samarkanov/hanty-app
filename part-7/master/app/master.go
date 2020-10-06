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

func notify_client(w * http.ResponseWriter, client_portno string, topic string, value string) {
    // TODO: send request towards client
    fmt.Fprintf(*w, "about to notify {client: %s, topic: %s, value: %s}", client_portno, topic, value)
}

func handle_notify(w http.ResponseWriter, r * http.Request) {
    db_portno, db_host := db_config()
    topic := r.FormValue("topic")
    message := r.FormValue("value")

    if len(topic) > 0 && len(message) > 0 {

        // retrieve all subscrubers for a given topic
        url_path := fmt.Sprintf("%s:%s/master/%s", db_host, db_portno, topic)
        client := &http.Client{}
        request, err := http.NewRequest("GET", url_path, nil)

        if err != nil {
            fmt.Fprintf(w, utils.ReplyError(err.Error()))
            return
        }

        resp, _ := client.Do(request)

        if resp.StatusCode == http.StatusOK {
            body, _ := ioutil.ReadAll(resp.Body)
            var reply []string
            json.Unmarshal(body, &reply)
            for _, client := range reply {
                notify_client(&w, client, topic, message)
            }
        }
    }
}

func main() {
    r := mux.NewRouter()

    // POST subscribe
    r.HandleFunc("/subscribe", handle_subscribe).Methods("POST")

    // POST notify
    r.HandleFunc("/notify", handle_notify).Methods("POST")

    // DELETE unsubscribe
    r.HandleFunc("/unsubscribe/{topic}/{client_id}", handle_unsubscribe).Methods("DELETE")

    http.Handle("/", r)
    log.Fatal(http.ListenAndServe(":"+utils.Portno("master"), nil))
}
