package main

import (
	"github.build.ge.com/aviation-intelligent-airport/configuration-manager-svc/util"
	"github.build.ge.com/aviation-predix-common/vcap-support"
	"net/http"
	"os"
	"runtime"
	"time"
)

func handleInfoGet(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	vcapAppMap, _ := vcap.LoadApplication()
	guid := vcapAppMap.ID
	predixSpace := util.GetPredixSpace()

	// Always set content type and status code
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)

	// And write your response to w
	var ts = time.Now().UTC()
	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)

	hostname := util.GetServiceHostName("uaa_url")

	oauth2_state := "Disabled"
	if OAUTH2_ENABLED {
		oauth2_state = "Enabled"
	}

	//rma := float64(dbg.GetRuntimeMaxAllocMem()) / 1024.0 / 1024.0 //  convert to MB
	ma := float64(mem.Alloc) / 1024.0 / 1024.0         //  convert to MB
	ta := float64(mem.TotalAlloc) / 1024.0 / 1024.0    //  convert to MB
	sm := float64(mem.Sys) / 1024.0 / 1024.0           //  convert to MB
	fm := float64(mem.Sys-mem.Alloc) / 1024.0 / 1024.0 //  convert to MB

	msg := ProcessInfo{
		Service:       "Configuration Manager Microservice - Copyright GE 2015-2016",
		CurrentTime:   ts,
		StartTime:     startTime,
		UpTime:        ts.Sub(startTime).String(),
		PredixSpace:   predixSpace,
		InstanceGUID:  guid,
		InstanceIP:    os.Getenv("CF_INSTANCE_ADDR"),
		InstanceIndex: os.Getenv("CF_INSTANCE_INDEX"),
		UAA:           oauth2_state,
		ConfigurationVariables: map[string]string{
			"PORT":          os.Getenv("PORT"),
			"URI":           responsePath,
			"UAA_HOST_NAME": hostname,
		},
	}

	msg.MemoryStats.Alloc = ma
	msg.MemoryStats.TotalAlloc = ta
	msg.MemoryStats.Sys = sm
	msg.MemoryStats.Free = fm
	msg.MemoryStats.Lookups = mem.Lookups
	msg.MemoryStats.Mallocs = mem.Mallocs
	msg.MemoryStats.Frees = mem.Frees

	//fmt.Sprintf("System Operation\n"),
	//fmt.Sprintf("  Number Usable CPUS:     %d\n", dbg.GetRuntimeNumCpus()),
	//fmt.Sprintf("  Number Go Threads:      %d\n", dbg.GetRuntimeGoThreads()),
	//fmt.Sprintf("  Number cGo Threads:     %d\n", dbg.GetRuntimeCGoThreads()),
	//fmt.Sprintf("  Max Number Go Threads:  %d\n", dbg.GetRuntimeMaxGoThreads()),
	//fmt.Sprintf("  Max Number cGo Threads: %d\n", dbg.GetRuntimeMaxCGoThreads()),
	//fmt.Sprintf("\n"),
	//fmt.Sprintf("  Max Alloc:   %f  (MB allocated and not yet freed)\n", rma),
	//fmt.Sprintf("  PostgreSQL:     %-s\n", db_state),
	//fmt.Sprintf("  Oauth Version:  %-s\n", au.GetVersion()),
	//fmt.Sprintf("  Logger Version: %-s\n", dbg.GetVersion()),
	//fmt.Sprintf("\n"),

	writeObjectAsJson(w, msg)
}

type ProcessInfo struct {
	CurrentTime            time.Time
	StartTime              time.Time
	UpTime                 string
	Service                string
	PredixSpace            string
	InstanceGUID           string
	InstanceIP             string
	InstanceIndex          string
	UAA                    string
	MemoryStats            ProcessInfoMemoryStats
	ConfigurationVariables map[string]string
}

type ProcessInfoMemoryStats struct {
	Alloc      float64
	TotalAlloc float64
	Sys        float64
	Free       float64
	Lookups    uint64
	Mallocs    uint64
	Frees      uint64
}