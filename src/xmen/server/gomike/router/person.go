package router

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"

	xModels "gomike/models"
	xSession "gomike/session"
	xDb "lib/dbchef"

	"github.com/google/uuid"
)

func GetPerson(dbSession *xDb.DBSession, log *slog.Logger) func(context.Context, string, io.ReadCloser) RespDetail {
	return func(reqCtx context.Context, reqObjID string, reqBody io.ReadCloser) RespDetail {
		if reqObjID == "" {
			log.Error("Person ID not provided in the URL", slog.String("request-id", reqCtx.Value(RequestIDKey("requestID")).(string)))
			errResponse := "Person ID not provided in the URL"
			return RespDetail{
				Statuscode: http.StatusBadRequest,
				Message:    []byte(errResponse),
			}
		}

		log.Info("Got person ID from request", slog.String("request-id", reqCtx.Value(RequestIDKey("requestID")).(string)), slog.String("reqObjID", reqObjID))

		personPtr, err := xSession.ReadRecord[xModels.Person](dbSession, reqObjID)
		if err != nil {
			if strings.Contains(strings.ToLower(err.Error()), "record not found") {
				log.Error("Person not found", slog.String("request-id", reqCtx.Value(RequestIDKey("requestID")).(string)))
				errResponse := fmt.Sprintf("Person  with ID %s not found", reqObjID)
				return RespDetail{
					Statuscode: http.StatusNotFound,
					Message:    []byte(errResponse),
				}
			}
			log.Error("Error retrieving person", slog.String("request-id", reqCtx.Value(RequestIDKey("requestID")).(string)), slog.String("error", err.Error()))
			errResponse := fmt.Sprintf("Error retrieving person: %s", err.Error())
			return RespDetail{
				Statuscode: http.StatusInternalServerError,
				Message:    []byte(errResponse),
			}
		}

		objDetails, err := json.Marshal(personPtr)
		if err != nil {
			log.Error("Failed to marshal person details", slog.String("request-id", reqCtx.Value(RequestIDKey("requestID")).(string)), slog.String("error", err.Error()))
			errResponse := fmt.Sprintf("Failed to marshal person details: %s", err.Error())
			return RespDetail{
				Statuscode: http.StatusInternalServerError,
				Message:    []byte(errResponse),
			}
		}
		log.Info("Person details retrieved successfully", slog.String("request-id", reqCtx.Value(RequestIDKey("requestID")).(string)))
		return RespDetail{
			Statuscode: http.StatusOK,
			Message:    objDetails,
			Type:       "application/json",
		}
	}
}

func UpdatePerson(dbSession *xDb.DBSession, log *slog.Logger) func(context.Context, string, io.ReadCloser) RespDetail {
	return func(reqCtx context.Context, reqObjID string, reqBody io.ReadCloser) RespDetail {
		if reqObjID == "" {
			log.Error("Person ID not provided in the URL", slog.String("request-id", reqCtx.Value(RequestIDKey("requestID")).(string)))
			errResponse := "Person ID not provided in the URL"
			return RespDetail{
				Statuscode: http.StatusBadRequest,
				Message:    []byte(errResponse),
			}
		}

		log.Info("Got person ID from request", slog.String("request-id", reqCtx.Value(RequestIDKey("requestID")).(string)), slog.String("reqObjID", reqObjID))

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

		personPtr, err := xSession.ReadRecord[xModels.Person](dbSession, reqObjID)
		if err != nil {
			if strings.Contains(strings.ToLower(err.Error()), "record not found") {
				log.Info("Person not found, Creating a new one", slog.String("request-id", reqCtx.Value(RequestIDKey("requestID")).(string)))
				// If the person is not found, create a new person with the provided ID
				person := xModels.Person{
					ID:            reqObjID,
					Kind:          "person",
					Cloned:        false,
					ClonedFromRef: "",
				} // Create a new person instance with the ID from the URL
				err = json.Unmarshal(body, &person)
				if err != nil {
					errResponse := fmt.Sprintf("Failed to unmarshal request body into person: %s", err.Error())
					log.Error("Failed to unmarshal request body into person", slog.String("request-id", reqCtx.Value(RequestIDKey("requestID")).(string)), slog.String("error", err.Error()))
					return RespDetail{
						Statuscode: http.StatusBadRequest,
						Message:    []byte(errResponse),
					}

				}
				err = xSession.CreateRecord(dbSession, person)
				if err != nil {
					errResponse := fmt.Sprintf("Failed to create person: %s", err.Error())
					log.Error("Failed to create person", slog.String("request-id", reqCtx.Value(RequestIDKey("requestID")).(string)), slog.String("error", err.Error()))
					return RespDetail{
						Statuscode: http.StatusInternalServerError,
						Message:    []byte(errResponse),
					}

				}
				log.Info("New person created successfully", slog.String("request-id", reqCtx.Value(RequestIDKey("requestID")).(string)))
				return RespDetail{
					Statuscode: http.StatusNoContent,
				}

			}
			log.Error("Error retrieving person", slog.String("request-id", reqCtx.Value(RequestIDKey("requestID")).(string)), slog.String("error", err.Error()))
			errResponse := fmt.Sprintf("Error retrieving person with ID %s: %s", reqObjID, err.Error())
			return RespDetail{
				Statuscode: http.StatusInternalServerError,
				Message:    []byte(errResponse),
			}

		}

		log.Info("Updating existing person", slog.String("request-id", reqCtx.Value(RequestIDKey("requestID")).(string)))
		err = json.Unmarshal(body, personPtr)
		if err != nil {
			log.Error("Failed to unmarshal request body into person", slog.String("request-id", reqCtx.Value(RequestIDKey("requestID")).(string)), slog.String("error", err.Error()))
			errResponse := fmt.Sprintf("Failed to unmarshal request body into person: %s", err.Error())
			return RespDetail{
				Statuscode: http.StatusBadRequest,
				Message:    []byte(errResponse),
			}

		}

		err = xSession.UpdateRecord(dbSession, *personPtr)
		if err != nil {
			log.Error("Failed to update person", slog.String("request-id", reqCtx.Value(RequestIDKey("requestID")).(string)), slog.String("error", err.Error()))
			errResponse := fmt.Sprintf("Failed to update person: %s", err.Error())
			return RespDetail{
				Statuscode: http.StatusInternalServerError,
				Message:    []byte(errResponse),
			}

		}

		log.Info("Person updated successfully", slog.String("request-id", reqCtx.Value(RequestIDKey("requestID")).(string)))
		return RespDetail{
			Statuscode: http.StatusNoContent,
		}
	}
}

func CreatePerson(dbSession *xDb.DBSession, log *slog.Logger) func(context.Context, string, io.ReadCloser) RespDetail {
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

		person := xModels.Person{
			ID:            uuid.New().String(),
			Kind:          "person",
			Cloned:        false,
			ClonedFromRef: "",
		} // Create a new person instance

		log.Info("Creating a new person with ID", slog.String("request-id", reqCtx.Value(RequestIDKey("requestID")).(string)), slog.String("reqObjID", person.ID))
		err = json.Unmarshal(body, &person)
		if err != nil {
			log.Error("Failed to unmarshal request body into person", slog.String("request-id", reqCtx.Value(RequestIDKey("requestID")).(string)), slog.String("error", err.Error()))
			errResponse := fmt.Sprintf("Failed to unmarshal request body into person: %s", err.Error())
			return RespDetail{
				Statuscode: http.StatusBadRequest,
				Message:    []byte(errResponse),
			}

		}

		err = xSession.CreateRecord(dbSession, person)
		if err != nil {
			log.Error("Failed to create person", slog.String("request-id", reqCtx.Value(RequestIDKey("requestID")).(string)), slog.String("error", err.Error()))
			errResponse := fmt.Sprintf("Failed to create person: %s", err.Error())
			return RespDetail{
				Statuscode: http.StatusInternalServerError,
				Message:    []byte(errResponse),
			}

		}

		res := fmt.Sprintf("Person with ID %s added", person.ID)
		log.Info("New Person created successfully", slog.String("request-id", reqCtx.Value(RequestIDKey("requestID")).(string)))
		return RespDetail{
			Statuscode: http.StatusCreated,
			Message:    []byte(res),
		}

	}
}

func PatchPerson(dbSession *xDb.DBSession, log *slog.Logger) func(context.Context, string, io.ReadCloser) RespDetail {
	return func(reqCtx context.Context, reqObjID string, reqBody io.ReadCloser) RespDetail {
		if reqObjID == "" {
			log.Error("Person ID not provided in the URL", slog.String("request-id", reqCtx.Value(RequestIDKey("requestID")).(string)))
			errResponse := "Person ID not provided in the URL"
			return RespDetail{
				Statuscode: http.StatusBadRequest,
				Message:    []byte(errResponse),
			}

		}

		log.Info("Got person ID from request", slog.String("request-id", reqCtx.Value(RequestIDKey("requestID")).(string)), slog.String("reqObjID", reqObjID))

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

		personPtr, err := xSession.ReadRecord[xModels.Person](dbSession, reqObjID)
		if err != nil {
			if strings.Contains(strings.ToLower(err.Error()), "record not found") {
				log.Warn("Patch method called on non-existing person, only allowed on existing records", slog.String("request-id", reqCtx.Value(RequestIDKey("requestID")).(string)))
				errResponse := fmt.Sprintf("Patch method is only allowed on existing records. Person with ID %s not found", reqObjID)
				return RespDetail{
					Statuscode: http.StatusNotFound,
					Message:    []byte(errResponse),
				}

			}
			log.Error("Error retrieving person", slog.String("request-id", reqCtx.Value(RequestIDKey("requestID")).(string)), slog.String("error", err.Error()))
			errResponse := fmt.Sprintf("Error retrieving person with ID %s: %s", reqObjID, err.Error())
			return RespDetail{
				Statuscode: http.StatusInternalServerError,
				Message:    []byte(errResponse),
			}

		}

		err = json.Unmarshal(body, personPtr)
		if err != nil {
			log.Error("Failed to unmarshal request body into person", slog.String("request-id", reqCtx.Value(RequestIDKey("requestID")).(string)), slog.String("error", err.Error()))
			errResponse := fmt.Sprintf("Failed to unmarshal request body into person: %s", err.Error())
			return RespDetail{
				Statuscode: http.StatusBadRequest,
				Message:    []byte(errResponse),
			}

		}

		log.Info("Patching existing person", slog.String("request-id", reqCtx.Value(RequestIDKey("requestID")).(string)))
		err = xSession.UpdateRecord(dbSession, *personPtr)
		if err != nil {
			log.Error("Failed to update person", slog.String("request-id", reqCtx.Value(RequestIDKey("requestID")).(string)), slog.String("error", err.Error()))
			errResponse := fmt.Sprintf("Failed to update person: %s", err.Error())
			return RespDetail{
				Statuscode: http.StatusInternalServerError,
				Message:    []byte(errResponse),
			}

		}

		log.Info("Person patched successfully", slog.String("request-id", reqCtx.Value(RequestIDKey("requestID")).(string)))
		return RespDetail{
			Statuscode: http.StatusNoContent,
		}

	}
}

func DeletePerson(dbSession *xDb.DBSession, log *slog.Logger) func(context.Context, string, io.ReadCloser) RespDetail {
	return func(reqCtx context.Context, reqObjID string, reqBody io.ReadCloser) RespDetail {
		if reqObjID == "" {
			log.Error("Person ID not provided in the URL", slog.String("request-id", reqCtx.Value(RequestIDKey("requestID")).(string)))
			errResponse := "Person ID not provided in the URL"
			return RespDetail{
				Statuscode: http.StatusBadRequest,
				Message:    []byte(errResponse),
			}

		}

		log.Info("Got person ID from request", slog.String("request-id", reqCtx.Value(RequestIDKey("requestID")).(string)), slog.String("reqObjID", reqObjID))

		personPtr, err := xSession.ReadRecord[xModels.Person](dbSession, reqObjID)
		if err != nil {
			if strings.Contains(strings.ToLower(err.Error()), "record not found") {
				log.Error("Person not found", slog.String("request-id", reqCtx.Value(RequestIDKey("requestID")).(string)))
				errResponse := fmt.Sprintf("Person with ID %s not found", reqObjID)
				return RespDetail{
					Statuscode: http.StatusNotFound,
					Message:    []byte(errResponse),
				}

			}
			log.Error("Error retrieving person", slog.String("request-id", reqCtx.Value(RequestIDKey("requestID")).(string)), slog.String("error", err.Error()))
			errResponse := fmt.Sprintf("Error retrieving person: %s", err.Error())
			return RespDetail{
				Statuscode: http.StatusInternalServerError,
				Message:    []byte(errResponse),
			}

		}

		log.Info("Deleting person", slog.String("request-id", reqCtx.Value(RequestIDKey("requestID")).(string)))
		err = xSession.DeleteRecord(dbSession, *personPtr)
		if err != nil {
			log.Error("Failed to delete person", slog.String("request-id", reqCtx.Value(RequestIDKey("requestID")).(string)), slog.String("error", err.Error()))
			errResponse := fmt.Sprintf("Failed to delete person: %s", err.Error())
			return RespDetail{
				Statuscode: http.StatusInternalServerError,
				Message:    []byte(errResponse),
			}

		}

		log.Info("Person deleted successfully", slog.String("request-id", reqCtx.Value(RequestIDKey("requestID")).(string)))
		res := fmt.Sprintf("Person with ID %s deleted", reqObjID)
		return RespDetail{
			Statuscode: http.StatusOK,
			Message:    []byte(res),
		}

	}
}
