package main

import (
    "fmt"
    "os"
    "log"
    "encoding/json"
    "io/ioutil"
    "net/http"
)

type HantyConfig struct {
    HantyServices []HantyServiceNode `json:"services"`
}

type HantyServiceNode struct {
    Name string `json:"name"`
    Host string `json:"host"`
    Port string `json:"port"`
}

type HantyTopicContext map[string] []string


func get_config(service_name string) (HantyServiceNode) {
    var config HantyConfig
    var res HantyServiceNode

    config_file, err := os.Open(os.Getenv("HANTY_CONFIG_FILE"))
    defer config_file.Close()
    if err != nil {
        panic(err)
    }

    byte_config, _ := ioutil.ReadAll(config_file)
    json.Unmarshal(byte_config, &config)

    for _, item := range config.HantyServices {
        if item.Name == service_name {
            res = item
        }
    }

    return res
}

func get_context() (HantyTopicContext) {
    var ctx HantyTopicContext = make(HantyTopicContext)

    ctx["ChangeColor"] = []string{"#5e72e4", "#f3a4b5", "#ffd600",
                                  "#2bffc6", "#fd5d93", "#ffffff"}

    ctx["SendMessage"] = []string{"SayHi", "SayStopIt"}

    return ctx
}

func handle_get(w http.ResponseWriter, r * http.Request) {
    w.Header().Set("Access-Control-Allow-Origin", "*")
    topic_, _ := r.URL.Query()["topic"]
    var ctx HantyTopicContext
    var reply = make(map[string] []string)
    ctx = get_context()

    if len(topic_) == 0 {
        // return all supported topics
        supported_topics := make([]string, 0, len(ctx))
        for k := range ctx {
           supported_topics = append(supported_topics, k)
        }
        reply["supported_topics"] = supported_topics
    } else {
        // return list of states for a given topic
        topic := topic_[0]
        reply[topic] = ctx[topic]
    }

    res, _ := json.Marshal(reply)
    fmt.Fprintf(w, string(res))
}

func main() {
    config := get_config(os.Getenv("THIS_SERVICE_NAME"))
    http.HandleFunc("/get", handle_get)
    log.Fatal(http.ListenAndServe(":"+config.Port, nil))
}
