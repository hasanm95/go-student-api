package student

import (
	"net/http"
	"log/slog"
	"encoding/json"
	"errors"
	"io"
	"fmt"


	"github.com/hasanm95/go-student-api/internal/types"
	"github.com/hasanm95/go-student-api/internal/utils/response"
	"github.com/go-playground/validator/v10"
)

func New() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request){

		var student types.Student

		err := json.NewDecoder(r.Body).Decode(&student)
		if errors.Is(err, io.EOF) {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(fmt.Errorf("request body is empty")))
			return
		}

		if err != nil {
			response.WriteJson(w, http.StatusBadRequest, response.GeneralError(err))
		}

		// Request validation
		if err := validator.New().Struct(student); err != nil {
			validateErrs := err.(validator.ValidationErrors)
			response.WriteJson(w, http.StatusBadRequest, response.ValidationError(validateErrs))
			return
		}


		slog.Info("Creating a new student")

		response.WriteJson(w, http.StatusCreated, map[string]string{"success": "OK"})
	}
}