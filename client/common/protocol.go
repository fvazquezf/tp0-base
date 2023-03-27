package common

import (
    "bytes"
    "encoding/binary"
)

func SerializeBet(bet *Bet) []byte {
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
    
    lenBuf := make([]byte, 2)
    binary.BigEndian.PutUint16(lenBuf, uint16(buf.Len()))

    // Write the original buffer to the new buffer
    newBuf := bytes.NewBuffer(lenBuf)
    newBuf.Write(buf.Bytes())
    return newBuf.Bytes()
}

const (
    SUCCESS_BYTE = 0
    SUCCESS = true
    FAIL = false
)

func ValidateResult(answer []byte) bool {
    if answer[0] == SUCCESS_BYTE {
        return SUCCESS
    } else {
        return FAIL
    }
}