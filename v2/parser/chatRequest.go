package parser

import (
	objs "github.com/fulviodenza/telego/v2/objects"
)

type chatRequestHandler struct {
	requestId int
	function  *func(*objs.Update)
}
