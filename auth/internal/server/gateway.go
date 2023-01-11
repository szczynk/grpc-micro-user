package server

import (
	"auth/pb"
	"context"
	"net/http"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func (s *Server) RunGateway() error {
	errCh := make(chan error)

	ctx, cancel := context.WithCancel(s.ctx)
	defer cancel()

	endpoint := s.cfg.Server.Port

	// Register gRPC server endpoint
	// Note: Make sure the gRPC server is running properly and accessible
	grpcMux := runtime.NewServeMux(
		runtime.WithHealthzEndpoint(s.healthClient),
	)
	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	err := pb.RegisterAuthServiceHandlerFromEndpoint(ctx, grpcMux, endpoint, opts)
	if err != nil {
		return err
	}

	mux := http.NewServeMux()
	// mount the gRPC HTTP gateway to the root
	mux.Handle("/", grpcMux)

	// mount a path to expose the generated OpenAPI specification on disk
	mux.HandleFunc("/swagger-ui/swagger.json", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./doc/swagger/auth.swagger.json")
	})

	// mount the Swagger UI that uses the OpenAPI specification path above
	mux.Handle("/swagger-ui/", http.StripPrefix("/swagger-ui/", http.FileServer(http.Dir("./doc/swagger-ui"))))

	server := &http.Server{Addr: s.cfg.Gateway.Port, Handler: mux}

	go func() {
		s.logger.Infof("Gateway available URL: %s", s.cfg.Gateway.Port)
		errCh <- server.ListenAndServe()
	}()

	select {
	case <-s.ctx.Done():
		c := make(chan bool)
		go func() {
			defer close(c)
			errCh <- server.Shutdown(s.ctx)
		}()
		select {
		case <-c:
		case <-time.After(5 * time.Second):
		}
		s.logger.Info("Gateway Exited Properly")
		return nil
	case err := <-errCh:
		return err
	}
}
