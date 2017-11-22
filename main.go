package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

var (
	enodes map[string]struct{}
)

func getEnodes() {
	ipList, err := ResolveAddressRecord("bootnode-service.default.svc.cluster.local")
	if err != nil {
		fmt.Println(err)
		return
	}

	for _, ip := range ipList {
		fmt.Println("ip:", ip)
		response, err := http.Get(fmt.Sprintf("%s:%s", ip, "8080"))
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

		if _, ok := enodes[string(contents)]; !ok {
			enodes[string(contents)] = struct{}{}
		}
	}
}

func startPolling() {
	for {
		<-time.After(2 * time.Second)
		go getEnodes()
	}
}

func webHandler(w http.ResponseWriter, r *http.Request) {
	for _, enode := range enodes {
		fmt.Println(w, enode)
	}
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func main() {
	go startPolling()
	http.HandleFunc("/", webHandler)
	http.ListenAndServe(":8080", nil)
}
