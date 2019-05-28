create table user(
	id int unsigned primary key,
	name varchar(40),
	number int unsigned,
	token varchar(40),
	email varchar(40),
    school varchar(40),
	type int unsigned default 0
);
create table course(
	id int unsigned auto_increment primary key,
	name varchar(40),
	content varchar(200),
	creator_id int unsigned,
	course_key varchar(40),
	img_path varchar(40),
	foreign key(creator_id) references user(id)
);

create table in_course(
	id int unsigned auto_increment primary key,
	course_id int unsigned,
	student_id int unsigned,
	foreign key(course_id) references course(id),
	foreign key(student_id) references user(id)
);
create table charge_course(
	id int unsigned auto_increment primary key,
	course_id int unsigned,
	ta_id int unsigned,
	foreign key(course_id) references course(id),
	foreign key(ta_id) references user(id)
);
create table ppt_file(
	id int unsigned auto_increment primary key,
	course_id int unsigned,
	name varchar(40),
	file_path varchar(100),
	foreign key(course_id) references course(id)
);

create table homework(
	id int unsigned auto_increment primary key,
	course_id int unsigned,
	title varchar(40),
	content varchar(200),
	foreign key(course_id) references course(id)
);

create table roll(
	id int unsigned auto_increment primary key,
	course_id int unsigned,
	title varchar(40),
	begin_time datetime,
	end_time datetime,
	foreign key(course_id) references course(id)
);

create table in_roll(
	id int unsigned auto_increment primary key,
	roll_id int unsigned,
	student_id int unsigned,
	time datetime,
	foreign key(roll_id) references roll(id),
	foreign key(student_id) references user(id)
);