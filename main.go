package main

import (
	"context"
	"database/sql"
	"net"
	"net/http"
	"os"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	_ "github.com/lib/pq"
	"github.com/rakyll/statik/fs"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/encoding/protojson"
	"tutorial.sqlc.dev/app/api"
	db "tutorial.sqlc.dev/app/db/sqlc"
	"tutorial.sqlc.dev/app/db/util"
	_ "tutorial.sqlc.dev/app/doc/statik"
	"tutorial.sqlc.dev/app/gapi"
	"tutorial.sqlc.dev/app/pb"
)

// const (
// 	dbDriver  = "postgres"
// 	dbSource  = "postgres://root:fHWFyt98gPR51h3NxjcroWoIscjt7QOb@dpg-cp649mmn7f5s73a6r8ag-a.oregon-postgres.render.com/simple_bank_7qc2"
// 	adrServer = "0.0.0.0:8080"
// )

func main() {
	var err error

	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal().Msg("cannot load config")
	}

	if config.Env == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal().Msg("cannot connect to db:")
	}

	store := db.NewsStore(conn)
	go runGatewayServer(config, store)
	rungRPCServer(config, store)

}

//use evans to call gRPC
//evans --host localhost --port 9090 repl

func rungRPCServer(config util.Config, store db.Store) {
	server, err := gapi.NewServer(config, store)
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

func runGatewayServer(config util.Config, store db.Store) {
	// Create a new server instance
	server, err := gapi.NewServer(config, store)
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
