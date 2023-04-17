
CREATE TABLE `user`
(
    `id`                INT(11)               NOT NULL AUTO_INCREMENT,
    `name`              VARCHAR(255)          NOT NULL,
    `telegram_id`       INT(16)               NOT NULL,
    `telegram_name`     VARCHAR(255)          NOT NULL,
    `created_at`        DATETIME              DEFAULT CURRENT_TIMESTAMP,
    `updated_at`        DATETIME              DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`)
)
    ENGINE = InnoDB
    DEFAULT CHARSET = utf8mb4
    COLLATE utf8mb4_unicode_ci;

CREATE TABLE `office`
(
    `id`                INT(11)               NOT NULL AUTO_INCREMENT,
    `name`              VARCHAR(255)          NOT NULL,
    `city`              VARCHAR(255)          NOT NULL,
    `created_at`        DATETIME              DEFAULT CURRENT_TIMESTAMP,
    `updated_at`        DATETIME              DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`)
)
    ENGINE = InnoDB
    DEFAULT CHARSET = utf8mb4
    COLLATE utf8mb4_unicode_ci;

CREATE TABLE `seat`
(
    `id`                INT(11)          NOT NULL AUTO_INCREMENT,
    `seat_number`       INT(11)          NOT NULL,
    `have_monitor`      BOOLEAN          DEFAULT FALSE,
    `office_id`         INT(11)          NOT NULL,
    `created_at`        DATETIME         DEFAULT CURRENT_TIMESTAMP,
    `updated_at`        DATETIME         DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`)
)
    ENGINE = InnoDB
    DEFAULT CHARSET = utf8mb4
    COLLATE utf8mb4_unicode_ci;

CREATE TABLE `book_seat`
(
    `id`                INT(11)          NOT NULL AUTO_INCREMENT,
    `office_id`         INT(11)          NOT NULL,
    `seat_id`           INT(11)          NOT NULL,
    `user_id`           INT(11)          DEFAULT NULL,
    `book_date`         DATETIME         NOT NULL,
    `book_start_time`   DATETIME         DEFAULT NULL,
    `book_end_time`     DATETIME         DEFAULT NULL,
    `created_at`        DATETIME         DEFAULT CURRENT_TIMESTAMP,
    `updated_at`        DATETIME         DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`)
)
    ENGINE = InnoDB
    DEFAULT CHARSET = utf8mb4
    COLLATE utf8mb4_unicode_ci;

alter table user add office_id INT(11) default NULL;

ALTER TABLE user add column notify_office_id int(16) default 0;

ALTER TABLE office add column time_zone VARCHAR(255) NOT NULL;

ALTER TABLE user add column chat_id INT(11) default 0;

CREATE TABLE `work_date`
(
    `id`                INT(11)          NOT NULL AUTO_INCREMENT,
    `status`            VARCHAR(255)     NOT NULL,
    `work_date`         DATETIME         NOT NULL,
    `created_at`        DATETIME         DEFAULT CURRENT_TIMESTAMP,
    `updated_at`        DATETIME         DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`)
)
    ENGINE = InnoDB
    DEFAULT CHARSET = utf8mb4
    COLLATE utf8mb4_unicode_ci;

ALTER TABLE book_seat add column confirm BOOLEAN DEFAULT FALSE;

ALTER TABLE seat add column seat_sign VARCHAR(255) NOT NULL;

ALTER TABLE book_seat add column `hold` BOOLEAN DEFAULT FALSE;

ALTER TABLE volt.user MODIFY COLUMN chat_id bigint;

ALTER TABLE volt.user MODIFY COLUMN telegram_id bigint;