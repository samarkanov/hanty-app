package main

import (
    "fmt"
    "log"
    "net/http"
    "net/url"
    // "io/ioutil"
    // "encoding/json"
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

func handle_notify(w http.ResponseWriter, r * http.Request) {

}



func main() {
    r := mux.NewRouter()

    // POST subscribe
    r.HandleFunc("/subscribe", handle_subscribe).Methods("POST")

    // POST notify
    r.HandleFunc("/notify", handle_notify).Methods("POST")

    // DELETE unsubscribe
    r.HandleFunc("/unsubscribe/{client_id}/{topic}", handle_unsubscribe).Methods("DELETE")

    http.Handle("/", r)
    log.Fatal(http.ListenAndServe(":"+utils.Portno("master"), nil))
}
