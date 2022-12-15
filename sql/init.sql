CREATE DATABASE IF NOT EXISTS tpms_prod;
USE tpms_prod;
CREATE TABLE IF NOT EXISTS `users`
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
    PRIMARY KEY (`id`)
);

CREATE TABLE IF NOT EXISTS `possible_matches`
(
    `dog_id`          BIGINT NOT NULL,
    `possible_dog_id` BIGINT NOT NULL,
    `ack`             INT    NOT NULL,
    PRIMARY KEY (`dog_id`, `possible_dog_id`)
);

CREATE TABLE IF NOT EXISTS `posts`
(
    `id`       BIGINT       NOT NULL AUTO_INCREMENT,
    `dog_id`   BIGINT       NOT NULL,
    `url`      varchar(500) NOT NULL,
    `title`    varchar(250) NOT NULL,
    `location` varchar(250) NOT NULL,
    PRIMARY KEY (`id`)
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
    PRIMARY KEY (`id`)
);
CREATE TABLE IF NOT EXISTS `possible_matches`
(
    `dog_id`          BIGINT NOT NULL,
    `possible_dog_id` BIGINT NOT NULL,
    `ack`             INT    NOT NULL,
    PRIMARY KEY (`dog_id`, `possible_dog_id`)
);

CREATE TABLE IF NOT EXISTS `posts`
(
    `id`       BIGINT       NOT NULL AUTO_INCREMENT,
    `dog_id`   BIGINT       NOT NULL,
    `url`      varchar(500) NOT NULL,
    `title`    varchar(250) NOT NULL,
    `location` varchar(250) NOT NULL,
    PRIMARY KEY (`id`)
);
