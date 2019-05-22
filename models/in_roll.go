package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type InRoll struct {
	Id        int       `orm:"column(id);auto"`
	RollId    *Roll     `orm:"column(roll_id);rel(fk)"`
	StudentId *Student  `orm:"column(student_id);rel(fk)"`
	Time      time.Time `orm:"column(time);type(datetime);null"`
}

func (t *InRoll) TableName() string {
	return "in_roll"
}

func init() {
	orm.RegisterModel(new(InRoll))
}

// AddInRoll insert a new InRoll into database and returns
// last inserted Id on success.
func AddInRoll(m *InRoll) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetInRollById retrieves InRoll by Id. Returns error if
// Id doesn't exist
func GetInRollById(id int) (v *InRoll, err error) {
	o := orm.NewOrm()
	v = &InRoll{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllInRoll retrieves all InRoll matches certain condition. Returns empty list if
// no records exist
func GetAllInRoll(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(InRoll))
	// query k=v
	for k, v := range query {
		// rewrite dot-notation to Object__Attribute
		k = strings.Replace(k, ".", "__", -1)
		if strings.Contains(k, "isnull") {
			qs = qs.Filter(k, (v == "true" || v == "1"))
		} else {
			qs = qs.Filter(k, v)
		}
	}
	// order by:
	var sortFields []string
	if len(sortby) != 0 {
		if len(sortby) == len(order) {
			// 1) for each sort field, there is an associated order
			for i, v := range sortby {
				orderby := ""
				if order[i] == "desc" {
					orderby = "-" + v
				} else if order[i] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
			qs = qs.OrderBy(sortFields...)
		} else if len(sortby) != len(order) && len(order) == 1 {
			// 2) there is exactly one order, all the sorted fields will be sorted by this order
			for _, v := range sortby {
				orderby := ""
				if order[0] == "desc" {
					orderby = "-" + v
				} else if order[0] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
		} else if len(sortby) != len(order) && len(order) != 1 {
			return nil, errors.New("Error: 'sortby', 'order' sizes mismatch or 'order' size is not 1")
		}
	} else {
		if len(order) != 0 {
			return nil, errors.New("Error: unused 'order' fields")
		}
	}

	var l []InRoll
	qs = qs.OrderBy(sortFields...)
	if _, err = qs.Limit(limit, offset).All(&l, fields...); err == nil {
		if len(fields) == 0 {
			for _, v := range l {
				ml = append(ml, v)
			}
		} else {
			// trim unused fields
			for _, v := range l {
				m := make(map[string]interface{})
				val := reflect.ValueOf(v)
				for _, fname := range fields {
					m[fname] = val.FieldByName(fname).Interface()
				}
				ml = append(ml, m)
			}
		}
		return ml, nil
	}
	return nil, err
}

// UpdateInRoll updates InRoll by Id and returns error if
// the record to be updated doesn't exist
func UpdateInRollById(m *InRoll) (err error) {
	o := orm.NewOrm()
	v := InRoll{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteInRoll deletes InRoll by Id and returns error if
// the record to be deleted doesn't exist
func DeleteInRoll(id int) (err error) {
	o := orm.NewOrm()
	v := InRoll{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&InRoll{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}

func QueryInRoll(roll_id int, student_id int, time_str string) []InRoll {
	var records []InRoll
	o := orm.NewOrm()
	qs := o.QueryTable(new(InRoll))
	cond := orm.NewCondition()
	if roll_id != -1 {
		roll, err := GetRollById(roll_id)
		if err != nil {
			return records
		}
		cond = cond.And("RollId", roll)
	}
	if student_id != -1 {
		student, err := GetStudentById(student_id)
		if err != nil {
			return records
		}
		cond = cond.And("StudentId", student)
	}
	loc, _ := time.LoadLocation("UTC")
	time, err := time.ParseInLocation("2006-01-02 15:04:05", time_str, loc)
	if err == nil {
		cond = cond.And("Time", time)
	}
	qs.SetCond(cond).All(&records)
	return records
}
