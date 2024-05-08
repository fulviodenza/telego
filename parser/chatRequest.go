package parser

import (
	objs "github.com/fulviodenza/telego/objects"
)

type chatRequestHandler struct {
	requestId int
	function  *func(*objs.Update)
}
