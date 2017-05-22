package actionstore

import (
	"fmt"

	"github.com/deadcheat/pedigree/chainer"
)

const (
	defaultRange = 10
)

// DefaultConfig デフォルト設定で実行するためのConfigstruct
var DefaultConfig = &Config{
	Object:    nil,
	SizeRange: defaultRange,
}

// ActionStore ExecutableなActionをStoreし、Chainableを実装する
type ActionStore struct {
	store       []chainer.Executable
	size        int
	pos         int
	sizeUpRange int
	Object      interface{}
}

// NewActionStore ActionStoreを生成（Objectはnil）
func NewActionStore() *ActionStore {
	return NewActionStoreWithConfig(DefaultConfig)
}

// Config ActionStore生成のためのコンフィグstruct
type Config struct {
	Object    interface{}
	SizeRange int
}

// NewActionStoreWithConfig ActionStoreを生成(Objectを指定する)
func NewActionStoreWithConfig(c *Config) *ActionStore {
	if c.SizeRange <= 0 {
		c.SizeRange = defaultRange
	}
	return &ActionStore{
		store:       make([]chainer.Executable, c.SizeRange),
		size:        0,
		pos:         0,
		sizeUpRange: c.SizeRange,
		Object:      c.Object,
	}
}

// Add storeにExecutableをAdd
func (a *ActionStore) Add(e chainer.Executable) {
	defer a.sizeUp()
	p := a.size
	if a.vacancy() {
		a.store[p] = e
		return
	}
	newStore := make([]chainer.Executable, a.size+a.sizeUpRange)
	for i := range a.store {
		newStore[i] = a.store[i]
	}
	newStore[p] = e
	a.store = newStore
}

// Next 次のActionを実行する
func (a *ActionStore) Next() error {
	// 全部終わった == Executableがnil入ってる
	s := a.store
	p := a.pos
	if p >= len(s) || s[p] == nil {
		return nil
	}
	if err := a.store[p].Execute(a, a.Object); err != nil {
		return fmt.Errorf("Error occured in Chainer #%d, err: %v", a.pos+1, err)
	}
	a.posUp()

	return a.Next()
}

// storeに空きがあるかをチェックする
func (a *ActionStore) vacancy() bool {
	return (a.size < len(a.store))
}

// sizeのincrement
func (a *ActionStore) sizeUp() {
	a.size++
}

// popのincrement
func (a *ActionStore) posUp() {
	a.pos++
}

// Store 保持しているStoreを取得
func (a *ActionStore) Store() []chainer.Executable {
	return a.store
}

// Size 現在のサイズを取得
func (a *ActionStore) Size() int {
	return a.size
}
