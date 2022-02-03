package trap_test

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"testing"
)
func TestTrap(t *testing.T) {
	var err error
	func () {
		// 此处 的 err 是新增变量而不是赋值给函数外的 err
		data, err := json.Marshal(func() {})
		assert.Nil(t,data, nil)
		assert.EqualError(t, err, "json: unsupported type: func()")
	}()
	// 函数外的 err 还是 nil,一直没有改变过
	assert.Nil(t, err)

	/*
	在这个场景中:

	data, err := json.Marshal(func() {})
	等于
	var data []byte
	var err error
	data, err = json.Marshal(func() {})

	改为:
	var data []byte
	data, err = json.Marshal(func() {})
	才能赋值到外部 err

	或者参考 TestAvoidTrap 代码避开陷阱
	*/
}

// 通过函数出参传递 err 避免陷阱
func TestAvoidTrap(t *testing.T) {
	var err error
	err = func() (err error) {
		// 此处 的 err 是新增变量而不是赋值给函数外的 err
		data, err := json.Marshal(func() {})
		assert.Nil(t, data, nil)
		assert.EqualError(t, err, "json: unsupported type: func()")
		return
	}()
	// 函数外的 err 还是 nil,一直没有改变过
	assert.EqualError(t, err, "json: unsupported type: func()")
}