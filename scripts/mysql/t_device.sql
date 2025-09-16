-- 创建选型方案表
CREATE TABLE "t_device" (
  "id" bigserial NOT NULL,
  "name" varchar(255) NOT NULL,
  "group_id" bigint NOT NULL,
  "details" jsonb NOT NULL,
  "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "deleted_at" timestamptz,
  PRIMARY KEY ("id")
);

-- 添加注释
COMMENT ON COLUMN "t_device"."name" IS '方案名称，来自Excel';
COMMENT ON COLUMN "t_device"."group_id" IS '所属分组ID';
COMMENT ON COLUMN "t_device"."details" IS '方案的动态详情(jsonb)，包含筛选条件、组件列表和公共参数';
COMMENT ON TABLE "t_device" IS '选型方案信息表';

-- 创建索引
CREATE INDEX "idx_t_device_group_id" ON "t_device" ("group_id");
CREATE INDEX "idx_t_device_deleted_at" ON "t_device" ("deleted_at");

-- 【关键】为 JSONB 字段创建 GIN 索引以加速筛选
CREATE INDEX "idx_t_device_details_gin" ON "t_device" USING GIN ("details");