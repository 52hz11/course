package controllers

import (
	"course/models"
	"time"

	"github.com/astaxie/beego"
	"github.com/bitly/go-simplejson"
)

type RollController struct {
	beego.Controller
}

func (this *RollController) Get() {
	id, err := this.GetInt("id")
	if err != nil {
		id = -1
	}
	course_id, err := this.GetInt("course_id")
	if err != nil {
		course_id = -1
	}
	title := this.GetString("title")
	begin_time := this.GetString("begin_time")
	end_time := this.GetString("end_time")
	rolls := models.QueryRoll(id, course_id, title, begin_time, end_time)
	bodyJSON := simplejson.New()
	bodyJSON.Set("status", "success")
	tmpMapArr := make([]interface{}, len(rolls))
	for i, r := range rolls {
		tmpMap := make(map[string]interface{})
		tmpMap["id"] = r.Id
		tmpMap["course_id"] = r.CourseId.Id
		tmpMap["title"] = r.Title
		tmpMap["begin_time"] = r.BeginTime.Format("2006-01-02 15:04:05")
		tmpMap["end_time"] = r.EndTime.Format("2006-01-02 15:04:05")
		tmpMapArr[i] = tmpMap
	}
	bodyJSON.Set("data", tmpMapArr)
	body, _ := bodyJSON.Encode()
	this.Ctx.Output.Body(body)
}

func (this *RollController) Post() {
	sess, _ := models.GlobalSessions.SessionStart(this.Ctx.ResponseWriter, this.Ctx.Request)
	defer sess.SessionRelease(this.Ctx.ResponseWriter)
	if inputJSON, err := simplejson.NewJson(this.Ctx.Input.RequestBody); err == nil {
		//fmt.Println("111")
		course_id := inputJSON.Get("course_id").MustInt()
		title := inputJSON.Get("title").MustString()
		begin_time_str := inputJSON.Get("begin_time").MustString()
		end_time_str := inputJSON.Get("end_time").MustString()
		loc, _ := time.LoadLocation("UTC")
		begin_time, err := time.ParseInLocation("2006-01-02 15:04:05", begin_time_str, loc)
		if err != nil {
			this.Abort(models.ErrJson("invalid data format"))
		}
		end_time, err := time.ParseInLocation("2006-01-02 15:04:05", end_time_str, loc)
		if err != nil {
			this.Abort(models.ErrJson("invalid data format"))
		}
		var roll models.Roll
		course, err := models.GetCourseById(course_id)
		if err != nil {
			this.Abort(models.ErrJson("invalid course id"))
		}
		if sess.Get("id") == nil || sess.Get("id").(int) != course.CreatorId.Id {
			this.Abort(models.ErrJson("login expired"))
		}
		roll.CourseId = course
		roll.Title = title
		roll.BeginTime = begin_time
		roll.EndTime = end_time
		id, err := models.AddRoll(&roll)
		if err != nil {
			this.Abort(models.ErrJson("add roll failed, database error"))
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

func (this *RollController) Put() {
	sess, _ := models.GlobalSessions.SessionStart(this.Ctx.ResponseWriter, this.Ctx.Request)
	defer sess.SessionRelease(this.Ctx.ResponseWriter)
	if inputJSON, err := simplejson.NewJson(this.Ctx.Input.RequestBody); err == nil {
		id := inputJSON.Get("id").MustInt()
		roll, err := models.GetRollById(id)
		if err != nil {
			this.Abort(models.ErrJson("invalid roll id"))
		}
		roll.CourseId, _ = models.GetCourseById(roll.CourseId.Id)
		if sess.Get("id") == nil || sess.Get("id").(int) != roll.CourseId.CreatorId.Id {
			this.Abort(models.ErrJson("login expired"))
		}
		title := inputJSON.Get("title").MustString()
		begin_time_str := inputJSON.Get("begin_time").MustString()
		end_time_str := inputJSON.Get("end_time").MustString()
		loc, _ := time.LoadLocation("UTC")
		begin_time, err := time.ParseInLocation("2006-01-02 15:04:05", begin_time_str, loc)
		if err != nil {
			this.Abort(models.ErrJson("invalid data format"))
		}
		end_time, err := time.ParseInLocation("2006-01-02 15:04:05", end_time_str, loc)
		if err != nil {
			this.Abort(models.ErrJson("invalid data format"))
		}
		roll.Title = title
		roll.BeginTime = begin_time
		roll.EndTime = end_time
		err = models.UpdateRollById(roll)
		if err != nil {
			this.Abort(models.ErrJson("fail to update roll, database error"))
		} else {
			this.Ctx.Output.Body(models.SuccessJson())
		}
	}
}

func (this *RollController) Delete() {
	sess, _ := models.GlobalSessions.SessionStart(this.Ctx.ResponseWriter, this.Ctx.Request)
	defer sess.SessionRelease(this.Ctx.ResponseWriter)
	id, err := this.GetInt("id")
	if err != nil {
		this.Abort(models.ErrJson("invalid id"))
	}
	roll, err := models.GetRollById(id)
	if err != nil {
		this.Abort(models.ErrJson("roll doesn't exist"))
	}
	roll.CourseId, _ = models.GetCourseById(roll.CourseId.Id)
	if sess.Get("id") == nil || sess.Get("id").(int) != roll.CourseId.CreatorId.Id {
		this.Abort(models.ErrJson("login expired"))
	}
	err = models.DeleteRoll(id)
	if err != nil {
		this.Abort(models.ErrJson("roll doesn't exist"))
	}
	this.Ctx.Output.Body(models.SuccessJson())
}
