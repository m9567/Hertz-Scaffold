CREATE TABLE g_platform_key
(
    id         bigint       NOT NULL AUTO_INCREMENT,
    created_at datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at datetime     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    is_delete  int          NOT NULL DEFAULT '0',
    code       varchar(255) NOT NULL,
    currency   varchar(255) NOT NULL,
    key_json   text         NOT NULL,
    url_json   text         NOT NULL,
    PRIMARY KEY (id)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

INSERT INTO g_platform_key (id, created_at, updated_at, is_delete, code, currency, key_json, url_json)
VALUES (1, '2025-06-05 00:12:57', '2025-06-10 23:18:24', 0, 'jdb', 'usd',
        '{\"agent\":\"wwyl_wp001t\",\"dc\":\"zfs\",\"key\":\"47e0cd2ece0883e2\",\"iv\":\"b87f2867577b68ce\"}',
        '{\"apiRequest.do\":\"https://api.jygrq.com/apiRequest.do\"}');
INSERT INTO g_platform_key (id, created_at, updated_at, is_delete, code, currency, key_json, url_json)
VALUES (2, '2025-06-10 01:25:23', '2025-06-10 23:50:42', 0, 'pg', 'usd', '{}',
        '{\"LoginProxy\":\"https://api.pg-bo.me/external/Login/v1/LoginProxy\"}');
