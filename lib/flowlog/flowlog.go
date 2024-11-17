package flowlog

import (
	"errors"
	"strconv"
)

// Expected number of fields for each version.
const (
	Version2       = 2
	Version2Fields = 14
)

// Reference: https://docs.aws.amazon.com/vpc/latest/userguide/flow-log-records.html
// Only Version 2 is supported in the current implementation. There is no Version 1.
// To implement higher versions, define a new struct type and embed the V2 type.
type V2 struct {
	Version     int32
	AccountID   string
	InterfaceID string
	SrcAddr     string
	DstAddr     string
	SrcPort     int32
	DstPort     int32
	Protocol    int32
	Packets     int64
	Bytes       int64
	Start       int64
	End         int64
	Action      string
	LogStatus   string
}

// Errors for parsing flow log records.
var (
	ErrInvalidVersion  = errors.New("invalid version")
	ErrNotEnoughFields = errors.New("not enough fields")
)

func (rec *V2) Parse(fields []string) error {
	var (
		tmp int
		err error
	)
	if len(fields) < Version2Fields {
		return ErrNotEnoughFields
	}
	// Version
	if tmp, err = strconv.Atoi(fields[0]); err != nil {
		return err
	}
	if tmp < Version2 {
		return ErrInvalidVersion
	}
	rec.Version = int32(tmp)
	// AccountID, InterfaceID, SrcAddr, DstAddr
	rec.AccountID = fields[1]
	rec.InterfaceID = fields[2]
	rec.SrcAddr = fields[3]
	rec.DstAddr = fields[4]
	// SrcPort, DstPort, Protocol
	if tmp, err = strconv.Atoi(fields[5]); err != nil {
		return err
	}
	rec.SrcPort = int32(tmp)
	if tmp, err = strconv.Atoi(fields[6]); err != nil {
		return err
	}
	rec.DstPort = int32(tmp)
	if tmp, err = strconv.Atoi(fields[7]); err != nil {
		return err
	}
	rec.Protocol = int32(tmp)
	// Packets, Bytes, Start, End
	if rec.Packets, err = strconv.ParseInt(fields[8], 10, 64); err != nil {
		return err
	}
	if rec.Bytes, err = strconv.ParseInt(fields[9], 10, 64); err != nil {
		return err
	}
	if rec.Start, err = strconv.ParseInt(fields[10], 10, 64); err != nil {
		return err
	}
	if rec.End, err = strconv.ParseInt(fields[11], 10, 64); err != nil {
		return err
	}
	// Action, LogStatus
	rec.Action = fields[12]
	rec.LogStatus = fields[13]
	return nil
}
