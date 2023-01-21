CREATE TABLE `user`
(
    `id`                INT(11)      NOT NULL AUTO_INCREMENT,
    `name`              VARCHAR(255) NOT NULL,
    `telegram_id`                INT(11)      NOT NULL AUTO_INCREMENT,
    `created_at`        DATETIME              DEFAULT CURRENT_TIMESTAMP,
    `updated_at`        DATETIME              DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`)
)
    ENGINE = InnoDB
    DEFAULT CHARSET = utf8mb4
    COLLATE utf8mb4_unicode_ci;