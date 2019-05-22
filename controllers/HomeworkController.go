package controllers

import (
	"course/models"

	"github.com/astaxie/beego"
	"github.com/bitly/go-simplejson"
)

type HomeworkController struct {
	beego.Controller
}

func (this *HomeworkController) Get() {
	id, err := this.GetInt("id")
	if err != nil {
		id = -1
	}
	course_id, err := this.GetInt("course_id")
	if err != nil {
		course_id = -1
	}
	title := this.GetString("title")
	content := this.GetString("content")
	homeworks := models.QueryHomework(id, course_id, title, content)
	bodyJSON := simplejson.New()
	bodyJSON.Set("status", "success")
	tmpMapArr := make([]interface{}, len(homeworks))
	for i, h := range homeworks {
		tmpMap := make(map[string]interface{})
		tmpMap["id"] = h.Id
		tmpMap["course_id"] = h.CourseId.Id
		tmpMap["title"] = h.Title
		tmpMap["content"] = h.Content
		tmpMapArr[i] = tmpMap
	}
	bodyJSON.Set("data", tmpMapArr)
	body, _ := bodyJSON.Encode()
	this.Ctx.Output.Body(body)
}

func (this *HomeworkController) Post() {
	if inputJSON, err := simplejson.NewJson(this.Ctx.Input.RequestBody); err == nil {
		course_id := inputJSON.Get("course_id").MustInt()
		title := inputJSON.Get("title").MustString()
		content := inputJSON.Get("content").MustString()
		var homework models.Homework
		homework.CourseId, _ = models.GetCourseById(course_id)
		homework.Title = title
		homework.Content = content
		id, err := models.AddHomework(&homework)
		if err != nil {
			this.Abort(models.ErrJson("invalid data, check course id is valid"))
		}
		bodyJSON := simplejson.New()
		bodyJSON.Set("status", "success")
		bodyJSON.Set("id", id)
		body, _ := bodyJSON.Encode()
		this.Ctx.Output.Body(body)
	} else {
		this.Abort(models.ErrJson("invalid data format"))
	}
}

func (this *HomeworkController) Put() {
	if inputJSON, err := simplejson.NewJson(this.Ctx.Input.RequestBody); err == nil {
		id := inputJSON.Get("id").MustInt()
		title := inputJSON.Get("title").MustString()
		content := inputJSON.Get("content").MustString()
		homework, err := models.GetHomeworkById(id)
		if err != nil {
			this.Abort(models.ErrJson("homework not exist"))
		}
		homework.Title = title
		homework.Content = content
		err = models.UpdateHomeworkById(homework)
		if err != nil {
			this.Abort(models.ErrJson("update homework failed, database error"))
		}
		this.Ctx.Output.Body(models.SuccessJson())
	}
}

func (this *HomeworkController) Delete() {
	id, err := this.GetInt("id")
	if err != nil {
		this.Abort(models.ErrJson("homework not exist"))
	}
	err = models.DeleteHomework(id)
	if err != nil {
		this.Abort(models.ErrJson("delete homework failed, databse error"))
	}
	this.Ctx.Output.Body(models.SuccessJson())
}
