CREATE TABLE IF NOT EXISTS `oauth_clients` (
    `client_id` VARCHAR(32) NOT NULL,
    `client_secret` VARCHAR(32) NOT NULL,
    `redirect_uri` VARCHAR(1000) NULL,
    `grant_types` VARCHAR(100) NOT NULL,
    `scope` VARCHAR(2000) NULL,
    `user_id` BIGINT(20) NULL,
    PRIMARY KEY (`client_id`)
) ENGINE=InnoDB
DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `oauth_access_tokens` (
    `access_token` VARCHAR(40) NOT NULL,
    `client_id` VARCHAR(32) NOT NULL,
    `user_id` VARCHAR(20) NULL,
    `expires` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    `scope` VARCHAR(2000) NULL,
    PRIMARY KEY (`access_token`)
) ENGINE=InnoDB
DEFAULT CHARSET=utf8;

INSERT INTO `oauth_clients`
(`client_id`, `client_secret`, `redirect_uri`, `grant_types`, `scope`, `user_id`)
VALUES
('client_web', '3v3rm0s', 'https://evermos.com/', 'client_credentials password refresh_token', 'user','10001');

INSERT INTO `oauth_access_tokens`
(`access_token`, `client_id`, `user_id`, `expires`, `scope`)
VALUES
('00000c708db9bdf1a70d5988a8f321a82970ceb6', 'client_web', '10001', '2024-08-24 20:57:59', NULL),
('4d1ed8aea16c7212886ecf91241aa2499f053bcf', 'client_web', NULL, '2021-08-24 20:57:59', NULL);