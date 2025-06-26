package main

import (
	"fmt"
	"net/http"

	xRouter "gomike/router"
	xSession "gomike/session"
	xDb "lib/dbchef"
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
- Create model independent, not member of model class ; CRUD operations indenpendent of models
- API field validation : Check field function (optional) ; Explore OpenAPI 3 schema ; Generate models from schema automatically
- Just pass object in db methods, no need to pass model DONE
- Multiline string in ToString method > use %v DONE
- To Status > no need to implement, just return object as json DONE
- Return resp of API as json marshall of object DONE
- Logging? log/Slog structured logging go package

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

	mux := http.NewServeMux()

	mux.Handle("GET /person/{personID}", handlerWithMiddleware(xRouter.GetPerson(dbSession), middleware))
	mux.Handle("PUT /person/{personID}", handlerWithMiddleware(xRouter.UpdatePerson(dbSession), middleware))
	mux.Handle("POST /person/", handlerWithMiddleware(xRouter.CreatePerson(dbSession), middleware))
	mux.Handle("PATCH /person/{personID}", handlerWithMiddleware(xRouter.PatchPerson(dbSession), middleware))
	mux.Handle("DELETE /person/{personID}", handlerWithMiddleware(xRouter.DeletePerson(dbSession), middleware))

	mux.Handle("GET /animal/{animalID}", handlerWithMiddleware(xRouter.GetAnimal(dbSession), middleware))
	mux.Handle("PUT /animal/{animalID}", handlerWithMiddleware(xRouter.UpdateAnimal(dbSession), middleware))
	mux.Handle("POST /animal/", handlerWithMiddleware(xRouter.CreateAnimal(dbSession), middleware))
	mux.Handle("PATCH /animal/{animalID}", handlerWithMiddleware(xRouter.PatchAnimal(dbSession), middleware))
	mux.Handle("DELETE /animal/{animalID}", handlerWithMiddleware(xRouter.DeleteAnimal(dbSession), middleware))

	fmt.Println("Starting server at port 8080")
	if err = http.ListenAndServe("localhost:8080", mux); err != nil {
		fmt.Println(err)
	}

}

func handlerWithMiddleware(handle func(http.ResponseWriter, *http.Request), middleware func(http.Handler) http.Handler) http.Handler {
	return middleware(http.HandlerFunc(handle))
}

func middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		authHeader := req.Header.Get("Authorization")
		if authHeader != "Basic bWl0ZXNoOk1pdGVzaC4xMjM=" {
			fmt.Println("Authorization failed")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized"))
			return
		}
		next.ServeHTTP(w, req)
	})
}
