/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8mb4 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;


# Dump of table favorite
# ------------------------------------------------------------

CREATE TABLE `favorite` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `userId` bigint(20) NOT NULL,
  `userName` varchar(100) NOT NULL,
  `tweetId` bigint(20) NOT NULL,
  `status` text NOT NULL,
  `favDate` datetime NOT NULL,
  `unfavDate` datetime DEFAULT NULL,
  `lastAction` datetime NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;



# Dump of table follow
# ------------------------------------------------------------

CREATE TABLE `follow` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `userId` bigint(20) NOT NULL,
  `userName` varchar(100) NOT NULL,
  `status` text,
  `followDate` datetime NOT NULL,
  `unfollowDate` datetime DEFAULT NULL,
  `lastAction` datetime NOT NULL,
  PRIMARY KEY (`id`),
  KEY `UserId` (`userId`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;



# Dump of table reply
# ------------------------------------------------------------

CREATE TABLE `reply` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `userId` bigint(20) NOT NULL,
  `userName` varchar(100) NOT NULL,
  `tweetId` bigint(20) NOT NULL,
  `status` text NOT NULL,
  `answer` text NOT NULL,
  `replyDate` datetime NOT NULL,
  PRIMARY KEY (`id`),
  KEY `TweetId` (`tweetId`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;



# Dump of table tweet
# ------------------------------------------------------------

CREATE TABLE `tweet` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `content` text NOT NULL,
  `date` datetime NOT NULL,
  PRIMARY KEY (`id`),
  KEY `Content` (`content`(255))
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;




/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;
/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
