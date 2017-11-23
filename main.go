package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	log "github.com/Sirupsen/logrus"
)

const (
	pollingDuration time.Duration = 60 * time.Second
	listeningPort                 = ":9898"
)

var (
	bootnodes = make(map[string]string)
	enodes    string
)

func getEnodes(addressRecord string) {
	ipAddresses, err := ResolveAddressRecord(addressRecord)
	if err != nil {
		log.Errorf("Error resolving A record: %s", err)
		return
	}

	log.Printf("A Record {%s} resolved to  %s", addressRecord, ipAddresses)
	for _, ipAddress := range ipAddresses {
		if _, ok := bootnodes[ipAddress]; ok {
			log.Printf("%s. Already exists.", ipAddresses)
			continue
		}

		resp, err := http.Get(fmt.Sprintf("http://%s:%s", ipAddress, "8080"))
		if err != nil {
			log.Errorf("Error on HTTP GET: %s", err)
			return
		}
		defer resp.Body.Close()
		contents, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Errorf("Error parsing response: %s", err)
			return
		}

		bootnodes[ipAddress] = strings.TrimSpace(string(contents))
		log.Infof("%s. Adding %s", ipAddress, bootnodes[ipAddress])
	}

	i := 0
	buffer := bytes.NewBufferString("")
	for _, v := range bootnodes {
		buffer.WriteString(v)

		if i < len(bootnodes)-1 {
			buffer.WriteString(",")
		}

		i++
	}
	enodes = buffer.String()
}

func startPollGetEnodes(addressRecord string) {
	for {
		go getEnodes(addressRecord)
		<-time.After(pollingDuration)
	}
}

func webHandler(w http.ResponseWriter, r *http.Request) {
	log.Infof("Request from %s", r.RemoteAddr)
	fmt.Fprintln(w, enodes)
}

func main() {
	bootNodeService := flag.String("service", os.Getenv("BOOTNODE_SERVICE"), "DNS A Record for `bootnode` services. Alternatively set `BOOTNODE_SERVICE` env.")
	flag.Parse()

	if *bootNodeService == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	log.Infof("Starting bootnode-registrar. {%s}.", *bootNodeService)

	go startPollGetEnodes(*bootNodeService)
	http.HandleFunc("/", webHandler)
	log.Fatal(http.ListenAndServe(listeningPort, nil))
}
