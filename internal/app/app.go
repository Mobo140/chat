package app

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/Mobo140/chat/internal/config"
	"github.com/Mobo140/chat/internal/interceptor"
	"github.com/Mobo140/chat/internal/ratelimiter"

	desc "github.com/Mobo140/chat/pkg/chat_v1"
	_ "github.com/Mobo140/chat/statik" // init statik
	"github.com/Mobo140/platform_common/pkg/closer"
	"github.com/Mobo140/platform_common/pkg/logger"
	"github.com/Mobo140/platform_common/pkg/tracing"
	grpcMiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/grpc-ecosystem/grpc-opentracing/go/otgrpc"
	"github.com/opentracing/opentracing-go"
	"github.com/rakyll/statik/fs"
	"github.com/rs/cors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	count           = 3
	countPerSecond  = 10
	logsMaxSize     = 10
	logsMaxBackups  = 3
	logsMaxAge      = 7
	chatServiceName = "chat_service"
	reqTimeout      = 5 * time.Second
)

type App struct {
	serviceProvider  *serviceProvider
	httpServer       *http.Server
	grpcServer       *grpc.Server
	grpcAccessClient *grpc.ClientConn
	swaggerServer    *http.Server
	configPath       string
	loggerLevel      string
}

func NewApp(ctx context.Context, configPath string, loggerLevel string) (*App, error) {
	a := &App{configPath: configPath, loggerLevel: loggerLevel}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initConfig,
		a.initLogger,
		a.initServiceProvider,
		a.initTracer,
		a.initHTTPServer,
		a.initGRPCAccessClient,
		a.initGRPCServer,
		a.initSwaggerServer,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initConfig(_ context.Context) error {
	flag.Parse()

	err := config.Load(a.configPath)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) initLogger(_ context.Context) error {
	logger.Init(getCore(getAtomicLevel(a.loggerLevel)))

	return nil
}

func (a *App) initTracer(_ context.Context) error {
	tracing.Init(logger.Logger(), chatServiceName, a.serviceProvider.JaegerConfig().Address())

	return nil
}

func getCore(level zap.AtomicLevel) zapcore.Core {
	stdout := zapcore.AddSync(os.Stdout)

	file := zapcore.AddSync(&lumberjack.Logger{
		Filename:   "logs/app.log",
		MaxSize:    logsMaxSize, // megabytes
		MaxBackups: logsMaxBackups,
		MaxAge:     logsMaxAge, // days
	})

	productionCfg := zap.NewProductionEncoderConfig()
	productionCfg.TimeKey = "timestamp"
	productionCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	developmentCfg := zap.NewDevelopmentEncoderConfig()
	developmentCfg.EncodeLevel = zapcore.CapitalColorLevelEncoder

	consoleEncoder := zapcore.NewConsoleEncoder(developmentCfg)
	fileEncoder := zapcore.NewJSONEncoder(productionCfg)

	return zapcore.NewTee(
		zapcore.NewCore(consoleEncoder, stdout, level),
		zapcore.NewCore(fileEncoder, file, level),
	)
}

func getAtomicLevel(logLevel string) zap.AtomicLevel {
	var level zapcore.Level
	if err := level.Set(logLevel); err != nil {
		log.Fatalf("failed to set log level: %v", err)
	}

	return zap.NewAtomicLevelAt(level)
}

func (a *App) initServiceProvider(_ context.Context) error {
	a.serviceProvider = newServiceProvider()
	return nil
}

func (a *App) initGRPCServer(ctx context.Context) error {
	creds, err := credentials.NewServerTLSFromFile("secure/service.pem", "secure/service.key")
	if err != nil {
		err = fmt.Errorf("failed to load TLS keys: %w", err)
		return err
	}

	rateLimiter := ratelimiter.NewTokenBucketLimiter(ctx, countPerSecond, time.Second)
	a.grpcServer = grpc.NewServer(grpc.Creds(creds),
		grpc.UnaryInterceptor(
			grpcMiddleware.ChainUnaryServer(
				interceptor.LogInterceptor,
				interceptor.ValidateInterceptor,
				interceptor.TimeoutUnaryServerInterceptor(reqTimeout),
				interceptor.NewRateLimiterInterceptor(rateLimiter).Unary,
				interceptor.ServerTracingInterceptor,
			),
		))

	reflection.Register(a.grpcServer)

	desc.RegisterChatV1Server(a.grpcServer, a.serviceProvider.ChatHandler(ctx, a.grpcAccessClient))

	return nil
}

func (a *App) initGRPCAccessClient(_ context.Context) error {
	creds, err := credentials.NewClientTLSFromFile("secure/auth.pem", "")
	if err != nil {
		log.Fatalf("failed to load TLS keys for access client: %v", err)
	}

	cl, err := grpc.NewClient(
		a.serviceProvider.AccessClientConfig().Address(),
		grpc.WithTransportCredentials(creds),
		grpc.WithUnaryInterceptor(
			otgrpc.OpenTracingClientInterceptor(opentracing.GlobalTracer()),
		),
	)

	closer.Add(cl.Close)

	a.grpcAccessClient = cl

	if err != nil {
		log.Fatalf("failed to dial gRPC access client: %v", err)
	}

	return nil
}

func (a *App) initHTTPServer(ctx context.Context) error {
	mux := runtime.NewServeMux()

	creds, err := credentials.NewClientTLSFromFile("secure/service.pem", "")
	if err != nil {
		log.Fatalf("failed to load TLS keys: %v", err)
	}

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(creds),
	}

	err = desc.RegisterChatV1HandlerFromEndpoint(ctx, mux, a.serviceProvider.GRPCConfig().Address(), opts)
	if err != nil {
		return err
	}

	corsMiddleware := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "Content-length"},
		AllowCredentials: true,
	})

	a.httpServer = &http.Server{
		Addr:    a.serviceProvider.HTTPConfig().Address(),
		Handler: corsMiddleware.Handler(mux),
	}

	return nil
}

func (a *App) initSwaggerServer(_ context.Context) error {
	statikFs, err := fs.New()
	if err != nil {
		return err
	}

	mux := http.NewServeMux()
	mux.Handle("/", http.StripPrefix("/", http.FileServer(statikFs)))
	mux.HandleFunc("/api.swagger.json", serveSwaggerFile("/api.swagger.json"))

	a.swaggerServer = &http.Server{
		Addr:    a.serviceProvider.SwaggerConfig().Address(),
		Handler: mux,
	}

	return nil
}

func serveSwaggerFile(path string) http.HandlerFunc {
	return func(w http.ResponseWriter, _ *http.Request) {
		log.Printf("Serving swagger file: %s", path)

		statikFs, err := fs.New()
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		log.Printf("Open swagger file: %s", path)

		file, err := statikFs.Open(path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		defer file.Close()

		log.Printf("Read swagger file: %s", path)

		content, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		log.Printf("Write swagger file: %s", path)

		w.Header().Set("Content-Type", "application/json")

		_, err = w.Write(content)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		log.Printf("Served swagger file: %s", path)
	}
}

func (a *App) Run() error {
	defer func() {
		closer.CloseAll()
		closer.Wait()
	}()

	wg := sync.WaitGroup{}
	wg.Add(count)

	go func() {
		defer wg.Done()

		err := a.runGRPCServer()
		if err != nil {
			log.Fatalf("failed to run GRPC server: %v", err)
		}
	}()

	go func() {
		defer wg.Done()

		err := a.runHTTPServer()
		if err != nil {
			log.Fatalf("failed to run HTTP server: %v", err)
		}
	}()

	go func() {
		defer wg.Done()

		err := a.runSwaggerServer()
		if err != nil {
			log.Fatalf("failed to run Swagger server: %v", err)
		}
	}()

	wg.Wait()

	return nil
}

func (a *App) runGRPCServer() error {
	log.Printf("GRPC server is running on: %s", a.serviceProvider.GRPCConfig().Address())

	list, err := net.Listen("tcp", a.serviceProvider.GRPCConfig().Address())
	if err != nil {
		return err
	}

	err = a.grpcServer.Serve(list)
	if err != nil {
		return err
	}

	return nil
}

func (a *App) runHTTPServer() error {
	log.Printf("HTTP server is running on: %s", a.serviceProvider.HTTPConfig().Address())

	err := a.httpServer.ListenAndServeTLS("secure/service.pem", "secure/service.key")
	if err != nil {
		return err
	}

	return nil
}

func (a *App) runSwaggerServer() error {
	log.Printf("Swagger server is running on: %s", a.serviceProvider.SwaggerConfig().Address())

	err := a.swaggerServer.ListenAndServeTLS("secure/service.pem", "secure/service.key")
	if err != nil {
		return err
	}

	return nil
}
