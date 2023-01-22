
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