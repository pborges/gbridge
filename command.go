package gbridge

import "github.com/pborges/gbridge/proto"

// Command belong to Traits
type Command interface {
	Name() string
	// passing maps in functions feels dirty
	// return the result (can be nil) and an error (if any)
	Execute(ctx Context, args map[string]interface{}) proto.CommandResponse
}
