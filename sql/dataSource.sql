CREATE SCHEMA IF NOT EXISTS data_source;

CREATE TABLE IF NOT EXISTS data_source.base_info(
    id serial PRIMARY KEY,
    uuid uuid NOT NULL UNIQUE,
    name varchar(255) NOT NULL,
    data_type integer NOT NULL,
    file_path varchar(255) NOT NULL,
    owner uuid NOT NULL,
    created timestamptz default now()
);

COMMENT ON TABLE data_source.base_info IS '存储数据源的基础信息表';

COMMENT ON COLUMN data_source.base_info.id IS '自增ID';
COMMENT ON COLUMN data_source.base_info.uuid IS '数据源的唯一标识';
COMMENT ON COLUMN data_source.base_info.name IS '数据源的名称';
COMMENT ON COLUMN data_source.base_info.data_type IS '数据源的数据类型：1-矢量，2-影像';
COMMENT ON COLUMN data_source.base_info.file_path IS '数据源的minio的存储路径';
COMMENT ON COLUMN data_source.base_info.owner IS '数据源的上传用户的id';
COMMENT ON COLUMN data_source.base_info.created IS '数据源的上传时间';

CREATE INDEX origin_info_id_index ON data_source.base_info(id);
CREATE INDEX origin_info_uuid_index ON data_source.base_info(uuid);

