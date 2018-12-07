DROP TABLE IF EXISTS `shop_user`;
CREATE TABLE IF NOT EXISTS `shop_user`(
  `uid` BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `username` VARCHAR(32) NOT NULL DEFAULT '',
  `password` CHAR(32) NOT NULL DEFAULT '',
  `email` VARCHAR(100) NOT NULL DEFAULT '',
  `create_time` INT UNSIGNED NOT NULL DEFAULT '0',
  `login_ip` int not null default '0',
  `login_time` int not null default '0',
  unique (`email`),
  PRIMARY KEY(`uid`)
)ENGINE=InnoDB DEFAULT CHARSET=utf8;