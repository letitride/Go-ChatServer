package trace

import "io"

import "fmt"

// Tracer はログ記録できるオブジェクトのインターフェース
type Tracer interface {
	//...interface{}は任意の型の引数を0個以上受け取れる
	Trace(...interface{})
}

type tracer struct {
	out io.Writer
}

// Tracerインターフェースの実装
func (t *tracer) Trace(a ...interface{}) {
	t.out.Write([]byte(fmt.Sprint(a...)))
	t.out.Write([]byte("\n"))
}

// New はインターフェースTracerを返すメソッド
func New(w io.Writer) Tracer {
	return &tracer{out: w}
}
