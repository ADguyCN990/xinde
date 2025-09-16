CREATE TABLE "t_device_type" (
  "id" bigserial NOT NULL,
  "name" varchar(255) NOT NULL, -- 用户在导入时输入的“设备名称”
  "group_id" bigint NOT NULL,   -- 它属于哪个分组
  -- 这个表的 attachment business_type 可以是 'device_type_main_image'
  "created_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "updated_at" timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,
  "deleted_at" timestamptz,
  PRIMARY KEY ("id"),
  UNIQUE("name", "group_id") -- 同一分组下设备类型名称唯一
);
-- 添加注释
COMMENT ON COLUMN "t_device_type"."name" IS '设备类型名称 (e.g., U钻)';
COMMENT ON COLUMN "t_device_type"."group_id" IS '设备类型所属的分组ID';
COMMENT ON TABLE "t_device_type" IS '设备类型信息表';