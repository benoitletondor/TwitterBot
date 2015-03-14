CREATE TABLE `favorite` (
  `idfavorite` INT NOT NULL AUTO_INCREMENT,
  `userId` INT NOT NULL,
  `userName` VARCHAR(100) NOT NULL,
  `tweetId` INT NOT NULL,
  `status` TEXT NOT NULL,
  `favDate` DATETIME NOT NULL,
  `unfavDate` DATETIME NULL,
  `lastAction` DATETIME NOT NULL,
  PRIMARY KEY (`idfavorite`)
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8;