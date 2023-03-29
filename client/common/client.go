package common

import (
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	DATASET_PATH = "dataset.csv"
)

// ClientConfig Configuration used by the client
type ClientConfig struct {
	Agency        string
	ServerAddress string
	LoopLapse     time.Duration
	LoopPeriod    time.Duration
}

// Client Entity that encapsulates how
type Client struct {
	config    ClientConfig
	batchSize int
	agencyId  int
}

// NewClient Initializes a new client receiving the configuration
// as a parameter and the desired bet
func NewClient(config ClientConfig, batchSize int, agencyId int) *Client {
	client := &Client{
		config:    config,
		batchSize: batchSize,
		agencyId:  agencyId,
	}
	return client
}

// Create a conection to the server and send bet, then wait response
func (c *Client) SendBetAndValidate() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM)

	datasetReader, err := NewCsvReader(DATASET_PATH)
	if err != nil {
		log.Fatalf(
			"action: openDataset | result: fail | Agency: %v | error: %v",
			c.config.Agency,
			err,
		)
		return
	}
	defer datasetReader.Close()

	batch := NewBatch(c.batchSize)
	for !datasetReader.IsAtEnd() {
		socket, err := NewConnectedSocket(c.config.ServerAddress)
		if err != nil {
			log.Fatalf(
				"action: connect | result: fail | Agency: %v | error: %v",
				c.config.Agency,
				err,
			)
			return
		}
		defer socket.Close()

		batch.BuildBatch(datasetReader, c.config.Agency)
		err = socket.SendAll(batch.buf.Bytes())
		if err != nil {
			log.Fatalf(
				"action: sendBatch | result: fail | Agency: %v | error: %v",
				c.config.Agency,
				err,
			)
			return
		}
		batch.Reset()
		answer := make([]byte, 1) //esto es protocolo deberia esta encapsulado :/
		err = socket.RecvAll(answer)
		if err == nil && ValidateResult(answer) {
			log.Infof("action: batch_sent | result: success | result: %v",
				answer,
			)
		} else {
			log.Fatalf(
				"action: batch_sent | result: fail | error: %v | answer: %v",
				err,
				answer,
			)
			return
		}
	}
}

func (c *Client) CheckWinner() {
	finished := NOT_FINISHED
	for finished == NOT_FINISHED {
		socket, err := NewConnectedSocket(c.config.ServerAddress)
		if err != nil {
			log.Fatalf(
				"action: connect | result: fail | Agency: %v | error: %v",
				c.config.Agency,
				err,
			)
			return
		}
		defer socket.Close()
		finished = CheckIfFinished(c.agencyId, socket)
		if finished == NOT_FINISHED {
			log.Infof("action: sleeping ")
			time.Sleep(1 * time.Second)
		} else {
			log.Infof("action: server finished ")
			ReceiveWinners(finished, socket)
			return
		}
	}
}

// esta funcion quedo un poco verbosa con el error handling, se podria manejar un poco mejor
