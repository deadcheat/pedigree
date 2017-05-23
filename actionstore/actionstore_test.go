package actionstore

import (
	"errors"
	"fmt"
	"testing"

	"github.com/deadcheat/pedigree/chainer"
	"github.com/golang/mock/gomock"

	mock_chainer "github.com/deadcheat/pedigree/mocks/chainer"
)

// Add()をテスト
func TestActionStore_Add(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	me := mock_chainer.NewMockExecutable(ctrl)

	as := NewActionStore()
	r := DefaultConfig.Expansion
	// １件追加した場合には配列サイズが拡張されないこと
	as.Add(me)
	if len(as.Store()) != r || as.Size() != 1 {
		t.Error("ActionStore.Add() may not work")
	}

	// そのまま10件目までは同様に追加されること
	for i := 1; i < r; i++ {
		as.Add(mock_chainer.NewMockExecutable(ctrl))
		if len(as.Store()) != 10 || as.Size() != i+1 {
			t.Errorf("ActionStore.Add() may not work, when size: %d", as.Size())
		}
	}

	// 11件目でサイズを超えるので範囲が拡張されること
	as.Add(mock_chainer.NewMockExecutable(ctrl))
	if len(as.Store()) != r*2 || as.Size() != r+1 {
		t.Errorf("ActionStore.Add() may not work, when size: %d", as.Size())
	}

	// Config指定で小さいものを定義したときに、defaultに戻されていること
	as = NewActionStoreWithConfig(&Config{
		Object:    nil,
		Expansion: -1,
	})
	if len(as.Store()) != 10 {
		t.Error("NewActionStoreWithConfig() will change their size-range to default-value(10) when argment is invalid")
	}
}

// TestActionStore_Next Next()をテスト
func TestActionStore_Next(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// 通常の実行（エラーが発生しない場合）
	m1 := mock_chainer.NewMockExecutable(ctrl)
	m2 := mock_chainer.NewMockExecutable(ctrl)
	test := "test"
	c := &Config{
		Object:    test,
		Expansion: 1,
	}
	var as chainer.Chainable = NewActionStoreWithConfig(c)
	as.Add(m1)
	as.Add(m2)
	m1.EXPECT().Execute(gomock.Any(), gomock.Any()).Times(1).Return(nil)
	m2.EXPECT().Execute(gomock.Any(), gomock.Any()).Times(1).Return(nil)

	if err := as.Next(); err != nil {
		t.Errorf("Next() should not return any error but occured. err : %v", err)
	}

	// 途中でエラーが発生する場合
	m3 := mock_chainer.NewMockExecutable(ctrl)
	m4 := mock_chainer.NewMockExecutable(ctrl)
	as = NewActionStoreWithConfig(c)
	as.Add(m3)
	as.Add(m4)
	argErr := errors.New("test-error")
	expectedErr := fmt.Errorf("Error occured in Chainer #%d, err: %v", 1, argErr)
	m3.EXPECT().Execute(gomock.Any(), gomock.Any()).Times(1).Return(argErr)
	m4.EXPECT().Execute(gomock.Any(), gomock.Any()).AnyTimes()

	if err := as.Next(); err.Error() != expectedErr.Error() {
		t.Errorf("Next() should return error: [%v] but occured. err : [%v]", expectedErr, err)
	}
}
