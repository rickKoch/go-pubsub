package main

import (
	"context"
	"fmt"
	"net"
	"os"

	"github.com/golang/protobuf/ptypes/empty"
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	go_pubsub "github.com/rickKoch/go-pubsub/genproto/pubsub"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func RunGRPCServer(logger *logrus.Entry, registerServer func(server *grpc.Server)) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	addr := fmt.Sprintf(":%s", port)
	RunGRPCServerOnAddr(logger, addr, registerServer)
}

func RunGRPCServerOnAddr(logger *logrus.Entry, addr string, registerServer func(server *grpc.Server)) {
	grpcServer := grpc.NewServer(
		grpc_middleware.WithUnaryServerChain(
			grpc_ctxtags.UnaryServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
			grpc_logrus.UnaryServerInterceptor(logger),
		),
		grpc_middleware.WithStreamServerChain(
			grpc_ctxtags.StreamServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
			grpc_logrus.StreamServerInterceptor(logger),
		),
	)
	registerServer(grpcServer)

	listen, err := net.Listen("tcp", addr)
	if err != nil {
		logrus.Fatal(err)
	}
	logrus.WithField("grpcEndpoint", addr).Info("Starting: gRPC Listener")
	logrus.Fatal(grpcServer.Serve(listen))
}

type GrpcServer struct {
	queueManager QueueManager
}

func NewGrpcServer(qm QueueManager) GrpcServer {
	return GrpcServer{queueManager: qm}
}

func (g GrpcServer) CreateTopic(ctx context.Context, request *go_pubsub.CreateTopicRequest) (*empty.Empty, error) {
	if err := g.queueManager.CreateTopic(request.Name); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &empty.Empty{}, nil
}

func (g GrpcServer) CreateSubscription(ctx context.Context, request *go_pubsub.CreateSubscriptionRequest) (*empty.Empty, error) {
	if err := g.queueManager.CreateSubscription(request.TopicName, request.SubscriptionName); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &empty.Empty{}, nil
}

func (g GrpcServer) PublishMessage(ctx context.Context, request *go_pubsub.PublishMessageRequest) (*empty.Empty, error) {
	g.queueManager.Emit(request.Topic, request.Message)

	return &empty.Empty{}, nil
}

func (g GrpcServer) PullMessages(ctx context.Context, request *go_pubsub.PullMessagesRequest) (*go_pubsub.PullMessagesResponse, error) {
	sub, err := g.queueManager.Subscription(request.Subscription)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var pullMessages []*go_pubsub.Message
	for _, msg := range sub.PullMessages() {
		pullMessages = append(pullMessages, &go_pubsub.Message{Data: msg.data})
	}

	return &go_pubsub.PullMessagesResponse{Messages: pullMessages}, nil
}

func (g GrpcServer) PullMessagesStreaming(request *go_pubsub.PullMessagesRequest, stream go_pubsub.PubSubService_PullMessagesStreamingServer) error {
	sub, err := g.queueManager.Subscription(request.Subscription)
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	con, err := sub.CreateConnection()
	if err != nil {
		return status.Error(codes.Internal, err.Error())
	}

	for _, msg := range sub.PullMessages() {
		if err := stream.Send(&go_pubsub.Message{Data: msg.data}); err != nil {
			return err
		}
	}

	for {
		select {
		case msg := <-con.message:
			if err := stream.Send(&go_pubsub.Message{Data: msg.data}); err != nil {
				return err
			}
		case <-stream.Context().Done():
			sub.RemoveConnection(con.id)
			return nil
		}
	}
}
