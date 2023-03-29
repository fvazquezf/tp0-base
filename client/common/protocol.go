package common

import (
	"bytes"
	"encoding/binary"

	log "github.com/sirupsen/logrus"
)

const (
	SUCCESS_BYTE    = 0
	SUCCESS         = true
	FAIL            = false
	MORE_BATCHES    = 1234
	NO_MORE_BATCHES = 1235
	NOT_FINISHED    = 1
)

func SerializeBet(bet *Bet) *bytes.Buffer {
	var buf bytes.Buffer

	// Write the length and string data of each string field to the buffer
	buf.WriteByte(uint8(len(bet.Name)))
	buf.WriteString(bet.Name)
	buf.WriteByte(uint8(len(bet.LastName)))
	buf.WriteString(bet.LastName)
	buf.WriteByte(uint8(len(bet.ID)))
	buf.WriteString(bet.ID)
	buf.WriteByte(uint8(len(bet.BirthDate)))
	buf.WriteString(bet.BirthDate)

	// Write the uint16 fields to the buffer
	binary.Write(&buf, binary.BigEndian, bet.Number)
	binary.Write(&buf, binary.BigEndian, bet.AgencyId)

	lenBuf := make([]byte, U16SIZE)
	binary.BigEndian.PutUint16(lenBuf, uint16(buf.Len()))

	// Write the original buffer to the new buffer
	newBuf := bytes.NewBuffer(lenBuf)
	newBuf.Write(buf.Bytes())
	return newBuf
}

func ReceiveWinners(numIds int, socket *Socket) {
	// Initialize a slice to store the ids
	idList := make([]string, numIds)

	// Start reading the ids after the first 2 bytes
	for i := 0; i < numIds; i++ {
		// Read the length of the id as a 2-byte big-endian integer
		buffer := make([]byte, 2)
		err := socket.RecvAll(buffer)
		if err != nil {
			log.Fatalf(
				"action: recieve_buffer | result: fail | error: %v",
				err,
			)
			return
		}
		idLength := binary.BigEndian.Uint16(buffer)

		buffer = make([]byte, idLength)
		err = socket.RecvAll(buffer)
		if err != nil {
			log.Fatalf(
				"action: recieve_Id | result: fail | error: %v",
				err,
			)
			return
		}
		id := string(buffer)
		idList[i] = id
	}
	log.Infof("action: consulta_ganadores | result: success | cant_ganadores: %v",
		len(idList),
	)
}

func addMoreBatchesBytes(batch bytes.Buffer) bytes.Buffer {
	lastBytes := make([]byte, U16SIZE)
	binary.BigEndian.PutUint16(lastBytes, uint16(MORE_BATCHES))
	batch.Write(lastBytes)
	return batch
}

func AddNoMoreBatchesBytes(batch bytes.Buffer) bytes.Buffer {
	log.Infof("action: finished all bets reads ")
	lastBytes := make([]byte, U16SIZE)
	binary.BigEndian.PutUint16(lastBytes, uint16(NO_MORE_BATCHES))
	batch.Write(lastBytes)
	return batch
}

func ValidateResult(answer []byte) bool {
	if answer[0] == SUCCESS_BYTE {
		return SUCCESS
	} else {
		return FAIL
	}
}

func CheckIfFinished(agencyId int, socket *Socket) int {
	msg := make([]byte, U16SIZE)
	binary.BigEndian.PutUint16(msg, uint16(agencyId))
	err := socket.SendAll(msg)
	if err != nil {
		log.Fatalf(
			"action: recieve_finish | result: fail | error: %v",
			err,
		)
		return NOT_FINISHED
	}
	buffer := make([]byte, 2)
	err = socket.RecvAll(buffer)
	if err != nil {
		log.Fatalf(
			"action: recieve_finish | result: fail | error: %v",
			err,
		)
		return NOT_FINISHED
	}
	return int(binary.BigEndian.Uint16(buffer))
}
