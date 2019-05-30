package controllers

import (
	"course/models"
	"fmt"

	"github.com/astaxie/beego"
	"github.com/bitly/go-simplejson"
)

type CourseController struct {
	beego.Controller
}

func (this *CourseController) Get() {
	sess, _ := models.GlobalSessions.SessionStart(this.Ctx.ResponseWriter, this.Ctx.Request)
	defer sess.SessionRelease(this.Ctx.ResponseWriter)
	fmt.Println(sess.Get("id").(int))
	method := this.GetString("method")
	if method == "data" {
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
			if sess.Get("id") == nil || sess.Get("id").(int) != c.CreatorId.Id {
				tmpMap["course_key"] = c.CourseKey
			}
			tmpMap["name"] = c.Name
			tmpMap["content"] = c.Content
			tmpMapArr[i] = tmpMap
		}
		bodyJSON.Set("data", tmpMapArr)
		body, _ := bodyJSON.Encode()
		this.Ctx.Output.Body(body)
	} else if method == "head" {
		id, err := this.GetInt("id")
		if err != nil {
			this.Abort(models.ErrJson("invalid id"))
		}
		course, err := models.GetCourseById(id)
		if err != nil {
			this.Abort(models.ErrJson("invalid id"))
		}
		this.Ctx.Output.Download(course.ImgPath)
	}

}

func (this *CourseController) Post() {
	sess, _ := models.GlobalSessions.SessionStart(this.Ctx.ResponseWriter, this.Ctx.Request)
	defer sess.SessionRelease(this.Ctx.ResponseWriter)
	method := this.GetString("method")
	if method == "data" {
		if inputJSON, err := simplejson.NewJson(this.Ctx.Input.RequestBody); err == nil {
			var course models.Course
			creator_id := inputJSON.Get("creator_id").MustInt()
			if sess.Get("id") == nil || sess.Get("id").(int) != creator_id {
				this.Abort(models.ErrJson("login expired"))
			}
			course.CreatorId, _ = models.GetUserById(creator_id)
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
	} else if method == "head" {
		id, err := this.GetInt("id")
		if err != nil {
			this.Abort(models.ErrJson("must provide a course id"))
		}
		course, err := models.GetCourseById(id)
		if err != nil {
			this.Abort(models.ErrJson("invalid course id"))
		}
		if sess.Get("id") == nil || sess.Get("id").(int) != course.CreatorId.Id {
			this.Abort(models.ErrJson("login expired"))
		}
		file, head, err := this.GetFile("file")
		defer file.Close()
		if err != nil {
			this.Abort(models.ErrJson("error when trying to get file"))
		}
		course.ImgPath = "./upload/" + models.GenerateKey() + "__" + head.Filename
		err = models.UpdateCourseById(course)
		if err != nil {
			this.Abort(models.ErrJson("update head image failed, database error"))
		}
		this.SaveToFile("file", course.ImgPath)
		this.Ctx.Output.Body(models.SuccessJson())
	} else {
		this.Abort(models.ErrJson("must provide a method"))
	}

}

func (this *CourseController) Put() {
	sess, _ := models.GlobalSessions.SessionStart(this.Ctx.ResponseWriter, this.Ctx.Request)
	defer sess.SessionRelease(this.Ctx.ResponseWriter)
	if inputJSON, err := simplejson.NewJson(this.Ctx.Input.RequestBody); err == nil {
		var course *models.Course
		id := inputJSON.Get("id").MustInt()
		course, err := models.GetCourseById(id)
		if err != nil {
			this.Abort(models.ErrJson("course not exist"))
		}
		if sess.Get("id") == nil || sess.Get("id").(int) != course.CreatorId.Id {
			this.Abort(models.ErrJson("login expired"))
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
	sess, _ := models.GlobalSessions.SessionStart(this.Ctx.ResponseWriter, this.Ctx.Request)
	defer sess.SessionRelease(this.Ctx.ResponseWriter)
	id, err := this.GetInt("id")
	if err != nil {
		this.Abort(models.ErrJson("invalid id"))
	} else {
		course, err := models.GetCourseById(id)
		if err != nil {
			this.Abort(models.ErrJson("course not exist"))
		}
		if sess.Get("id") == nil || sess.Get("id").(int) != course.CreatorId.Id {
			this.Abort(models.ErrJson("login expired"))
		}
		err = models.DeleteCourse(id)
		if err != nil {
			this.Abort(models.ErrJson("course not exist"))
		}
		this.Ctx.Output.Body(models.SuccessJson())
	}
}
