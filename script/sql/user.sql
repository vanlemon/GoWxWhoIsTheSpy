-- drop database if exists user_db;
-- create database user_db;
use user_db;
drop table if exists `user`;

create table `user`(
    `id` int (64) not null auto_increment,
    `openid` varchar (64) not null,
    `nick_name` varchar (32),
    `gender` int (2) comment '性别 0：未知、1：男、2：女',
    `avatar_url`  varchar (512),
    `city`  varchar (16),
    `province`  varchar (16),
    `country`  varchar (16),
    `language`  varchar (16),
    `create_time` timestamp not null DEFAULT CURRENT_TIMESTAMP,
    `update_time` timestamp not null DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    primary key (`id`),
    unique key (`openid`)
) ENGINE=InnoDB default charset=utf8mb4;