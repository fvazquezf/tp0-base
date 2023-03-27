package common

import (
	"time"
    "os"
    "os/signal"
    "syscall"
	log "github.com/sirupsen/logrus"
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
	connWrap  Socket
	bet       *Bet
}

// NewClient Initializes a new client receiving the configuration
// as a parameter and the desired bet
func NewClient(config ClientConfig, bet *Bet) *Client {
	client := &Client{
		config: config,
		bet: bet,
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
	        "action: connect | result: fail | Agendy: %v | error: %v",
			c.config.Agency,
			err,
		)
        return
    }
	defer socket.Close()

	betByteArray := SerializeBet(c.bet)
	err = socket.SendAll(betByteArray)
    if err != nil {
		log.Fatalf(
	        "action: sendBet | result: fail | Agendy: %v | error: %v",
			c.config.Agency,
			err,
		)
        return
    }

	answer := make([]byte, 1)
	err = socket.RecvAll(answer)
    log.Infof("Received data: %d result", ValidateResult(answer))
    log.Infof("Received data: %d error", err)

	if err == nil && ValidateResult(answer) {
		log.Infof("action: bet_sent | result: success | ID: %v | number: %v",
			c.bet.ID,
			c.bet.Number,
		)
    } else {
		log.Fatalf(
	        "action: bet_sent | result: fail | ID: %v | error: %v",
			c.bet.ID,
			err,
		)
	}
}
// esta funcion quedo un poco verbosa con el error handling, se podria manejar un poco mejor
