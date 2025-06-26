package router

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/google/uuid"

	xModels "gomike/models"
	xSession "gomike/session"
	xDb "lib/dbchef"
)

func GetAnimal(dbSession *xDb.DBSession) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/plain")

		animalID := req.PathValue("animalID")
		if animalID == "" {
			errResponse := "Animal ID not provided in the URL"
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(errResponse))
			return
		}

		animalPtr, err := xSession.ReadRecord[xModels.Animal](dbSession, animalID)
		if err != nil {
			if strings.Contains(strings.ToLower(err.Error()), "record not found") {
				errResponse := fmt.Sprintf("Animal with ID %s not found", animalID)
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
		objDetails, err := json.Marshal(animalPtr)
		if err != nil {
			errResponse := fmt.Sprintf("Failed to marshal animal details: %s", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errResponse))
			return
		}
		w.Write(objDetails)
	}
}

func UpdateAnimal(dbSession *xDb.DBSession) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		animalID := req.PathValue("animalID")
		if animalID == "" {
			errResponse := "Animal ID not provided in the URL"
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

		animalPtr, err := xSession.ReadRecord[xModels.Animal](dbSession, animalID)
		if err != nil {
			if strings.Contains(strings.ToLower(err.Error()), "record not found") {
				// Create the animal
				animal := xModels.Animal{
					ID:            animalID,
					Kind:          "animal",
					Cloned:        false,
					ClonedFromRef: "",
				} // Create a new animal instance with the ID from the URL
				err = json.Unmarshal(body, &animal)
				if err != nil {
					errResponse := fmt.Sprintf("Failed to unmarshal request body: %s", err.Error())
					w.WriteHeader(http.StatusBadRequest)
					w.Write([]byte(errResponse))
					return
				}
				err = xSession.CreateRecord(dbSession, animal)
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

		err = json.Unmarshal(body, animalPtr)
		if err != nil {
			errResponse := fmt.Sprintf("Failed to unmarshal request body into animal: %s", err.Error())
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(errResponse))
			return
		}

		err = xSession.UpdateRecord(dbSession, animalPtr)
		if err != nil {
			errResponse := fmt.Sprintf("Failed to update animal: %s", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errResponse))
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func CreateAnimal(dbSession *xDb.DBSession) func(http.ResponseWriter, *http.Request) {
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

		animal := xModels.Animal{
			ID:            uuid.New().String(), // Generate a new UUID for the animal
			Kind:          "animal",
			Cloned:        false,
			ClonedFromRef: "",
		} // Create a new animal instance with a generated ID

		err = json.Unmarshal(body, &animal)
		if err != nil {
			errResponse := fmt.Sprintf("Failed to unmarshal request body into animal: %s", err.Error())
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(errResponse))
			return
		}

		err = xSession.CreateRecord(dbSession, animal)
		if err != nil {
			errResponse := fmt.Sprintf("Failed to create animal: %s", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errResponse))
			return
		}

		w.WriteHeader(http.StatusCreated)
		res := fmt.Sprintf("Animal with ID %s added", animal.ID)
		w.Write([]byte(res))
	}
}

func PatchAnimal(dbSession *xDb.DBSession) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		animalID := req.PathValue("animalID")
		if animalID == "" {
			errResponse := "Animal ID not provided in the URL"
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

		animalPtr, err := xSession.ReadRecord[xModels.Animal](dbSession, animalID)
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

		// Unmarshal the request body into a map to allow partial updates
		err = json.Unmarshal(body, &animalPtr)
		if err != nil {
			errResponse := fmt.Sprintf("Failed to unmarshal request body into animal: %s", err.Error())
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(errResponse))
			return
		}

		err = xSession.UpdateRecord(dbSession, animalPtr)
		if err != nil {
			errResponse := fmt.Sprintf("Failed to update animal: %s", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errResponse))
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func DeleteAnimal(dbSession *xDb.DBSession) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		animalID := req.PathValue("animalID")
		if animalID == "" {
			errResponse := "Animal ID not provided in the URL"
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(errResponse))
			return
		}

		animalPtr, err := xSession.ReadRecord[xModels.Animal](dbSession, animalID)
		if err != nil {
			if strings.Contains(strings.ToLower(err.Error()), "record not found") {
				errResponse := fmt.Sprintf("Animal with ID %s not found", animalID)
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte(errResponse))
				return
			}

			errResponse := fmt.Sprintf("Error retrieving animal: %s", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errResponse))
			return
		}

		err = xSession.DeleteRecord(dbSession, animalPtr)
		if err != nil {
			errResponse := fmt.Sprintf("Failed to delete animal: %s", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errResponse))
			return
		}

		w.WriteHeader(http.StatusOK)
		res := fmt.Sprintf("Animal with ID %s deleted", animalID)
		w.Write([]byte(res))
	}
}
