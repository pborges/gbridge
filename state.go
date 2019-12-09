package gbridge

import "github.com/pborges/gbridge/proto"

type State struct {
	Name  string
	Value interface{}
	Error proto.ErrorCode
}
