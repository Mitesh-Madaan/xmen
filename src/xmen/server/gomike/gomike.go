package main

import (
	"context"
	"fmt"
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
	fmt.Println("This is a placeholder for the main function")
	run()
	fmt.Println("End of the main function")
}

/*
TODO:
- Add UTs
- Minimum viable product (binaries, tests, documentation)
- Publish as docker image
- Intergration tests

Code Review:
- Auth: Pending middleware for authentication DONE
- Change UUID to string > Check GORM documentation DONE
- Different handlers for different methods DONE
- Fix ID parsing from URL {placeholders} DONE
- Parse object from request body directly DONE
- Create model independent, not member of model class ; CRUD operations indenpendent of models DONE
- API field validation : Check field function (optional) ; Explore OpenAPI 3 schema ; Generate models from schema automatically DEFERRED
- Just pass object in db methods, no need to pass model DONE
- Multiline string in ToString method > use %v DONE
- To Status > no need to implement, just return object as json DONE
- Return resp of API as json marshall of object DONE
- Logging? log/Slog structured logging go package DONE
- Context? > Use context for request ID, timeout, cancellation, error handling , and logging

Next:
- UT: Check ways to unit test http calls
- MVP
*/

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

	mux.Handle("GET /person/{personID}", handleWithLogger(logger, handleWithMiddleware(authMiddleware, http.HandlerFunc(xRouter.GetPerson(dbSession, logger)), logger)))
	mux.Handle("PUT /person/{personID}", handleWithLogger(logger, handleWithMiddleware(authMiddleware, http.HandlerFunc(xRouter.UpdatePerson(dbSession, logger)), logger)))
	mux.Handle("POST /person/", handleWithLogger(logger, handleWithMiddleware(authMiddleware, http.HandlerFunc(xRouter.CreatePerson(dbSession, logger)), logger)))
	mux.Handle("PATCH /person/{personID}", handleWithLogger(logger, handleWithMiddleware(authMiddleware, http.HandlerFunc(xRouter.PatchPerson(dbSession, logger)), logger)))
	mux.Handle("DELETE /person/{personID}", handleWithLogger(logger, handleWithMiddleware(authMiddleware, http.HandlerFunc(xRouter.DeletePerson(dbSession, logger)), logger)))

	mux.Handle("GET /animal/{animalID}", handleWithLogger(logger, handleWithMiddleware(authMiddleware, http.HandlerFunc(xRouter.GetAnimal(dbSession, logger)), logger)))
	mux.Handle("PUT /animal/{animalID}", handleWithLogger(logger, handleWithMiddleware(authMiddleware, http.HandlerFunc(xRouter.UpdateAnimal(dbSession, logger)), logger)))
	mux.Handle("POST /animal/", handleWithLogger(logger, handleWithMiddleware(authMiddleware, http.HandlerFunc(xRouter.CreateAnimal(dbSession, logger)), logger)))
	mux.Handle("PATCH /animal/{animalID}", handleWithLogger(logger, handleWithMiddleware(authMiddleware, http.HandlerFunc(xRouter.PatchAnimal(dbSession, logger)), logger)))
	mux.Handle("DELETE /animal/{animalID}", handleWithLogger(logger, handleWithMiddleware(authMiddleware, http.HandlerFunc(xRouter.DeleteAnimal(dbSession, logger)), logger)))

	fmt.Println("Starting server at port 8080")
	if err = http.ListenAndServe("localhost:8080", mux); err != nil {
		fmt.Println(err)
	}

}

func handleWithLogger(log *slog.Logger, middlewareHandler http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, req *http.Request) {
			requestID := uuid.New().String()
			ctx, _ := context.WithTimeout(req.Context(), 2*time.Second)
			ctx = context.WithValue(ctx, xRouter.RequestIDKey("requestID"), requestID)
			req = req.WithContext(ctx)
			log.Info("Handling request method", slog.String("request-id", req.Context().Value(xRouter.RequestIDKey("requestID")).(string)), slog.String("method", req.Method), slog.String("Path", req.URL.Path))
			responseChan := make(chan struct{})
			go func() {
				defer func() {
					close(responseChan)
				}()
				middlewareHandler.ServeHTTP(w, req)
				responseChan <- struct{}{}
			}()

			select {
			case <-req.Context().Done():
				if req.Context().Err() == context.DeadlineExceeded {
					w.WriteHeader(http.StatusRequestTimeout)
					w.Write([]byte("Request timed out"))
					log.Error("Request timed out due to deadline exceeded", slog.String("request-id", req.Context().Value(xRouter.RequestIDKey("requestID")).(string)))
					return
				} else {
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte("Request failed due to context cancellation"))
					log.Error("Request failed due to context cancellation", slog.String("request-id", req.Context().Value(xRouter.RequestIDKey("requestID")).(string)))
					return
				}
			case <-responseChan:
				log.Info("Request completed successfully", slog.String("request-id", req.Context().Value(xRouter.RequestIDKey("requestID")).(string)))
				return
			}
		},
	)
}

func handleWithMiddleware(middleware func(http.Handler, *slog.Logger) http.Handler, next http.Handler, log *slog.Logger) http.Handler {
	return middleware(next, log)
}

func authMiddleware(next http.Handler, log *slog.Logger) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, req *http.Request) {
			authHeader := req.Header.Get("Authorization")
			if authHeader != "Basic bWl0ZXNoOk1pdGVzaC4xMjM=" {
				log.Error("Unauthorized", slog.String("request-id", req.Context().Value(xRouter.RequestIDKey("requestID")).(string)))
				w.WriteHeader(http.StatusUnauthorized)
				w.Write([]byte("Unauthorized"))
				return
			}
			log.Info("Authorization successful", slog.String("request-id", req.Context().Value(xRouter.RequestIDKey("requestID")).(string)))
			next.ServeHTTP(w, req)
		},
	)
}
