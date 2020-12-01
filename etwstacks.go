package etwstacks

import (
	"runtime"
	"strings"

	"github.com/Microsoft/go-winio/pkg/etw"
	"github.com/Microsoft/go-winio/pkg/guid"
)

var p *etw.Provider

func init() {
	// Provider ID: fc009c55-f069-5ee5-b70f-fa2fa7b2dd20
	p, _ = etw.NewProvider("etwstacks", etwCallback)
}

func writeError(err error) {
	p.WriteEvent(
		"StackDumpError",
		etw.WithEventOpts(
			etw.WithLevel(etw.LevelError),
		),
		etw.WithFields(
			etw.StringField("err", err.Error()),
		),
	)
}

func etwCallback(sourceID guid.GUID, state etw.ProviderState, level etw.Level, matchAnyKeyword uint64, matchAllKeyword uint64, filterData uintptr) {
	if state == etw.ProviderStateCaptureState {
		stacks := getStacks()
		g, err := guid.NewV4()
		if err != nil {
			writeError(err)
			return
		}
		p.WriteEvent(
			"BeginStackDump",
			etw.WithEventOpts(
				etw.WithOpcode(etw.OpcodeStart),
				etw.WithActivityID(g),
			),
			etw.WithFields(
				etw.IntField("stackCount", len(stacks)),
			),
		)
		for i, stack := range stacks {
			p.WriteEvent(
				"StackDump",
				etw.WithEventOpts(
					etw.WithActivityID(g),
				),
				etw.WithFields(
					etw.IntField("stackIndex", i),
					etw.StringField("stack", stack),
				),
			)
		}
		p.WriteEvent(
			"EndStackDump",
			etw.WithEventOpts(
				etw.WithOpcode(etw.OpcodeStop),
				etw.WithActivityID(g),
			),
			nil,
		)
	}
}

func getStacks() []string {
	var (
		buf       []byte
		stackSize int
	)
	bufferLen := 16384
	for stackSize == len(buf) {
		buf = make([]byte, bufferLen)
		stackSize = runtime.Stack(buf, true)
		bufferLen *= 2
	}
	buf = buf[:stackSize]
	return strings.Split(string(buf), "\n\n")
}
