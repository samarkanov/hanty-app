package main

import (
    "fmt"
    "log"
    "net/http"
    "encoding/json"
    "khantydb/model"
    "github.com/samarkanov/khanty-app/utils"
    "github.com/gorilla/mux"
)

func print_success(w * http.ResponseWriter) {
    res_, _ := json.Marshal(struct{
        Success bool `json:"success"`
    }{
        true,
    })

    fmt.Fprintf(*w, string(res_))
}

func handle_get(w http.ResponseWriter, r * http.Request) {
    db := model.Cursor()
    vars := mux.Vars(r)
    client_id  := vars["client_id"]
    topic  := vars["topic"]

    query := &model.Query {
        Client: client_id,
        Topic: topic,
    }

    res := db.Get(query)
    fmt.Fprintf(w, res)
}

func handle_delete(w http.ResponseWriter, r * http.Request) {
    db := model.Cursor()
    vars := mux.Vars(r)
    client_id  := vars["client_id"]
    topic := vars["topic"]
    entry := vars["entry"]

    query := &model.Query {
        Client: client_id,
        Topic: topic,
        Entry: entry,
    }

    db.Delete(query)
    print_success(&w)
}

func handle_post(w http.ResponseWriter, r * http.Request) {
    db := model.Cursor()
    vars := mux.Vars(r)
    client_id  := vars["client_id"]
    topic := r.FormValue("topic")
    value := r.FormValue("value")

    if len(topic) > 0 && len(value) > 0 {
        query := &model.Query {
            Client: client_id,
            Topic: topic,
            Entry: value,
        }
        db.Set(query)
    }
    print_success(&w)
}

func main() {
    r := mux.NewRouter()

    // GET requests
    r.HandleFunc("/{client_id}", handle_get).Methods("GET")
    r.HandleFunc("/{client_id}/{topic}", handle_get).Methods("GET")

    // DELETE requests
    r.HandleFunc("/{client_id}", handle_delete).Methods("DELETE")
    r.HandleFunc("/{client_id}/{topic}", handle_delete).Methods("DELETE")
    r.HandleFunc("/{client_id}/{topic}/{entry}", handle_delete).Methods("DELETE")

    // POST requests
    r.HandleFunc("/{client_id}", handle_post).Methods("POST")

    http.Handle("/", r)
    log.Fatal(http.ListenAndServe(":"+utils.Portno("db"), nil))
}
