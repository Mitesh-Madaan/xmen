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
		// Create a response channel
		respChan := make(chan RespDetail)
		defer close(respChan)

		go func() {
			animalID := req.PathValue("animalID")
			if animalID == "" {
				errResponse := "Animal ID not provided in the URL"
				respChan <- RespDetail{
					statusCode: http.StatusBadRequest,
					message:    []byte(errResponse),
				}
				log.Error("Animal ID not provided in the URL", slog.String("request-id", req.Context().Value(RequestIDKey("requestID")).(string)))
				return
			}

			log.Info("Got animal ID from request", slog.String("request-id", req.Context().Value(RequestIDKey("requestID")).(string)), slog.String("animalID", animalID))

			animalPtr, err := xSession.ReadRecord[xModels.Animal](dbSession, animalID)
			if err != nil {
				if strings.Contains(strings.ToLower(err.Error()), "record not found") {
					errResponse := fmt.Sprintf("Animal  with ID %s not found", animalID)
					respChan <- RespDetail{
						statusCode: http.StatusNotFound,
						message:    []byte(errResponse),
					}
					log.Error("Animal not found", slog.String("request-id", req.Context().Value(RequestIDKey("requestID")).(string)))
					return
				}
				errResponse := fmt.Sprintf("Error retrieving animal: %s", err.Error())
				respChan <- RespDetail{
					statusCode: http.StatusInternalServerError,
					message:    []byte(errResponse),
				}
				log.Error("Error retrieving animal", slog.String("request-id", req.Context().Value(RequestIDKey("requestID")).(string)), slog.String("error", err.Error()))
				return
			}

			objDetails, err := json.Marshal(animalPtr)
			if err != nil {
				errResponse := fmt.Sprintf("Failed to marshal animal details: %s", err.Error())
				respChan <- RespDetail{
					statusCode: http.StatusInternalServerError,
					message:    []byte(errResponse),
				}
				log.Error("Failed to marshal animal details", slog.String("request-id", req.Context().Value(RequestIDKey("requestID")).(string)), slog.String("error", err.Error()))
				return
			}
			respChan <- RespDetail{
				statusCode: http.StatusOK,
				message:    objDetails,
			}
			log.Info("Animal details retrieved successfully", slog.String("request-id", req.Context().Value(RequestIDKey("requestID")).(string)))
		}()

		resp := <-respChan
		w.WriteHeader(resp.statusCode)
		w.Write(resp.message)
		log.Info("Response sent", slog.String("request-id", req.Context().Value(RequestIDKey("requestID")).(string)), slog.Int("status-code", resp.statusCode))
	}
}

func UpdateAnimal(dbSession *xDb.DBSession, log *slog.Logger) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		// Create a response channel
		respChan := make(chan RespDetail)
		defer close(respChan)

		go func() {
			animalID := req.PathValue("animalID")
			if animalID == "" {
				errResponse := "Animal ID not provided in the URL"
				respChan <- RespDetail{
					statusCode: http.StatusBadRequest,
					message:    []byte(errResponse),
				}
				log.Error("Animal ID not provided in the URL", slog.String("request-id", req.Context().Value(RequestIDKey("requestID")).(string)), slog.String("error", errResponse))
				return
			}

			log.Info("Got animal ID from request", slog.String("request-id", req.Context().Value(RequestIDKey("requestID")).(string)), slog.String("animalID", animalID))

			body, err := ioutil.ReadAll(req.Body)
			if err != nil {
				errResponse := fmt.Sprintf("Failed to read request body: %s", err.Error())
				respChan <- RespDetail{
					statusCode: http.StatusInternalServerError,
					message:    []byte(errResponse),
				}
				log.Error("Failed to read request body", slog.String("request-id", req.Context().Value(RequestIDKey("requestID")).(string)), slog.String("error", err.Error()))
				return
			}
			if len(body) == 0 {
				errResponse := "Request body is empty"
				respChan <- RespDetail{
					statusCode: http.StatusBadRequest,
					message:    []byte(errResponse),
				}
				log.Error("Request body is empty", slog.String("request-id", req.Context().Value(RequestIDKey("requestID")).(string)))
				return
			}

			animalPtr, err := xSession.ReadRecord[xModels.Animal](dbSession, animalID)
			if err != nil {
				if strings.Contains(strings.ToLower(err.Error()), "record not found") {
					log.Info("Animal not found, Creating a new one", slog.String("request-id", req.Context().Value(RequestIDKey("requestID")).(string)))
					// If the animal is not found, create a new animal with the provided ID
					animal := xModels.Animal{
						ID:            animalID,
						Kind:          "animal",
						Cloned:        false,
						ClonedFromRef: "",
					} // Create a new animal instance with the ID from the URL
					err = json.Unmarshal(body, &animal)
					if err != nil {
						errResponse := fmt.Sprintf("Failed to unmarshal request body into animal: %s", err.Error())
						respChan <- RespDetail{
							statusCode: http.StatusBadRequest,
							message:    []byte(errResponse),
						}
						log.Error("Failed to unmarshal request body into animal", slog.String("request-id", req.Context().Value(RequestIDKey("requestID")).(string)), slog.String("error", err.Error()))
						return
					}
					err = xSession.CreateRecord(dbSession, animal)
					if err != nil {
						errResponse := fmt.Sprintf("Failed to create animal: %s", err.Error())
						respChan <- RespDetail{
							statusCode: http.StatusInternalServerError,
							message:    []byte(errResponse),
						}
						log.Error("Failed to create animal", slog.String("request-id", req.Context().Value(RequestIDKey("requestID")).(string)), slog.String("error", err.Error()))
						return
					}

					respChan <- RespDetail{
						statusCode: http.StatusNoContent,
						message:    []byte(""),
					}
					log.Info("New animal created successfully", slog.String("request-id", req.Context().Value(RequestIDKey("requestID")).(string)))
					return
				}
				errResponse := fmt.Sprintf("Error retrieving animal with ID %s: %s", animalID, err.Error())
				respChan <- RespDetail{
					statusCode: http.StatusInternalServerError,
					message:    []byte(errResponse),
				}
				log.Error("Error retrieving animal", slog.String("request-id", req.Context().Value(RequestIDKey("requestID")).(string)), slog.String("error", err.Error()))
				return
			}

			log.Info("Updating existing animal", slog.String("request-id", req.Context().Value(RequestIDKey("requestID")).(string)))
			err = json.Unmarshal(body, animalPtr)
			if err != nil {
				errResponse := fmt.Sprintf("Failed to unmarshal request body into animal: %s", err.Error())
				respChan <- RespDetail{
					statusCode: http.StatusBadRequest,
					message:    []byte(errResponse),
				}
				log.Error("Failed to unmarshal request body into animal", slog.String("request-id", req.Context().Value(RequestIDKey("requestID")).(string)), slog.String("error", err.Error()))
				return
			}

			err = xSession.UpdateRecord(dbSession, *animalPtr)
			if err != nil {
				errResponse := fmt.Sprintf("Failed to update animal: %s", err.Error())
				respChan <- RespDetail{
					statusCode: http.StatusInternalServerError,
					message:    []byte(errResponse),
				}
				log.Error("Failed to update animal", slog.String("request-id", req.Context().Value(RequestIDKey("requestID")).(string)), slog.String("error", err.Error()))
				return
			}

			respChan <- RespDetail{
				statusCode: http.StatusNoContent,
				message:    []byte(""),
			}
			log.Info("Animal updated successfully", slog.String("request-id", req.Context().Value(RequestIDKey("requestID")).(string)))
		}()

		resp := <-respChan
		w.WriteHeader(resp.statusCode)
		w.Write(resp.message)
		log.Info("Response sent", slog.String("request-id", req.Context().Value(RequestIDKey("requestID")).(string)), slog.Int("status-code", resp.statusCode))
	}
}

func CreateAnimal(dbSession *xDb.DBSession, log *slog.Logger) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		// Create a response channel
		respChan := make(chan RespDetail)
		defer close(respChan)

		go func() {
			body, err := ioutil.ReadAll(req.Body)
			if err != nil {
				errResponse := fmt.Sprintf("Failed to read request body: %s", err.Error())
				respChan <- RespDetail{
					statusCode: http.StatusInternalServerError,
					message:    []byte(errResponse),
				}
				log.Error("Failed to read request body", slog.String("request-id", req.Context().Value(RequestIDKey("requestID")).(string)), slog.String("error", err.Error()))
				return
			}

			if len(body) == 0 {
				errResponse := "Request body is empty"
				respChan <- RespDetail{
					statusCode: http.StatusBadRequest,
					message:    []byte(errResponse),
				}
				log.Error("Request body is empty", slog.String("request-id", req.Context().Value(RequestIDKey("requestID")).(string)))
				return
			}

			animal := xModels.Animal{
				ID:            uuid.New().String(),
				Kind:          "animal",
				Cloned:        false,
				ClonedFromRef: "",
			} // Create a new animal instance

			log.Info("Creating a new animal with ID", slog.String("request-id", req.Context().Value(RequestIDKey("requestID")).(string)), slog.String("animalID", animal.ID))
			err = json.Unmarshal(body, &animal)
			if err != nil {
				errResponse := fmt.Sprintf("Failed to unmarshal request body into animal: %s", err.Error())
				respChan <- RespDetail{
					statusCode: http.StatusBadRequest,
					message:    []byte(errResponse),
				}
				log.Error("Failed to unmarshal request body into animal", slog.String("request-id", req.Context().Value(RequestIDKey("requestID")).(string)), slog.String("error", err.Error()))
				return
			}

			err = xSession.CreateRecord(dbSession, animal)
			if err != nil {
				errResponse := fmt.Sprintf("Failed to create animal: %s", err.Error())
				respChan <- RespDetail{
					statusCode: http.StatusInternalServerError,
					message:    []byte(errResponse),
				}
				log.Error("Failed to create animal", slog.String("request-id", req.Context().Value(RequestIDKey("requestID")).(string)), slog.String("error", err.Error()))
				return
			}

			res := fmt.Sprintf("Animal with ID %s added", animal.ID)
			respChan <- RespDetail{
				statusCode: http.StatusCreated,
				message:    []byte(res),
			}
			log.Error("New Animal created successfully", slog.String("request-id", req.Context().Value(RequestIDKey("requestID")).(string)))
		}()

		resp := <-respChan
		w.WriteHeader(resp.statusCode)
		w.Write(resp.message)
		log.Info("Response sent", slog.String("request-id", req.Context().Value(RequestIDKey("requestID")).(string)), slog.Int("status-code", resp.statusCode))
	}
}

func PatchAnimal(dbSession *xDb.DBSession, log *slog.Logger) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		// Create a response channel
		respChan := make(chan RespDetail)
		defer close(respChan)

		go func() {
			animalID := req.PathValue("animalID")
			if animalID == "" {
				errResponse := "Animal ID not provided in the URL"
				respChan <- RespDetail{
					statusCode: http.StatusBadRequest,
					message:    []byte(errResponse),
				}
				log.Error("Animal ID not provided in the URL", slog.String("request-id", req.Context().Value(RequestIDKey("requestID")).(string)), slog.String("error", errResponse))
				return
			}

			log.Info("Got animal ID from request", slog.String("request-id", req.Context().Value(RequestIDKey("requestID")).(string)), slog.String("animalID", animalID))

			body, err := ioutil.ReadAll(req.Body)
			if err != nil {
				errResponse := fmt.Sprintf("Failed to read request body: %s", err.Error())
				respChan <- RespDetail{
					statusCode: http.StatusInternalServerError,
					message:    []byte(errResponse),
				}
				log.Error("Failed to read request body", slog.String("request-id", req.Context().Value(RequestIDKey("requestID")).(string)), slog.String("error", err.Error()))
				return
			}

			if len(body) == 0 {
				errResponse := "Request body is empty"
				respChan <- RespDetail{
					statusCode: http.StatusBadRequest,
					message:    []byte(errResponse),
				}
				log.Error("Request body is empty", slog.String("request-id", req.Context().Value(RequestIDKey("requestID")).(string)))
				return
			}

			animalPtr, err := xSession.ReadRecord[xModels.Animal](dbSession, animalID)
			if err != nil {
				if strings.Contains(strings.ToLower(err.Error()), "record not found") {
					errResponse := fmt.Sprintf("Patch method is only allowed on existing records. Animal with ID %s not found", animalID)
					respChan <- RespDetail{
						statusCode: http.StatusNotFound,
						message:    []byte(errResponse),
					}
					log.Warn("Patch method called on non-existing animal, only allowed on existing records", slog.String("request-id", req.Context().Value(RequestIDKey("requestID")).(string)))
					return
				}
				errResponse := fmt.Sprintf("Error retrieving animal with ID %s: %s", animalID, err.Error())
				respChan <- RespDetail{
					statusCode: http.StatusInternalServerError,
					message:    []byte(errResponse),
				}
				log.Error("Error retrieving animal", slog.String("request-id", req.Context().Value(RequestIDKey("requestID")).(string)), slog.String("error", err.Error()))
				return
			}

			err = json.Unmarshal(body, animalPtr)
			if err != nil {
				errResponse := fmt.Sprintf("Failed to unmarshal request body into animal: %s", err.Error())
				respChan <- RespDetail{
					statusCode: http.StatusBadRequest,
					message:    []byte(errResponse),
				}
				log.Error("Failed to unmarshal request body into animal", slog.String("request-id", req.Context().Value(RequestIDKey("requestID")).(string)), slog.String("error", err.Error()))
				return
			}

			log.Info("Patching existing animal", slog.String("request-id", req.Context().Value(RequestIDKey("requestID")).(string)))
			err = xSession.UpdateRecord(dbSession, *animalPtr)
			if err != nil {
				errResponse := fmt.Sprintf("Failed to update animal: %s", err.Error())
				respChan <- RespDetail{
					statusCode: http.StatusInternalServerError,
					message:    []byte(errResponse),
				}
				log.Error("Failed to update animal", slog.String("request-id", req.Context().Value(RequestIDKey("requestID")).(string)), slog.String("error", err.Error()))
				return
			}

			respChan <- RespDetail{
				statusCode: http.StatusNoContent,
				message:    []byte(""),
			}
			log.Info("Animal patched successfully", slog.String("request-id", req.Context().Value(RequestIDKey("requestID")).(string)))
		}()

		resp := <-respChan
		w.WriteHeader(resp.statusCode)
		w.Write(resp.message)
		log.Info("Response sent", slog.String("request-id", req.Context().Value(RequestIDKey("requestID")).(string)), slog.Int("status-code", resp.statusCode))
	}
}

func DeleteAnimal(dbSession *xDb.DBSession, log *slog.Logger) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		// Create a response channel
		respChan := make(chan RespDetail)
		defer close(respChan)

		go func() {
			animalID := req.PathValue("animalID")
			if animalID == "" {
				errResponse := "Animal ID not provided in the URL"
				respChan <- RespDetail{
					statusCode: http.StatusBadRequest,
					message:    []byte(errResponse),
				}
				log.Error("Animal ID not provided in the URL", slog.String("request-id", req.Context().Value(RequestIDKey("requestID")).(string)), slog.String("error", errResponse))
				return
			}

			log.Info("Got animal ID from request", slog.String("request-id", req.Context().Value(RequestIDKey("requestID")).(string)), slog.String("animalID", animalID))

			animalPtr, err := xSession.ReadRecord[xModels.Animal](dbSession, animalID)
			if err != nil {
				if strings.Contains(strings.ToLower(err.Error()), "record not found") {
					errResponse := fmt.Sprintf("Animal with ID %s not found", animalID)
					respChan <- RespDetail{
						statusCode: http.StatusNotFound,
						message:    []byte(errResponse),
					}
					log.Error("Animal not found", slog.String("request-id", req.Context().Value(RequestIDKey("requestID")).(string)))
					return
				}
				errResponse := fmt.Sprintf("Error retrieving animal: %s", err.Error())
				respChan <- RespDetail{
					statusCode: http.StatusInternalServerError,
					message:    []byte(errResponse),
				}
				log.Error("Error retrieving animal", slog.String("request-id", req.Context().Value(RequestIDKey("requestID")).(string)), slog.String("error", err.Error()))
				return
			}

			log.Info("Deleting animal", slog.String("request-id", req.Context().Value(RequestIDKey("requestID")).(string)))
			err = xSession.DeleteRecord(dbSession, *animalPtr)
			if err != nil {
				errResponse := fmt.Sprintf("Failed to delete animal: %s", err.Error())
				respChan <- RespDetail{
					statusCode: http.StatusInternalServerError,
					message:    []byte(errResponse),
				}
				log.Error("Failed to delete animal", slog.String("request-id", req.Context().Value(RequestIDKey("requestID")).(string)), slog.String("error", err.Error()))
				return
			}

			res := fmt.Sprintf("Animal with ID %s deleted", animalID)
			respChan <- RespDetail{
				statusCode: http.StatusOK,
				message:    []byte(res),
			}
			log.Info("Animal deleted successfully", slog.String("request-id", req.Context().Value(RequestIDKey("requestID")).(string)))
		}()

		resp := <-respChan
		w.WriteHeader(resp.statusCode)
		w.Write(resp.message)
		log.Info("Response sent", slog.String("request-id", req.Context().Value(RequestIDKey("requestID")).(string)), slog.Int("status-code", resp.statusCode))
	}
}
