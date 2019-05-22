package controllers

import (
	"course/models"

	"github.com/astaxie/beego"
	"github.com/bitly/go-simplejson"
)

type CourseController struct {
	beego.Controller
}

func (this *CourseController) Get() {
	id, err := this.GetInt("id")
	if err != nil {
		id = -1
	}
	name := this.GetString("name")
	content := this.GetString("content")
	creator_id, err := this.GetInt("creator_id")
	if err != nil {
		creator_id = -1
	}
	offset, err := this.GetInt("offset")
	if err != nil {
		offset = 0
	}
	limit, err := this.GetInt("limit")
	if err != nil {
		limit = 100
	}
	courses := models.QueryCourse(id, name, content, creator_id, offset, limit)
	bodyJSON := simplejson.New()
	bodyJSON.Set("status", "success")
	tmpMapArr := make([]interface{}, len(courses))
	for i, c := range courses {
		tmpMap := make(map[string]interface{})
		tmpMap["id"] = c.Id
		tmpMap["creator_id"] = c.CreatorId.Id
		tmpMap["course_key"] = c.CourseKey
		tmpMap["name"] = c.Name
		tmpMap["content"] = c.Content
		tmpMapArr[i] = tmpMap
	}
	bodyJSON.Set("data", tmpMapArr)
	body, _ := bodyJSON.Encode()
	this.Ctx.Output.Body(body)
}

func (this *CourseController) Post() {
	if inputJSON, err := simplejson.NewJson(this.Ctx.Input.RequestBody); err == nil {
		var course models.Course
		creator_id := inputJSON.Get("creator_id").MustInt()
		course.CreatorId, _ = models.GetTeacherById(creator_id)
		course.CourseKey = models.GenerateKey()
		course.Name = inputJSON.Get("name").MustString()
		course.Content = inputJSON.Get("content").MustString()
		id, err := models.AddCourse(&course)
		if err == nil {
			bodyJSON := simplejson.New()
			bodyJSON.Set("status", "success")
			bodyJSON.Set("id", id)
			body, _ := bodyJSON.Encode()
			this.Ctx.Output.Body(body)
		} else {
			this.Abort(models.ErrJson("fail to add course, database error"))
		}
	} else {
		this.Abort(models.ErrJson("invalid data format"))
	}
}

func (this *CourseController) Put() {
	if inputJSON, err := simplejson.NewJson(this.Ctx.Input.RequestBody); err == nil {
		var course *models.Course
		id := inputJSON.Get("id").MustInt()
		course, err := models.GetCourseById(id)
		if err != nil {
			this.Abort(models.ErrJson("course not exist"))
		}
		course.Name = inputJSON.Get("name").MustString()
		course.Content = inputJSON.Get("content").MustString()
		err = models.UpdateCourseById(course)
		if err == nil {
			this.Ctx.Output.Body(models.SuccessJson())
		} else {
			this.Abort(models.ErrJson("fail to update course, database error"))
		}
	} else {
		this.Abort(models.ErrJson("invalid data format"))
	}
}

func (this *CourseController) Delete() {
	id, err := this.GetInt("id")
	if err != nil {
		this.Abort(models.ErrJson("invalid id"))
	} else {
		_, err := models.GetCourseById(id)
		if err != nil {
			this.Abort(models.ErrJson("course not exist"))
		}
		err = models.DeleteCourse(id)
		if err != nil {
			this.Abort(models.ErrJson("course not exist"))
		}
		this.Ctx.Output.Body(models.SuccessJson())
	}
}
