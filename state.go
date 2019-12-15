package gbridge

import "github.com/pborges/gbridge/proto"

// State represents the state of a device 
type State struct {
	Name  string
	Value interface{}
	Error proto.ErrorCode
}
