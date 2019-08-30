
DROP TABLE IF EXISTS `chat`;

CREATE TABLE `chat` (
  `chat_id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `chat_type` tinyint(1) unsigned NOT NULL DEFAULT '1' COMMENT '1:私聊 2:群聊',
  `create_time` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
  PRIMARY KEY (`chat_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='聊天索引表，新增数据时，需要同时创建chat_user表中用户的chat_list记录';


INSERT INTO `chat` (`chat_id`, `chat_type`, `create_time`)
VALUES
	(13,1,1567134019),
	(14,1,1567134216),
	(15,1,1567153001);


# Dump of table chat_list
# ------------------------------------------------------------

DROP TABLE IF EXISTS `chat_list`;

CREATE TABLE `chat_list` (
  `list_id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `chat_id` int(11) unsigned NOT NULL DEFAULT '0',
  `chat_type` tinyint(1) unsigned NOT NULL DEFAULT '1' COMMENT '1:私聊 2:群聊',
  `uid` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '所属用户uid',
  `relate_id` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '私聊类型->用户id 群聊类型->群组id',
  `new` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '新消息数量',
  `active_time` int(10) unsigned NOT NULL DEFAULT '0' COMMENT '最后活跃时间',
  `create_time` int(11) unsigned NOT NULL DEFAULT '0',
  PRIMARY KEY (`list_id`),
  KEY `uid` (`uid`,`chat_id`),
  KEY `active_time` (`active_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户的聊天对话列表';


INSERT INTO `chat_list` (`list_id`, `chat_id`, `chat_type`, `uid`, `relate_id`, `new`, `active_time`, `create_time`)
VALUES
	(21,13,1,4,12,0,1567134165,1567134019),
	(22,13,1,12,4,0,1567134165,1567134019),
	(23,14,1,10,4,0,1567146709,1567134216),
	(24,14,1,4,10,0,1567146709,1567134216),
	(25,15,1,4,11,0,1567153570,1567153001),
	(26,15,1,11,4,0,1567153570,1567153001);


# Dump of table chat_message
# ------------------------------------------------------------

DROP TABLE IF EXISTS `chat_message`;

CREATE TABLE `chat_message` (
  `msg_id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `chat_id` int(11) unsigned NOT NULL DEFAULT '0',
  `uid` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '发送人uid',
  `msg_type` tinyint(1) NOT NULL DEFAULT '1' COMMENT '1:文字 2:图片',
  `content` text NOT NULL COMMENT '消息内容',
  `create_time` int(11) NOT NULL,
  PRIMARY KEY (`msg_id`),
  KEY `chat_id` (`chat_id`),
  KEY `create_time` (`create_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;


INSERT INTO `chat_message` (`msg_id`, `chat_id`, `uid`, `msg_type`, `content`, `create_time`)
VALUES
	(82,13,4,1,'你好！',1567134029),
	(83,13,12,1,'hi',1567134036),
	(85,13,4,1,'我去',1567134152),
	(86,13,12,1,'？？？？',1567134165),
	(87,14,10,1,'加我干嘛？？？',1567134227),
	(88,14,4,1,'你妹',1567146709),
	(92,15,11,1,'哈哈哈哈',1567153565),
	(93,15,4,1,'哈哈哈哈',1567153570);


# Dump of table chat_user
# ------------------------------------------------------------

DROP TABLE IF EXISTS `chat_user`;

CREATE TABLE `chat_user` (
  `chat_id` int(11) unsigned NOT NULL DEFAULT '0',
  `uid` int(11) unsigned NOT NULL DEFAULT '0',
  UNIQUE KEY `uid` (`uid`,`chat_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='聊天成员';


INSERT INTO `chat_user` (`chat_id`, `uid`)
VALUES
	(13,4),
	(14,4),
	(15,4),
	(14,10),
	(15,11),
	(13,12);



# Dump of table user
# ------------------------------------------------------------

DROP TABLE IF EXISTS `user`;

CREATE TABLE `user` (
  `uid` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `open` tinyint(1) unsigned NOT NULL DEFAULT '1' COMMENT '1:公开',
  `nickname` varchar(20) NOT NULL DEFAULT '' COMMENT '昵称',
  `email` varchar(100) NOT NULL DEFAULT '' COMMENT '邮箱',
  `password` char(32) NOT NULL DEFAULT '' COMMENT '密码',
  `salt` char(6) NOT NULL DEFAULT '',
  `avatar` varchar(255) NOT NULL DEFAULT '' COMMENT '头像',
  `signature` varchar(100) NOT NULL DEFAULT '' COMMENT '个性签名',
  `create_ip` varchar(50) NOT NULL DEFAULT '',
  `create_time` int(11) NOT NULL,
  PRIMARY KEY (`uid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表';

LOCK TABLES `user` WRITE;
/*!40000 ALTER TABLE `user` DISABLE KEYS */;

INSERT INTO `user` (`uid`, `open`, `nickname`, `email`, `password`, `salt`, `avatar`, `signature`, `create_ip`, `create_time`)
VALUES
	(4,1,'黄三岁','abc@qq.com','5477e13bec922caabeb448545d95502f','ijniyh','https://wx.qlogo.cn/mmopen/vi_32/KMzUH2aj5qvOostTYcJC1AgQSQzPTKRT80U0WxyFN3TmLqHGveFicTXs0W8jq94avzwvsvI84jLoFowMNYtL7zg/132','','127.0.0.1:63751',1566547242),
	(10,1,'八戒的钉耙','123@qq.com','2e6bee7e39748b894d5416f89c046e73','8pmg4t','https://timgsa.baidu.com/timg?image&quality=80&size=b9999_10000&sec=1566971313030&di=a4bcb03f1ccad5d829afa41f17d65c37&imgtype=0&src=http%3A%2F%2Fb-ssl.duitang.com%2Fuploads%2Fitem%2F201711%2F10%2F20171110225150_ym2jw.jpeg','大师兄，师傅被妖怪抓走了','127.0.0.1:61388',1566550759),
	(11,1,'小河弯弯','111@111.com','8a67e2d4e3f0da0a67298bc057207d6a','89ga7b','https://timgsa.baidu.com/timg?image&quality=80&size=b9999_10000&sec=1567566373&di=89b0c0a32cc4a21b90866fbe72df4078&imgtype=jpg&er=1&src=http%3A%2F%2Fb-ssl.duitang.com%2Fuploads%2Fitem%2F201809%2F26%2F20180926162125_vjbwi.jpg','','127.0.0.1',1566896410),
	(12,1,'马总','222@qq.com','797ea3722c46b04dc1a5ae961b0661b8','t0nd6b','http://img5.imgtn.bdimg.com/it/u=2499107440,3201198713&fm=26&gp=0.jpg','','127.0.0.1',1567053925);



# Dump of table user_friend
# ------------------------------------------------------------

DROP TABLE IF EXISTS `user_friend`;

CREATE TABLE `user_friend` (
  `uid` int(11) unsigned NOT NULL,
  `friend_id` int(11) unsigned NOT NULL COMMENT '好友id',
  `remark` varchar(100) NOT NULL DEFAULT '' COMMENT '备注',
  `create_time` int(10) unsigned NOT NULL COMMENT '添加时间',
  UNIQUE KEY `idx_u_f` (`uid`,`friend_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;



INSERT INTO `user_friend` (`uid`, `friend_id`, `remark`, `create_time`)
VALUES
	(4,10,'',1567134216),
	(4,11,'',1567153001),
	(4,12,'',1567134019),
	(10,4,'',1567134216),
	(11,4,'',1567153001),
	(12,4,'',1567134019);



# Dump of table user_friend_request
# ------------------------------------------------------------

DROP TABLE IF EXISTS `user_friend_request`;

CREATE TABLE `user_friend_request` (
  `req_id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `uid` int(10) unsigned NOT NULL,
  `from_uid` int(11) unsigned NOT NULL,
  `message` varchar(100) NOT NULL DEFAULT '' COMMENT '消息',
  `status` tinyint(1) unsigned NOT NULL DEFAULT '0' COMMENT '0:未处理 1:已通过 2:已拒绝',
  `create_time` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '申请时间',
  PRIMARY KEY (`req_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='好友申请表';


INSERT INTO `user_friend_request` (`req_id`, `uid`, `from_uid`, `message`, `status`, `create_time`)
VALUES
	(45,4,12,'',1,1567132680),
	(46,10,4,'',1,1567134177),
	(48,4,11,'',1,1567152991);


