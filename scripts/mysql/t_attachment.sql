CREATE TABLE `t_attachment`
(
    `id`              int unsigned                            NOT NULL AUTO_INCREMENT COMMENT '附件主键ID',
    `filename`        varchar(255) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '原始文件名',
    `storage_path`    varchar(512) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '文件在存储系统中的相对路径或Key',
    `file_type`       varchar(100) COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '文件MIME类型 (e.g., application/vnd.ms-excel)',
    `file_size`       bigint unsigned                         NOT NULL COMMENT '文件大小 (bytes)',
    `storage_driver`  varchar(20) COLLATE utf8mb4_unicode_ci  NOT NULL DEFAULT 'local' COMMENT '存储驱动 (local, s3, oss, etc.)',
    `uploaded_by_uid` int unsigned                            NOT NULL COMMENT '上传者用户ID',
    `business_type`   varchar(50) COLLATE utf8mb4_unicode_ci           DEFAULT NULL COMMENT '业务类型 (e.g., group_icon, price_import)',
    `business_id`     int unsigned                                     DEFAULT NULL COMMENT '业务ID，关联具体业务表的主键',
    `created_at`      timestamp                               NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '上传时间',
    `updated_at`      timestamp                               NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    `deleted_at`      timestamp                               NULL     DEFAULT NULL,
    PRIMARY KEY (`id`),
    KEY `idx_uploaded_by_uid` (`uploaded_by_uid`),
    KEY `idx_business_lookup` (`business_type`, `business_id`),
    KEY `idx_deleted_at` (`deleted_at`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci COMMENT ='附件信息表';