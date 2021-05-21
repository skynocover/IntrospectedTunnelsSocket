package api

import (
	"context"
	"itgserver/src/proto"
)

type Server struct{}

// GetKey 回傳key
func (s *Server) GetKey(ctx context.Context, in *proto.Request) (*proto.Reply, error) {

	return &proto.Reply{Result: []byte("")}, nil
}
