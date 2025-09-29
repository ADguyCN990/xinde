CREATE TABLE `t_device_access_log` (
  `id` int unsigned NOT NULL AUTO_INCREMENT,
  `user_id` int unsigned NOT NULL COMMENT '访问的用户ID',
  `company_id` int unsigned DEFAULT NULL COMMENT '用户所属公司ID',
  `device_type_id` int unsigned NOT NULL COMMENT '访问的设备类型ID',
  `accessed_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '访问时间',
  PRIMARY KEY (`id`),
  KEY `idx_user_id` (`user_id`),
  KEY `idx_company_id` (`company_id`),
  KEY `idx_device_type_id` (`device_type_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci COMMENT='设备类型访问记录表';