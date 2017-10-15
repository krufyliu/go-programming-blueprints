package trace

import (
	"fmt"
	"io"
)

// Tracer defines a interface for tracing information.
type Tracer interface {
	Trace(...interface{})
}

type tracer struct {
	out io.Writer
}

func (t tracer) Trace(a ...interface{}) {
	fmt.Fprint(t.out, a...)
	fmt.Fprintln(t.out)
}

// New creates a Tracer that will output the trace info to the io.writer.
func New(w io.Writer) Tracer {
	return &tracer{
		out: w,
	}
}

type nilTracer struct{}

func (t nilTracer) Trace(a ...interface{}) {}

// Off creates a Tracer that will ignore calls to Trace.
func Off() Tracer {
	return &nilTracer{}
}
