package controllers

import (
	"course/models"
	"encoding/json"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/bitly/go-simplejson"
)

type UserController struct {
	beego.Controller
}

func (this *UserController) Get() {
	method := this.GetString("method")
	if method == "id" {
		id, err := strconv.Atoi(this.Ctx.Input.Param(":id"))
		if err == nil {
			bodyJSON := simplejson.New()
			user, err := models.GetUserById(id)
			if err == nil {
				bodyJSON.Set("status", "success")
				dataMap := make(map[string]interface{})
				dataMap["id"] = id
				dataMap["name"] = user.Name
				dataMap["number"] = user.Number
				dataMap["email"] = user.Email
				dataMap["school"] = user.School
				dataMap["type"] = user.Type
				bodyJSON.Set("data", dataMap)
				body, _ := bodyJSON.Encode()
				this.Ctx.Output.Body(body)
			} else {
				this.Abort(models.ErrJson("user not exist"))
			}
		} else {
			this.Abort(models.ErrJson("invalid id format"))
		}
	} else if method == "token" {
		token := this.GetString("token")
		bodyJSON := simplejson.New()
		user, err := models.GetUserByToken(token)
		if err == nil {
			bodyJSON.Set("status", "success")
			dataMap := make(map[string]interface{})
			dataMap["id"] = user.Id
			dataMap["name"] = user.Name
			dataMap["number"] = user.Number
			dataMap["token"] = user.Token
			dataMap["email"] = user.Email
			dataMap["school"] = user.School
			dataMap["type"] = user.Type
			bodyJSON.Set("data", dataMap)
			body, _ := bodyJSON.Encode()
			this.Ctx.Output.Body(body)
		} else {
			this.Abort(models.ErrJson("user not exist"))
		}
	}
}

func (this *UserController) Post() {
	sess, _ := models.GlobalSessions.SessionStart(this.Ctx.ResponseWriter, this.Ctx.Request)
	defer sess.SessionRelease(this.Ctx.ResponseWriter)
	var user models.User
	if err := json.Unmarshal(this.Ctx.Input.RequestBody, &user); err == nil {
		_, err := models.GetUserByToken(user.Token)
		if err == nil {
			this.Abort(models.ErrJson("token duplicated"))
		}
		id, err := models.AddUser(&user)
		if err == nil {
			sess.Set("id", int(id))
			bodyJSON := simplejson.New()
			bodyJSON.Set("status", "success")
			bodyJSON.Set("id", id)
			body, _ := bodyJSON.Encode()
			this.Ctx.Output.Body(body)
		} else {
			this.Abort(models.ErrJson("fail to register, database error"))
		}
	} else {
		this.Abort(models.ErrJson("invalid data format"))
	}
}

func (this *UserController) Put() {
	sess, _ := models.GlobalSessions.SessionStart(this.Ctx.ResponseWriter, this.Ctx.Request)
	defer sess.SessionRelease(this.Ctx.ResponseWriter)
	var user models.User
	if err := json.Unmarshal(this.Ctx.Input.RequestBody, &user); err == nil {
		if sess.Get("id") == nil || sess.Get("id").(int) != user.Id {
			this.Abort(models.ErrJson("login expired"))
		}
		err := models.UpdateUserById(&user)
		if err == nil {
			this.Ctx.Output.Body(models.SuccessJson())
		} else {
			this.Abort(models.ErrJson("invlaid user id"))
		}
	} else {
		this.Abort(models.ErrJson("invalid data format"))
	}
}

func (this *UserController) Delete() {
	sess, _ := models.GlobalSessions.SessionStart(this.Ctx.ResponseWriter, this.Ctx.Request)
	defer sess.SessionRelease(this.Ctx.ResponseWriter)
	id, err := strconv.Atoi(this.Ctx.Input.Param(":id"))
	if sess.Get("id") == nil || sess.Get("id").(int) != id {
		this.Abort(models.ErrJson("login expired"))
	}
	if err != nil {
		this.Abort(models.ErrJson("invalid id"))
	}
	if err := models.DeleteUser(id); err == nil {
		this.Ctx.Output.Body(models.SuccessJson())
	} else {
		this.Abort(models.ErrJson("invalid id"))
	}
}
