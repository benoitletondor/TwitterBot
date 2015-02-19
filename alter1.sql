ALTER TABLE `follow` ADD COLUMN `lastAction` DATETIME NOT NULL AFTER `unfollowDate`;
UPDATE `follow` SET `lastAction` = NOW();