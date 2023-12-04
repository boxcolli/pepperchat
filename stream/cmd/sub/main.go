package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	pb "github.com/boxcolli/go-transistor/api/gen/transistor/v1"
	"github.com/boxcolli/go-transistor/base/basicbase"
	"github.com/boxcolli/go-transistor/collector/basiccollector"
	"github.com/boxcolli/go-transistor/emitter/basicemitter"
	"github.com/boxcolli/go-transistor/index/basicindex"
	"github.com/boxcolli/go-transistor/server/grpcserver"
	"github.com/boxcolli/go-transistor/transistor"
	"github.com/boxcolli/go-transistor/transistor/simpletransistor"
	"github.com/boxcolli/go-transistor/types"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/peterbourgon/ff/v4"
	"google.golang.org/grpc"
)

func main() {
	// Parse flag, env var
	fs := flag.NewFlagSet("sub", flag.ContinueOnError)
	var (
		addr = fs.String("addr", ":50050", "address of the publisher server")
		topic = fs.String("topic", "chat", "static topic to subscribe on")
		port = fs.String("port", "50051", "listen port")
		cmqs = fs.Int("cmqs", 100, "collector message queue size")
		bcqs = fs.Int("bcqs", 100, "base change queue size")
		emqs = fs.Int("emqs", 100, "emitter message queue size")
	)
	ff.Parse(fs, os.Args[1:],
		ff.WithEnvVars(),
	)

	// Client
	var client pb.TransistorServiceClient
	{
		var conn *grpc.ClientConn
		var err error
		for {
			conn, err = grpc.Dial(*addr, dialOpts...)
			if err != nil {
				time.Sleep(time.Second * 2)
				continue
			}
			break
		}
		defer conn.Close()
		client = pb.NewTransistorServiceClient(conn)
		fmt.Println("connected with pub transistor")
	}

	// Subscribe
	{
		var opts = []grpc.CallOption{}
		stream, err := client.Subscribe(context.Background(), opts...)
		if err != nil {
			panic(err)
		}

		// Send initial change
		cg := types.Change{ Op: types.OperationAdd, Topic: types.Topic{*topic} }
		err = stream.Send(&pb.SubscribeRequest{
			Change: cg.Marshal(),
		})
		if err != nil {
			panic(err)
		}

		fmt.Println("now listening..")
		for {
			res, err := stream.Recv()
			if err != nil {
				log.Fatalf("Subscribe() received error: %v\n", err)
				break
			}

			msg := new(types.Message)
			msg.Unmarshal(res.GetMsg())
			log.Printf("Subscribe() receivd: %s\n", msg.String())
		}
	}

	// Transistor
	var tr transistor.Transistor
	{
		tr = simpletransistor.NewSimpleCore(
			transistor.Component{
				Base: basicbase.NewBasicBase(basicindex.NewBasicIndex, *bcqs),
				Collector: basiccollector.NewBasicCollector(*cmqs),
				Emitter: basicemitter.NewBasicEmitter(*emqs),
			},
			simpletransistor.Option{},
		)
		log.Println("tr started.")
	}

	// Server
	{
		lis, err := net.Listen("tcp", fmt.Sprintf(":%s", *port))
		if err != nil {
			logger.Fatalf("failed to listen: %v", err)
		}
		grpcServer := grpc.NewServer(
			grpc.ChainUnaryInterceptor(
				logging.UnaryServerInterceptor(InterceptorLogger(logger), logOpts...),
			),
			grpc.ChainStreamInterceptor(
				logging.StreamServerInterceptor(InterceptorLogger(logger), logOpts...),
			),
		)
		pb.RegisterTransistorServiceServer(grpcServer, grpcserver.NewGrpcServer(tr))
		logger.Printf("server listening at %v", lis.Addr())
		if err := grpcServer.Serve(lis); err != nil {
			logger.Fatalf("failed to serve: %v", err)
		}
	}
}
