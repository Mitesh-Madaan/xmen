package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	xModels "gomike/models"
	xBase "lib/base"
)

func main() {
	fmt.Println("This is a placeholder for the main function")
	run()
	fmt.Println("End of the main function")
}

/*
TODO:
- Singleton DB session management
- Add UTs
- Minimum viable product (binaries, tests, documentation)
- Publish as docker image
- Intergration tests
- Remote DB PostgresQL
*/

func run() {
	var err error

	if err = InitDirectory(); err != nil {
		panic(err)
	}

	if err = SeedDirectory(); err != nil {
		panic(err)
	}

	if err = PrintDirectory(); err != nil {
		panic(err)
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/person/", handlePerson())

	mux.HandleFunc("/animal/", handleAnimal())

	fmt.Println("Starting server at port 8080")
	if err := http.ListenAndServe("localhost:8080", mux); err != nil {
		fmt.Println(err)
	}

	if err = StoreDirectory(); err != nil {
		panic(err)
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
			person, err := xModels.GetObjectFromDB(personID, "Person")
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
			person, err := xModels.GetObjectFromDB(personID, "Person")
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

			person.Edit(editConfig)
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

			var b *xBase.Base
			err = json.Unmarshal(body, &b)
			if err != nil {
				errResponse := fmt.Sprintf("Failed to unmarshal request body: %s", err.Error())
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(errResponse))
				return
			}

			person := xModels.NewPerson(b)
			err = person.Save()
			if err != nil {
				errResponse := fmt.Sprintf("Failed to save person: %s", err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(errResponse))
				return
			}

			w.WriteHeader(http.StatusOK)
			res := fmt.Sprintf("Person with ID %s added", person.GetBase().ID.String())
			w.Write([]byte(res))
			return

		case http.MethodDelete:
			personID := req.URL.Path[len("/person/"):]
			person, err := xModels.GetObjectFromDB(personID, "Person")
			if err != nil {
				errResponse := fmt.Sprintf("Person not found: %s", err.Error())
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte(errResponse))
				return
			}

			err = person.Delete()
			if err != nil {
				errResponse := fmt.Sprintf("Failed to delete person: %s", err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(errResponse))
				return
			}

			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Person deleted"))
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
			animal, err := xModels.GetObjectFromDB(animalID, "Animal")
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
			animal, err := xModels.GetObjectFromDB(animalID, "Animal")
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

			animal.Edit(editConfig)
			err = animal.Save()
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

			var b *xBase.Base
			err = json.Unmarshal(body, &b)
			if err != nil {
				errResponse := fmt.Sprintf("Failed to unmarshal request body: %s", err.Error())
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(errResponse))
				return
			}

			animal := xModels.NewPerson(b)
			err = animal.Save()
			if err != nil {
				errResponse := fmt.Sprintf("Failed to save person: %s", err.Error())
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(errResponse))
				return
			}

			w.WriteHeader(http.StatusOK)
			res := fmt.Sprintf("Animal with ID %s added", animal.GetBase().ID.String())
			w.Write([]byte(res))
			return

		case http.MethodDelete:
			animalID := req.URL.Path[len("/animal/"):]
			animal, err := xModels.GetObjectFromDB(animalID, "Animal")
			if err != nil {
				errResponse := fmt.Sprintf("Animal not found: %s", err.Error())
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte(errResponse))
				return
			}

			err = animal.Delete()
			if err != nil {
				errResponse := fmt.Sprintf("Failed to delete person: %v", err)
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(errResponse))
				return
			}

			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Animal deleted"))
			return

		default:
			errResponse := fmt.Sprintf("Method %s not allowed", method)
			w.WriteHeader(http.StatusMethodNotAllowed)
			w.Write([]byte(errResponse))
			return
		}
	}
}
