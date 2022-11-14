CREATE DATABASE IF NOT EXISTS tpms_prod;
USE tpms_prod;
CREATE TABLE IF NOT EXISTS `users`
(
    `id`        VARCHAR(255) NOT NULL,
    `email`     VARCHAR(255) NOT NULL,
    `phone`     VARCHAR(40)  NOT NULL,
    `fmt_token` VARCHAR(255),
    `name`      VARCHAR(255) NOT NULL,
    `optout`    BOOLEAN DEFAULT FALSE,
    `latitude`  DOUBLE,
    `longitude` DOUBLE,
    PRIMARY KEY (`id`)
);
CREATE TABLE IF NOT EXISTS `dogs`
(
    `id`              BIGINT      NOT NULL AUTO_INCREMENT,
    `name`            VARCHAR(50) NOT NULL,
    `breed`           int         NOT NULL,
    `age`             int         NOT NULL,
    `size`            int         NOT NULL,
    `coat_color`      int         NOT NULL,
    `coat_length`     int         NOT NULL,
    `tail_length`     int         NOT NULL,
    `ear`             int         NOT NULL,
    `is_lost`         BOOLEAN,
    `owner_id`        VARCHAR(255),
    `host_id`         VARCHAR(255),
    `latitude`        DOUBLE,
    `longitude`       DOUBLE,
    `img_url`         LONGTEXT,
    `additional_info` varchar(500),
    `embedding`       LONGTEXT,
    `created_at`      DATETIME    NOT NULL,
    `deleted_at`      DATETIME,
    PRIMARY KEY (`id`),
    CONSTRAINT `fk_owner_id` FOREIGN KEY (`owner_id`) REFERENCES users (`id`) ON UPDATE CASCADE ON DELETE SET NULL,
    CONSTRAINT `fk_host_id` FOREIGN KEY (`host_id`) REFERENCES users (`id`) ON UPDATE CASCADE ON DELETE SET NULL
);

CREATE TABLE IF NOT EXISTS `possible_matches`
(
    `dog_id`          BIGINT NOT NULL,
    `possible_dog_id` BIGINT NOT NULL,
    `ack`             INT    NOT NULL,
    PRIMARY KEY (`dog_id`, `possible_dog_id`)
);

CREATE DATABASE IF NOT EXISTS tpms_test;
USE tpms_test;
CREATE TABLE IF NOT EXISTS users
(
    `id`        VARCHAR(255) NOT NULL,
    `email`     VARCHAR(255) NOT NULL,
    `phone`     VARCHAR(40)  NOT NULL,
    `fcm_token` VARCHAR(255),
    `name`      VARCHAR(255) NOT NULL,
    `optout`    BOOLEAN DEFAULT FALSE,
    `latitude`  DOUBLE,
    `longitude` DOUBLE,
    PRIMARY KEY (`id`)
);
CREATE TABLE IF NOT EXISTS `dogs`
(
    `id`              BIGINT      NOT NULL AUTO_INCREMENT,
    `name`            VARCHAR(50) NOT NULL,
    `breed`           int         NOT NULL,
    `age`             int         NOT NULL,
    `size`            int         NOT NULL,
    `coat_color`      int         NOT NULL,
    `coat_length`     int         NOT NULL,
    `tail_length`     int         NOT NULL,
    `ear`             int         NOT NULL,
    `is_lost`         boolean,
    `owner_id`        VARCHAR(255),
    `host_id`         VARCHAR(255),
    `latitude`        DOUBLE,
    `longitude`       DOUBLE,
    `img_url`         LONGTEXT,
    `additional_info` varchar(500),
    `embedding`       LONGTEXT,
    `created_at`      DATETIME    NOT NULL,
    `deleted_at`      DATETIME,
    PRIMARY KEY (`id`),
    CONSTRAINT `fk_owner_id` FOREIGN KEY (`owner_id`) REFERENCES users (`id`) ON UPDATE CASCADE ON DELETE SET NULL,
    CONSTRAINT `fk_host_id` FOREIGN KEY (`host_id`) REFERENCES users (`id`) ON UPDATE CASCADE ON DELETE SET NULL
);
CREATE TABLE IF NOT EXISTS `possible_matches`
(
    `dog_id`          BIGINT NOT NULL,
    `possible_dog_id` BIGINT NOT NULL,
    `ack`             INT    NOT NULL,
    PRIMARY KEY (`dog_id`, `possible_dog_id`)
);
