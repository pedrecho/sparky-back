package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/uptrace/bunrouter"
	"net/http"
	"sparky-back/internal/convert"
	"sparky-back/internal/logic"
	"strconv"
)

type Controller struct {
	logic *logic.Logic
}

func New(l *logic.Logic) *Controller {
	return &Controller{
		logic: l,
	}
}

func (c *Controller) AddUser(w http.ResponseWriter, req bunrouter.Request) error {
	err := req.ParseMultipartForm(1 << 22)
	if err != nil {
		return fmt.Errorf("big multipartform size: %w", err)
	}
	user, err := convert.FormToUser(req.PostForm)
	if err != nil {
		return fmt.Errorf("parsing post form: %w", err)
	}
	file, handler, err := req.FormFile("img")
	if err != nil {
		return fmt.Errorf("gettin form file img: %w", err)
	}
	defer file.Close()
	imageName := handler.Filename
	imagePath, err := c.logic.SaveImg(file, imageName)
	if err != nil {
		return fmt.Errorf("saving image: %w", err)
	}
	user.ImgPath = imagePath
	err = c.logic.SaveUser(context.TODO(), user)
	if err != nil {
		return fmt.Errorf("saving user: %w", err)
	}
	w.Write([]byte("user added"))
	return nil
}

func (c *Controller) Login(w http.ResponseWriter, req bunrouter.Request) error {
	err := req.ParseMultipartForm(1 << 22)
	if err != nil {
		return fmt.Errorf("big multipartform size: %w", err)
	}
	user, err := convert.FormToUser(req.PostForm)
	if err != nil {
		return fmt.Errorf("parsing post form: %w", err)
	}
	id, err := c.logic.LogIn(context.TODO(), user.Email, user.Password)
	if err != nil {
		return fmt.Errorf("login: %w", err)
	}
	w.Write([]byte(strconv.FormatInt(id, 10)))
	return nil
}

func (c *Controller) GetUserByID(w http.ResponseWriter, req bunrouter.Request) error {
	idStr, ok := req.Params().Get("id")
	if !ok {
		return fmt.Errorf("no id param")
	}
	id, err := strconv.ParseInt(idStr, 10, 64)
	user, err := c.logic.GetUserByID(context.TODO(), id)
	if err != nil {
		return fmt.Errorf("getting user: %w", err)
	}
	jsonData, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf("marshaling json: %w", err)
	}
	w.Write(jsonData)
	return nil
}

func (c *Controller) GetFile(w http.ResponseWriter, req bunrouter.Request) error {
	filename, ok := req.Params().Get("filename")
	if !ok {
		return fmt.Errorf("no filename param")
	}
	data, err := c.logic.GetFile(filename)
	if err != nil {
		return fmt.Errorf("getting file: %w", err)
	}
	w.Write(data)
	return nil
}

func (c *Controller) SetReaction(w http.ResponseWriter, req bunrouter.Request) error {
	err := req.ParseMultipartForm(1 << 22)
	if err != nil {
		return fmt.Errorf("big multipartform size: %w", err)
	}
	if err != nil {
		return fmt.Errorf("big multipartform size: %w", err)
	}
	reaction, err := convert.FormToReaction(req.PostForm)
	if err != nil {
		return fmt.Errorf("parsing post form: %w", err)
	}
	err = c.logic.SetReaction(context.TODO(), reaction)
	if err != nil {
		return fmt.Errorf("login: %w", err)
	}
	return nil
}
