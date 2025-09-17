-- 创建选型方案表
CREATE TABLE "t_device" (
  "id" bigserial NOT NULL,
  "device_type_id" bigint NOT NULL, -- 【新增】外键，关联到 t_device_type
  "name" varchar(255) NOT NULL,    -- 方案名称 ("方案1", "方案2"...)
  "details" jsonb NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "deleted_at" timestamptz DEFAULT NULL,
  PRIMARY KEY ("id")
);

-- 添加注释
COMMENT ON COLUMN "t_device"."device_type_id" IS '外键，关联到 t_device_type.id';
COMMENT ON COLUMN "t_device"."name" IS '方案名称 (e.g., 方案1, 方案2)';
COMMENT ON COLUMN "t_device"."details" IS '方案的动态详情(jsonb)，包含筛选条件、组件列表和公共参数';
COMMENT ON TABLE "t_device" IS '选型方案信息表';


-- 创建索引
CREATE INDEX "idx_t_device_device_type_id" ON "t_device" ("device_type_id");
CREATE INDEX "idx_t_device_deleted_at" ON "t_device" ("deleted_at");

-- 【关键】为 JSONB 字段创建 GIN 索引以加速筛选
CREATE INDEX "idx_t_device_details_gin" ON "t_device" USING GIN ("details");
