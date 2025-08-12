-- phpxxxd.t_company definition

CREATE TABLE `t_company`
(
    `id`          int unsigned                                                  NOT NULL AUTO_INCREMENT COMMENT '主键默认id',
    `name`        varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '公司名称',
    `address`     varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci          DEFAULT NULL COMMENT '公司地址',
    `price_level` varchar(31) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci  NOT NULL DEFAULT 'price_1' COMMENT '该公司查看产品的价格等级，默认为1，一共4级',

    -- 新增的字段
    `created_at`  timestamp                                                     NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '记录创建时间',
    `updated_at`  timestamp                                                     NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '记录最后更新时间',
    `deleted_at`  timestamp                                                     NULL     DEFAULT NULL COMMENT '软删除时间戳', -- NULL 表示未删除                                                    NOT NULL DEFAULT '0' COMMENT '是否删除，0为未删除，1为已删除',

    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_name` (`name`),
    KEY `idx_deleted_at` (`deleted_at`)                                                                                  -- 为软删除标记添加索引
) ENGINE = InnoDB
  AUTO_INCREMENT = 1
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci COMMENT ='公司信息表';