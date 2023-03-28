package common

import (
	"time"
    "os"
    "os/signal"
    "syscall"
    "strings"
    "bytes"
	log "github.com/sirupsen/logrus"
)
const (DATASET_PATH = "dataset.csv")

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

	datasetReader, err := NewCsvParser(DATASET_PATH)
	if err != nil {
		log.Fatalf(
	        "action: openDataset | result: fail | Agency: %v | error: %v",
			c.config.Agency,
			err,
		)
        return
	}
	defer datasetReader.Close()

	record, err := ReadLine()
	var batch bytes.Buffer
	packetsInBatch = 0
	for err == nil {
		bet := &common.Bet{
			Name:      record[0],
			LastName:  record[1],
			BirthDate: record[3],
			ID:        record[2],
			Number:    record[4],
			AgencyId:  agencyId,
		}
		betByteBuffer := SerializeBet(bet)
		batch, batchState = AddBetToBatch(bet, batch, packetsInBatch, c.batchSize)

		if batchStatus != BATCH_PARTWAY {
			err = socket.SendAll(batch)
			if err != nil {
				log.Fatalf(
					"action: sendBatch | result: fail | Agency: %v | error: %v",
					c.config.Agency,
					err,
				)
				return
			}
			batch.Truncate(0)
			packetsInBatch = 0
		}
		if batchStatus == BATCH_FULL {
			continue
		} else {
			record, err = ReadLine()
		}
	}
	if err.Error() == "EOF" {
		batch = AddNoMoreBatchesBytes(batch)
		err = socket.SendAll(batch)
		if err != nil {
			log.Fatalf(
				"action: sendBatch | result: fail | Agency: %v | error: %v",
				c.config.Agency,
				err,
			)
			return
		}
	} else {
		log.Fatalf(
			"action: readCDV | result: fail | Agency: %v | error: %v",
			c.config.Agency,
			err,
		)
		return
	}

	answer := make([]byte, 1)
	err = socket.RecvAll(answer)
	if err == nil && ValidateResult(answer) {
		log.Infof("action: bet_sent | result: success | ID: %v | number: %v | result: %v",
			c.bet.ID,
			c.bet.Number,
			answer,
		)
    } else {
		log.Fatalf(
	        "action: bet_sent | result: fail | ID: %v | error: %v | answer: %v",
			c.bet.ID,
			err,
			answer,
		)
	}
}
// esta funcion quedo un poco verbosa con el error handling, se podria manejar un poco mejor
