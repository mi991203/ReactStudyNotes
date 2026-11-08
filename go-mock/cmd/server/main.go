// go-mock 入口：同时启动 gRPC 服务与 grpc-gateway(REST) 服务。
package main

import (
	"context"
	"flag"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/encoding/protojson"

	"go-mock/api"
	"go-mock/dao"
	"go-mock/logic"
	commentv1 "go-mock/protocols/comment/v1"
)

func main() {
	var (
		dbPath   = flag.String("db", "data/mock.db", "SQLite 数据库文件路径")
		grpcAddr = flag.String("grpc", "127.0.0.1:9090", "gRPC 监听地址")
		httpAddr = flag.String("http", ":3004", "HTTP(gateway) 监听地址")
		origin   = flag.String("cors-origin", "*", "允许的 CORS 来源")
	)
	flag.Parse()

	db, err := dao.Open(*dbPath)
	if err != nil {
		log.Fatalf("open db: %v", err)
	}
	defer db.Close()

	// ---- gRPC server ----
	lis, err := net.Listen("tcp", *grpcAddr)
	if err != nil {
		log.Fatalf("listen grpc: %v", err)
	}
	grpcSrv := grpc.NewServer()
	commentv1.RegisterCommentServiceServer(grpcSrv,
		api.NewServer(logic.NewCommentLogic(dao.NewCommentDAO(db))))
	go func() {
		log.Printf("gRPC listening on %s", *grpcAddr)
		if err := grpcSrv.Serve(lis); err != nil {
			log.Fatalf("grpc serve: %v", err)
		}
	}()

	// ---- grpc-gateway ----
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	mux := runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{
				UseProtoNames:   true, // 字段名与 proto 一致（rpid/ctime/like）
				EmitUnpopulated: true, // 零值也输出，如 like: 0
			},
			UnmarshalOptions: protojson.UnmarshalOptions{DiscardUnknown: true},
		}),
	)
	err = commentv1.RegisterCommentServiceHandlerFromEndpoint(ctx, mux, *grpcAddr,
		[]grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())})
	if err != nil {
		log.Fatalf("register gateway: %v", err)
	}

	httpSrv := &http.Server{
		Addr:              *httpAddr,
		Handler:           cors(*origin, mux),
		ReadHeaderTimeout: 5 * time.Second,
	}
	go func() {
		log.Printf("HTTP gateway listening on %s (db=%s)", *httpAddr, *dbPath)
		if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("http serve: %v", err)
		}
	}()

	// ---- 优雅退出 ----
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	<-sig
	log.Println("shutting down...")

	shutdownCtx, c := context.WithTimeout(context.Background(), 5*time.Second)
	defer c()
	_ = httpSrv.Shutdown(shutdownCtx)
	grpcSrv.GracefulStop()
}

// cors 允许前端(:3000)跨域访问，并直接应答预检请求。
func cors(origin string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := w.Header()
		h.Set("Access-Control-Allow-Origin", origin)
		h.Set("Access-Control-Allow-Methods", "GET, POST, DELETE, PUT, PATCH, OPTIONS")
		h.Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		next.ServeHTTP(w, r)
	})
}
