package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

var (
	enodes = make(map[string]struct{})
)

func getEnodes() {
	ipList, err := ResolveAddressRecord("bootnode-service.default.svc.cluster.local")
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, ip := range ipList {
		response, err := http.Get(fmt.Sprintf("http://%s:%s", ip, "8080"))
		if err != nil {
			fmt.Println(err)
			return
		}

		defer response.Body.Close()
		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Printf("%s", err)
			return
		}

		enode := strings.TrimSpace(string(contents))
		if _, ok := enodes[enode]; !ok {
			enodes[enode] = struct{}{}
		}
	}
}

func startPolling() {
	for {
		go getEnodes()
		<-time.After(10 * time.Second)
	}
}

func webHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "[")

	i := 1
	for k := range enodes {
		fmt.Fprintf(w, "\"%s\"", k)

		if i < len(enodes) {
			fmt.Fprintf(w, ",")
		}

		fmt.Fprintln(w, "")
		i++
	}
	fmt.Fprintln(w, "]")
}

func main() {
	go startPolling()
	http.HandleFunc("/", webHandler)
	http.ListenAndServe(":8080", nil)
}
