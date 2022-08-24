CREATE DATABASE IF NOT EXISTS tpms_prod;
USE tpms_prod;
CREATE TABLE IF NOT EXISTS `users`
(
    `id`         VARCHAR(255) PRIMARY KEY,
    `first_name` VARCHAR(50)  NOT NULL,
    `last_name`  VARCHAR(50)  NOT NULL,
    `email`      VARCHAR(255) NOT NULL,
    `phone`      VARCHAR(40)  NOT NULL,
    `city`       VARCHAR(255) NOT NULL
);
CREATE TABLE IF NOT EXISTS `dogs`
(
    `id`         VARCHAR(255),
    `name`       VARCHAR(50) NOT NULL,
    `breed`      int         NOT NULL,
    `age`        int         NOT NULL,
    `size`       int         NOT NULL,
    `is_lost`    boolean,
    `owner_id`   VARCHAR(255),
    `host_id`    VARCHAR(255),
    `latitude`   DECIMAL,
    `longitude`  DECIMAL,
    `img_url`    LONGTEXT,
    `created_at` DATETIME    NOT NULL,
    `updated_at` DATETIME    NOT NULL,
    `deleted_at` DATETIME    NOT NULL,
    PRIMARY KEY (`id`),
    CONSTRAINT `fk_owner_id` FOREIGN KEY (`owner_id`) REFERENCES users (`id`) ON UPDATE CASCADE ON DELETE SET NULL,
    CONSTRAINT `fk_host_id` FOREIGN KEY (`host_id`) REFERENCES users (`id`) ON UPDATE CASCADE ON DELETE SET NULL
);

CREATE DATABASE IF NOT EXISTS tpms_test;
USE tpms_test;
CREATE TABLE IF NOT EXISTS users
(
    `id`         VARCHAR(255) PRIMARY KEY,
    `first_name` VARCHAR(50)  NOT NULL,
    `last_name`  VARCHAR(50)  NOT NULL,
    `email`      VARCHAR(255) NOT NULL,
    `phone`      VARCHAR(40)  NOT NULL,
    `city`       VARCHAR(255) NOT NULL
);
CREATE TABLE IF NOT EXISTS `dogs`
(
    `id`         VARCHAR(255),
    `name`       VARCHAR(50) NOT NULL,
    `breed`      int         NOT NULL,
    `age`        int         NOT NULL,
    `size`       int         NOT NULL,
    `is_lost`    boolean,
    `owner_id`   VARCHAR(255),
    `host_id`    VARCHAR(255),
    `latitude`   DECIMAL,
    `longitude`  DECIMAL,
    `img_url`    LONGTEXT,
    `created_at` DATETIME    NOT NULL,
    `updated_at` DATETIME    NOT NULL,
    `deleted_at` DATETIME    NOT NULL,
    PRIMARY KEY (`id`),
    CONSTRAINT `fk_owner_id` FOREIGN KEY (`owner_id`) REFERENCES users (`id`) ON UPDATE CASCADE ON DELETE SET NULL,
    CONSTRAINT `fk_host_id` FOREIGN KEY (`host_id`) REFERENCES users (`id`) ON UPDATE CASCADE ON DELETE SET NULL
);
