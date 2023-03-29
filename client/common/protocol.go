package common

import (
	"bytes"
	"encoding/binary"
)

const (
	SUCCESS_BYTE    = 0
	SUCCESS         = true
	FAIL            = false
	MORE_BATCHES    = 1234
	NO_MORE_BATCHES = 1235
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

func addMoreBatchesBytes(batch bytes.Buffer) bytes.Buffer {
	lastBytes := make([]byte, U16SIZE)
	binary.BigEndian.PutUint16(lastBytes, uint16(NO_MORE_BATCHES))
	batch.Write(lastBytes)
	return batch
}

func AddNoMoreBatchesBytes(batch bytes.Buffer) bytes.Buffer {
	lastBytes := make([]byte, U16SIZE)
	binary.BigEndian.PutUint16(lastBytes, uint16(MORE_BATCHES))
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
