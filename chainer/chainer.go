package chainer

// Chainable Execute-Chainを実装するためのもの
type Chainable interface {
	Add(e Executable)
	Next() error
}

// Executable Chainerbleに格納され、順次Executeされる
type Executable interface {
	Execute(c Chainable, o interface{}) error
}
