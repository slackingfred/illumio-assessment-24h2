# Illumio Technical Assessment: Flow Log Analyzer

This work sample takes an Amazon VPC Version 2 Flow Log from standard input and
counts tag matches and port–protocol combinations in the log. Results can be
printed to standard output or CSV files.

To minimize external dependencies (the current code has none), Go's `flag`
package is used for command line argument parsing. This means the program
follows Go's command line convention which differs from the GNU style. Please
note the usage in the examples below.

The input part is kept simple and supports only stdin. For the output, because
there are two result sets, options are provided to write them into two files.

## Running the program

On UNIX-based platforms, just run this in any shell (tested with Go 1.21):

``` bash
go run cmd/filter/main.go < testdata/flow.log
```

It will print the results into the console like this:

``` csv
Tag Counts:
Tag,Count
sv_P1,2
email,3
Untagged,8
sv_P2,1

Port/Protocol Combination Counts:
Port,Protocol,Count
80,tcp,1
443,tcp,1
23,tcp,1
49155,tcp,1
49154,tcp,1
1024,tcp,1
993,tcp,1
49153,tcp,1
49157,tcp,1
110,tcp,1
49156,tcp,1
25,tcp,1
143,tcp,1
49158,tcp,1
```

If you want one output file, just add `> output-file` to redirect it. See the
comments in [main.go](cmd/filter/main.go) for full usage.

You can also build the program into binary:

``` bash
cd cmd/filter
go build
./filter -lookup ../../testdata/lookup_table.csv < ../../testdata/flow.log
```

Windows is not recommended because PowerShell's I/O redirection may add an UTF-8
BOM and interfere with parsing.

## Performance analysis

This program supports arbitrarily large flow logs because it processes the input
in a streaming manner. The lookup table is kept in memory, but for day-to-day
usage (<=10K mappings), memory usage is not an issue.

## Tests

Unit tests are performed on `flowlog` and `lookup` packages with normal and edge
cases.

## Limitations

- Only ICMP, TCP, and UDP are supported. To extend this, add more protocols to
[iana_proto_numbers.go](pkg/flowlog/iana_proto_numbers.go).
- Errors in IP addresses and ports (e.g. out of range) are not checked.
