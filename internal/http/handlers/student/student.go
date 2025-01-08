package student

import (
	"net/http"
	"log/slog"
	"encoding/json"
	"errors"
	"io"
	"fmt"
	"strconv"

	"github.com/hasanm95/go-student-api/internal/types"
	"github.com/hasanm95/go-student-api/internal/utils/response"
	"github.com/hasanm95/go-student-api/internal/storage/sqlite"
	"github.com/go-playground/validator/v10"
)

func New(storage *sqlite.Storage) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var student types.Student

		err := json.NewDecoder(r.Body).Decode(&student)
		if errors.Is(err, io.EOF) {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("request body is empty")))
			return
		}

		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
			return
		}

		// Request validation
		if err := validator.New().Struct(student); err != nil {
			validateErrs := err.(validator.ValidationErrors)
			response.WriteJson(w, http.StatusBadRequest, response.ValidationError(validateErrs))
			return
		}

		// Convert to storage model
		storageStudent := &sqlite.Student{
			Name:  student.Name,
			Email: student.Email,
			Age:   student.Age,
		}

		// Save to database
		if err := storage.CreateStudent(storageStudent); err != nil {
			slog.Error("failed to create student", 
				slog.String("error", err.Error()),
				slog.String("email", student.Email),
			)
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(fmt.Errorf("failed to create student")))
			return
		}

		slog.Info("student created successfully", 
			slog.String("email", student.Email),
			slog.String("id", fmt.Sprint(storageStudent.ID)),
		)

		response.WriteJson(w, http.StatusCreated, map[string]interface{}{
			"success": true,
			"data": map[string]interface{}{
				"id":    storageStudent.ID,
				"name":  storageStudent.Name,
				"email": storageStudent.Email,
				"age":   storageStudent.Age,
			},
		})
	}
}

func GetStudentById(storage *sqlite.Storage) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request){
		id := r.PathValue("id")
		slog.Info("get student by id", slog.String("id", id))

		intId, err := strconv.ParseUint(id, 10, 64)

		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("invalid id")))
			return
		}

		student, err := storage.GetStudent(intId)

		if err != nil {
			slog.Error("failed to get student", slog.String("error", err.Error()))
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(fmt.Errorf("failed to get student")))
			return
		}

		if student == nil {
			response.WriteJson(w, http.StatusNotFound, response.GeneralError(fmt.Errorf("student not found")))
			return
		}

		response.WriteJson(w, http.StatusOK, map[string]interface{}{
			"success": true,
			"data": map[string]interface{}{
				"id":    student.ID,
				"name":  student.Name,
				"email": student.Email,
				"age":   student.Age,
			},
		})
	}
}

func GetStudents(storage *sqlite.Storage) http.HandlerFunc{
	return func(w http.ResponseWriter, r *http.Request){
		slog.Info("get all students")

		students, err := storage.GetStudents()

		if err != nil {
			slog.Error("failed to get students", slog.String("error", err.Error()))
			response.WriteJson(w, http.StatusInternalServerError, response.GeneralError(fmt.Errorf("failed to get students")))
			return
		}

		response.WriteJson(w, http.StatusOK, map[string]interface{}{
			"success": true,
			"data": students,
		})
	}
}