package common

import (
	"bytes"
	"strconv"

	log "github.com/sirupsen/logrus"
)

const (
	MAX_BATCH_SIZE = 8192
	U16SIZE        = 2
)

type Batch struct {
	buf         bytes.Buffer
	packets     int
	maxPackets  int
	isFull      bool
	savedPacket bytes.Buffer
}

func NewBatch(batchMaxSize int) *Batch {
	batch := &Batch{
		maxPackets: batchMaxSize,
		packets:    0,
		isFull:     false,
	}
	return batch
}

func (b *Batch) AddBetToBatch(bet bytes.Buffer) {
	if (bet.Len() + b.buf.Len()) > MAX_BATCH_SIZE-U16SIZE {
		b.buf = addMoreBatchesBytes(b.buf)
		b.isFull = true
		b.savedPacket = bet
	}
	b.buf.Write(bet.Bytes())
	b.packets++
	if b.packets == b.maxPackets {
		b.buf = addMoreBatchesBytes(b.buf)
		b.isFull = true
	}
}

func (b *Batch) Reset() {
	b.buf.Truncate(0)
	b.isFull = false
	b.packets = 0
	if b.savedPacket.Len() != 0 {
		b.buf.Write(b.savedPacket.Bytes())
		b.savedPacket.Truncate(0)
		b.packets = 1
	}
}

func (b *Batch) BuildBatch(datasetReader *CsvReader, agency string) {
	for !b.isFull {
		// time.Sleep(5 * time.Second)
		record, err := datasetReader.ReadLine()
		if err != nil {
			if err.Error() == "EOF" {
				b.buf = AddNoMoreBatchesBytes(b.buf)
				b.isFull = true
				continue
			}
			log.Fatalf(
				"action: readCDV | result: fail | error: %v",
				err,
			)
			return
		}
		number, _ := strconv.ParseUint(record[4], 10, 16)
		agencyId, _ := strconv.ParseUint(agency, 10, 16)
		bet := &Bet{
			Name:      record[0],
			LastName:  record[1],
			BirthDate: record[3],
			ID:        record[2],
			Number:    uint16(number),
			AgencyId:  uint16(agencyId),
		}
		betByteBuffer := SerializeBet(bet)
		b.AddBetToBatch(*betByteBuffer)
	}
}
