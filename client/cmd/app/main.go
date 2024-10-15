package main

import (
	"fmt"
	"net"
	"net/http"
	"text/template"
	"time"

	"github.com/gorilla/mux"
	"github.com/shlmvgleb/mtg-task/client/internal/config"
	"github.com/shlmvgleb/mtg-task/client/internal/handlers"
	"github.com/shlmvgleb/mtg-task/client/pkg/form"
	log "github.com/sirupsen/logrus"
)

func startServer(cntrl *handlers.Controller, config *config.AppConfig) {
	r := mux.NewRouter()

	r.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		t := template.New("form")
		t, _ = t.Parse(form.FormTmpl)
		t.Execute(w, nil)
	})

	r.Methods(http.MethodPost).Path("/submit").HandlerFunc(cntrl.SubmitNewConnections)

	srv := &http.Server{
		Addr:         fmt.Sprintf("0.0.0.0:%d", config.ClientPort),
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}

	log.Fatal(srv.ListenAndServe())
}

func main() {
	config := config.ReadFromEnv()
	ping, err := net.Dial("tcp", fmt.Sprintf("%s:%d", config.ServerHost, config.ServerPort))
	if err != nil {
		log.Fatalf("failed to connect to server: %v\n", err)
	}
	defer ping.Close()

	log.Infof("successfully connected to the server on port %d", config.ServerPort)

	cntrl := handlers.NewController(config)
	startServer(cntrl, config)
}
