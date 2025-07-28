package router

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/google/uuid"

	xModels "gomike/models"
	xSession "gomike/session"
	xDb "lib/dbchef"
)

func GetAnimal(dbSession *xDb.DBSession, log *slog.Logger) func(context.Context, string, io.ReadCloser) RespDetail {
	return func(reqCtx context.Context, reqObjID string, reqBody io.ReadCloser) RespDetail {
		if reqObjID == "" {
			errResponse := "Animal ID not provided in the URL"
			log.Error("Animal ID not provided in the URL", slog.String("request-id", reqCtx.Value(RequestIDKey("requestID")).(string)))
			return RespDetail{
				Statuscode: http.StatusBadRequest,
				Message:    []byte(errResponse),
			}
		}

		log.Info("Got animal ID from request", slog.String("request-id", reqCtx.Value(RequestIDKey("requestID")).(string)), slog.String("reqObjID", reqObjID))

		animalPtr, err := xSession.ReadRecord[xModels.Animal](dbSession, reqObjID)
		if err != nil {
			if strings.Contains(strings.ToLower(err.Error()), "record not found") {
				log.Error("Animal not found", slog.String("request-id", reqCtx.Value(RequestIDKey("requestID")).(string)))
				errResponse := fmt.Sprintf("Animal  with ID %s not found", reqObjID)
				return RespDetail{
					Statuscode: http.StatusNotFound,
					Message:    []byte(errResponse),
				}
			}
			log.Error("Error retrieving animal", slog.String("request-id", reqCtx.Value(RequestIDKey("requestID")).(string)), slog.String("error", err.Error()))
			errResponse := fmt.Sprintf("Error retrieving animal: %s", err.Error())
			return RespDetail{
				Statuscode: http.StatusInternalServerError,
				Message:    []byte(errResponse),
			}
		}

		objDetails, err := json.Marshal(animalPtr)
		if err != nil {
			log.Error("Failed to marshal animal details", slog.String("request-id", reqCtx.Value(RequestIDKey("requestID")).(string)), slog.String("error", err.Error()))
			errResponse := fmt.Sprintf("Failed to marshal animal details: %s", err.Error())
			return RespDetail{
				Statuscode: http.StatusInternalServerError,
				Message:    []byte(errResponse),
			}
		}
		log.Info("Animal details retrieved successfully", slog.String("request-id", reqCtx.Value(RequestIDKey("requestID")).(string)))
		return RespDetail{
			Statuscode: http.StatusOK,
			Message:    objDetails,
			Type:       "application/json",
		}
	}
}

func UpdateAnimal(dbSession *xDb.DBSession, log *slog.Logger) func(context.Context, string, io.ReadCloser) RespDetail {
	return func(reqCtx context.Context, reqObjID string, reqBody io.ReadCloser) RespDetail {
		if reqObjID == "" {
			log.Error("Animal ID not provided in the URL", slog.String("request-id", reqCtx.Value(RequestIDKey("requestID")).(string)))
			errResponse := "Animal ID not provided in the URL"
			return RespDetail{
				Statuscode: http.StatusBadRequest,
				Message:    []byte(errResponse),
			}
		}

		log.Info("Got animal ID from request", slog.String("request-id", reqCtx.Value(RequestIDKey("requestID")).(string)), slog.String("reqObjID", reqObjID))

		body, err := io.ReadAll(reqBody)
		if err != nil {
			log.Error("Failed to read request body", slog.String("request-id", reqCtx.Value(RequestIDKey("requestID")).(string)), slog.String("error", err.Error()))
			errResponse := fmt.Sprintf("Failed to read request body: %s", err.Error())
			return RespDetail{
				Statuscode: http.StatusInternalServerError,
				Message:    []byte(errResponse),
			}
		}
		if len(body) == 0 {
			log.Error("Request body is empty", slog.String("request-id", reqCtx.Value(RequestIDKey("requestID")).(string)))
			errResponse := "Request body is empty"
			return RespDetail{
				Statuscode: http.StatusBadRequest,
				Message:    []byte(errResponse),
			}
		}

		animalPtr, err := xSession.ReadRecord[xModels.Animal](dbSession, reqObjID)
		if err != nil {
			if strings.Contains(strings.ToLower(err.Error()), "record not found") {
				log.Info("Animal not found, Creating a new one", slog.String("request-id", reqCtx.Value(RequestIDKey("requestID")).(string)))
				// If the animal is not found, create a new animal with the provided ID
				animal := xModels.Animal{
					ID:            reqObjID,
					Kind:          "animal",
					Cloned:        false,
					ClonedFromRef: "",
				} // Create a new animal instance with the ID from the URL
				err = json.Unmarshal(body, &animal)
				if err != nil {
					log.Error("Failed to unmarshal request body into animal", slog.String("request-id", reqCtx.Value(RequestIDKey("requestID")).(string)), slog.String("error", err.Error()))
					errResponse := fmt.Sprintf("Failed to unmarshal request body into animal: %s", err.Error())
					return RespDetail{
						Statuscode: http.StatusBadRequest,
						Message:    []byte(errResponse),
					}
				}
				err = xSession.CreateRecord(dbSession, animal)
				if err != nil {
					log.Error("Failed to create animal", slog.String("request-id", reqCtx.Value(RequestIDKey("requestID")).(string)), slog.String("error", err.Error()))
					errResponse := fmt.Sprintf("Failed to create animal: %s", err.Error())
					return RespDetail{
						Statuscode: http.StatusInternalServerError,
						Message:    []byte(errResponse),
					}
				}

				log.Info("New animal created successfully", slog.String("request-id", reqCtx.Value(RequestIDKey("requestID")).(string)))
				return RespDetail{
					Statuscode: http.StatusNoContent,
				}
			}
			log.Error("Error retrieving animal", slog.String("request-id", reqCtx.Value(RequestIDKey("requestID")).(string)), slog.String("error", err.Error()))
			errResponse := fmt.Sprintf("Error retrieving animal with ID %s: %s", reqObjID, err.Error())
			return RespDetail{
				Statuscode: http.StatusInternalServerError,
				Message:    []byte(errResponse),
			}
		}

		log.Info("Updating existing animal", slog.String("request-id", reqCtx.Value(RequestIDKey("requestID")).(string)))
		err = json.Unmarshal(body, animalPtr)
		if err != nil {
			log.Error("Failed to unmarshal request body into animal", slog.String("request-id", reqCtx.Value(RequestIDKey("requestID")).(string)), slog.String("error", err.Error()))
			errResponse := fmt.Sprintf("Failed to unmarshal request body into animal: %s", err.Error())
			return RespDetail{
				Statuscode: http.StatusBadRequest,
				Message:    []byte(errResponse),
			}
		}

		err = xSession.UpdateRecord(dbSession, *animalPtr)
		if err != nil {
			log.Error("Failed to update animal", slog.String("request-id", reqCtx.Value(RequestIDKey("requestID")).(string)), slog.String("error", err.Error()))
			errResponse := fmt.Sprintf("Failed to update animal: %s", err.Error())
			return RespDetail{
				Statuscode: http.StatusInternalServerError,
				Message:    []byte(errResponse),
			}
		}

		log.Info("Animal updated successfully", slog.String("request-id", reqCtx.Value(RequestIDKey("requestID")).(string)))
		return RespDetail{
			Statuscode: http.StatusNoContent,
		}
	}
}

func CreateAnimal(dbSession *xDb.DBSession, log *slog.Logger) func(context.Context, string, io.ReadCloser) RespDetail {
	return func(reqCtx context.Context, reqObjID string, reqBody io.ReadCloser) RespDetail {
		body, err := io.ReadAll(reqBody)
		if err != nil {
			log.Error("Failed to read request body", slog.String("request-id", reqCtx.Value(RequestIDKey("requestID")).(string)), slog.String("error", err.Error()))
			errResponse := fmt.Sprintf("Failed to read request body: %s", err.Error())
			return RespDetail{
				Statuscode: http.StatusInternalServerError,
				Message:    []byte(errResponse),
			}
		}

		if len(body) == 0 {
			log.Error("Request body is empty", slog.String("request-id", reqCtx.Value(RequestIDKey("requestID")).(string)))
			errResponse := "Request body is empty"
			return RespDetail{
				Statuscode: http.StatusBadRequest,
				Message:    []byte(errResponse),
			}
		}

		animal := xModels.Animal{
			ID:            uuid.New().String(),
			Kind:          "animal",
			Cloned:        false,
			ClonedFromRef: "",
		} // Create a new animal instance

		log.Info("Creating a new animal with ID", slog.String("request-id", reqCtx.Value(RequestIDKey("requestID")).(string)), slog.String("reqObjID", animal.ID))
		err = json.Unmarshal(body, &animal)
		if err != nil {
			log.Error("Request body is empty", slog.String("request-id", reqCtx.Value(RequestIDKey("requestID")).(string)))
			errResponse := fmt.Sprintf("Failed to unmarshal request body into animal: %s", err.Error())
			return RespDetail{
				Statuscode: http.StatusBadRequest,
				Message:    []byte(errResponse),
			}
		}

		err = xSession.CreateRecord(dbSession, animal)
		if err != nil {
			log.Error("Request body is empty", slog.String("request-id", reqCtx.Value(RequestIDKey("requestID")).(string)))
			errResponse := fmt.Sprintf("Failed to create animal: %s", err.Error())
			return RespDetail{
				Statuscode: http.StatusInternalServerError,
				Message:    []byte(errResponse),
			}
		}

		log.Error("New Animal created successfully", slog.String("request-id", reqCtx.Value(RequestIDKey("requestID")).(string)))
		res := fmt.Sprintf("Animal with ID %s added", animal.ID)
		return RespDetail{
			Statuscode: http.StatusCreated,
			Message:    []byte(res),
		}

	}
}

func PatchAnimal(dbSession *xDb.DBSession, log *slog.Logger) func(context.Context, string, io.ReadCloser) RespDetail {
	return func(reqCtx context.Context, reqObjID string, reqBody io.ReadCloser) RespDetail {
		if reqObjID == "" {
			log.Error("Animal ID not provided in the URL", slog.String("request-id", reqCtx.Value(RequestIDKey("requestID")).(string)))
			errResponse := "Animal ID not provided in the URL"
			return RespDetail{
				Statuscode: http.StatusBadRequest,
				Message:    []byte(errResponse),
			}
		}

		log.Info("Got animal ID from request", slog.String("request-id", reqCtx.Value(RequestIDKey("requestID")).(string)), slog.String("reqObjID", reqObjID))

		body, err := io.ReadAll(reqBody)
		if err != nil {
			log.Error("Failed to read request body", slog.String("request-id", reqCtx.Value(RequestIDKey("requestID")).(string)), slog.String("error", err.Error()))
			errResponse := fmt.Sprintf("Failed to read request body: %s", err.Error())
			return RespDetail{
				Statuscode: http.StatusInternalServerError,
				Message:    []byte(errResponse),
			}
		}

		if len(body) == 0 {
			log.Error("Request body is empty", slog.String("request-id", reqCtx.Value(RequestIDKey("requestID")).(string)))
			errResponse := "Request body is empty"
			return RespDetail{
				Statuscode: http.StatusBadRequest,
				Message:    []byte(errResponse),
			}
		}

		animalPtr, err := xSession.ReadRecord[xModels.Animal](dbSession, reqObjID)
		if err != nil {
			if strings.Contains(strings.ToLower(err.Error()), "record not found") {
				log.Warn("Patch method called on non-existing animal, only allowed on existing records", slog.String("request-id", reqCtx.Value(RequestIDKey("requestID")).(string)))
				errResponse := fmt.Sprintf("Patch method is only allowed on existing records. Animal with ID %s not found", reqObjID)
				return RespDetail{
					Statuscode: http.StatusNotFound,
					Message:    []byte(errResponse),
				}
			}
			log.Error("Error retrieving animal", slog.String("request-id", reqCtx.Value(RequestIDKey("requestID")).(string)), slog.String("error", err.Error()))
			errResponse := fmt.Sprintf("Error retrieving animal with ID %s: %s", reqObjID, err.Error())
			return RespDetail{
				Statuscode: http.StatusInternalServerError,
				Message:    []byte(errResponse),
			}
		}

		err = json.Unmarshal(body, animalPtr)
		if err != nil {
			log.Error("Failed to unmarshal request body into animal", slog.String("request-id", reqCtx.Value(RequestIDKey("requestID")).(string)), slog.String("error", err.Error()))
			errResponse := fmt.Sprintf("Failed to unmarshal request body into animal: %s", err.Error())
			return RespDetail{
				Statuscode: http.StatusBadRequest,
				Message:    []byte(errResponse),
			}
		}

		log.Info("Patching existing animal", slog.String("request-id", reqCtx.Value(RequestIDKey("requestID")).(string)))
		err = xSession.UpdateRecord(dbSession, *animalPtr)
		if err != nil {
			log.Error("Failed to update animal", slog.String("request-id", reqCtx.Value(RequestIDKey("requestID")).(string)), slog.String("error", err.Error()))
			errResponse := fmt.Sprintf("Failed to update animal: %s", err.Error())
			return RespDetail{
				Statuscode: http.StatusInternalServerError,
				Message:    []byte(errResponse),
			}
		}

		log.Info("Animal patched successfully", slog.String("request-id", reqCtx.Value(RequestIDKey("requestID")).(string)))
		return RespDetail{
			Statuscode: http.StatusNoContent,
		}

	}
}

func DeleteAnimal(dbSession *xDb.DBSession, log *slog.Logger) func(context.Context, string, io.ReadCloser) RespDetail {
	return func(reqCtx context.Context, reqObjID string, reqBody io.ReadCloser) RespDetail {
		if reqObjID == "" {
			log.Error("Animal ID not provided in the URL", slog.String("request-id", reqCtx.Value(RequestIDKey("requestID")).(string)))
			errResponse := "Animal ID not provided in the URL"
			return RespDetail{
				Statuscode: http.StatusBadRequest,
				Message:    []byte(errResponse),
			}
		}

		log.Info("Got animal ID from request", slog.String("request-id", reqCtx.Value(RequestIDKey("requestID")).(string)), slog.String("reqObjID", reqObjID))

		animalPtr, err := xSession.ReadRecord[xModels.Animal](dbSession, reqObjID)
		if err != nil {
			if strings.Contains(strings.ToLower(err.Error()), "record not found") {
				log.Error("Animal ID not provided in the URL", slog.String("request-id", reqCtx.Value(RequestIDKey("requestID")).(string)))
				errResponse := fmt.Sprintf("Animal with ID %s not found", reqObjID)
				return RespDetail{
					Statuscode: http.StatusNotFound,
					Message:    []byte(errResponse),
				}
			}
			log.Error("Error retrieving animal", slog.String("request-id", reqCtx.Value(RequestIDKey("requestID")).(string)), slog.String("error", err.Error()))
			errResponse := fmt.Sprintf("Error retrieving animal: %s", err.Error())
			return RespDetail{
				Statuscode: http.StatusInternalServerError,
				Message:    []byte(errResponse),
			}
		}

		log.Info("Deleting animal", slog.String("request-id", reqCtx.Value(RequestIDKey("requestID")).(string)))
		err = xSession.DeleteRecord(dbSession, *animalPtr)
		if err != nil {
			log.Error("Failed to delete animal", slog.String("request-id", reqCtx.Value(RequestIDKey("requestID")).(string)), slog.String("error", err.Error()))
			errResponse := fmt.Sprintf("Failed to delete animal: %s", err.Error())
			return RespDetail{
				Statuscode: http.StatusInternalServerError,
				Message:    []byte(errResponse),
			}
		}

		res := fmt.Sprintf("Animal with ID %s deleted", reqObjID)
		log.Info("Animal deleted successfully", slog.String("request-id", reqCtx.Value(RequestIDKey("requestID")).(string)))
		return RespDetail{
			Statuscode: http.StatusOK,
			Message:    []byte(res),
		}

	}
}
