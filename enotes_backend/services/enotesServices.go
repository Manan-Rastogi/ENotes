package services

import (
	"fmt"
	"net/http"

	"github.com/Manan-Rastogi/enotes/configs"
	"github.com/Manan-Rastogi/enotes/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	SignUpService(mongodb models.Mongo, user models.User) (status int, msg string, err error)
	LogInService(mongodb models.Mongo, email, password string) (primitive.ObjectID, error)
	SaveNotesService(mongodb models.Mongo, user primitive.ObjectID, note models.Notes) (code, status int, msg string)
}

type service struct {
}

func NewService() Service {
	return &service{}
}

func (s *service) SignUpService(mongodb models.Mongo, user models.User) (status int, msg string, err error) {
	user.Date = configs.CurrTime()

	password, err := bcrypt.GenerateFromPassword([]byte(user.Password), 14)
	if err != nil {
		if err == bcrypt.ErrPasswordTooLong {
			msg = "Password too long."

		} else {
			msg = "Please use a different Password."
		}
		status = 1003
		return
	}

	user.Password = string(password)

	if err = mongodb.CreateNewUser(user, "E-Notes", "users"); err != nil {
		msg = "Error Creating New User"
		status = 1002
		return
	}

	msg = "User Created Successfully"
	status = 1
	err = nil

	return
}

func (s *service) LogInService(mongodb models.Mongo, email, password string) (primitive.ObjectID, error) {
	user, err := mongodb.ValidateEmailPassword(email, password, "E-Notes", "users")
	if err != nil {
		fmt.Printf("Invalid Username or password: %v", err.Error())
		return primitive.NilObjectID, fmt.Errorf("Invalid Username or Password.")
	}

	return user.ID, nil
}

func (s *service) SaveNotesService(mongodb models.Mongo, user primitive.ObjectID, note models.Notes) (code, status int, msg string) {
	note.Date = configs.CurrTime()
	note.User = user

	err := mongodb.CreateNewNote(note, user, "E-Notes", "notes")
	if err != nil {
		fmt.Println("Error Creating New Note: ", err.Error())
		return http.StatusInternalServerError, 1007, "Unable to save note for now."
	}

	return http.StatusOK, 1, "Notes Saved Successfully!"
}
