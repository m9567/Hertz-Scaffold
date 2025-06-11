CREATE TABLE g_platform_tenant
(
    id            bigint                                                        NOT NULL AUTO_INCREMENT,
    created_at    datetime                                                      NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at    datetime                                                      NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    is_delete     int                                                           NOT NULL DEFAULT '0',
    platform_code varchar(255)                                                  NOT NULL,
    prefix        varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
    host          varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_0900_ai_ci NOT NULL,
    tenant_code   varchar(255)                                                  NOT NULL,
    PRIMARY KEY (id)
) ENGINE=InnoDB AUTO_INCREMENT=3 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
INSERT INTO g_platform_tenant (id, created_at, updated_at, is_delete, platform_code, prefix, host,
                                 tenant_code)
VALUES (1, '2025-06-06 23:22:06', '2025-06-10 00:34:37', 0, 'jdb', '9wyl', 'https://testapi.abcvip.website', 'dv');
INSERT INTO g_platform_tenant (id, created_at, updated_at, is_delete, platform_code, prefix, host,
                                 tenant_code)
VALUES (2, '2025-06-10 01:25:46', '2025-06-10 22:09:44', 0, 'pg', 'test', 'https://testapi.abcvip.website', 'dv');
