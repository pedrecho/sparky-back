package controllers

import (
	"context"
	"encoding/json"
	"errors"
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
		if errors.Is(err, http.ErrMissingFile) {
			return fmt.Errorf("file img does not exist in the form: %w", err)
		} else {
			return fmt.Errorf("getting form file img: %w", err)
		}
	}
	defer file.Close()
	imageName := handler.Filename
	imagePath, err := c.logic.SaveImg(file, imageName)
	if err != nil {
		return fmt.Errorf("saving image: %w", err)
	}
	user.ImgPath = imagePath
	id, err := c.logic.AddUser(context.TODO(), user)
	if err != nil {
		return fmt.Errorf("adding user: %w", err)
	}
	w.Write([]byte(fmt.Sprintf("%d", id)))
	return nil
}

func (c *Controller) UpdateUser(w http.ResponseWriter, req bunrouter.Request) error {
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
		if !errors.Is(err, http.ErrMissingFile) {
			return fmt.Errorf("getting form file img: %w", err)
		}
	} else {
		defer file.Close()
		imageName := handler.Filename
		imagePath, err := c.logic.SaveImg(file, imageName)
		if err != nil {
			return fmt.Errorf("saving image: %w", err)
		}
		user.ImgPath = imagePath
	}
	id, err := c.logic.UpdateUser(context.TODO(), user)
	if err != nil {
		return fmt.Errorf("updating user: %w", err)
	}
	w.Write([]byte(fmt.Sprintf("%d", id)))
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

func (c *Controller) GetUserByEmail(w http.ResponseWriter, req bunrouter.Request) error {
	emailStr, ok := req.Params().Get("email")
	if !ok {
		return fmt.Errorf("no id param")
	}
	user, err := c.logic.GetUserByEmail(context.TODO(), emailStr)
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
	reaction, err := convert.FormToReaction(req.PostForm)
	if err != nil {
		return fmt.Errorf("parsing post form: %w", err)
	}
	err = c.logic.SetReaction(context.TODO(), reaction)
	if err != nil {
		return fmt.Errorf("setting reaction: %w", err)
	}
	return nil
}

func (c *Controller) ClientConnection(w http.ResponseWriter, req bunrouter.Request) error {
	err := req.ParseMultipartForm(1 << 22)
	if err != nil {
		return fmt.Errorf("big multipartform size: %w", err)
	}
	msg, err := convert.FormToMessage(req.PostForm)
	if err != nil {
		return fmt.Errorf("parsing post form: %w", err)
	}

	//TODO it's sse((
	flusher, ok := w.(http.Flusher)
	if !ok {
		return fmt.Errorf("streaming unsupported: %w", err)
	}
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	sse := func(data []byte) {
		fmt.Fprintf(w, "data: %s\n\n", string(data))
		flusher.Flush()
	}
	err = c.logic.SendMessages(req.Context(), sse, msg)
	if err != nil {
		return fmt.Errorf("sending messages: %w", err)
	}

	<-req.Context().Done()
	return nil
}

func (c *Controller) NewMessage(w http.ResponseWriter, req bunrouter.Request) error {
	err := req.ParseMultipartForm(1 << 22)
	if err != nil {
		return fmt.Errorf("big multipartform size: %w", err)
	}
	msg, err := convert.FormToMessage(req.PostForm)
	if err != nil {
		return fmt.Errorf("parsing post form: %w", err)
	}
	return c.logic.NewMessage(msg)
}

func (c *Controller) GetRecommendations(w http.ResponseWriter, req bunrouter.Request) error {
	err := req.ParseMultipartForm(1 << 22)
	if err != nil {
		return fmt.Errorf("big multipartform size: %w", err)
	}
	filter, err := convert.FormToFilter(req.PostForm)
	if err != nil {
		return fmt.Errorf("parsing post form: %w", err)
	}
	users, err := c.logic.GetRecommendations(context.TODO(), filter)
	if err != nil {
		return fmt.Errorf("getting recomendations: %w", err)
	}
	jsonData, err := json.Marshal(users)
	if err != nil {
		return fmt.Errorf("marshaling json: %w", err)
	}
	w.Write(jsonData)
	return nil
}
