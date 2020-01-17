package trace

import "io"

// Tracer はログ記録できるオブジェクトのインターフェース
type Tracer interface {
	//...interface{}は任意の型の引数を0個以上受け取れる
	Trace(...interface{})
}

// New はメソッド
func New(w io.Writer) Tracer {
	return nil
}
