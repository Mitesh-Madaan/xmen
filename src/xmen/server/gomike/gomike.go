package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	xModels "gomike/models"
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
- Auth: Pending middleware for authentication, user login mechanism
- Change UUID to string > Check GORM documentation
- Different handlers for different methods
- Fix ID parsing from URL {placeholders}
- Parse object from request body directly
- Create model independent, not member of model class ; CRUD operations indenpendent of models
- API field validation : Check field function (optional) ; Explore OpenAPI 3 schema ; Generate models from schema automatically
- Just pass object in db methods, no need to pass model, Use generic
- Multiline string in ToString method > use %v
- To Status > no need to implement, just return object as json
- Return resp of API as json marshall of object
- Logging? log/Slog structured logging go package

Next:
- UT: Check ways to unit test http calls
- MVP
*/

var dbSession *xDb.DBSession
var connStr = "host=localhost user=postgres password=Postsql.123 dbname=postgres port=5432 sslmode=disable TimeZone=Asia/Shanghai"

func run() {
	var err error
	if dbSession == nil {
		dbSession = xDb.GetSession(connStr)
	}
	fmt.Printf("DB Session: %v\n", dbSession)

	err = SeedTables(dbSession)
	if err != nil {
		fmt.Printf("Error seeding tables: %v\n", err)
		return
	}

	mux := http.NewServeMux()

	// mux.HandleFunc("/person/", handlePerson())
	// mux.Handle("/person/", handler(handlePerson()))
	mux.Handle("/person/", handlerWithMiddleware(handlePerson(), middleware))

	// mux.HandleFunc("/animal/", handleAnimal())
	// mux.Handle("/animal/", handler(handleAnimal()))
	mux.Handle("/animal/", handlerWithMiddleware(handleAnimal(), middleware))

	fmt.Println("Starting server at port 8080")
	if err := http.ListenAndServe("localhost:8080", mux); err != nil {
		fmt.Println(err)
	}

}

func handler(handle func(http.ResponseWriter, *http.Request)) http.Handler {
	return http.HandlerFunc(handle)
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

func handlePerson() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		method := req.Method
		w.Header().Set("Content-Type", "text/plain")
		fmt.Printf("Request method: %s with body: %s\n", method, req.Body)

		switch method {
		case http.MethodGet:
			personID := req.URL.Path[len("/person/"):]
			personIDUint, err := strconv.ParseUint(personID, 10, 32)
			if err != nil {
				errResponse := fmt.Sprintf("Invalid person ID: %s", err.Error())
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(errResponse))
				return
			}
			person, err := xModels.GetPersonByID(dbSession, personIDUint)
			if err != nil {
				if strings.Contains(strings.ToLower(err.Error()), "record not found") {
					errResponse := fmt.Sprintf("Person  with ID %d not found: %s", personIDUint, err.Error())
					w.WriteHeader(http.StatusNotFound)
					w.Write([]byte(errResponse))
					return
				}
				errResponse := fmt.Sprintf("Error retrieving person: %s", err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(errResponse))
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(person.ToString()))
			return

		case http.MethodPut:
			personID := req.URL.Path[len("/person/"):]
			body, err := ioutil.ReadAll(req.Body)
			if err != nil {
				errResponse := fmt.Sprintf("Failed to read request body: %s", err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(errResponse))
				return
			}

			var objMap map[string]interface{}
			err = json.Unmarshal(body, &objMap)
			if err != nil {
				errResponse := fmt.Sprintf("Failed to unmarshal request body: %s", err.Error())
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(errResponse))
				return
			}

			personIDUint, err := strconv.ParseUint(personID, 10, 32)
			if err != nil {
				errResponse := fmt.Sprintf("Invalid person ID: %s", err.Error())
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(errResponse))
				return
			}
			person, err := xModels.GetPersonByID(dbSession, personIDUint)
			if err != nil {
				if strings.Contains(strings.ToLower(err.Error()), "record not found") {
					objMap["id"] = personIDUint // Ensure the ID is set for creation
					// Create the person
					person := xModels.Person{}
					err = person.Create(dbSession, objMap)
					if err != nil {
						errResponse := fmt.Sprintf("Failed to create person: %s", err.Error())
						w.WriteHeader(http.StatusInternalServerError)
						w.Write([]byte(errResponse))
						return
					}

					w.WriteHeader(http.StatusNoContent)
					return
				}
				errResponse := fmt.Sprintf("Error retrieving person with ID %s: %s", personID, err.Error())
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte(errResponse))
				return
			}

			// Update the person
			err = person.Update(dbSession, objMap)
			if err != nil {
				errResponse := fmt.Sprintf("Failed to update person: %s", err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(errResponse))
				return
			}

			w.WriteHeader(http.StatusNoContent)
			return

		case http.MethodPost:
			body, err := ioutil.ReadAll(req.Body)
			if err != nil {
				errResponse := fmt.Sprintf("Failed to read request body: %s", err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(errResponse))
				return
			}

			var objMap map[string]interface{}
			err = json.Unmarshal(body, &objMap)
			if err != nil {
				errResponse := fmt.Sprintf("Failed to unmarshal request body: %s", err.Error())
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(errResponse))
				return
			}

			person := xModels.Person{}
			err = person.Create(dbSession, objMap)
			if err != nil {
				errResponse := fmt.Sprintf("Failed to create person: %s", err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(errResponse))
				return
			}

			w.WriteHeader(http.StatusCreated)
			res := fmt.Sprintf("Person with ID %d added", person.ID)
			w.Write([]byte(res))
			return

		case http.MethodPatch:
			personID := req.URL.Path[len("/person/"):]
			body, err := ioutil.ReadAll(req.Body)
			if err != nil {
				errResponse := fmt.Sprintf("Failed to read request body: %s", err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(errResponse))
				return
			}

			var objMap map[string]interface{}
			err = json.Unmarshal(body, &objMap)
			if err != nil {
				errResponse := fmt.Sprintf("Failed to unmarshal request body: %s", err.Error())
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(errResponse))
				return
			}

			personIDUint, err := strconv.ParseUint(personID, 10, 32)
			if err != nil {
				errResponse := fmt.Sprintf("Invalid person ID: %s", err.Error())
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(errResponse))
				return
			}
			person, err := xModels.GetPersonByID(dbSession, personIDUint)
			if err != nil {
				if strings.Contains(strings.ToLower(err.Error()), "record not found") {
					errResponse := fmt.Sprintf("Patch method is only allowed on existing records. Person with ID %s not found", personID)
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte(errResponse))
					return
				}
				errResponse := fmt.Sprintf("Error retrieving person with ID %s: %s", personID, err.Error())
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte(errResponse))
				return
			}

			// Update the person
			err = person.Update(dbSession, objMap)
			if err != nil {
				errResponse := fmt.Sprintf("Failed to update person: %s", err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(errResponse))
				return
			}

			w.WriteHeader(http.StatusNoContent)
			return

		case http.MethodDelete:
			personID := req.URL.Path[len("/person/"):]
			personIDUint, err := strconv.ParseUint(personID, 10, 32)
			if err != nil {
				errResponse := fmt.Sprintf("Invalid person ID: %s", err.Error())
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(errResponse))
				return
			}
			person, err := xModels.GetPersonByID(dbSession, personIDUint)
			if err != nil {
				errResponse := fmt.Sprintf("Person not found: %s", err.Error())
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte(errResponse))
				return
			}

			err = person.Delete(dbSession)
			if err != nil {
				errResponse := fmt.Sprintf("Failed to delete person: %s", err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(errResponse))
				return
			}

			w.WriteHeader(http.StatusOK)
			res := fmt.Sprintf("Person with ID %d deleted", person.ID)
			w.Write([]byte(res))
			return

		default:
			errResponse := fmt.Sprintf("Method %s not allowed", method)
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte(errResponse))
			return
		}
	}
}

func handleAnimal() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		method := req.Method
		w.Header().Set("Content-Type", "text/plain")
		fmt.Printf("Request method: %s with body: %s\n", method, req.Body)

		switch method {
		case http.MethodGet:
			animalID := req.URL.Path[len("/animal/"):]
			animalIDUint, err := strconv.ParseUint(animalID, 10, 32)
			if err != nil {
				errResponse := fmt.Sprintf("Invalid animal ID: %s", err.Error())
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(errResponse))
				return
			}
			// Retrieve the animal by ID
			animal, err := xModels.GetAnimalByID(dbSession, animalIDUint)
			if err != nil {
				if strings.Contains(strings.ToLower(err.Error()), "record not found") {
					errResponse := fmt.Sprintf("Animal  with ID %d not found: %s", animalIDUint, err.Error())
					w.WriteHeader(http.StatusNotFound)
					w.Write([]byte(errResponse))
					return
				}
				errResponse := fmt.Sprintf("Error retrieving animal: %s", err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(errResponse))
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(animal.ToString()))
			return

		case http.MethodPut:
			animalID := req.URL.Path[len("/animal/"):]
			body, err := ioutil.ReadAll(req.Body)
			if err != nil {
				errResponse := fmt.Sprintf("Failed to read request body: %s", err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(errResponse))
				return
			}

			var objMap map[string]interface{}
			err = json.Unmarshal(body, &objMap)
			if err != nil {
				errResponse := fmt.Sprintf("Failed to unmarshal request body: %s", err.Error())
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(errResponse))
				return
			}

			animalIDUint, err := strconv.ParseUint(animalID, 10, 32)
			if err != nil {
				errResponse := fmt.Sprintf("Invalid animal ID: %s", err.Error())
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(errResponse))
			}
			animal, err := xModels.GetAnimalByID(dbSession, animalIDUint)
			if err != nil {
				if strings.Contains(strings.ToLower(err.Error()), "record not found") {
					objMap["id"] = animalIDUint // Ensure the ID is set for creation
					// Create the animal
					animal := xModels.Animal{}
					err = animal.Create(dbSession, objMap)
					if err != nil {
						errResponse := fmt.Sprintf("Failed to create animal: %s", err.Error())
						w.WriteHeader(http.StatusInternalServerError)
						w.Write([]byte(errResponse))
						return
					}

					w.WriteHeader(http.StatusNoContent)
					return
				}
				errResponse := fmt.Sprintf("Error retrieving animal with ID %s: %s", animalID, err.Error())
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte(errResponse))
				return
			}

			// Update the animal
			err = animal.Update(dbSession, objMap)
			if err != nil {
				errResponse := fmt.Sprintf("Failed to update animal: %s", err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(errResponse))
				return
			}

			w.WriteHeader(http.StatusNoContent)
			return

		case http.MethodPost:
			body, err := ioutil.ReadAll(req.Body)
			if err != nil {
				errResponse := fmt.Sprintf("Failed to read request body: %s", err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(errResponse))
				return
			}

			var objMap map[string]interface{}
			err = json.Unmarshal(body, &objMap)
			if err != nil {
				errResponse := fmt.Sprintf("Failed to unmarshal request body: %s", err.Error())
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(errResponse))
				return
			}

			animal := xModels.Animal{}
			err = animal.Create(dbSession, objMap)
			if err != nil {
				errResponse := fmt.Sprintf("Failed to create animal: %s", err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(errResponse))
				return
			}

			w.WriteHeader(http.StatusCreated)
			res := fmt.Sprintf("Animal with ID %d added", animal.ID)
			w.Write([]byte(res))
			return

		case http.MethodPatch:
			animalID := req.URL.Path[len("/animal/"):]
			body, err := ioutil.ReadAll(req.Body)
			if err != nil {
				errResponse := fmt.Sprintf("Failed to read request body: %s", err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(errResponse))
				return
			}

			var objMap map[string]interface{}
			err = json.Unmarshal(body, &objMap)
			if err != nil {
				errResponse := fmt.Sprintf("Failed to unmarshal request body: %s", err.Error())
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(errResponse))
				return
			}

			animalIDUint, err := strconv.ParseUint(animalID, 10, 32)
			if err != nil {
				errResponse := fmt.Sprintf("Invalid animal ID: %s", err.Error())
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(errResponse))
				return
			}

			animal, err := xModels.GetAnimalByID(dbSession, animalIDUint)
			if err != nil {
				if strings.Contains(strings.ToLower(err.Error()), "record not found") {
					errResponse := fmt.Sprintf("Patch method is only allowed on existing records. Animal with ID %s not found", animalID)
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte(errResponse))
					return
				}
				errResponse := fmt.Sprintf("Error retrieving animal with ID %s: %s", animalID, err.Error())
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte(errResponse))
				return
			}

			// Update the animal
			err = animal.Update(dbSession, objMap)
			if err != nil {
				errResponse := fmt.Sprintf("Failed to update animal: %s", err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(errResponse))
				return
			}

			w.WriteHeader(http.StatusNoContent)
			return

		case http.MethodDelete:
			animalID := req.URL.Path[len("/animal/"):]
			animalIDUint, err := strconv.ParseUint(animalID, 10, 32)
			if err != nil {
				errResponse := fmt.Sprintf("Invalid animal ID: %s", err.Error())
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(errResponse))
				return
			}
			animal, err := xModels.GetAnimalByID(dbSession, animalIDUint)
			if err != nil {
				errResponse := fmt.Sprintf("Animal not found: %s", err.Error())
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte(errResponse))
				return
			}

			err = animal.Delete(dbSession)
			if err != nil {
				errResponse := fmt.Sprintf("Failed to delete animal: %s", err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(errResponse))
				return
			}

			w.WriteHeader(http.StatusOK)
			res := fmt.Sprintf("Animal with ID %d deleted", animal.ID)
			w.Write([]byte(res))
			return

		default:
			errResponse := fmt.Sprintf("Method %s not allowed", method)
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte(errResponse))
			return
		}
	}
}
