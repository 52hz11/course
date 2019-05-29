注：对于所有的错误信息，格式如下：
```
{
    "status" : "failed",
    "msg" : "错误信息"
}
```
对于所有无需返回id的api，返回的正确信息格式如下：
```
{
    "status" : "success"
}
```
对于所有创建类而且需要返回id的api，运行正确返回信息格式如下：
```
{
    "status" : "success",
    "id" : id
}
```
对于一些查询类api，返回的格式如下
```
{
    "status" : "success",
    "data" : [{数组内容}]
}
```
对于认证失败的返回格式如下：
```
{
    "status" : "failed",
    "msg" : "login expired"
}
```
数组具体内容可以参见数据库结构，所有的字段名均与数据库表字段名统一，并且注意请求体中提供和url参数提供的区别，请求体统一为json格式

api全部遵循RESTFul规范，GET方法获取信息，POST方法上传信息，PUT修改信息，DELETE删除信息

因为有时候有时间信息需要写在url中，可能包含空格，所以在发送前要进行转义，空格转义后为%20
### /user/?:id
本路由支持GET,POST,PUT,DELETE类型请求
#### GET
如果是GET请求，则必需提供查询参数method，如果method为id，则为按id获取用户信息，必须提供?:id部分，即/student/1的格式来访问，否则报错，返回的data包括除了token外的全部信息，如果method为token，则为按token获取用户信息，不需要提供?:id，返回的data包括用户的全部信息。
#### POST
如果是POST请求，则请求体格式如下：
```
{
    "id" : id
    "name" : "name",
    ···
}
```
#### PUT
供修改用户信息用，和POST类似，但是请求体中需要提供用户id，而且返回的json中不会包括id，只有当前用户可以修改当前用户账号
#### DELETE
必须在url中提供参数id，返回内容没有额外信息,只有当前用户可以删除当前用户账号

### /course
本路由支持GET,POST,PUT,DELETE类型请求
#### GET
供查询course使用，可选参数包括id,name,content,creator_id,offset,limit, 直接附在url尾部参数部分即可，返回的data中包括course表中所有的信息（也是唯一可以获取course_key的接口），只有确认目前用户是创建者时才会返回course_key
#### POST
供创建course使用，必须提供method，如果method为data,则请求体中需要提供除了id和course_key以外的所有course信息，创建成功会返回id，如果method为head，则为上传头像，需要在参数中同时提供id，也就是/course?method=head&id=1类似的地址，表单文件栏的的name必须设为file，如下：
```xml
<input type="file" name="file">
```
只有创建者可以上传（修改）头像
#### PUT
和POST的method=data时类似（头像上传依旧用POST方法，这里也不需要再提供method了），但是请求体中需要提供id，不会返回额外信息
#### DELETE
把要删除的id放在参数列表中，例如/course?id=1，无额外返回信息，只有创建者可以删除课程

### /in_course
本路由支持GET,POST,DELETE类型请求
#### GET
get请求可供选择的url参数有course_id和student_id，返回的data也只包括course_id和student_id
#### POST
在请求体中提供course_id和student_id，不会返回id（因为这个id没有意义），只有当前用户可以加入课程
#### DELETE
在url可选参数有course_id和student_id，如果没有任何信息被删除会返回错误信息，只有当前用户可以退出课程

### /roll
本路由支持GET,POST,PUT,DELETE类型请求
#### GET
支持的url参数有id,course_id,title,begin_time,end_time, 其中begin_time和end_time是一个区间，返回的data的签到开放时间会在这个区间内，返回的信息包括roll表中的所有信息
#### POST
除了id外需要提供其他全部信息，创建成功会返回id，这里请求体中提供的course_id必须是当前用户创建的课程
#### PUT
需要提供除了course_id外的所有信息，也就是不可以更改所属课程，其他非id信息都可以更改，只有course_id为当前用户创建的课程时允许修改
#### DELETE
在url参数中提供id，无额外返回信息，只有course_id为当前用户创建的课程时允许删除

### /in_roll
本路由支持GET,POST类型请求
#### GET
get请求可供选择的url参数有roll_id和student_id和time，返回的data包括roll_id和student_id和time
#### POST
在请求体中提供course_id、student_id和time，不会返回id（因为这个id没有意义），只有roll所属课程的拥有者有此操作权限


### /homework
本路由支持GET,POST,PUT,DELETE类型请求
#### GET
支持的url参数有id,course_id,title,content, 返回的data包括所有homework表中的信息
#### POST
请求体中需要提供除了id外其他所有信息，成功则会返回id，这里请求体中提供的course_id必须是当前用户创建的课程
#### PUT
和POST类似，但是需要提供id，不需要提供course_id（同时也不允许修改course_id）,无额外返回信息，只有course_id为当前用户创建的课程时允许修改
#### DELETE
在url参数中提供id即可，无额外返回信息，只有course_id为当前用户创建的课程时允许删除

### /ppt_file
支持GET,POST,DELETE方法
#### GET
必须在url参数中包含method，method可选list和getfile

如果method为list，则为获取文件列表，可供选择的url查询参数有id,name,course_id，返回的data包括除了file_path外其他信息（file_path不暴露给前端）

如果method为getfile，则只需在url参数中提供id即可，会返回文件（注意文件名可能和上传的不一样，为了保证文件名不重复加了一些乱码进去，在前端可能要重命名一下）
#### POST
直接在url列表中提供course_key, body中传文件即可(和上面上传头像一样)，成功会返回id
#### DELETE
在url参数中提供id即可，无额外返回信息，只有course_id为当前用户创建的课程时允许删除

### /charge_course
本路由支持GET,POST,DELETE类型请求
#### GET
get请求可供选择的url参数有course_id和ta_id，返回的data也只包括course_id和ta_id
#### POST
在请求体中提供course_id和ta_id，不会返回id（因为这个id没有意义），只有course_id为当前用户创建的课程时允许操作
#### DELETE
在url可选参数有course_id和ta_id，如果没有任何信息被删除会返回错误信息，只有course_id为当前用户创建的课程时允许删除
