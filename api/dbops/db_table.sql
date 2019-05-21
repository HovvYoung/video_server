CREATE TABLE users
(
    id INT UNSIGNED PRIMARY KEY AUTO_INCREMENT,
    login_name VARCHAR(64) UNIQUE KEY,
	pwd TEXT        
);

CREATE TABLE video_info
(
	id VARCHAR(64) PRIMARY KEY NOT NULL,
	author_id INT UNSIGNED,
	name TEXT,
	display_ctime TEXT,
	create_time DATETIME
);

CREATE TABLE comments
(
	id VARCHAR(64) PRIMARY KEY NOT NULL,
	video_id VARCHAR(64),
	author_id INT UNSIGNED,
	content TEXT,
	time DATETIME
);

CREATE TABLE sessions
(
	session_id TINYTEXT NOT NULL,
	TTL TINYTEXT,
	login_name TEXT
);
alter table sessions add PRIMARY KEY (session_id(64));

CREATE TABLE video_del_rec
(
	video_id VARCHAR(64) NOT NULL PRIMARY KEY
);
