CREATE TABLE `user`
(
    `id`                INT(11)         NOT NULL AUTO_INCREMENT,
    `name`              VARCHAR(255)    NOT NULL,
    `telegram_id`       BIGINT          NOT NULL,
    `telegram_name`     VARCHAR(255)    NOT NULL,
    `office_id`     	INT(11)         NOT NULL,
    `notify_office_id`  INT(16)			NOT NULL,
    `chat_id`     		BIGINT          NOT NULL,
    `created_at`        DATETIME        DEFAULT CURRENT_TIMESTAMP,
    `updated_at`        DATETIME        DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`)
)
    ENGINE = InnoDB
    DEFAULT CHARSET = utf8mb4
    COLLATE utf8mb4_unicode_ci;

CREATE TABLE `office`
(
    `id`                INT(11)          NOT NULL AUTO_INCREMENT,
    `name`              VARCHAR(255)     NOT NULL,
    `city`              VARCHAR(255)     NOT NULL,
    `time_zone`         VARCHAR(255)     NOT NULL,
    `created_at`        DATETIME         DEFAULT CURRENT_TIMESTAMP,
    `updated_at`        DATETIME         DEFAULT CURRENT_TIMESTAMP,
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
    `seat_sign`         VARCHAR(255)     NOT NULL,
    `created_at`        DATETIME         DEFAULT CURRENT_TIMESTAMP,
    `updated_at`        DATETIME         DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`)
)
    ENGINE = InnoDB
    DEFAULT CHARSET = utf8mb4
    COLLATE utf8mb4_unicode_ci;

CREATE TABLE `book_seat`
(
    `id`                INT(11)         NOT NULL AUTO_INCREMENT,
    `office_id`         INT(11)         NOT NULL,
    `seat_id`           INT(11)         NOT NULL,
    `user_id`           INT(11)         DEFAULT NULL,
    `book_date`         DATETIME        NOT NULL,
    `book_start_time`   DATETIME        DEFAULT NULL,
    `book_end_time`     DATETIME        DEFAULT NULL,
    `confirm`     		BOOLEAN         DEFAULT FALSE,
    `hold`     			BOOLEAN         DEFAULT FALSE,
    `created_at`        DATETIME        DEFAULT CURRENT_TIMESTAMP,
    `updated_at`        DATETIME        DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`)
)
    ENGINE = InnoDB
    DEFAULT CHARSET = utf8mb4
    COLLATE utf8mb4_unicode_ci;

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

CREATE INDEX idx_user_id_date ON book_seat (user_id, book_date);
CREATE INDEX idx_office_id ON book_seat (office_id, book_date);

CREATE INDEX id_telegram_id ON user (telegram_id);

CREATE INDEX idx_work_date ON work_date (work_date);