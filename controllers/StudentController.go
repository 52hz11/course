package controllers

import (
	"course/models"
	"encoding/json"
	"strconv"

	"github.com/astaxie/beego"
	"github.com/bitly/go-simplejson"
)

type StudentController struct {
	beego.Controller
}

func (this *StudentController) Get() {
	id, err := strconv.Atoi(this.Ctx.Input.Param(":id"))
	if err == nil {
		bodyJSON := simplejson.New()
		student, err := models.GetStudentById(id)
		if err == nil {
			bodyJSON.Set("status", "success")
			dataMap := make(map[string]interface{})
			dataMap["id"] = id
			dataMap["name"] = student.Name
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

func (this *StudentController) Post() {
	var student models.Student
	if err := json.Unmarshal(this.Ctx.Input.RequestBody, &student); err == nil {
		id, err := models.AddStudent(&student)
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

func (this *StudentController) Put() {
	var student models.Student
	if err := json.Unmarshal(this.Ctx.Input.RequestBody, &student); err == nil {
		err := models.UpdateStudentById(&student)
		if err == nil {
			this.Ctx.Output.Body(models.SuccessJson())
		} else {
			this.Abort(models.ErrJson("invlaid user id"))
		}
	} else {
		this.Abort(models.ErrJson("invalid data format"))
	}
}

func (this *StudentController) Delete() {
	id, err := strconv.Atoi(this.Ctx.Input.Param(":id"))
	if err != nil {
		this.Abort(models.ErrJson("invalid id"))
	}
	if err := models.DeleteStudent(id); err == nil {
		this.Ctx.Output.Body(models.SuccessJson())
	} else {
		this.Abort(models.ErrJson("invalid id"))
	}
}
