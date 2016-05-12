create database ryo;

	use ryo;

	create table UserInfo(

		UserID int primary key AUTO_INCREMENT,
		UserAccount varchar(16) not null,
		UserPwd varchar(16) not null

	)

	create table BookInfo(
		TableID int primary key AUTO_INCREMENT,
		BookCode varchar(32) not null,
		BookName varchar(32),
		UserID int
	)
