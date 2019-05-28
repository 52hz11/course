// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"course/controllers"

	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/user/?:id", &controllers.UserController{})
	beego.Router("/course", &controllers.CourseController{})
	beego.Router("/in_course", &controllers.InCourseController{})
	beego.Router("/roll", &controllers.RollController{})
	beego.Router("/in_roll", &controllers.InRollController{})
	beego.Router("homework", &controllers.HomeworkController{})
	beego.Router("/ppt_file", &controllers.PPTFileController{})
	beego.Router("/charge_course", &controllers.ChargeCourseController{})
}
