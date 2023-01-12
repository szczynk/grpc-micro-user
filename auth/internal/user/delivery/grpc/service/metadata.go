package service

import (
	"auth/pkg/grpc_errors"
	"context"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

const (
	grpcGatewayUserAgentHeader = "grpcgateway-user-agent"
	userAgentHeader            = "user-agent"
	xForwardedForHeader        = "x-forwarded-for"
)

type Metadata struct {
	UserAgent string
	ClientIP  string
}

func (u *usersService) ExtractMetadata(ctx context.Context) (*Metadata, error) {
	mtdt := &Metadata{}

	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, grpc_errors.ErrNoCtxMetaData
	}

	if userAgents := md.Get(grpcGatewayUserAgentHeader); len(userAgents) > 0 {
		mtdt.UserAgent = userAgents[0]
	}

	if userAgents := md.Get(userAgentHeader); len(userAgents) > 0 {
		mtdt.UserAgent = userAgents[0]
	}

	if clientIPs := md.Get(xForwardedForHeader); len(clientIPs) > 0 {
		mtdt.ClientIP = clientIPs[0]
	}

	if p, ok := peer.FromContext(ctx); ok {
		mtdt.ClientIP = p.Addr.String()
	}

	return mtdt, nil
}
