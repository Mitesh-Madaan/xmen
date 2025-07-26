package main

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"time"

	xRouter "gomike/router"
	xSession "gomike/session"
	xDb "lib/dbchef"

	"github.com/google/uuid"
)

func main() {
	run()
}

var dbSession *xDb.DBSession

const (
	connStr = "host=localhost user=postgres password=Postsql.123 dbname=postgres port=5432 sslmode=disable TimeZone=Asia/Shanghai"
)

func run() {
	var err error
	dbSession = xSession.GetDBSession(connStr)
	if dbSession == nil {
		fmt.Println("Failed to get DB session")
		return
	}
	fmt.Printf("DB Session: %v\n", dbSession)

	err = xSession.SeedTables(dbSession)
	if err != nil {
		fmt.Printf("Error seeding tables: %v\n", err)
		return
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	mux := http.NewServeMux()

	mux.Handle("GET /person/{reqObjID}", handleRequest(handleWithLogger(logger, handleWithAuth(authMiddleware(logger), handleWithRouter(logger, xRouter.GetPerson(dbSession, logger))))))
	mux.Handle("PUT /person/{reqObjID}", handleRequest(handleWithLogger(logger, handleWithAuth(authMiddleware(logger), handleWithRouter(logger, xRouter.UpdatePerson(dbSession, logger))))))
	mux.Handle("POST /person/", handleRequest(handleWithLogger(logger, handleWithAuth(authMiddleware(logger), handleWithRouter(logger, xRouter.CreatePerson(dbSession, logger))))))
	mux.Handle("PATCH /person/{reqObjID}", handleRequest(handleWithLogger(logger, handleWithAuth(authMiddleware(logger), handleWithRouter(logger, xRouter.PatchPerson(dbSession, logger))))))
	mux.Handle("DELETE /person/{reqObjID}", handleRequest(handleWithLogger(logger, handleWithAuth(authMiddleware(logger), handleWithRouter(logger, xRouter.DeletePerson(dbSession, logger))))))

	mux.Handle("GET /animal/{reqObjID}", handleRequest(handleWithLogger(logger, handleWithAuth(authMiddleware(logger), handleWithRouter(logger, xRouter.GetAnimal(dbSession, logger))))))
	mux.Handle("PUT /animal/{reqObjID}", handleRequest(handleWithLogger(logger, handleWithAuth(authMiddleware(logger), handleWithRouter(logger, xRouter.UpdateAnimal(dbSession, logger))))))
	mux.Handle("POST /animal/", handleRequest(handleWithLogger(logger, handleWithAuth(authMiddleware(logger), handleWithRouter(logger, xRouter.CreateAnimal(dbSession, logger))))))
	mux.Handle("PATCH /animal/{reqObjID}", handleRequest(handleWithLogger(logger, handleWithAuth(authMiddleware(logger), handleWithRouter(logger, xRouter.PatchAnimal(dbSession, logger))))))
	mux.Handle("DELETE /animal/{reqObjID}", handleRequest(handleWithLogger(logger, handleWithAuth(authMiddleware(logger), handleWithRouter(logger, xRouter.DeleteAnimal(dbSession, logger))))))

	fmt.Println("Starting server at port 8080")
	if err = http.ListenAndServe("localhost:8080", mux); err != nil {
		fmt.Println(err)
	}

}

func handleRequest(nextHandler http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, req *http.Request) {
			requestID := uuid.New().String()
			ctx := context.WithValue(req.Context(), xRouter.RequestIDKey("requestID"), requestID)
			req = req.WithContext(ctx)
			nextHandler.ServeHTTP(w, req)
		},
	)
}

func handleWithLogger(log *slog.Logger, nextHandler http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, req *http.Request) {
			log.Info("Received request",
				slog.String("request-id", req.Context().Value(xRouter.RequestIDKey("requestID")).(string)),
				slog.String("method", req.Method),
				slog.String("url", req.URL.String()),
			)
			nextHandler.ServeHTTP(w, req)
		},
	)
}

func handleWithAuth(authFunc func(context.Context, string) bool, nextHandler http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, req *http.Request) {
			isAuthorized := authFunc(req.Context(), req.Header.Get("Authorization"))
			if !isAuthorized {
				w.Header().Set("Content-Type", "text/plain; charset=utf-8")
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Unauthorized"))
				return
			}
			nextHandler.ServeHTTP(w, req)
		},
	)
}

func handleWithRouter(log *slog.Logger, routerFunc func(context.Context, string, io.ReadCloser) xRouter.RespDetail) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, req *http.Request) {
			timeoutCtx, cancelCtx := context.WithTimeout(req.Context(), 2*time.Second)
			defer cancelCtx()
			req = req.WithContext(timeoutCtx)
			respChan := make(chan xRouter.RespDetail, 1)

			go func() {
				defer close(respChan)
				respChan <- routerFunc(req.Context(), req.PathValue("reqObjID"), req.Body)
			}()

			select {
			case <-timeoutCtx.Done():
				log.Error("Request timed out", slog.String("request-id", req.Context().Value(xRouter.RequestIDKey("requestID")).(string)))
				w.Header().Set("Content-Type", "text/plain; charset=utf-8")
				w.WriteHeader(http.StatusRequestTimeout)
				w.Write([]byte("Request timed out"))
			case resp := <-respChan:
				log.Info("Response received", slog.String("request-id", req.Context().Value(xRouter.RequestIDKey("requestID")).(string)))
				switch resp.Type {
				case "text/plain":
					w.Header().Set("Content-Type", "text/plain; charset=utf-8")
				case "application/json":
					w.Header().Set("Content-Type", "application/json; charset=utf-8")
				}
				w.WriteHeader(resp.Statuscode)
				w.Write(resp.Message)
			}
		},
	)
}

func authMiddleware(log *slog.Logger) func(context.Context, string) bool {
	return func(reqCtx context.Context, authHeader string) bool {
		if authHeader != "Basic bWl0ZXNoOk1pdGVzaC4xMjM=" {
			log.Error("Unauthorized", slog.String("request-id", reqCtx.Value(xRouter.RequestIDKey("requestID")).(string)))
			return false
		}
		log.Info("Authorization successful", slog.String("request-id", reqCtx.Value(xRouter.RequestIDKey("requestID")).(string)))
		return true
	}
}
