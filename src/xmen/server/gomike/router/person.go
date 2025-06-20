package router

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/google/uuid"

	xModels "gomike/models"
	xDb "lib/dbchef"
)

func GetPerson(dbSession *xDb.DBSession) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		fmt.Printf("Request GET with body: %s\n", req.Body)

		personID, err := ParseIDFromURL(req.URL)
		if err != nil {
			errResponse := fmt.Sprintf("Failed to parse ID from URL: %s", err.Error())
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(errResponse))
			return
		}

		person, err := xModels.GetPersonByID(dbSession, personID)
		if err != nil {
			if strings.Contains(strings.ToLower(err.Error()), "record not found") {
				errResponse := fmt.Sprintf("Person  with ID %s not found: %s", personID, err.Error())
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
		objDetails, err := json.Marshal(person)
		if err != nil {
			errResponse := fmt.Sprintf("Failed to marshal person details: %s", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errResponse))
			return
		}
		w.Write(objDetails)
	}
}

func UpdatePerson(dbSession *xDb.DBSession) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		personID, err := ParseIDFromURL(req.URL)
		if err != nil {
			errResponse := fmt.Sprintf("Failed to parse ID from URL: %s", err.Error())
			w.WriteHeader(http.StatusBadRequest)
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

		person, err := xModels.GetPersonByID(dbSession, personID)
		if err != nil {
			if strings.Contains(strings.ToLower(err.Error()), "record not found") {
				// Create the person
				person := xModels.Person{ID: personID} // Create a new person instance with the ID from the URL
				err = person.Create(dbSession, body)
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
		err = person.Update(dbSession, body)
		if err != nil {
			errResponse := fmt.Sprintf("Failed to update person: %s", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errResponse))
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func CreatePerson(dbSession *xDb.DBSession) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			errResponse := fmt.Sprintf("Failed to read request body: %s", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errResponse))
			return
		}

		person := xModels.Person{ID: uuid.New().String()} // Create a new person instance
		err = person.Create(dbSession, body)
		if err != nil {
			errResponse := fmt.Sprintf("Failed to create person: %s", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errResponse))
			return
		}

		w.WriteHeader(http.StatusCreated)
		res := fmt.Sprintf("Person with ID %s added", person.ID)
		w.Write([]byte(res))
	}
}

func PatchPerson(dbSession *xDb.DBSession) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		personID, err := ParseIDFromURL(req.URL)
		if err != nil {
			errResponse := fmt.Sprintf("Failed to parse ID from URL: %s", err.Error())
			w.WriteHeader(http.StatusBadRequest)
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

		person, err := xModels.GetPersonByID(dbSession, personID)
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
		err = person.Update(dbSession, body)
		if err != nil {
			errResponse := fmt.Sprintf("Failed to update person: %s", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errResponse))
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func DeletePerson(dbSession *xDb.DBSession) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		personID, err := ParseIDFromURL(req.URL)
		if err != nil {
			errResponse := fmt.Sprintf("Failed to parse ID from URL: %s", err.Error())
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(errResponse))
			return
		}

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
		res := fmt.Sprintf("Person with ID %s deleted", personID)
		w.Write([]byte(res))
	}
}
