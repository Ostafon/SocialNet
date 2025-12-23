package main

import (
	"log"
	"net/http"

	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"google.golang.org/grpc"
)

func StartGrpcWebWrapper(grpcServer *grpc.Server) {
	wrapped := grpcweb.WrapServer(
		grpcServer,
		grpcweb.WithOriginFunc(func(origin string) bool { return true }),
		grpcweb.WithAllowNonRootResource(true),
		grpcweb.WithWebsockets(true),
	)

	httpServer := &http.Server{
		Addr: "0.0.0.0:8081",
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			w.Header().Set("access-control-allow-origin", "*")
			w.Header().Set("access-control-allow-credentials", "true")
			w.Header().Set("access-control-allow-methods", "GET, POST, OPTIONS")
			w.Header().Set("access-control-allow-headers",
				"content-type, x-grpc-web, x-user-agent, x-user-id, user-id, authorization, grpc-timeout",
			)

			w.Header().Set("access-control-expose-headers",
				"grpc-status, grpc-message, grpc-status-details-bin",
			)

			if r.Method == "POST" {
				w.Header().Set("Content-Type", "application/grpc-web+proto")
			}

			if r.Method == http.MethodOptions {
				w.WriteHeader(200)
				return
			}

			if wrapped.IsGrpcWebRequest(r) ||
				wrapped.IsAcceptableGrpcCorsRequest(r) {
				wrapped.ServeHTTP(w, r)
				return
			}

			w.WriteHeader(404)
		}),
	}

	log.Println("🌐 grpc-web wrapper started on :8081")
	if err := httpServer.ListenAndServe(); err != nil {
		log.Fatalf("grpc-web wrapper error: %v", err)
	}
}
