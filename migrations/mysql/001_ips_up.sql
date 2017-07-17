

CREATE TABLE IF NOT EXISTS `my_schema`.`ips` (
  `id` BIGINT NOT NULL AUTO_INCREMENT,
  `ip` BIGINT NOT NULL,
  PRIMARY KEY (`id`))
  ENGINE = InnoDB;

INSERT INTO `my_schema`.`ips` (`ip`) VALUES (0), (1), (23456), (-9876543210), (6789012345);
