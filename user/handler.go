package user

import (
	"log"
	"net/http"

	"github.com/harkce/beasiswakita"
	"github.com/harkce/beasiswakita/authentication"
	"github.com/harkce/beasiswakita/errors"
	"github.com/harkce/beasiswakita/server/response"
	"github.com/harkce/beasiswakita/utils"

	"github.com/julienschmidt/httprouter"
)

type UserHandler struct{}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var user beasiswakita.User
	err := utils.Decode(r, &user)
	if err != nil {
		response.Error(w, errors.InternalServerError)
		log.Println(err)
		return
	}

	beasiswakita.Transaction, err = beasiswakita.DbMap.Begin()
	if err != nil {
		response.Error(w, errors.InternalServerError)
		log.Println(err)
		return
	}

	newUser, err := AddUser(user)
	if err != nil {
		beasiswakita.Transaction.Rollback()
		response.Error(w, errors.InternalServerError)
		log.Println(err)
		return
	}

	err = beasiswakita.Transaction.Commit()
	if err != nil {
		beasiswakita.Transaction.Rollback()
		response.Error(w, errors.InternalServerError)
		log.Println(err)
		return
	}

	token, err := authentication.Authorize(authentication.Owner{
		ID:           newUser.ID,
		EmailAddress: newUser.EmailAddress,
		Role:         newUser.Role,
	}, newUser.EmailAddress)
	if err != nil {
		response.Error(w, errors.Unauthorized)
	}

	loginData := struct {
		User  beasiswakita.User `json:"user"`
		Token string            `json:"token"`
	}{newUser, token}

	response.Created(w, loginData)
	return
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	var user beasiswakita.User

	err := utils.Decode(r, &user)
	if err != nil {
		log.Println(err)
		response.Error(w, errors.InternalServerError)
		return
	}

	var checkUser beasiswakita.User
	err = beasiswakita.DbMap.SelectOne(&checkUser, "select * from users where email_address = ?", user.EmailAddress)
	if err != nil {
		log.Println(err)
		response.Error(w, errors.Unauthorized)
		return
	}

	err = authentication.Authenticate(user.Password, checkUser.Password)
	if err != nil {
		log.Println(err)
		response.Error(w, errors.Unauthorized)
		return
	}

	token, err := authentication.Authorize(authentication.Owner{
		checkUser.ID,
		checkUser.EmailAddress,
		checkUser.Role,
	}, checkUser.EmailAddress)

	response.OK(w, struct {
		Token string `json:"token"`
	}{
		token,
	})
}