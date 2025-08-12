-- phpxxxd.t_user definition

CREATE TABLE `t_user`
(
    `uid`           int                                                     NOT NULL AUTO_INCREMENT,
    `username`      varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '用户账号',
    `password`      varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '用户密码',
    `phone`    varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '用户电话号码',
    `is_admin`      tinyint                                                 NOT NULL DEFAULT '0' COMMENT '是否为管理员',
    `remarks`       varchar(64)                                                      DEFAULT NULL COMMENT '备注',
    `recent_search_at`      timestamp                                                        DEFAULT NULL COMMENT '上次访问时间',
    `search_device` varchar(100)                                                     DEFAULT NULL COMMENT '上次访问的设备',
    `company_name`       varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '公司名称',
    `company_address`       varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci          DEFAULT NULL COMMENT '公司地址',
    `name`          varchar(32)                                             NOT NULL COMMENT '用户真实姓名',
    `user_email`    varchar(255) CHARACTER SET utf8 COLLATE utf8_general_ci          DEFAULT NULL COMMENT '用户邮箱',
    `created_at`    timestamp                                               NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '记录创建时间',
    `updated_at`    timestamp                                               NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '记录最后更新时间',
    `is_user`       int                                                     NOT NULL COMMENT '是否审核通过，1为通过',
    `why`           varchar(255)                                                     DEFAULT NULL COMMENT '审核拒绝的原因',
    `handled_at`    timestamp                                                        DEFAULT NULL COMMENT '注册申请通过的时间',
    `company_id`    int unsigned                                                     DEFAULT NULL COMMENT '用户对应的公司ID',
    `deleted_at`    timestamp                                               NULL     DEFAULT NULL COMMENT '软删除时间戳',
    PRIMARY KEY (`uid`),
    KEY `idx_company_id` (`company_id`), -- 公司ID作为外键，应该有索引
    KEY `idx_deleted_at` (`deleted_at`)  -- 为软删除字段添加索引
) ENGINE = InnoDB
  AUTO_INCREMENT = 1
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_unicode_ci COMMENT ='用户信息表';