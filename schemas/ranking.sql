# ************************************************************
# Sequel Pro SQL dump
# Version 481
#
# https://www.sequelpro.com/
# https://github.com/sequelpro/sequelpro
#
# Host: 127.0.0.1 (MySQL 5.7.18-15-log)
# Database: ranking
# Generation Time: 2024-08-29 07:43:24 +0000
# ************************************************************


/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
SET NAMES utf8mb4;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;


# Dump of table activity
# ------------------------------------------------------------

DROP TABLE IF EXISTS `activity`;

CREATE TABLE `activity` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `activity_name` varchar(50) NOT NULL DEFAULT '',
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

LOCK TABLES `activity` WRITE;
/*!40000 ALTER TABLE `activity` DISABLE KEYS */;

INSERT INTO `activity` (`id`, `activity_name`, `created_at`, `updated_at`)
VALUES
	(1,'六一活动','2024-06-01 00:00:00','2024-06-01 00:00:00');

/*!40000 ALTER TABLE `activity` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table player
# ------------------------------------------------------------

DROP TABLE IF EXISTS `player`;

CREATE TABLE `player` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `activity_id` int(11) NOT NULL,
  `player_id` int(11) DEFAULT NULL,
  `player_name` varchar(25) DEFAULT NULL,
  `avatar` varchar(25) DEFAULT NULL,
  `desc` text,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

LOCK TABLES `player` WRITE;
/*!40000 ALTER TABLE `player` DISABLE KEYS */;

INSERT INTO `player` (`id`, `activity_id`, `player_id`, `player_name`, `avatar`, `desc`, `created_at`, `updated_at`)
VALUES
	(1,1,1,'张三','/images/11.png','我是参赛选手，请为我投票','2024-06-02 00:00:00','2024-08-12 15:01:28'),
	(2,1,2,'李四','/images/12.png','我是参赛选手，请为我投票','2024-06-02 00:00:00','2024-08-12 15:01:30'),
	(3,1,3,'王五','/images/13.png','我是参赛选手，请为我投票','2024-06-02 00:00:00','2024-08-12 15:01:33'),
	(4,1,4,'陈六','/images/14.png','我是参赛选手，请为我投票','2024-06-02 00:00:00','2024-08-07 19:18:29');

/*!40000 ALTER TABLE `player` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table score
# ------------------------------------------------------------

DROP TABLE IF EXISTS `score`;

CREATE TABLE `score` (
  `activity_id` int(11) unsigned NOT NULL,
  `player_id` int(11) unsigned NOT NULL,
  `score` int(11) NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`activity_id`,`player_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

LOCK TABLES `score` WRITE;
/*!40000 ALTER TABLE `score` DISABLE KEYS */;

INSERT INTO `score` (`activity_id`, `player_id`, `score`, `created_at`, `updated_at`)
VALUES
	(1,1,1,'2024-08-28 16:07:22','2024-08-28 16:07:22'),
	(1,2,1,'2024-08-28 17:34:26','2024-08-28 17:34:26'),
	(1,4,1,'2024-08-28 16:07:45','2024-08-28 16:07:45');

/*!40000 ALTER TABLE `score` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table user
# ------------------------------------------------------------

DROP TABLE IF EXISTS `user`;

CREATE TABLE `user` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `user_name` varchar(25) NOT NULL DEFAULT '',
  `password` varchar(50) NOT NULL DEFAULT '',
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

LOCK TABLES `user` WRITE;
/*!40000 ALTER TABLE `user` DISABLE KEYS */;

INSERT INTO `user` (`id`, `user_name`, `password`, `created_at`, `updated_at`)
VALUES
	(1,'张山','96e79218965eb72c92a549dd5a330112','2024-07-18 00:00:00','2024-07-25 17:45:23'),
	(2,'1111','96e79218965eb72c92a549dd5a330112','2024-07-18 18:01:55','2024-07-18 18:01:55'),
	(7,'222','96e79218965eb72c92a549dd5a330112','2024-07-25 17:03:14','2024-07-25 17:03:14'),
	(8,'王五','96e79218965eb72c92a549dd5a330112','2024-07-25 17:52:54','2024-07-25 17:55:17'),
	(9,'陈六','96e79218965eb72c92a549dd5a330112','2024-07-25 17:54:32','2024-07-25 17:54:32'),
	(13,'111222','96e79218965eb72c92a549dd5a330112','2024-08-05 11:14:17','2024-08-05 11:14:17'),
	(14,'123456','96e79218965eb72c92a549dd5a330112','2024-08-12 15:01:01','2024-08-12 15:01:01');

/*!40000 ALTER TABLE `user` ENABLE KEYS */;
UNLOCK TABLES;


# Dump of table vote
# ------------------------------------------------------------

DROP TABLE IF EXISTS `vote`;

CREATE TABLE `vote` (
  `id` int(11) unsigned NOT NULL AUTO_INCREMENT,
  `user_id` int(11) NOT NULL,
  `activity_id` int(11) NOT NULL,
  `player_id` int(11) NOT NULL,
  `created_at` datetime NOT NULL,
  `updated_at` datetime NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

LOCK TABLES `vote` WRITE;
/*!40000 ALTER TABLE `vote` DISABLE KEYS */;

INSERT INTO `vote` (`id`, `user_id`, `activity_id`, `player_id`, `created_at`, `updated_at`)
VALUES
	(1,2,1,1,'2024-08-28 16:07:21','2024-08-28 16:07:21'),
	(2,2,1,4,'2024-08-28 16:07:45','2024-08-28 16:07:45'),
	(3,13,1,1,'2024-08-28 17:34:17','2024-08-28 17:34:17'),
	(4,13,1,2,'2024-08-28 17:34:26','2024-08-28 17:34:26');

/*!40000 ALTER TABLE `vote` ENABLE KEYS */;
UNLOCK TABLES;



/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;
/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
