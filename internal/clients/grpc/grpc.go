package grpc

import (
	"context"
	"fmt"
	grpclog "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpcretry "github.com/grpc-ecosystem/go-grpc-middleware/retry"
	ssov1 "github.com/markgregr/FruitfulFriends-protos/gen/go/sso"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"time"
)

type Client struct {
	AuthService ssov1.AuthClient
	log         *logrus.Entry
}

func New(ctx context.Context, log *logrus.Entry, targetAddr string, timeout time.Duration, retriesCount int) (*Client, error) {
	const op = "grpc.New"
	log = log.WithField("operation", op)
	log.Infof("dialing to %s", targetAddr)

	retryOpts := []grpcretry.CallOption{
		grpcretry.WithCodes(codes.Unavailable, codes.DeadlineExceeded, codes.Internal, codes.ResourceExhausted, codes.PermissionDenied, codes.Unimplemented),
		grpcretry.WithMax(uint(retriesCount)),
		grpcretry.WithBackoff(grpcretry.BackoffLinear(timeout)),
		grpcretry.WithPerRetryTimeout(timeout),
	}

	log.Infof("retry options: %+v", retryOpts)

	conn, err := grpc.DialContext(ctx, targetAddr, grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithChainUnaryInterceptor(
			grpclog.UnaryClientInterceptor(log),
			grpcretry.UnaryClientInterceptor(retryOpts...),
		))
	if err != nil {
		log.WithError(err).Errorf("%s: dial: failed to dial to %s", op, targetAddr)
		return nil, fmt.Errorf("%s: dial: %w", op, err)
	}

	return &Client{
		AuthService: ssov1.NewAuthClient(conn),
		log:         log,
	}, nil
}
