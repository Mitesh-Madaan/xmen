package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"

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

	mux.Handle("GET /person/{personID}", handleRequest(handleWithLogger(logger, handleWithAuth(authMiddleware(logger), handleWithRouter(xRouter.GetPerson(dbSession, logger))))))
	mux.Handle("PUT /person/{personID}", handleRequest(handleWithLogger(logger, handleWithAuth(authMiddleware(logger), handleWithRouter(xRouter.UpdatePerson(dbSession, logger))))))
	mux.Handle("POST /person/", handleRequest(handleWithLogger(logger, handleWithAuth(authMiddleware(logger), handleWithRouter(xRouter.CreatePerson(dbSession, logger))))))
	mux.Handle("PATCH /person/{personID}", handleRequest(handleWithLogger(logger, handleWithAuth(authMiddleware(logger), handleWithRouter(xRouter.PatchPerson(dbSession, logger))))))
	mux.Handle("DELETE /person/{personID}", handleRequest(handleWithLogger(logger, handleWithAuth(authMiddleware(logger), handleWithRouter(xRouter.DeletePerson(dbSession, logger))))))

	mux.Handle("GET /animal/{animalID}", handleRequest(handleWithLogger(logger, handleWithAuth(authMiddleware(logger), handleWithRouter(xRouter.GetAnimal(dbSession, logger))))))
	mux.Handle("PUT /animal/{animalID}", handleRequest(handleWithLogger(logger, handleWithAuth(authMiddleware(logger), handleWithRouter(xRouter.UpdateAnimal(dbSession, logger))))))
	mux.Handle("POST /animal/", handleRequest(handleWithLogger(logger, handleWithAuth(authMiddleware(logger), handleWithRouter(xRouter.CreateAnimal(dbSession, logger))))))
	mux.Handle("PATCH /animal/{animalID}", handleRequest(handleWithLogger(logger, handleWithAuth(authMiddleware(logger), handleWithRouter(xRouter.PatchAnimal(dbSession, logger))))))
	mux.Handle("DELETE /animal/{animalID}", handleRequest(handleWithLogger(logger, handleWithAuth(authMiddleware(logger), handleWithRouter(xRouter.DeleteAnimal(dbSession, logger))))))

	fmt.Println("Starting server at port 8080")
	if err = http.ListenAndServe("localhost:8080", mux); err != nil {
		fmt.Println(err)
	}

}

func handleRequest(logHandler http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, req *http.Request) {
			requestID := uuid.New().String()
			ctx := context.WithValue(req.Context(), xRouter.RequestIDKey("requestID"), requestID)
			req = req.WithContext(ctx)
			logHandler.ServeHTTP(w, req)
		},
	)
}

func handleWithLogger(log *slog.Logger, authHandler http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, req *http.Request) {
			log.Info("Received request",
				slog.String("request-id", req.Context().Value(xRouter.RequestIDKey("requestID")).(string)),
				slog.String("method", req.Method),
				slog.String("url", req.URL.String()),
			)
			authHandler.ServeHTTP(w, req)
		},
	)
}

func handleWithAuth(authFunc func(http.ResponseWriter, *http.Request), nextHandler http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, req *http.Request) {
			authFunc(w, req)
			nextHandler.ServeHTTP(w, req)
		},
	)
}

func handleWithRouter(routerFunc func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, req *http.Request) {
			routerFunc(w, req)
		},
	)
}

func authMiddleware(log *slog.Logger) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		authHeader := req.Header.Get("Authorization")
		if authHeader != "Basic bWl0ZXNoOk1pdGVzaC4xMjM=" {
			log.Error("Unauthorized", slog.String("request-id", req.Context().Value(xRouter.RequestIDKey("requestID")).(string)))
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized"))
			return
		}
		log.Info("Authorization successful", slog.String("request-id", req.Context().Value(xRouter.RequestIDKey("requestID")).(string)))
	}
}
