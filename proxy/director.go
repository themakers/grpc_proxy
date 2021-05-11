// Copyright 2017 Michal Witkowski. All Rights Reserved.
// See LICENSE for licensing terms.

package proxy

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// StreamDirector returns a gRPC ClientConn to be used to forward the call to.
//
// The presence of the `Context` allows for rich filtering, e.g. based on Metadata (headers).
// If no handling is meant to be done, a `codes.NotImplemented` gRPC error should be returned.
//
// The context returned from this function should be the context for the *outgoing* (to backend) call. In case you want
// to forward any Metadata between the inbound request and outbound requests, you should do it manually. However, you
// *must* propagate the cancel function (`context.WithCancel`) of the inbound context to the one returned.
//
// It is worth noting that the StreamDirector will be fired *after* all server-side stream interceptors
// are invoked. So decisions around authorization, monitoring etc. are better to be handled there.
//
// See the rather rich example.
//type StreamDirector func(ctx context.Context, fullMethodName string) (context.Context, *grpc.ClientConn, error)

type StreamDirector interface {
	// Connect returns a connection to use for the given method,
	// or an error if the call should not be handled.
	//
	// The provided context may be inspected for filtering on request
	// metadata.
	//
	// Method is the gRPC request path, which is in the form "/service/method".
	//
	// The returned context is used as the basis for the outgoing connection.
	Connect(ctx context.Context, method string) (context.Context, *grpc.ClientConn, error)

	// Release is called when a connection is longer being used.  This is called
	// once for every call to Connect that does not return an error.
	//
	// The provided context is the one returned from Connect.
	//
	// This can be used by the director to pool connections or close unused
	// connections.
	Release(ctx context.Context, conn *grpc.ClientConn)
}
