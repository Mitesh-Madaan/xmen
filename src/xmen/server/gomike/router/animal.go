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

func GetAnimal(dbSession *xDb.DBSession) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		fmt.Printf("Request GET with body: %s\n", req.Body)

		animalID, err := ParseIDFromURL(req.URL)
		if err != nil {
			errResponse := fmt.Sprintf("Failed to parse ID from URL: %s", err.Error())
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(errResponse))
			return
		}

		animal, err := xModels.GetAnimalByID(dbSession, animalID)
		if err != nil {
			if strings.Contains(strings.ToLower(err.Error()), "record not found") {
				errResponse := fmt.Sprintf("Animal with ID %s not found: %s", animalID, err.Error())
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
	}
}

func UpdateAnimal(dbSession *xDb.DBSession) func(http.ResponseWriter, *http.Request) {

	return func(w http.ResponseWriter, req *http.Request) {
		animalID, err := ParseIDFromURL(req.URL)
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

		animal, err := xModels.GetAnimalByID(dbSession, animalID)
		if err != nil {
			if strings.Contains(strings.ToLower(err.Error()), "record not found") {
				// Create the animal
				animal := xModels.Animal{ID: animalID} // Create a new animal instance with the ID from the URL
				err = animal.Create(dbSession, body)
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
		err = animal.Update(dbSession, body)
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

		var objMap map[string]interface{}
		err = json.Unmarshal(body, &objMap)
		if err != nil {
			errResponse := fmt.Sprintf("Failed to unmarshal request body: %s", err.Error())
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(errResponse))
			return
		}

		animal := xModels.Animal{ID: uuid.New().String()} // Create a new animal instance
		err = animal.Create(dbSession, body)
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
		animalID, err := ParseIDFromURL(req.URL)
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

		var objMap map[string]interface{}
		err = json.Unmarshal(body, &objMap)
		if err != nil {
			errResponse := fmt.Sprintf("Failed to unmarshal request body: %s", err.Error())
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(errResponse))
			return
		}

		animal, err := xModels.GetAnimalByID(dbSession, animalID)
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
		err = animal.Update(dbSession, body)
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
		animalID, err := ParseIDFromURL(req.URL)
		if err != nil {
			errResponse := fmt.Sprintf("Failed to parse ID from URL: %s", err.Error())
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(errResponse))
			return
		}

		animal, err := xModels.GetAnimalByID(dbSession, animalID)
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
		res := fmt.Sprintf("Animal with ID %s deleted", animalID)
		w.Write([]byte(res))
	}
}
