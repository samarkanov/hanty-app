package main

import (
    "fmt"
    "log"
    "net/http"
    "github.com/samarkanov/khanty-app/utils"
    "github.com/gorilla/mux"
)

func handle_get(w http.ResponseWriter, r * http.Request) {
    vars := mux.Vars(r)
    client_id  := vars["client_id"]
    topic  := vars["topic"]

    if len(topic) == 0 {
        // dump all DB
        fmt.Fprint(w, "returning everything I have in DB for client = ", client_id)
    } else {
        // get data for specific topic
        fmt.Fprint(w, "returning data for client = ", client_id, " and topic = ", topic)
    }
}

func handle_delete(w http.ResponseWriter, r * http.Request) {
    vars := mux.Vars(r)
    client_id  := vars["client_id"]
    topic := vars["topic"]

     if len(topic) > 0 {
        // deleting everything for specified client and topic
        fmt.Fprint(w, "deleting all entries for client = ", client_id, " and topic = ", topic)
    } else {
        // delete all entries for client
        fmt.Fprint(w, "deleting all entries for client = ", client_id)
    }
}

func handle_post(w http.ResponseWriter, r * http.Request) {
    vars := mux.Vars(r)
    client_id  := vars["client_id"]
    topic := r.FormValue("topic")
    value := r.FormValue("value")

    if len(topic) > 0 && len(value) > 0 {
        fmt.Fprint(w, "about to store an entry in database: ", client_id, ":", topic, ":", value)
    }
}

func main() {
    r := mux.NewRouter()

    // GET requests
    r.HandleFunc("/{client_id}", handle_get).Methods("GET")
    r.HandleFunc("/{client_id}/{topic}", handle_get).Methods("GET")

    // DELETE requests
    r.HandleFunc("/{client_id}", handle_delete).Methods("DELETE")
    r.HandleFunc("/{client_id}/{topic}", handle_delete).Methods("DELETE")

    // POST requests
    r.HandleFunc("/{client_id}", handle_post).Methods("POST")

    http.Handle("/", r)
    log.Fatal(http.ListenAndServe(":"+utils.Portno("db"), nil))
}
