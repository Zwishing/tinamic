-- 基础信息的模式，存放矢量、影像数据的基础信息
CREATE SCHEMA IF NOT EXISTS spatial_info;
-- 矢量图层
CREATE SCHEMA IF NOT EXISTS vectors;
-- 影像图层
CREATE SCHEMA IF NOT EXISTS images;

-- 所有上传到minio的数据数据记录
CREATE TABLE spatial_info.sources(
    id serial PRIMARY KEY,
    uuid uuid NOT NULL UNIQUE,
    name varchar(255) NOT NULL,
    data_type varchar(255) NOT NULL,
    file_path varchar(255)[] NOT NULL,
    create_at timestamptz

)

CREATE TABLE spatial_info.vectors(
    id serial PRIMARY KEY,
    uuid uuid NOT NULL UNIQUE,
    schema varchar(255) NOT NULL,
    data_type varchar(255) NOT NULL,
    file_path varchar(255) NOT NULL,
    create_at timestamptz
)

CREATE TABLE spatial_info.table_layer(
		id  serial PRIMARY KEY,
		uid uuid NOT NULL UNIQUE,
		schema varchar(255) NOT NULL,
		name varchar(255) NOT NULL,
		attr json NOT NULL,
        geometry_type varchar(255),
        id_column varchar(255),
        geometry_column varchar(255),
        srid int,
        center float[],
        bounds float[],
        min_zoom int,
        max_zoom int,
        tile_url varchar(255),
		description text,
        thumbnail bytea,
		create_at timestamptz,
		update_at timestamptz
);
CREATE TABLE IF NOT EXISTS spatial_info.esri_wkt(
    srid integer PRIMARY KEY NOT NULL,
    type varchar(10) NOT NULL,
    auth_srid integer,
    esri_wkt varchar(2048)
);


CREATE TABLE IF NOT EXISTS spatial_info.spatial_data(
    id serial PRIMARY KEY,
    uid uuid NOT NULL UNIQUE,
    name varchar(255) NOT NULL,
    is_publish bool NOT NULL,
    file_type varchar(255) NOT NULL,
    size integer NOT NULL,
    file_path varchar(255) NOT NULL,
    create_at timestamptz,
    update_at timestamptz
)

-- 存放所有影像图层的基础信息的数据的表
CREATE TABLE spatial_info.images(
    id  serial PRIMARY KEY,
    uid uuid NOT NULL UNIQUE,
    schema varchar(255) NOT NULL,
    name varchar(255) NOT NULL,
    srid int,
    center float[],
    bounds float[],
    min_zoom int,
    max_zoom int,
    tile_url varchar(255),
    description text,
    thumbnail bytea,
    create_at timestamptz,
    update_at timestamptz
);

-- 影像图层数据，一个影像对应一张表，表名用uuid标识
CREATE TABLE IF NOT EXISTS images.uuid(

)

