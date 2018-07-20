CREATE TABLE `userregistry`.`user` (
  `id` INT(10) NOT NULL AUTO_INCREMENT,
  `name` VARCHAR(16) NOT NULL,
  `email` VARCHAR(256) NOT NULL,
  `password` BINARY(64) NOT NULL,
  `create_time` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `role` INT(10) NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `fk_role_idx` (`role` ASC),
  CONSTRAINT `fk_role`
    FOREIGN KEY (`role`)
    REFERENCES `userregistry`.`role` (`id`)
    ON DELETE NO ACTION
    ON UPDATE NO ACTION);

INSERT INTO `userregistry`.`user` (`name`, `email`, `password`, `role`) VALUES ('John Doe', 'johndoe1@user.com', 'MQ==', '2');
INSERT INTO `userregistry`.`user` (`name`, `email`, `password`, `role`) VALUES ('John Doe Jr.', 'johndoe2@user.com', 'MQ==', '2');
INSERT INTO `userregistry`.`user` (`name`, `email`, `password`, `role`) VALUES ('Joan Doe', 'joandoe1@user.com', 'MQ==', '2');
INSERT INTO `userregistry`.`user` (`name`, `email`, `password`, `role`) VALUES ('Joan Doe Jr.', 'joandoe2@user.com', 'MQ==', '2');
INSERT INTO `userregistry`.`user` (`name`, `email`, `password`, `role`) VALUES ('John Doe A', 'johndoe1@admin.com', 'MQ==', '1');
INSERT INTO `userregistry`.`user` (`name`, `email`, `password`, `role`) VALUES ('John Doe Jr. A', 'johndoe2@admin.com', 'MQ==', '1');
