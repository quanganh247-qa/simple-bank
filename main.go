package main

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rakyll/statik/fs"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/techschool/simplebank/mail"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
	"tutorial.sqlc.dev/app/api"
	db "tutorial.sqlc.dev/app/db/sqlc"
	"tutorial.sqlc.dev/app/db/util"
	_ "tutorial.sqlc.dev/app/doc/statik"
	"tutorial.sqlc.dev/app/gapi"
	"tutorial.sqlc.dev/app/pb"
	"tutorial.sqlc.dev/app/worker"
)

var interruptSignals = []os.Signal{
	os.Interrupt,
	syscall.SIGTERM,
	syscall.SIGINT,
}

func main() {
	var err error

	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal().Msg("cannot load config")
	}

	if config.Env == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
	//********************************************************************
	//********************************************************************
	ctx, stop := signal.NotifyContext(context.Background(), interruptSignals...)
	defer stop()
	connPool, err := pgxpool.New(ctx, config.DBSource)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot connect to db")
	}
	store := db.NewStore(connPool)
	//********************************************************************
	//********************************************************************
	redisOpt := asynq.RedisClientOpt{
		Addr: config.RedisAddress,
	}

	taskDistributor := worker.NewRedisTaskDistributor(redisOpt)

	go runTaskProccessor(config, redisOpt, store)
	go runGatewayServer(config, store, taskDistributor)
	rungRPCServer(config, store, taskDistributor)

}

//use evans to call gRPC
//evans --host localhost --port 9090 repl

func rungRPCServer(config util.Config, store db.Store, taskDistributor worker.TaskDistributor) {
	server, err := gapi.NewServer(config, store, taskDistributor)
	if err != nil {
		log.Fatal().Msg("cannot connect server")
	}

	grpcLogger := grpc.UnaryInterceptor(gapi.GrpcLogger)
	grpcServer := grpc.NewServer(grpcLogger)
	pb.RegisterSimpleBankServer(grpcServer, server)

	//Allow grpc to ealisy explore what rpc availiable in server
	reflection.Register(grpcServer)

	// Create a new listener on a specified port (9090)
	listener, err := net.Listen("tcp", config.GRPCServerAddress)
	if err != nil {
		log.Fatal().Msg("cannot listen server")
	}

	log.Info().Msgf("start gRPC server at %s", listener.Addr().String())

	//Start serving requests
	err = grpcServer.Serve(listener)

	if err != nil {
		log.Fatal().Msg("cannot start gRPC server")
	}

}

func runGinServer(config util.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal().Msg("cannot create server: ")
	}
	err = server.Starts(config.HTTPServerAddress)
	if err != nil {
		log.Fatal().Msg("cannot start server")
	}
}

func runGatewayServer(config util.Config, store db.Store, taskDistributor worker.TaskDistributor) {
	// Create a new server instance
	server, err := gapi.NewServer(config, store, taskDistributor)
	if err != nil {
		log.Fatal().Msg("cannot connect server")
	}

	jsonOption := runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
		MarshalOptions: protojson.MarshalOptions{
			UseProtoNames: true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	})

	// Create a gRPC-Gateway mux
	grpcMux := runtime.NewServeMux(jsonOption)

	// Create a context with cancel
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Register gRPC server handler to the gRPC-Gateway mux
	err = pb.RegisterSimpleBankHandlerServer(ctx, grpcMux, server)
	if err != nil {
		log.Fatal().Msg("cannot register handler server: ")
	}

	// Create an HTTP mux
	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	statikFS, err := fs.New()
	if err != nil {
		log.Fatal().Msg("cannot create statik fs")
	}

	swaggerHandler := http.StripPrefix("/swagger/", http.FileServer(statikFS))

	mux.Handle("/swagger/", swaggerHandler)
	// Create a new listener on a specified port (9090)
	listener, err := net.Listen("tcp", config.HTTPServerAddress)
	if err != nil {
		log.Fatal().Msg("cannot start HTTP gateway server: ")
	}

	log.Printf("start HTTP server at %s", listener.Addr().String())
	handler := gapi.HttpLogger(mux)
	//Start serving requests
	err = http.Serve(listener, handler)

	if err != nil {
		log.Fatal().Msg("cannot start gRPC server")
	}

}

func runTaskProccessor(config util.Config, redisOpt asynq.RedisClientOpt, store db.Store) {
	mailer := mail.NewGmailSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword)
	taskProccessor := worker.NewRedisTaskProccessor(redisOpt, store, mailer)
	log.Info().Msg("start task poccessor")
	err := taskProccessor.Start()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to start task proccesor")
	}
}
