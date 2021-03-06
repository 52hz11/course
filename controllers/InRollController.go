package controllers

import (
	"course/models"
	"time"

	"github.com/astaxie/beego"
	"github.com/bitly/go-simplejson"
	"github.com/tealeg/xlsx"
)

type InRollController struct {
	beego.Controller
}

func (this *InRollController) Get() {
	sess, _ := models.GlobalSessions.SessionStart(this.Ctx.ResponseWriter, this.Ctx.Request)
	defer sess.SessionRelease(this.Ctx.ResponseWriter)
	roll_id, err := this.GetInt("roll_id")
	if err != nil {
		roll_id = -1
	}
	student_id, err := this.GetInt("student_id")
	if err != nil {
		student_id = -1
	}
	time := this.GetString("time")
	file := xlsx.NewFile()
	sheet, _ := file.AddSheet("Sheet1")
	records := models.QueryInRoll(roll_id, student_id, time)
	for _, r := range records {
		row := sheet.AddRow()
		cell := row.AddCell()
		r.StudentId, _ = models.GetUserById(r.StudentId.Id)
		cell.Value = r.StudentId.Name
		cell = row.AddCell()
		cell.Value = r.Time.Format("2006-01-02 15:04:05")
	}
	file.Save("tmp.xlsx")
	/*
		bodyJSON := simplejson.New()
		bodyJSON.Set("status", "success")
		tmpMapArr := make([]interface{}, len(records))
		for i, r := range records {
			tmpMap := make(map[string]interface{})
			tmpMap["roll_id"] = r.RollId.Id
			tmpMap["student_id"] = r.StudentId.Id
			tmpMap["time"] = r.Time.Format("2006-01-02 15:04:05")
			tmpMapArr[i] = tmpMap
		}
		bodyJSON.Set("data", tmpMapArr)
		body, _ := bodyJSON.Encode()
		this.Ctx.Output.Body(body)
	*/
	this.Ctx.Output.Download("tmp.xlsx")
}

func (this *InRollController) Post() {
	sess, _ := models.GlobalSessions.SessionStart(this.Ctx.ResponseWriter, this.Ctx.Request)
	defer sess.SessionRelease(this.Ctx.ResponseWriter)
	if inputJSON, err := simplejson.NewJson(this.Ctx.Input.RequestBody); err == nil {
		roll_id := inputJSON.Get("roll_id").MustInt()
		student_id := inputJSON.Get("student_id").MustInt()
		time_str := inputJSON.Get("time").MustString()
		//这里的local是因为明明是按照UTC存的，但是取出来后又会变成当前时区，怀疑是orm对这方面有所遗漏，所以
		//比较的时候用本地时区，存的时候再转为UTC
		loc, _ := time.LoadLocation("Local")
		Time, err := time.ParseInLocation("2006-01-02 15:04:05", time_str, loc)
		var inroll models.InRoll
		roll, err := models.GetRollById(roll_id)
		if err != nil {
			this.Abort(models.ErrJson("invalid roll id"))
		}
		roll.CourseId, _ = models.GetCourseById(roll.CourseId.Id)
		if sess.Get("id") == nil || sess.Get("id").(int) != roll.CourseId.CreatorId.Id {
			this.Abort(models.ErrJson("login expired"))
		}
		student, err := models.GetUserById(student_id)
		if err != nil {
			this.Abort(models.ErrJson("invalid student id"))
		}
		if !(Time.After(roll.BeginTime) && Time.Before(roll.EndTime)) {
			this.Abort(models.ErrJson("the roll is not open now"))
		}
		inroll.RollId = roll
		inroll.StudentId = student
		loc, _ = time.LoadLocation("UTC")
		Time, _ = time.ParseInLocation("2006-01-02 15:04:05", time_str, loc)
		inroll.Time = Time
		_, err = models.AddInRoll(&inroll)
		if err != nil {
			this.Abort(models.ErrJson("add in roll record failed, database error"))
		}
		this.Ctx.Output.Body(models.SuccessJson())
	} else {
		this.Abort(models.ErrJson("invalid data format"))
	}
}
