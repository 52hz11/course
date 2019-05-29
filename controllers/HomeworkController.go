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
		tmpMap["deadline"] = h.Deadline
		tmpMapArr[i] = tmpMap
	}
	bodyJSON.Set("data", tmpMapArr)
	body, _ := bodyJSON.Encode()
	this.Ctx.Output.Body(body)
}

func (this *HomeworkController) Post() {
	sess, _ := models.GlobalSessions.SessionStart(this.Ctx.ResponseWriter, this.Ctx.Request)
	defer sess.SessionRelease(this.Ctx.ResponseWriter)
	if inputJSON, err := simplejson.NewJson(this.Ctx.Input.RequestBody); err == nil {
		course_id := inputJSON.Get("course_id").MustInt()
		title := inputJSON.Get("title").MustString()
		content := inputJSON.Get("content").MustString()
		deadline := inputJSON.Get("deadline").MustString()
		var homework models.Homework
		course, err := models.GetCourseById(course_id)
		if err != nil {
			this.Abort(models.ErrJson("invalid course id"))
		}
		if sess.Get("id") == nil || sess.Get("id").(int) != course.CreatorId.Id {
			this.Abort(models.ErrJson("login expired"))
		}
		homework.CourseId = course
		homework.Title = title
		homework.Content = content
		homework.Deadline = deadline
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
	sess, _ := models.GlobalSessions.SessionStart(this.Ctx.ResponseWriter, this.Ctx.Request)
	defer sess.SessionRelease(this.Ctx.ResponseWriter)
	if inputJSON, err := simplejson.NewJson(this.Ctx.Input.RequestBody); err == nil {
		id := inputJSON.Get("id").MustInt()
		title := inputJSON.Get("title").MustString()
		content := inputJSON.Get("content").MustString()
		deadline := inputJSON.Get("deadline").MustString()
		homework, err := models.GetHomeworkById(id)
		if err != nil {
			this.Abort(models.ErrJson("homework not exist"))
		}
		homework.CourseId, _ = models.GetCourseById(homework.CourseId.Id)
		if sess.Get("id") == nil || sess.Get("id").(int) != homework.CourseId.CreatorId.Id {
			this.Abort(models.ErrJson("login expired"))
		}
		homework.Title = title
		homework.Content = content
		homework.Deadline = deadline
		err = models.UpdateHomeworkById(homework)
		if err != nil {
			this.Abort(models.ErrJson("update homework failed, database error"))
		}
		this.Ctx.Output.Body(models.SuccessJson())
	}
}

func (this *HomeworkController) Delete() {
	sess, _ := models.GlobalSessions.SessionStart(this.Ctx.ResponseWriter, this.Ctx.Request)
	defer sess.SessionRelease(this.Ctx.ResponseWriter)
	id, err := this.GetInt("id")
	if err != nil {
		this.Abort(models.ErrJson("homework not exist"))
	}
	homework, err := models.GetHomeworkById(id)
	if err != nil {
		this.Abort(models.ErrJson("homework not exist"))
	}
	homework.CourseId, _ = models.GetCourseById(homework.CourseId.Id)
	if sess.Get("id") == nil || sess.Get("id").(int) != homework.CourseId.CreatorId.Id {
		this.Abort(models.ErrJson("login expired"))
	}
	err = models.DeleteHomework(id)
	if err != nil {
		this.Abort(models.ErrJson("delete homework failed, databse error"))
	}
	this.Ctx.Output.Body(models.SuccessJson())
}
