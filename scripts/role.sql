CREATE TABLE `userregistry`.`role` (
  `id` INT NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(5) NOT NULL,
  PRIMARY KEY (`id`));

INSERT INTO `userregistry`.`role` (`name`) VALUES ('admin');
INSERT INTO `userregistry`.`role` (`name`) VALUES ('user');
