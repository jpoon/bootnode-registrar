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
	enodes string
)

func getEnodes(addressRecord string) {
	ipAddresses, err := ResolveAddressRecord(addressRecord)
	if err != nil {
		log.Errorf("Error resolving DNS address record: %s", err)
		return
	}

	log.Printf("%s resolved to %s", addressRecord, ipAddresses)

	var bootnodes = make(map[string]struct{})
	for _, ipAddress := range ipAddresses {
		resp, err := http.Get(fmt.Sprintf("http://%s:%s", ipAddress, "8080"))
		if err != nil {
			log.Errorf("Error retrieving enode address: %s", err)
			return
		}
		defer resp.Body.Close()
		contents, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Errorf("Error parsing response: %s", err)
			return
		}

		var enodeAddress = strings.TrimSpace(string(contents))
		bootnodes[enodeAddress] = struct{}{}
		log.Infof("%s with enode address %s", ipAddress, enodeAddress)
	}

	i := 0
	buffer := bytes.NewBufferString("")
	for k := range bootnodes {
		buffer.WriteString(k)

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
	log.Infof("handling request from %s", r.RemoteAddr)
	fmt.Fprintln(w, enodes)
}

func main() {
	bootNodeService := flag.String("service", os.Getenv("BOOTNODE_SERVICE"), "DNS A Record for `bootnode` services. Alternatively set `BOOTNODE_SERVICE` env.")
	flag.Parse()

	if *bootNodeService == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	log.Infof("starting bootnode-registrar. {%s}.", *bootNodeService)

	go startPollGetEnodes(*bootNodeService)
	http.HandleFunc("/", webHandler)
	log.Fatal(http.ListenAndServe(listeningPort, nil))
}
