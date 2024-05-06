-- 基础信息的模式，存放矢量、影像数据的基础信息
CREATE SCHEMA IF NOT EXISTS source_info;

-- 所有上传到minio的数据数据记录
CREATE TABLE source_info.vectors(
    id serial PRIMARY KEY,
    uuid uuid NOT NULL UNIQUE,
    name varchar(255) NOT NULL,
    data_type smallint,
    size bigint,
    layers varchar(255)[] NOT NULL,
    file_path varchar(255) NOT NULL,
    created timestamptz,
    creator varchar(72),
    edited timestamptz,
    editor varchar(72),
    deleted bool DEFAULT FALSE
);

-- 注释
COMMENT ON TABLE source_info.vectors IS '矢量源数据表';

COMMENT ON COLUMN source_info.vectors.id IS '源数据ID';
COMMENT ON COLUMN source_info.vectors.uuid IS '唯一标识ID';
COMMENT ON COLUMN source_info.vectors.name IS '源文件名称';
COMMENT ON COLUMN source_info.vectors.data_type IS '数据类别,1=shapefile，2=geojson';
COMMENT ON COLUMN source_info.vectors.size IS '文件大小以kb为单位';
COMMENT ON COLUMN source_info.vectors.layers IS '图层名称,支持多图层';
COMMENT ON COLUMN source_info.vectors.file_path IS '文件存储路径';
COMMENT ON COLUMN source_info.vectors.created IS '创建时间';
COMMENT ON COLUMN source_info.vectors.creator IS '创建人';
COMMENT ON COLUMN source_info.vectors.edited IS '修改时间';
COMMENT ON COLUMN source_info.vectors.editor IS '修改人';
COMMENT ON COLUMN source_info.vectors.deleted IS '逻辑删除:true=删除,false=未删除';
-- 索引
CREATE INDEX vectors_id_index ON source_info.vectors(id);
CREATE INDEX vectors_uuid_index ON source_info.vectors(uuid);