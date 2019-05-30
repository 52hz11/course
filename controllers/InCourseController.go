package controllers

import (
	"course/models"

	"github.com/astaxie/beego"
	"github.com/bitly/go-simplejson"
)

type InCourseController struct {
	beego.Controller
}

func (this *InCourseController) Get() {
	sess, _ := models.GlobalSessions.SessionStart(this.Ctx.ResponseWriter, this.Ctx.Request)
	defer sess.SessionRelease(this.Ctx.ResponseWriter)
	course_id, err := this.GetInt("course_id")
	if err != nil {
		course_id = -1
	}
	student_id, err := this.GetInt("student_id")
	if err != nil {
		student_id = -1
	}
	records := models.QueryInCourse(course_id, student_id)
	bodyJSON := simplejson.New()
	bodyJSON.Set("status", "success")
	tmpMapArr := make([]interface{}, len(records))
	for i, r := range records {
		tmpMap := make(map[string]interface{})
		tmpMap["course_id"] = r.CourseId.Id
		tmpMap["student_id"] = r.StudentId.Id
		tmpMapArr[i] = tmpMap
	}
	bodyJSON.Set("data", tmpMapArr)
	body, _ := bodyJSON.Encode()
	this.Ctx.Output.Body(body)
}

func (this *InCourseController) Post() {
	sess, _ := models.GlobalSessions.SessionStart(this.Ctx.ResponseWriter, this.Ctx.Request)
	defer sess.SessionRelease(this.Ctx.ResponseWriter)
	if inputJSON, err := simplejson.NewJson(this.Ctx.Input.RequestBody); err == nil {
		course_id := inputJSON.Get("course_id").MustInt()
		student_id := inputJSON.Get("student_id").MustInt()
		course_key := inputJSON.Get("course_key").MustString()
		var incourse models.InCourse
		course, err := models.GetCourseById(course_id)
		if err != nil {
			this.Abort(models.ErrJson("course not exist"))
		}
		if course.CourseKey != course_key {
			this.Abort("invalid course key")
		}
		student, err := models.GetUserById(student_id)
		if err != nil {
			this.Abort(models.ErrJson("student not exist"))
		}
		if sess.Get("id") == nil || sess.Get("id").(int) != student.Id {
			this.Abort(models.ErrJson("login expired"))
		}
		incourse.CourseId = course
		incourse.StudentId = student
		_, err = models.AddInCourse(&incourse)
		if err != nil {
			this.Abort(models.ErrJson("this record already exist"))
		} else {
			this.Ctx.Output.Body(models.SuccessJson())
		}
	} else {
		this.Abort(models.ErrJson("invalid data format"))
	}
}

func (this *InCourseController) Delete() {
	sess, _ := models.GlobalSessions.SessionStart(this.Ctx.ResponseWriter, this.Ctx.Request)
	defer sess.SessionRelease(this.Ctx.ResponseWriter)
	course_id, err := this.GetInt("course_id")
	if err != nil {
		course_id = -1
	}
	student_id, err := this.GetInt("student_id")
	if err != nil {
		student_id = -1
	}
	records := models.QueryInCourse(course_id, student_id)
	if len(records) == 0 {
		this.Abort(models.ErrJson("invalid in-course record"))
	} else {
		if sess.Get("id") == nil || sess.Get("id").(int) != student_id {
			this.Abort(models.ErrJson("login expired"))
		}
		for _, r := range records {
			models.DeleteInCourse(r.Id)
		}
		this.Ctx.Output.Body(models.SuccessJson())
	}
}
