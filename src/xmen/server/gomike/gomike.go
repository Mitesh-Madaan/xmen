package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

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

	mux.HandleFunc("/person/", handlePerson())

	mux.HandleFunc("/animal/", handleAnimal())

	fmt.Println("Starting server at port 8080")
	if err := http.ListenAndServe("localhost:8080", mux); err != nil {
		fmt.Println(err)
	}

}

func handlePerson() func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		auth_header := req.Header.Get("Authorization")
		if auth_header == "Basic bWl0ZXNoOk1pdGVzaC4xMjM=" {
			fmt.Println("Authorization successful")
		} else {
			fmt.Println("Authorization failed")
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Unauthorized"))
			return
		}

		method := req.Method
		w.Header().Set("Content-Type", "text/plain")
		fmt.Printf("Request method: %s with body: %s\n", method, req.Body)

		switch method {
		case http.MethodGet:
			personID := req.URL.Path[len("/person/"):]
			person, err := xModels.GetPersonByID(dbSession, personID)
			if err != nil {
				errResponse := fmt.Sprintf("error saving person: %s", err.Error())
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte(errResponse))
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(person.ToString()))
			return

		case http.MethodPost:
			personID := req.URL.Path[len("/person/"):]
			person, err := xModels.GetPersonByID(dbSession, personID)
			if err != nil {
				errResponse := fmt.Sprintf("Person Not found: %s", err.Error())
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte(errResponse))
				return
			}

			body, err := ioutil.ReadAll(req.Body)
			if err != nil {
				errResponse := fmt.Sprintf("Failed to read request body: %s", err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(errResponse))
				return
			}

			editConfig := map[string]interface{}{}
			err = json.Unmarshal(body, &editConfig)
			if err != nil {
				errResponse := fmt.Sprintf("Failed to unmarshal request body: %s", err.Error())
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(errResponse))
				return
			}

			err = person.Update(dbSession, editConfig)
			if err != nil {
				errResponse := fmt.Sprintf("Failed to save person: %s", err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(errResponse))
				return
			}

			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Person edited"))
			return

		case http.MethodPut:
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
				errResponse := fmt.Sprintf("Failed to save person: %s", err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(errResponse))
				return
			}

			w.WriteHeader(http.StatusOK)
			res := fmt.Sprintf("Person with ID %d added", person.ID)
			w.Write([]byte(res))
			return

		case http.MethodDelete:
			personID := req.URL.Path[len("/person/"):]
			person, err := xModels.GetPersonByID(dbSession, personID)
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
			animal, err := xModels.GetAnimalByID(dbSession, animalID)
			if err != nil {
				errResponse := fmt.Sprintf("Animal not found: %s", err.Error())
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte(errResponse))
				return
			}
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(animal.ToString()))
			return

		case http.MethodPost:
			animalID := req.URL.Path[len("/animal/"):]
			animal, err := xModels.GetAnimalByID(dbSession, animalID)
			if err != nil {
				errResponse := fmt.Sprintf("Animal not found: %s", err.Error())
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte(errResponse))
				return
			}

			body, err := ioutil.ReadAll(req.Body)
			if err != nil {
				errResponse := fmt.Sprintf("Failed to read request body: %s", err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(errResponse))
				return
			}

			editConfig := map[string]interface{}{}
			err = json.Unmarshal(body, &editConfig)
			if err != nil {
				errResponse := fmt.Sprintf("Failed to unmarshal request body: %s", err.Error())
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(errResponse))
				return
			}

			err = animal.Update(dbSession, editConfig)
			if err != nil {
				errResponse := fmt.Sprintf("Failed to save person: %s", err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(errResponse))
				return
			}

			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Animal edited"))
			return

		case http.MethodPut:
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
				errResponse := fmt.Sprintf("Failed to save person: %s", err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(errResponse))
				return
			}

			w.WriteHeader(http.StatusOK)
			res := fmt.Sprintf("Animal with ID %d added", animal.ID)
			w.Write([]byte(res))
			return

		case http.MethodDelete:
			animalID := req.URL.Path[len("/animal/"):]
			animal, err := xModels.GetAnimalByID(dbSession, animalID)
			if err != nil {
				errResponse := fmt.Sprintf("Animal not found: %s", err.Error())
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte(errResponse))
				return
			}

			err = animal.Delete(dbSession)
			if err != nil {
				errResponse := fmt.Sprintf("Failed to delete person: %v", err)
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
