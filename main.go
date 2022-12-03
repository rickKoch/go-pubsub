package main

import (
	"context"

	go_pubsub "github.com/rickKoch/go-pubsub/genproto/pubsub"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	logInit()
	logrusEntry := logrus.NewEntry(logrus.StandardLogger())
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	qm := NewQueueManager(ctx, logrusEntry)

	RunGRPCServer(logrusEntry, func(server *grpc.Server) {
		svc := NewGrpcServer(qm)
		go_pubsub.RegisterPubSubServiceServer(server, svc)
    go svc.queueManager.Run()
		reflection.Register(server)
	})
}
