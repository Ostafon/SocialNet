package main

import (
	"log"
	"net/http"
	"time"

	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"google.golang.org/grpc"
)

func ChatGrpcWebWrapper(grpcServer *grpc.Server) {
	wrapped := grpcweb.WrapServer(
		grpcServer,
		grpcweb.WithOriginFunc(func(origin string) bool { return true }),
		grpcweb.WithWebsockets(true),
		grpcweb.WithWebsocketPingInterval(30*time.Second),
	)

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers",
			"Content-Type, X-Grpc-Web, X-User-Agent, User-Agent, Authorization, user-id, x-user-id, grpc-timeout",
		)

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		if wrapped.IsGrpcWebRequest(r) ||
			wrapped.IsGrpcWebSocketRequest(r) {
			wrapped.ServeHTTP(w, r)
			return
		}

		http.NotFound(w, r)
	})

	srv := &http.Server{
		Addr:    "0.0.0.0:8082",
		Handler: handler,
	}

	log.Println("🌐 GRPC-WEB running on :8082")
	log.Fatal(srv.ListenAndServe())
}
