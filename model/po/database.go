package po

/**
create table t_database (
	id int(11) primary key,
    name varchar(100) unique not null
);

create table t_table (
	id int(11) primary key,
    name varchar(100) not null,
	database_id int(11)
);

create table t_table_col (
	id int(11) primary key,
    name varchar(100) not null,
    table_id int(11) not null,
	index_id int(11) not null
);

create table t_table_index (
	id int(11) primary key,
    name varchar(100) not null,
    table_id int(11) not null
);

create table
*/

type DataBase struct {
}
