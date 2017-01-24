package clients

import (
    "github.build.ge.com/aviation-predix-common/vcap-support.git"

    "os"
    "encoding/json"
    "fmt"
    "log"
    "strings"
)

const AppKey = "VCAP_APPLICATION"

func GetServiceName(name string) (world string) {

    world = os.Getenv(name)

    return
}

func GetServiceClass(sType string) (fType string, cnt int32) {

    cnt = 0
    vcapServices, _ := vcap.LoadServices()
    for key,_ := range vcapServices {
        if strings.Index(strings.ToLower(key), strings.ToLower(sType)) >= 0 {
            if (cnt == 0) { fType = key }
            cnt++
        }
    }

    return
}

func GetServiceHostName(name string) (world string) {

    vcapServices, _ := vcap.LoadServices()
    for i := range vcapServices["user-provided"] {
        if vcapServices["user-provided"][i].Name == "aviation-dca-services" {
            vmap := vcapServices["user-provided"][i].Credentials
            world = vmap[name].(string)
        }
    }

    return
}

func GetVcapCredential(class, name, field string) (value string) {

    value = ""

    vcapServices, _ := vcap.LoadServices()
    for i := range vcapServices[class] {
        if vcapServices[class][i].Name == name {
            vmap := vcapServices[class][i].Credentials
            if vmap[field] != nil {
                value = vmap[field].(string)
            }
        }
    }

    return
}

func GetPredixSpace() (space string) {
    var v map[string]interface{}

    vcap := os.Getenv("VCAP_APPLICATION")

    err := json.Unmarshal([]byte(vcap), &v)
    if err != nil {
        if DEBUG { fmt.Println("DBG-> Vcap: ", vcap) }
        log.Println("ERROR: Could not convert Vcap Services json data; msg: ", err.Error())
    }

    if (v != nil) {
        space = v["space_name"].(string)
        if DEBUG { fmt.Println("DBG-> space: ", space) }
    } else {
        space = "unknown"
    }

    return
}
