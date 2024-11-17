package flowlog

import (
	"reflect"
	"strings"
	"testing"
)

func TestParseV2(t *testing.T) {
	// Normal input
	{
		fields := strings.Split("2 123456789012 eni-0a1b2c3d 10.0.1.201 198.51.100.2 443 49153 6 25 20000 1620140761 1620140821 ACCEPT OK", " ")
		rec := &V2{}
		if err := rec.Parse(fields); err != nil {
			t.Errorf("Unexpected error for normal input: %v", err)
		}
		if !reflect.DeepEqual(rec, &V2{
			Version:     2,
			AccountID:   "123456789012",
			InterfaceID: "eni-0a1b2c3d",
			SrcAddr:     "10.0.1.201",
			DstAddr:     "198.51.100.2",
			SrcPort:     443,
			DstPort:     49153,
			Protocol:    6,
			Packets:     25,
			Bytes:       20000,
			Start:       1620140761,
			End:         1620140821,
			Action:      "ACCEPT",
			LogStatus:   "OK",
		}) {
			t.Errorf("Incorrect parse result: %v", rec)
		}
	}
	// Invalid version
	{
		fields := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13", "14"}
		rec := &V2{}
		if err := rec.Parse(fields); err != ErrInvalidVersion {
			t.Errorf("Unexpected error for invalid version: %v", err)
		}
	}
	// Not enough fields
	{
		fields := []string{"2", "3", "4", "5", "6", "7", "8", "9", "10", "11", "12", "13"}
		rec := &V2{}
		if err := rec.Parse(fields); err != ErrNotEnoughFields {
			t.Errorf("Unexpected error for not enough fields: %v", err)
		}
	}
	// Higher version
	{
		fields := strings.Split("3 123456789012 eni-4d3c2b1a 192.168.1.100 203.0.113.101 23 49154 6 15 12000 1620140761 1620140821 REJECT OK 1 1 1 1 1 1 1", " ")
		rec := &V2{}
		if err := rec.Parse(fields); err != nil {
			t.Errorf("Unexpected error for higher version: %v", err)
		}
		if rec.Action != "REJECT" {
			t.Errorf("Incorrect action for higher version: %v", rec)
		}
	}
}
