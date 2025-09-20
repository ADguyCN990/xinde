-- 在 PostgreSQL 数据库中执行
CREATE TABLE "t_filter_image" (
"id" bigserial NOT NULL,
"device_type_id" bigint NOT NULL,
"filter_value" varchar(255) NOT NULL,
"created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
"updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
"deleted_at" timestamptz,
PRIMARY KEY ("id")
);
-- 添加注释
COMMENT ON COLUMN "t_filter_image"."device_type_id" IS '关联的设备类型ID';
COMMENT ON COLUMN "t_filter_image"."filter_value" IS '筛选条件的值 (e.g., 博世)';
COMMENT ON TABLE "t_filter_image" IS '筛选条件图片配置表';
-- 创建普通索引
CREATE INDEX "idx_t_filter_image_device_type_id" ON "t_filter_image" ("device_type_id");
CREATE INDEX "idx_t_filter_image_deleted_at" ON "t_filter_image" ("deleted_at");