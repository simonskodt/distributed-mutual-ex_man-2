package logic

import (
	proto "distributed-mutual-exclusion_mand-2/service"
)

type Server struct {
	Clients []Client
	proto.UnimplementedServiceServer
}

type Client struct {
	Id string
	Port string
	proto.UnimplementedReplyToAllServiceServer
}