package router

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log/slog"
	"net/http"
	"strings"

	"github.com/google/uuid"

	xModels "gomike/models"
	xSession "gomike/session"
	xDb "lib/dbchef"
)

func GetAnimal(dbSession *xDb.DBSession, log *slog.Logger) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		animalID := req.PathValue("animalID")
		if animalID == "" {
			errResponse := "Animal ID not provided in the URL"
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(errResponse))
			log.Error("Animal ID not provided in the URL", slog.String("request-id", req.Context().Value(RequestIDKey("requestID")).(string)))
			return
		}

		log.Info("Got animal ID from request", slog.String("request-id", req.Context().Value(RequestIDKey("requestID")).(string)), slog.String("animalID", animalID))

		animalPtr, err := xSession.ReadRecord[xModels.Animal](dbSession, animalID)
		if err != nil {
			if strings.Contains(strings.ToLower(err.Error()), "record not found") {
				errResponse := fmt.Sprintf("Animal with ID %s not found", animalID)
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte(errResponse))
				log.Error("Animal not found", slog.String("request-id", req.Context().Value(RequestIDKey("requestID")).(string)))
				return
			}
			errResponse := fmt.Sprintf("Error retrieving animal: %s", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errResponse))
			log.Error("Error retrieving animal", slog.String("request-id", req.Context().Value(RequestIDKey("requestID")).(string)), slog.String("error", err.Error()))
			return
		}
		w.WriteHeader(http.StatusOK)
		objDetails, err := json.Marshal(animalPtr)
		if err != nil {
			errResponse := fmt.Sprintf("Failed to marshal animal details: %s", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errResponse))
			log.Error("Failed to marshal animal details", slog.String("request-id", req.Context().Value(RequestIDKey("requestID")).(string)), slog.String("error", err.Error()))
			return
		}
		w.Write(objDetails)
		log.Info("Animal details retrieved successfully", slog.String("request-id", req.Context().Value(RequestIDKey("requestID")).(string)))
	}
}

func UpdateAnimal(dbSession *xDb.DBSession, log *slog.Logger) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		animalID := req.PathValue("animalID")
		if animalID == "" {
			errResponse := "Animal ID not provided in the URL"
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(errResponse))
			log.Error("Animal ID not provided in the URL", slog.String("request-id", req.Context().Value(RequestIDKey("requestID")).(string)))
			return
		}

		log.Info("Got animal ID from request", slog.String("request-id", req.Context().Value(RequestIDKey("requestID")).(string)), slog.String("animalID", animalID))

		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			errResponse := fmt.Sprintf("Failed to read request body: %s", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errResponse))
			log.Error("Failed to read request body", slog.String("request-id", req.Context().Value(RequestIDKey("requestID")).(string)), slog.String("error", err.Error()))
			return
		}
		if len(body) == 0 {
			errResponse := "Request body is empty"
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(errResponse))
			log.Error("Request body is empty", slog.String("request-id", req.Context().Value(RequestIDKey("requestID")).(string)))
			return
		}

		animalPtr, err := xSession.ReadRecord[xModels.Animal](dbSession, animalID)
		if err != nil {
			if strings.Contains(strings.ToLower(err.Error()), "record not found") {
				log.Info("Animal not found, creating a new one", slog.String("request-id", req.Context().Value(RequestIDKey("requestID")).(string)))
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
					log.Error("Failed to unmarshal request body", slog.String("request-id", req.Context().Value(RequestIDKey("requestID")).(string)), slog.String("error", err.Error()))
					return
				}
				err = xSession.CreateRecord(dbSession, animal)
				if err != nil {
					errResponse := fmt.Sprintf("Failed to create animal: %s", err.Error())
					w.WriteHeader(http.StatusInternalServerError)
					w.Write([]byte(errResponse))
					log.Error("Failed to create animal", slog.String("request-id", req.Context().Value(RequestIDKey("requestID")).(string)), slog.String("error", err.Error()))
					return
				}

				w.WriteHeader(http.StatusNoContent)
				log.Info("New Animal created successfully", slog.String("request-id", req.Context().Value(RequestIDKey("requestID")).(string)))
				return
			}
			errResponse := fmt.Sprintf("Error retrieving animal with ID %s: %s", animalID, err.Error())
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(errResponse))
			log.Error("Error retrieving animal", slog.String("request-id", req.Context().Value(RequestIDKey("requestID")).(string)), slog.String("error", err.Error()))
			return
		}

		log.Info("Updating existing animal", slog.String("request-id", req.Context().Value(RequestIDKey("requestID")).(string)))
		err = json.Unmarshal(body, animalPtr)
		if err != nil {
			errResponse := fmt.Sprintf("Failed to unmarshal request body into animal: %s", err.Error())
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(errResponse))
			log.Error("Failed to unmarshal request body into animal", slog.String("request-id", req.Context().Value(RequestIDKey("requestID")).(string)), slog.String("error", err.Error()))
			return
		}

		err = xSession.UpdateRecord(dbSession, animalPtr)
		if err != nil {
			errResponse := fmt.Sprintf("Failed to update animal: %s", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errResponse))
			log.Error("Failed to update animal", slog.String("request-id", req.Context().Value(RequestIDKey("requestID")).(string)), slog.String("error", err.Error()))
			return
		}

		w.WriteHeader(http.StatusNoContent)
		log.Info("Animal updated successfully", slog.String("request-id", req.Context().Value(RequestIDKey("requestID")).(string)))
	}
}

func CreateAnimal(dbSession *xDb.DBSession, log *slog.Logger) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			errResponse := fmt.Sprintf("Failed to read request body: %s", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errResponse))
			log.Error("Failed to read request body", slog.String("request-id", req.Context().Value(RequestIDKey("requestID")).(string)), slog.String("error", err.Error()))
			return
		}

		if len(body) == 0 {
			errResponse := "Request body is empty"
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(errResponse))
			log.Error("Request body is empty", slog.String("request-id", req.Context().Value(RequestIDKey("requestID")).(string)))
			return
		}

		animal := xModels.Animal{
			ID:            uuid.New().String(), // Generate a new UUID for the animal
			Kind:          "animal",
			Cloned:        false,
			ClonedFromRef: "",
		} // Create a new animal instance with a generated ID
		log.Info("Creating new animal with ID", slog.String("request-id", req.Context().Value(RequestIDKey("requestID")).(string)), slog.String("animalID", animal.ID))

		err = json.Unmarshal(body, &animal)
		if err != nil {
			errResponse := fmt.Sprintf("Failed to unmarshal request body into animal: %s", err.Error())
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(errResponse))
			log.Error("Failed to unmarshal request body into animal", slog.String("request-id", req.Context().Value(RequestIDKey("requestID")).(string)), slog.String("error", err.Error()))
			return
		}

		err = xSession.CreateRecord(dbSession, animal)
		if err != nil {
			errResponse := fmt.Sprintf("Failed to create animal: %s", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errResponse))
			log.Error("Failed to create animal", slog.String("request-id", req.Context().Value(RequestIDKey("requestID")).(string)), slog.String("error", err.Error()))
			return
		}

		w.WriteHeader(http.StatusCreated)
		res := fmt.Sprintf("Animal with ID %s added", animal.ID)
		w.Write([]byte(res))
		log.Info("New Animal created successfully", slog.String("request-id", req.Context().Value(RequestIDKey("requestID")).(string)))
	}
}

func PatchAnimal(dbSession *xDb.DBSession, log *slog.Logger) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		animalID := req.PathValue("animalID")
		if animalID == "" {
			errResponse := "Animal ID not provided in the URL"
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(errResponse))
			log.Error("Animal ID not provided in the URL", slog.String("request-id", req.Context().Value(RequestIDKey("requestID")).(string)))
			return
		}

		log.Info("Got animal ID from request", slog.String("request-id", req.Context().Value(RequestIDKey("requestID")).(string)), slog.String("animalID", animalID))

		body, err := ioutil.ReadAll(req.Body)
		if err != nil {
			errResponse := fmt.Sprintf("Failed to read request body: %s", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errResponse))
			log.Error("Failed to read request body", slog.String("request-id", req.Context().Value(RequestIDKey("requestID")).(string)), slog.String("error", err.Error()))
			return
		}

		if len(body) == 0 {
			errResponse := "Request body is empty"
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(errResponse))
			log.Error("Request body is empty", slog.String("request-id", req.Context().Value(RequestIDKey("requestID")).(string)))
			return
		}

		animalPtr, err := xSession.ReadRecord[xModels.Animal](dbSession, animalID)
		if err != nil {
			if strings.Contains(strings.ToLower(err.Error()), "record not found") {
				errResponse := fmt.Sprintf("Patch method is only allowed on existing records. Animal with ID %s not found", animalID)
				w.WriteHeader(http.StatusBadRequest)
				w.Write([]byte(errResponse))
				log.Warn("Patch method called on non-existing animal, only allowed on existing records", slog.String("request-id", req.Context().Value(RequestIDKey("requestID")).(string)))
				return
			}
			errResponse := fmt.Sprintf("Error retrieving animal with ID %s: %s", animalID, err.Error())
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(errResponse))
			log.Error("Error retrieving animal", slog.String("request-id", req.Context().Value(RequestIDKey("requestID")).(string)), slog.String("error", err.Error()))
			return
		}

		// Unmarshal the request body into a map to allow partial updates
		err = json.Unmarshal(body, &animalPtr)
		if err != nil {
			errResponse := fmt.Sprintf("Failed to unmarshal request body into animal: %s", err.Error())
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(errResponse))
			log.Error("Failed to unmarshal request body into animal", slog.String("request-id", req.Context().Value(RequestIDKey("requestID")).(string)), slog.String("error", err.Error()))
			return
		}

		log.Info("Patching existing animal", slog.String("request-id", req.Context().Value(RequestIDKey("requestID")).(string)))
		err = xSession.UpdateRecord(dbSession, animalPtr)
		if err != nil {
			errResponse := fmt.Sprintf("Failed to update animal: %s", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errResponse))
			log.Error("Failed to update animal", slog.String("request-id", req.Context().Value(RequestIDKey("requestID")).(string)), slog.String("error", err.Error()))
			return
		}

		w.WriteHeader(http.StatusNoContent)
		log.Info("Animal patched successfully", slog.String("request-id", req.Context().Value(RequestIDKey("requestID")).(string)))
	}
}

func DeleteAnimal(dbSession *xDb.DBSession, log *slog.Logger) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		animalID := req.PathValue("animalID")
		if animalID == "" {
			errResponse := "Animal ID not provided in the URL"
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(errResponse))
			log.Error("Animal ID not provided in the URL", slog.String("request-id", req.Context().Value(RequestIDKey("requestID")).(string)))
			return
		}

		log.Info("Got animal ID from request", slog.String("request-id", req.Context().Value(RequestIDKey("requestID")).(string)), slog.String("animalID", animalID))

		animalPtr, err := xSession.ReadRecord[xModels.Animal](dbSession, animalID)
		if err != nil {
			if strings.Contains(strings.ToLower(err.Error()), "record not found") {
				errResponse := fmt.Sprintf("Animal with ID %s not found", animalID)
				w.WriteHeader(http.StatusNotFound)
				w.Write([]byte(errResponse))
				log.Error("Animal not found", slog.String("request-id", req.Context().Value(RequestIDKey("requestID")).(string)))
				return
			}

			errResponse := fmt.Sprintf("Error retrieving animal: %s", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errResponse))
			log.Error("Error retrieving animal", slog.String("request-id", req.Context().Value(RequestIDKey("requestID")).(string)), slog.String("error", err.Error()))
			return
		}

		log.Info("Deleting animal", slog.String("request-id", req.Context().Value(RequestIDKey("requestID")).(string)))
		err = xSession.DeleteRecord(dbSession, animalPtr)
		if err != nil {
			errResponse := fmt.Sprintf("Failed to delete animal: %s", err.Error())
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(errResponse))
			log.Error("Failed to delete animal", slog.String("request-id", req.Context().Value(RequestIDKey("requestID")).(string)), slog.String("error", err.Error()))
			return
		}

		w.WriteHeader(http.StatusOK)
		res := fmt.Sprintf("Animal with ID %s deleted", animalID)
		w.Write([]byte(res))
		log.Info("Animal deleted successfully", slog.String("request-id", req.Context().Value(RequestIDKey("requestID")).(string)))
	}
}
