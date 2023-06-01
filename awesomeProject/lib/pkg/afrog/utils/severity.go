package utils

type Severity int

const (
	Undefined Severity = iota

	INFO

	LOW

	MEDIUM

	HIGH

	CRITICAL
)

var SeverityMap = map[string]Severity{
	"info":     INFO,
	"low":      LOW,
	"medium":   MEDIUM,
	"high":     HIGH,
	"critical": CRITICAL,
}
