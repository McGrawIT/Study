package OAuth2

import (
    "os"
    "github.build.ge.com/aviation-predix-common/vcap-support.git"
)

const AppKey = "VCAP_APPLICATION"

func GetServiceName(name string) (world string) {

    world = os.Getenv(name)

    return
}

func GetServiceHostName(name string) (world string) {

    vcapServices, _ := vcap.LoadServices()
    for i := range vcapServices["user-provided"] {
        if vcapServices["user-provided"][i].Name == "aviation-dca-services" {
            vmap := vcapServices["user-provided"][i].Credentials
            if (vmap[name] != nil) {
                world = vmap[name].(string)
            } else { world = "" }
        }
    }

    return
}
