package handlers

import (
	"encoding/base64"
	"fmt"
	"math/rand"
	"net"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/shlmvgleb/mtg-task/client/internal/config"
	log "github.com/sirupsen/logrus"
)

const letterBytes = "1234567890" // from task description

type Controller struct {
	config *config.AppConfig
}

func NewController(config *config.AppConfig) *Controller {
	return &Controller{
		config: config,
	}
}

func (cntrl *Controller) SubmitNewConnections(w http.ResponseWriter, r *http.Request) {
	numThreadsStr := r.FormValue("numThreads")
	numThreads, err := strconv.Atoi(numThreadsStr)
	if err != nil || numThreads < 1 {
		http.Error(w, "Неверное количество потоков", http.StatusBadRequest)
		return
	}

	log.Infoln("successfully connected to server")

	for i := 1; i <= numThreads; i++ {
		go cntrl.sendData(uuid.NewString())
	}

	fmt.Fprintf(w, "Запущено %d потоков для отправки данных", numThreads)
}

func (cntrl *Controller) sendData(id string) {
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", cntrl.config.ServerHost, cntrl.config.ServerPort))
	if err != nil {
		log.Errorf("failed to connect to server: %v\n", err)
	}

	defer conn.Close()

	for {
		data, err := generateRandomBase64String(1000 + rand.Intn(9000)) // от 1000 до 10000 байт
		if err != nil {
			log.Errorf("client %s: error generating random data: %v\n", id, err)
			return
		}

		_, err = conn.Write([]byte(data))
		if err != nil {
			log.Errorf("client %s: error sending data: %v\n", id, err)
			return
		}

		log.Infof("client %s sent data: %s\n", id, data[:50])
	}
}

func generateRandomBase64String(size int) (string, error) {
	bytes := randStringBytes(size)
	return base64.StdEncoding.EncodeToString(bytes), nil
}

func randStringBytes(n int) []byte {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return b
}
