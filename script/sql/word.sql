# drop database if exists user_db;
# create database user_db;
use user_db;
drop table if exists `word`;

create table `word`(
    `id` int (64) not null auto_increment,
    `normal_word` varchar (8) not null comment '平民词汇',
    `spy_word` varchar (8) not null comment '卧底词汇',
    `blank_word` varchar (8) not null comment '空白词汇',
    `class` varchar (8) not null comment '词汇类别：NotYet（未分类）',
    primary key (`id`),
    unique key (`normal_word`, `spy_word`)
) ENGINE=InnoDB default charset=utf8mb4;

insert into `word` (`normal_word`, `spy_word`, `blank_word`, `class`) values
('内裤', '内衣', '', 'NotYet'),
('浴缸', '鱼缸', '', 'NotYet'),
('筷子', '牙签', '', 'NotYet'),
('饺子', '包子', '', 'NotYet'),
('眉毛', '睫毛', '', 'NotYet');

insert into `word` (`normal_word`, `spy_word`, `blank_word`, `class`) values
('内衣', '内裤', '', 'NotYet'),
('鱼缸', '浴缸', '', 'NotYet'),
('牙签', '筷子', '', 'NotYet'),
('包子', '饺子', '', 'NotYet'),
('睫毛', '眉毛', '', 'NotYet');

insert into `word` (`normal_word`, `spy_word`, `blank_word`, `class`) values
('丑小鸭', '灰姑娘', '', 'NotYet'),
('灰姑娘', '丑小鸭', '', 'NotYet'),
('男/女朋友', '前男/女友', '', 'NotYet'),
('前/女男友', '男/女朋友', '', 'NotYet'),
('坐在你左边的人', '坐在你右边的人', '', 'NotYet'),
('坐在你右边的人', '坐在你左边的人', '', 'NotYet'),
('甄嬛传', '红楼梦', '', 'NotYet'),
('红楼梦', '甄嬛传', '', 'NotYet');
