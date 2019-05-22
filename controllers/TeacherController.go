package controllers

import (
	"course/models"
	"encoding/json"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/bitly/go-simplejson"
)

type TeacherController struct {
	beego.Controller
}

func (this *TeacherController) Get() {
	id, err := strconv.Atoi(this.Ctx.Input.Param(":id"))
	if err == nil {
		bodyJSON := simplejson.New()
		teacher, err := models.GetTeacherById(id)
		if err == nil {
			bodyJSON.Set("status", "success")
			dataMap := make(map[string]interface{})
			dataMap["id"] = id
			dataMap["name"] = teacher.Name
			bodyJSON.Set("data", dataMap)
			body, _ := bodyJSON.Encode()
			this.Ctx.Output.Body(body)
		} else {
			this.Abort(models.ErrJson("user not exist"))
		}
	} else {
		this.Abort(models.ErrJson("invalid id format"))
	}
}

func (this *TeacherController) Post() {
	var teacher models.Teacher
	if err := json.Unmarshal(this.Ctx.Input.RequestBody, &teacher); err == nil {
		id, err := models.AddTeacher(&teacher)
		if err == nil {
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

func (this *TeacherController) Put() {
	var teacher models.Teacher
	if err := json.Unmarshal(this.Ctx.Input.RequestBody, &teacher); err == nil {
		err := models.UpdateTeacherById(&teacher)
		if err == nil {
			this.Ctx.Output.Body(models.SuccessJson())
		} else {
			this.Abort(models.ErrJson("invlaid user id"))
		}
	} else {
		this.Abort(models.ErrJson("invalid data format"))
	}
}

func (this *TeacherController) Delete() {
	id, err := strconv.Atoi(this.Ctx.Input.Param(":id"))
	if err != nil {
		this.Abort(models.ErrJson("invalid id"))
	}
	if err := models.DeleteTeacher(id); err == nil {
		this.Ctx.Output.Body(models.SuccessJson())
	} else {
		this.Abort(models.ErrJson("invalid id"))
	}
}
