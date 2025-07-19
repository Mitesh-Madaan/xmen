package router

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log/slog"
	"net/http"
	"strings"

	xModels "gomike/models"
	xSession "gomike/session"
	xDb "lib/dbchef"

	"github.com/google/uuid"
)

func GetPerson(dbSession *xDb.DBSession, log *slog.Logger) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		// Create a response channel
		respChan := make(chan RespDetail)
		defer close(respChan)

		go func() {
			personID := req.PathValue("personID")
			if personID == "" {
				errResponse := "Person ID not provided in the URL"
				respChan <- RespDetail{
					statusCode: http.StatusBadRequest,
					message:    []byte(errResponse),
				}
				log.Error("Person ID not provided in the URL", slog.String("request-id", req.Context().Value(RequestIDKey("requestID")).(string)))
				return
			}

			log.Info("Got person ID from request", slog.String("request-id", req.Context().Value(RequestIDKey("requestID")).(string)), slog.String("personID", personID))

			personPtr, err := xSession.ReadRecord[xModels.Person](dbSession, personID)
			if err != nil {
				if strings.Contains(strings.ToLower(err.Error()), "record not found") {
					errResponse := fmt.Sprintf("Person  with ID %s not found", personID)
					respChan <- RespDetail{
						statusCode: http.StatusNotFound,
						message:    []byte(errResponse),
					}
					log.Error("Person not found", slog.String("request-id", req.Context().Value(RequestIDKey("requestID")).(string)))
					return
				}
				errResponse := fmt.Sprintf("Error retrieving person: %s", err.Error())
				respChan <- RespDetail{
					statusCode: http.StatusInternalServerError,
					message:    []byte(errResponse),
				}
				log.Error("Error retrieving person", slog.String("request-id", req.Context().Value(RequestIDKey("requestID")).(string)), slog.String("error", err.Error()))
				return
			}

			objDetails, err := json.Marshal(personPtr)
			if err != nil {
				errResponse := fmt.Sprintf("Failed to marshal person details: %s", err.Error())
				respChan <- RespDetail{
					statusCode: http.StatusInternalServerError,
					message:    []byte(errResponse),
				}
				log.Error("Failed to marshal person details", slog.String("request-id", req.Context().Value(RequestIDKey("requestID")).(string)), slog.String("error", err.Error()))
				return
			}
			respChan <- RespDetail{
				statusCode: http.StatusOK,
				message:    objDetails,
			}
			log.Info("Person details retrieved successfully", slog.String("request-id", req.Context().Value(RequestIDKey("requestID")).(string)))
		}()

		resp := <-respChan
		w.WriteHeader(resp.statusCode)
		w.Write(resp.message)
		log.Info("Response sent", slog.String("request-id", req.Context().Value(RequestIDKey("requestID")).(string)), slog.Int("status-code", resp.statusCode))
	}
}

func UpdatePerson(dbSession *xDb.DBSession, log *slog.Logger) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		// Create a response channel
		respChan := make(chan RespDetail)
		defer close(respChan)

		go func() {
			personID := req.PathValue("personID")
			if personID == "" {
				errResponse := "Person ID not provided in the URL"
				respChan <- RespDetail{
					statusCode: http.StatusBadRequest,
					message:    []byte(errResponse),
				}
				log.Error("Person ID not provided in the URL", slog.String("request-id", req.Context().Value(RequestIDKey("requestID")).(string)), slog.String("error", errResponse))
				return
			}

			log.Info("Got person ID from request", slog.String("request-id", req.Context().Value(RequestIDKey("requestID")).(string)), slog.String("personID", personID))

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

			personPtr, err := xSession.ReadRecord[xModels.Person](dbSession, personID)
			if err != nil {
				if strings.Contains(strings.ToLower(err.Error()), "record not found") {
					log.Info("Person not found, Creating a new one", slog.String("request-id", req.Context().Value(RequestIDKey("requestID")).(string)))
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
						respChan <- RespDetail{
							statusCode: http.StatusBadRequest,
							message:    []byte(errResponse),
						}
						log.Error("Failed to unmarshal request body into person", slog.String("request-id", req.Context().Value(RequestIDKey("requestID")).(string)), slog.String("error", err.Error()))
						return
					}
					err = xSession.CreateRecord(dbSession, person)
					if err != nil {
						errResponse := fmt.Sprintf("Failed to create person: %s", err.Error())
						respChan <- RespDetail{
							statusCode: http.StatusInternalServerError,
							message:    []byte(errResponse),
						}
						log.Error("Failed to create person", slog.String("request-id", req.Context().Value(RequestIDKey("requestID")).(string)), slog.String("error", err.Error()))
						return
					}

					respChan <- RespDetail{
						statusCode: http.StatusNoContent,
						message:    []byte(""),
					}
					log.Info("New person created successfully", slog.String("request-id", req.Context().Value(RequestIDKey("requestID")).(string)))
					return
				}
				errResponse := fmt.Sprintf("Error retrieving person with ID %s: %s", personID, err.Error())
				respChan <- RespDetail{
					statusCode: http.StatusInternalServerError,
					message:    []byte(errResponse),
				}
				log.Error("Error retrieving person", slog.String("request-id", req.Context().Value(RequestIDKey("requestID")).(string)), slog.String("error", err.Error()))
				return
			}

			log.Info("Updating existing person", slog.String("request-id", req.Context().Value(RequestIDKey("requestID")).(string)))
			err = json.Unmarshal(body, personPtr)
			if err != nil {
				errResponse := fmt.Sprintf("Failed to unmarshal request body into person: %s", err.Error())
				respChan <- RespDetail{
					statusCode: http.StatusBadRequest,
					message:    []byte(errResponse),
				}
				log.Error("Failed to unmarshal request body into person", slog.String("request-id", req.Context().Value(RequestIDKey("requestID")).(string)), slog.String("error", err.Error()))
				return
			}

			err = xSession.UpdateRecord(dbSession, *personPtr)
			if err != nil {
				errResponse := fmt.Sprintf("Failed to update person: %s", err.Error())
				respChan <- RespDetail{
					statusCode: http.StatusInternalServerError,
					message:    []byte(errResponse),
				}
				log.Error("Failed to update person", slog.String("request-id", req.Context().Value(RequestIDKey("requestID")).(string)), slog.String("error", err.Error()))
				return
			}

			respChan <- RespDetail{
				statusCode: http.StatusNoContent,
				message:    []byte(""),
			}
			log.Info("Person updated successfully", slog.String("request-id", req.Context().Value(RequestIDKey("requestID")).(string)))
		}()

		resp := <-respChan
		w.WriteHeader(resp.statusCode)
		w.Write(resp.message)
		log.Info("Response sent", slog.String("request-id", req.Context().Value(RequestIDKey("requestID")).(string)), slog.Int("status-code", resp.statusCode))
	}
}

func CreatePerson(dbSession *xDb.DBSession, log *slog.Logger) func(http.ResponseWriter, *http.Request) {
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

			person := xModels.Person{
				ID:            uuid.New().String(),
				Kind:          "person",
				Cloned:        false,
				ClonedFromRef: "",
			} // Create a new person instance

			log.Info("Creating a new person with ID", slog.String("request-id", req.Context().Value(RequestIDKey("requestID")).(string)), slog.String("personID", person.ID))
			err = json.Unmarshal(body, &person)
			if err != nil {
				errResponse := fmt.Sprintf("Failed to unmarshal request body into person: %s", err.Error())
				respChan <- RespDetail{
					statusCode: http.StatusBadRequest,
					message:    []byte(errResponse),
				}
				log.Error("Failed to unmarshal request body into person", slog.String("request-id", req.Context().Value(RequestIDKey("requestID")).(string)), slog.String("error", err.Error()))
				return
			}

			err = xSession.CreateRecord(dbSession, person)
			if err != nil {
				errResponse := fmt.Sprintf("Failed to create person: %s", err.Error())
				respChan <- RespDetail{
					statusCode: http.StatusInternalServerError,
					message:    []byte(errResponse),
				}
				log.Error("Failed to create person", slog.String("request-id", req.Context().Value(RequestIDKey("requestID")).(string)), slog.String("error", err.Error()))
				return
			}

			res := fmt.Sprintf("Person with ID %s added", person.ID)
			respChan <- RespDetail{
				statusCode: http.StatusCreated,
				message:    []byte(res),
			}
			log.Error("New Person created successfully", slog.String("request-id", req.Context().Value(RequestIDKey("requestID")).(string)))
		}()

		resp := <-respChan
		w.WriteHeader(resp.statusCode)
		w.Write(resp.message)
		log.Info("Response sent", slog.String("request-id", req.Context().Value(RequestIDKey("requestID")).(string)), slog.Int("status-code", resp.statusCode))
	}
}

func PatchPerson(dbSession *xDb.DBSession, log *slog.Logger) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		// Create a response channel
		respChan := make(chan RespDetail)
		defer close(respChan)

		go func() {
			personID := req.PathValue("personID")
			if personID == "" {
				errResponse := "Person ID not provided in the URL"
				respChan <- RespDetail{
					statusCode: http.StatusBadRequest,
					message:    []byte(errResponse),
				}
				log.Error("Person ID not provided in the URL", slog.String("request-id", req.Context().Value(RequestIDKey("requestID")).(string)), slog.String("error", errResponse))
				return
			}

			log.Info("Got person ID from request", slog.String("request-id", req.Context().Value(RequestIDKey("requestID")).(string)), slog.String("personID", personID))

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

			personPtr, err := xSession.ReadRecord[xModels.Person](dbSession, personID)
			if err != nil {
				if strings.Contains(strings.ToLower(err.Error()), "record not found") {
					errResponse := fmt.Sprintf("Patch method is only allowed on existing records. Person with ID %s not found", personID)
					respChan <- RespDetail{
						statusCode: http.StatusNotFound,
						message:    []byte(errResponse),
					}
					log.Warn("Patch method called on non-existing person, only allowed on existing records", slog.String("request-id", req.Context().Value(RequestIDKey("requestID")).(string)))
					return
				}
				errResponse := fmt.Sprintf("Error retrieving person with ID %s: %s", personID, err.Error())
				respChan <- RespDetail{
					statusCode: http.StatusInternalServerError,
					message:    []byte(errResponse),
				}
				log.Error("Error retrieving person", slog.String("request-id", req.Context().Value(RequestIDKey("requestID")).(string)), slog.String("error", err.Error()))
				return
			}

			err = json.Unmarshal(body, personPtr)
			if err != nil {
				errResponse := fmt.Sprintf("Failed to unmarshal request body into person: %s", err.Error())
				respChan <- RespDetail{
					statusCode: http.StatusBadRequest,
					message:    []byte(errResponse),
				}
				log.Error("Failed to unmarshal request body into person", slog.String("request-id", req.Context().Value(RequestIDKey("requestID")).(string)), slog.String("error", err.Error()))
				return
			}

			log.Info("Patching existing person", slog.String("request-id", req.Context().Value(RequestIDKey("requestID")).(string)))
			err = xSession.UpdateRecord(dbSession, *personPtr)
			if err != nil {
				errResponse := fmt.Sprintf("Failed to update person: %s", err.Error())
				respChan <- RespDetail{
					statusCode: http.StatusInternalServerError,
					message:    []byte(errResponse),
				}
				log.Error("Failed to update person", slog.String("request-id", req.Context().Value(RequestIDKey("requestID")).(string)), slog.String("error", err.Error()))
				return
			}

			respChan <- RespDetail{
				statusCode: http.StatusNoContent,
				message:    []byte(""),
			}
			log.Info("Person patched successfully", slog.String("request-id", req.Context().Value(RequestIDKey("requestID")).(string)))
		}()

		resp := <-respChan
		w.WriteHeader(resp.statusCode)
		w.Write(resp.message)
		log.Info("Response sent", slog.String("request-id", req.Context().Value(RequestIDKey("requestID")).(string)), slog.Int("status-code", resp.statusCode))
	}
}

func DeletePerson(dbSession *xDb.DBSession, log *slog.Logger) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		// Create a response channel
		respChan := make(chan RespDetail)
		defer close(respChan)

		go func() {
			personID := req.PathValue("personID")
			if personID == "" {
				errResponse := "Person ID not provided in the URL"
				respChan <- RespDetail{
					statusCode: http.StatusBadRequest,
					message:    []byte(errResponse),
				}
				log.Error("Person ID not provided in the URL", slog.String("request-id", req.Context().Value(RequestIDKey("requestID")).(string)), slog.String("error", errResponse))
				return
			}

			log.Info("Got person ID from request", slog.String("request-id", req.Context().Value(RequestIDKey("requestID")).(string)), slog.String("personID", personID))

			personPtr, err := xSession.ReadRecord[xModels.Person](dbSession, personID)
			if err != nil {
				if strings.Contains(strings.ToLower(err.Error()), "record not found") {
					errResponse := fmt.Sprintf("Person with ID %s not found", personID)
					respChan <- RespDetail{
						statusCode: http.StatusNotFound,
						message:    []byte(errResponse),
					}
					log.Error("Person not found", slog.String("request-id", req.Context().Value(RequestIDKey("requestID")).(string)))
					return
				}
				errResponse := fmt.Sprintf("Error retrieving person: %s", err.Error())
				respChan <- RespDetail{
					statusCode: http.StatusInternalServerError,
					message:    []byte(errResponse),
				}
				log.Error("Error retrieving person", slog.String("request-id", req.Context().Value(RequestIDKey("requestID")).(string)), slog.String("error", err.Error()))
				return
			}

			log.Info("Deleting person", slog.String("request-id", req.Context().Value(RequestIDKey("requestID")).(string)))
			err = xSession.DeleteRecord(dbSession, *personPtr)
			if err != nil {
				errResponse := fmt.Sprintf("Failed to delete person: %s", err.Error())
				respChan <- RespDetail{
					statusCode: http.StatusInternalServerError,
					message:    []byte(errResponse),
				}
				log.Error("Failed to delete person", slog.String("request-id", req.Context().Value(RequestIDKey("requestID")).(string)), slog.String("error", err.Error()))
				return
			}

			res := fmt.Sprintf("Person with ID %s deleted", personID)
			respChan <- RespDetail{
				statusCode: http.StatusOK,
				message:    []byte(res),
			}
			log.Info("Person deleted successfully", slog.String("request-id", req.Context().Value(RequestIDKey("requestID")).(string)))
		}()

		resp := <-respChan
		w.WriteHeader(resp.statusCode)
		w.Write(resp.message)
		log.Info("Response sent", slog.String("request-id", req.Context().Value(RequestIDKey("requestID")).(string)), slog.Int("status-code", resp.statusCode))
	}
}
