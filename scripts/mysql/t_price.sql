CREATE TABLE `t_price`
(
    `id`           int unsigned                                                 NOT NULL AUTO_INCREMENT COMMENT '价格记录主键ID',
    `product_code` varchar(31) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '产品编码',
    `unit`         varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '产品单位',
    `spec_code`    varchar(31) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '规格型号',

    -- 价格字段: decimal 类型非常适合存储精确的货币值
    `price_1`      decimal(10, 2)                                               NOT NULL DEFAULT '0.00' COMMENT '价格等级1',
    `price_2`      decimal(10, 2)                                               NOT NULL DEFAULT '0.00' COMMENT '价格等级2',
    `price_3`      decimal(10, 2)                                               NOT NULL DEFAULT '0.00' COMMENT '价格等级3',
    `price_4`      decimal(10, 2)                                               NOT NULL DEFAULT '0.00' COMMENT '价格等级4',

    -- 时间戳与软删除 (与 t_user 和 t_company 保持一致)
    `created_at`   timestamp                                                    NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '记录创建时间',
    `updated_at`   timestamp                                                    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '记录最后更新时间',
    `deleted_at`   timestamp                                                    NULL     DEFAULT NULL COMMENT '软删除时间戳',

    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_product_code` (`product_code`), -- 使用 uk_ 前缀表示唯一索引
    KEY `idx_deleted_at` (`deleted_at`)            -- 为软删除字段添加索引
) ENGINE = InnoDB
  AUTO_INCREMENT = 1
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci COMMENT ='产品价格信息表';