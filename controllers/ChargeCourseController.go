package controllers

import (
	"course/models"

	"github.com/astaxie/beego"
	"github.com/bitly/go-simplejson"
)

type ChargeCourseController struct {
	beego.Controller
}

func (this *ChargeCourseController) Get() {
	course_id, err := this.GetInt("course_id")
	if err != nil {
		course_id = -1
	}
	ta_id, err := this.GetInt("ta_id")
	if err != nil {
		ta_id = -1
	}
	records := models.QueryChargeCourse(course_id, ta_id)
	bodyJSON := simplejson.New()
	bodyJSON.Set("status", "success")
	tmpMapArr := make([]interface{}, len(records))
	for i, r := range records {
		tmpMap := make(map[string]interface{})
		tmpMap["course_id"] = r.CourseId.Id
		tmpMap["ta_id"] = r.StudentId.Id
		tmpMapArr[i] = tmpMap
	}
	bodyJSON.Set("data", tmpMapArr)
	body, _ := bodyJSON.Encode()
	this.Ctx.Output.Body(body)
}

func (this *ChargeCourseController) Post() {
	if inputJSON, err := simplejson.NewJson(this.Ctx.Input.RequestBody); err == nil {
		course_id := inputJSON.Get("course_id").MustInt()
		ta_id := inputJSON.Get("ta_id").MustInt()
		var chargecourse models.ChargeCourse
		course, err := models.GetCourseById(course_id)
		if err != nil {
			this.Abort(models.ErrJson("course not exist"))
		}
		ta, err := models.GetUserById(ta_id)
		if err != nil {
			this.Abort(models.ErrJson("student not exist"))
		}
		chargecourse.CourseId = course
		chargecourse.TaId = ta
		_, err = models.AddChargeCourse(&chargecourse)
		if err != nil {
			this.Abort(models.ErrJson("this record already exist"))
		} else {
			this.Ctx.Output.Body(models.SuccessJson())
		}
	} else {
		this.Abort(models.ErrJson("invalid data format"))
	}
}

func (this *ChargeCourseController) Delete() {
	course_id, err := this.GetInt("course_id")
	if err != nil {
		course_id = -1
	}
	ta_id, err := this.GetInt("ta_id")
	if err != nil {
		ta_id = -1
	}
	records := models.QueryChargeCourse(course_id, ta_id)
	if len(records) == 0 {
		this.Abort(models.ErrJson("invalid charge-course record"))
	} else {
		for _, r := range records {
			models.DeleteChargeCourse(r.Id)
		}
		this.Ctx.Output.Body(models.SuccessJson())
	}
}
