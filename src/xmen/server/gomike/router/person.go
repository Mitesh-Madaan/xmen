package router

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	xModels "gomike/models"
	xSession "gomike/session"
	xDb "lib/dbchef"

	"github.com/google/uuid"
)

func GetPerson(dbSession *xDb.DBSession) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/plain")

		personID := req.PathValue("personID")
		if personID == "" {
			errResponse := "Person ID not provided in the URL"
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(errResponse))
			return
		}

		personPtr, err := xSession.ReadRecord[xModels.Person](dbSession, personID)
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
		objDetails, err := json.Marshal(personPtr)
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
		personID := req.PathValue("personID")
		if personID == "" {
			errResponse := "Person ID not provided in the URL"
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
		if len(body) == 0 {
			errResponse := "Request body is empty"
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(errResponse))
			return
		}

		personPtr, err := xSession.ReadRecord[xModels.Person](dbSession, personID)
		if err != nil {
			if strings.Contains(strings.ToLower(err.Error()), "record not found") {
				// If the person is not found, create a new person with the provided ID
				person := xModels.Person{
					ID:            personID,
					Kind:          "person",
					Cloned:        false,
					ClonedFromRef: "",
				} // Create a new person instance with the ID from the URL
				err = json.Unmarshal(body, &person)
				if err != nil {
					errResponse := fmt.Sprintf("Failed to unmarshal request body into person: %s", err.Error())
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte(errResponse))
					return
				}
				err = xSession.CreateRecord(dbSession, person)
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
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errResponse))
			return
		}

		err = json.Unmarshal(body, personPtr)
		if err != nil {
			errResponse := fmt.Sprintf("Failed to unmarshal request body into person: %s", err.Error())
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(errResponse))
			return
		}

		err = xSession.UpdateRecord(dbSession, *personPtr)
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

		if len(body) == 0 {
			errResponse := "Request body is empty"
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(errResponse))
			return
		}

		person := xModels.Person{
			ID:            uuid.New().String(),
			Kind:          "person",
			Cloned:        false,
			ClonedFromRef: "",
		} // Create a new person instance

		err = json.Unmarshal(body, &person)
		if err != nil {
			errResponse := fmt.Sprintf("Failed to unmarshal request body into person: %s", err.Error())
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(errResponse))
			return
		}

		err = xSession.CreateRecord(dbSession, person)
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
		personID := req.PathValue("personID")
		if personID == "" {
			errResponse := "Person ID not provided in the URL"
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

		personPtr, err := xSession.ReadRecord[xModels.Person](dbSession, personID)
		if err != nil {
			if strings.Contains(strings.ToLower(err.Error()), "record not found") {
				errResponse := fmt.Sprintf("Patch method is only allowed on existing records. Person with ID %s not found", personID)
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(errResponse))
				return
			}
			errResponse := fmt.Sprintf("Error retrieving person with ID %s: %s", personID, err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errResponse))
			return
		}

		err = json.Unmarshal(body, personPtr)
		if err != nil {
			errResponse := fmt.Sprintf("Failed to unmarshal request body into person: %s", err.Error())
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(errResponse))
			return
		}

		err = xSession.UpdateRecord(dbSession, *personPtr)
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
		personID := req.PathValue("personID")
		if personID == "" {
			errResponse := "Person ID not provided in the URL"
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(errResponse))
			return
		}

		personPtr, err := xSession.ReadRecord[xModels.Person](dbSession, personID)
		if err != nil {
			if strings.Contains(strings.ToLower(err.Error()), "record not found") {
				errResponse := fmt.Sprintf("Person with ID %s not found", personID)
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte(errResponse))
				return
			}
			errResponse := fmt.Sprintf("Error retrieving person: %s", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errResponse))
			return
		}

		err = xSession.DeleteRecord(dbSession, *personPtr)
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
