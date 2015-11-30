package streaming

import (
	"time"

	"github.com/myshkin5/netspel/logs"
	"github.com/myshkin5/netspel/utils"
)

var ReporterLogger Logger

func init() {
	ReporterLogger = logs.Logger
}

type ReporterImpl struct {
	expectedMessagesPerSecond int
	secondsPerCycle           float64
}

type Report struct {
	MessageCount uint32
	ByteCount    uint64
	ErrorCount   uint32
}

type Logger interface {
	Info(format string, args ...interface{})
}

func (r *ReporterImpl) Init(expectedMessagesPerSecond int, reportCycle time.Duration) {
	r.expectedMessagesPerSecond = expectedMessagesPerSecond
	r.secondsPerCycle = float64(reportCycle) / float64(time.Second)
}

func (r *ReporterImpl) Report(report Report) {
	messagesPerSecond := float64(report.MessageCount) / r.secondsPerCycle
	percent := messagesPerSecond / float64(r.expectedMessagesPerSecond) * 100.0
	errorsPerSecond := float64(report.ErrorCount) / r.secondsPerCycle
	bytesPerSecond := utils.ByteSize(report.ByteCount) / utils.ByteSize(r.secondsPerCycle)
	ReporterLogger.Info("%8d messages/s (%6.2f%%), %8d errors/s, %s/s", uint64(messagesPerSecond), percent, uint64(errorsPerSecond), bytesPerSecond.String())
}
