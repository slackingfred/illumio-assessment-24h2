package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/slackingfred/illumio-assessment-24h2/pkg/flowlog"
	"github.com/slackingfred/illumio-assessment-24h2/pkg/lookup"
)

// Command-line arguments.
// Example:
//   - Use specified lookup table and write to stdout
//     go run cmd/filter/main.go -lookup testdata/lookup_table.csv < testdata/flow.log
//   - Use default lookup table and write to specified files
//     go run cmd/filter/main.go -tag-output 1.log -combination-output 2.log < testdata/flow_non_tcp.log
var (
	lookupTablePath      = flag.String("lookup", "testdata/lookup_table.csv", "path to the lookup table")
	tagCountFile         = flag.String("tag-output", "", "path to the tag count output file - omit for stdout")
	combinationCountFile = flag.String("combination-output", "", "path to the combination count output file - omit for stdout")
)

type PortProtocol struct {
	Port     int32
	Protocol string
}

func main() {
	flag.Parse()
	var (
		err         error
		tagTable    lookup.Table
		flowLogLine flowlog.V2
		rec         []string

		tag              string
		combination      PortProtocol
		tagCount         = make(map[string]int)
		combinationCount = make(map[PortProtocol]int)
	)
	// Load config
	if err = tagTable.LoadFile(*lookupTablePath); err != nil {
		log.Fatalf("failed to load lookup table: %v", err)
	}
	// Read flow logs from stdin
	rdr := csv.NewReader(os.Stdin)
	rdr.Comma = ' ' // Flow log is space-separated
	rdr.TrimLeadingSpace = true
	for {
		if rec, err = rdr.Read(); err != nil {
			break
		}
		if err = flowLogLine.Parse(rec); err != nil {
			break
		}
		// Extract fields and convert protocol number to string
		combination.Port = flowLogLine.DstPort
		combination.Protocol = flowlog.IANAProtoNumberToString(flowLogLine.Protocol)
		// Try matching the tag
		tag = tagTable.GetTag(combination.Port, combination.Protocol)
		if tag == "" {
			tag = "Untagged"
		}
		// Update counters
		tagCount[tag]++
		combinationCount[combination]++
	}
	if err != nil && err != io.EOF {
		log.Fatalf("failed to read flow log: %v", err)
	}
	outputStats(tagCount, combinationCount)
}

func outputStats(tagCount map[string]int, combinationCount map[PortProtocol]int) {
	var (
		err error

		tagOut, combinationOut       *os.File
		tagWriter, combinationWriter *csv.Writer
	)
	// If output files are specified, write to them
	// Otherwise, write to stdout
	if *tagCountFile != "" {
		if tagOut, err = os.Create(*tagCountFile); err != nil {
			log.Fatalf("failed to open tag count file: %v", err)
		}
		defer tagOut.Close()
		tagWriter = csv.NewWriter(tagOut)
	} else {
		tagWriter = csv.NewWriter(os.Stdout)
	}
	if *combinationCountFile != "" {
		if combinationOut, err = os.Create(*combinationCountFile); err != nil {
			log.Fatalf("failed to open combination count file: %v", err)
		}
		defer combinationOut.Close()
		combinationWriter = csv.NewWriter(combinationOut)
	} else {
		combinationWriter = csv.NewWriter(os.Stdout)
	}
	defer tagWriter.Flush()
	defer combinationWriter.Flush()

	// If writing tag counts to stdout, add an extra header line
	if *tagCountFile == "" {
		fmt.Println("Tag Counts:")
	}
	tagWriter.Write([]string{"Tag", "Count"})
	// Write tag counts
	for tag, count := range tagCount {
		if err = tagWriter.Write([]string{tag, fmt.Sprint(count)}); err != nil {
			log.Fatalf("failed to write tag count: %v", err)
		}
		if *tagCountFile == "" {
			tagWriter.Flush() // Flush every line when using stdout
		}
	}

	// If both files are stdout, add a newline between the two outputs
	if *tagCountFile == "" && *combinationCountFile == "" {
		fmt.Println()
	}
	// If writing combination counts to stdout, add an extra header line
	if *combinationCountFile == "" {
		fmt.Println("Port/Protocol Combination Counts:")
	}
	combinationWriter.Write([]string{"Port", "Protocol", "Count"})
	// Write combination counts
	for combination, count := range combinationCount {
		if err = combinationWriter.Write([]string{fmt.Sprintf("%d", combination.Port), combination.Protocol, fmt.Sprint(count)}); err != nil {
			log.Fatalf("failed to write combination count: %v", err)
		}
		if *combinationCountFile == "" {
			combinationWriter.Flush() // Flush every line when using stdout
		}
	}
}
