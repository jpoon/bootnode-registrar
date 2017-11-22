package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

const (
	pollingDuration time.Duration = 60 * time.Second
)

var (
	bootnodes = make(map[string]string)
	enodes    string
)

func getEnodes(addressRecord string) {
	ipAddresses, err := ResolveAddressRecord(addressRecord)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("A Record [%s] resolved to  %s", addressRecord, ipAddresses)
	for _, ipAddress := range ipAddresses {
		if _, ok := bootnodes[ipAddress]; ok {
			log.Printf("%s. Already exists.", ipAddresses)
			continue
		}

		resp, err := http.Get(fmt.Sprintf("http://%s:%s", ipAddress, "8080"))
		if err != nil {
			log.Fatal(err)
		}
		defer resp.Body.Close()
		contents, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		bootnodes[ipAddress] = strings.TrimSpace(string(contents))
		log.Printf("%s. Adding %s", ipAddress, bootnodes[ipAddress])
	}

	i := 0
	buffer := bytes.NewBufferString("[\n")
	for _, v := range bootnodes {
		buffer.WriteString(fmt.Sprintf("\"%s\"", v))

		if i < len(enodes)-1 {
			buffer.WriteString(fmt.Sprintf(","))
		}

		buffer.WriteString(fmt.Sprintf("\n"))
		i++
	}
	buffer.WriteString(fmt.Sprintf("]\n"))
	enodes = buffer.String()
}

func startPollGetEnodes(addressRecord string) {
	for {
		go getEnodes(addressRecord)
		<-time.After(pollingDuration)
	}
}

func webHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Request from %s", r.RemoteAddr)
	fmt.Fprintln(w, enodes)
}

func main() {
	go startPollGetEnodes("bootnode-service.default.svc.cluster.local")
	http.HandleFunc("/", webHandler)
	http.ListenAndServe(":8080", nil)
}
