create schema if not exists mockserver;

CREATE TABLE `mockserver`.`rules`
(
    `key`                 varchar(255) NOT NULL,
    `group`               varchar(255) NOT NULL,
    `name`                varchar(255) NOT NULL,
    `path`                varchar(255) NOT NULL,
    `strategy`            varchar(255) NOT NULL,
    `method`              varchar(45)  NOT NULL,
    `status`              varchar(255) NOT NULL,
    `pattern`             varchar(255) NOT NULL,
    `next_response_index` int NOT NULL DEFAULT '0' ,
        PRIMARY KEY (`key`),
    UNIQUE KEY `key_UNIQUE` (`key`)
) ENGINE = InnoDB
  DEFAULT CHARSET = latin1;

CREATE TABLE `mockserver`.`responses`
(
    `id`           bigint       NOT NULL AUTO_INCREMENT,
    `body`         longtext     NOT NULL,
    `content_type` varchar(255) NOT NULL,
    `http_status`  int          NOT NULL,
    `delay`        int          DEFAULT '0',
    `scene`        varchar(255) DEFAULT NULL,
    `rule_key`     varchar(255) NOT NULL,
    PRIMARY KEY (`id`),
    KEY `rules_idx` (`rule_key`),
    CONSTRAINT `rules` FOREIGN KEY (`rule_key`) REFERENCES `rules` (`key`)
) ENGINE = InnoDB
  DEFAULT CHARSET = latin1;

CREATE TABLE `mockserver`.`variables`
(
    `id`         bigint       NOT NULL AUTO_INCREMENT,
    `type`       varchar(255) NOT NULL,
    `name`       varchar(255) NOT NULL,
    `key`        varchar(255) NOT NULL,
    `rule_key`   varchar(255) NOT NULL,
    `assertions` json DEFAULT NULL,
    PRIMARY KEY (`id`),
    KEY `rule_idx` (`rule_key`),
    CONSTRAINT `rule` FOREIGN KEY (`rule_key`) REFERENCES `rules` (`key`)
) ENGINE = InnoDB
  DEFAULT CHARSET = latin1;
