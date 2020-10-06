package config

import (
    "os"
    "sync"
    "errors"
    "io/ioutil"
    "encoding/json"
)

/* Type defs */
type configItem struct {
    Host string `json:"host"`
    Port string `json:"port"`
}

type serviceName string
type Config map[serviceName]*configItem


/* Singleton */
var (
    conf * Config
    once sync.Once
)

func get_config() * Config {
    once.Do(func() {
        conf = init_config()
    })

    return conf
}

//* Reading/Parsing of the config file */
func init_config() * Config {
    var res Config

    type HantyServiceNode struct {
        Name serviceName `json:"name"`
        Host string `json:"host"`
        Port string `json:"port"`
    }

    type HantyConfig struct {
        HantyServices []HantyServiceNode `json:"services"`
    }

    var config HantyConfig

    config_file, err := os.Open(os.Getenv("KHANTY_CONFIG_FILE"))
    // config_file, err := os.Open("/mount/.config/config.json")

    defer config_file.Close()
    if err != nil {
        panic(err)
    }

    byte_config, _ := ioutil.ReadAll(config_file)
    json.Unmarshal(byte_config, &config)

    res = make(map[serviceName]*configItem)

    for _, item := range config.HantyServices {
        res[item.Name] = &configItem {
            Host: item.Host,
            Port: item.Port}
    }

    return &res
}

/* Config Constructor */
func New() * Config {
    return get_config()
}

/* Public getter of config data */
func (c Config) GetData(service string) (string, string, error) {
    var port string
    var host string

    if item, ok := c[serviceName(service)]; ok {
        port = item.Port
        host = item.Host
    } else {
        str := "service " + service + " unknown"
        return port, host, errors.New(str)
    }
    return port, host, nil
}
