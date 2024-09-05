package main

import "structs"

// NOTE:
// CGOのために必要になった。そもそもそのパターンって稀なのでは？そこから
// すなわち、HostLayoutは付けたものを変えるんじゃなくてつけない物を変える下準備です！
// 壊れてはいけないものにつける
type Example struct {
	_     structs.HostLayout
	Value string
}
