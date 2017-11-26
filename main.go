package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	log "github.com/Sirupsen/logrus"
)

const (
	pollingDuration time.Duration = 20 * time.Second
	listeningPort                 = ":9898"
)

var (
	mu            sync.RWMutex
	ethereumNodes string
)

func updateEthereumNodes(addressRecord string) {
	buffer := bytes.NewBufferString("")

	// Resolve DNS A record to a set of IP addresses
	ipAddresses, err := ResolveAddressRecord(addressRecord)
	if err != nil {
		log.Errorf("Error resolving DNS address record: %s", err)
		goto updateNodes
	}

	log.Printf("%s resolved to %s", addressRecord, ipAddresses)

	// Retrieve enode from each IP address
	for _, ipAddress := range ipAddresses {
		resp, err := http.Get(fmt.Sprintf("http://%s:%s", ipAddress, "8080"))
		if err != nil {
			log.Errorf("Error retrieving enode address: %s", err)
			goto updateNodes
		}
		defer resp.Body.Close()
		contents, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Errorf("Error parsing response: %s", err)
			goto updateNodes
		}

		var enodeAddress = strings.TrimSpace(string(contents))
		if buffer.Len() > 0 {
			buffer.WriteString(",")
		}
		buffer.WriteString(enodeAddress)
		log.Infof("%s with enode address %s", ipAddress, enodeAddress)
	}

updateNodes:
	// Update list
	mu.Lock()
	defer mu.Unlock()
	ethereumNodes = buffer.String()
}

func startPollUpdateEthereumNodes(addressRecord string) {
	for {
		go updateEthereumNodes(addressRecord)
		<-time.After(pollingDuration)
	}
}

func webHandler(w http.ResponseWriter, r *http.Request) {
	log.Debugf("handling request from %s", r.RemoteAddr)

	mu.RLock()
	defer mu.RUnlock()
	fmt.Fprintln(w, ethereumNodes)
}

func main() {
	bootNodeService := flag.String("service", os.Getenv("BOOTNODE_SERVICE"), "DNS A Record for `bootnode` services. Alternatively set `BOOTNODE_SERVICE` env.")
	flag.Parse()

	if *bootNodeService == "" {
		flag.PrintDefaults()
		os.Exit(1)
	}

	log.Infof("starting bootnode-registrar: %s.", *bootNodeService)
	go startPollUpdateEthereumNodes(*bootNodeService)
	http.HandleFunc("/", webHandler)
	log.Fatal(http.ListenAndServe(listeningPort, nil))
}
