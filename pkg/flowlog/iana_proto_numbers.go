package flowlog

var (
	// Protocol numbers defined by IANA (partial list).
	// Full list is available at https://www.iana.org/assignments/protocol-numbers/protocol-numbers.xhtml
	numToProto = map[int32]string{
		1:  "icmp",
		6:  "tcp",
		17: "udp",
	}
)

func IANAProtoNumberToString(num int32) string {
	if proto, ok := numToProto[num]; ok {
		return proto
	}
	return "unknown"
}
