CREATE SCHEMA IF NOT EXISTS user_info;

-- 账号表:记录登录账号信息
CREATE TABLE IF NOT EXISTS user_info.account(
    id serial PRIMARY KEY,
    user_id integer,
    login_account varchar(255),
    category smallint,
    created timestamptz,
    creator varchar(72),
    edited timestamptz,
    editor varchar(72),
    deleted bool DEFAULT FALSE
);
-- 注释
COMMENT ON TABLE user_info.account IS '账号表:记录登录账号信息';

COMMENT ON COLUMN user_info.account.id IS '账号ID';
COMMENT ON COLUMN user_info.account.user_id IS '用户唯一ID';
COMMENT ON COLUMN user_info.account.login_account IS '登录账号';
COMMENT ON COLUMN user_info.account.category IS '账号类别,1=用户名，2=邮箱，3=手机号';
COMMENT ON COLUMN user_info.account.created IS '创建时间';
COMMENT ON COLUMN user_info.account.creator IS '创建人';
COMMENT ON COLUMN user_info.account.edited IS '修改时间';
COMMENT ON COLUMN user_info.account.editor IS '修改人';
COMMENT ON COLUMN user_info.account.deleted IS '逻辑删除:true=删除,false=未删除';
-- 索引
CREATE INDEX account_id_index ON user_info.account(id);
CREATE INDEX account_user_id_index ON user_info.account(user_id);


-- 用户表:记录用户基本信息和密码
CREATE TABLE IF NOT EXISTS user_info.user(
    id serial PRIMARY KEY,
    state bool DEFAULT TRUE,
    name varchar(255),
    avatar bytea,
    cell_phone varchar(11),
    salt varchar(64),
    password varchar(64),
    created timestamptz,
    creator varchar(72),
    edited timestamptz,
    editor varchar(72),
    deleted bool DEFAULT FALSE
);
-- 注释
COMMENT ON TABLE user_info.user IS '用户表:记录用户基本信息和密码';

COMMENT ON COLUMN user_info.user.id IS '用户ID';
COMMENT ON COLUMN user_info.user.state IS '用户状态,true=正常,false=禁用';
COMMENT ON COLUMN user_info.user.name IS '姓名';
COMMENT ON COLUMN user_info.user.avatar IS '用户头像图片';
COMMENT ON COLUMN user_info.user.cell_phone IS '手机号码';
COMMENT ON COLUMN user_info.user.salt IS '密码加盐';
COMMENT ON COLUMN user_info.user.password IS '登录密码';
COMMENT ON COLUMN user_info.user.created IS '创建时间';
COMMENT ON COLUMN user_info.user.creator IS '创建人';
COMMENT ON COLUMN user_info.user.edited IS '修改时间';
COMMENT ON COLUMN user_info.user.editor IS '修改人';
COMMENT ON COLUMN user_info.user.deleted IS '逻辑删除:true=删除,false=未删除';
-- 索引
CREATE INDEX user_id_index ON user_info.user(id);

-- 权限表:记录权限信息
CREATE TABLE IF NOT EXISTS user_info.permission(
    id serial PRIMARY KEY,
    parent_id integer,
    code varchar(255),
    name varchar(255),
    introduction varchar(255),
    category smallint,
    uri integer,
    created timestamptz,
    creator varchar(72),
    edited timestamptz,
    editor varchar(72),
    deleted bool DEFAULT FALSE
);
COMMENT ON TABLE user_info.permission IS '权限表:记录权限信息';

COMMENT ON COLUMN user_info.permission.id IS '权限ID';
COMMENT ON COLUMN user_info.permission.parent_id IS '所属父级权限ID';
COMMENT ON COLUMN user_info.permission.code IS '权限唯一CODE代码';
COMMENT ON COLUMN user_info.permission.name IS '权限名称';
COMMENT ON COLUMN user_info.permission.introduction IS '权限介绍';
COMMENT ON COLUMN user_info.permission.category IS '权限类别,1=编辑,2=查看';
COMMENT ON COLUMN user_info.permission.uri IS 'URL规则';
COMMENT ON COLUMN user_info.permission.created IS '创建时间';
COMMENT ON COLUMN user_info.permission.creator IS '创建人';
COMMENT ON COLUMN user_info.permission.edited IS '修改时间';
COMMENT ON COLUMN user_info.permission.editor IS '修改人';
COMMENT ON COLUMN user_info.permission.deleted IS '逻辑删除:true=删除,false=未删除';
-- 索引
CREATE INDEX permission_id_index ON user_info.permission(id);
CREATE INDEX permission_parent_id_index ON user_info.permission(parent_id);
CREATE INDEX permission_code_index ON user_info.permission(code);

-- 角色表:记录角色信息，即定义权限组
CREATE TABLE IF NOT EXISTS user_info.role(
    id serial PRIMARY KEY,
    parent_id integer,
    code varchar(255),
    name varchar(255),
    introduction varchar(255),
    created timestamptz,
    creator varchar(72),
    edited timestamptz,
    editor varchar(72),
    deleted bool DEFAULT FALSE
);

COMMENT ON TABLE user_info.role IS '角色表:记录角色信息，即定义权限组';

COMMENT ON COLUMN user_info.role.id IS '角色ID';
COMMENT ON COLUMN user_info.role.parent_id IS '所属父级角色ID';
COMMENT ON COLUMN user_info.role.code IS '角色唯一CODE代码';
COMMENT ON COLUMN user_info.role.name IS '角色名称';
COMMENT ON COLUMN user_info.role.introduction IS '角色介绍';
COMMENT ON COLUMN user_info.role.created IS '创建时间';
COMMENT ON COLUMN user_info.role.creator IS '创建人';
COMMENT ON COLUMN user_info.role.edited IS '修改时间';
COMMENT ON COLUMN user_info.role.editor IS '修改人';
COMMENT ON COLUMN user_info.role.deleted IS '逻辑删除:true=删除,false=未删除';
-- 索引
CREATE INDEX role_id_index ON user_info.role(id);
CREATE INDEX role_parent_id_index ON user_info.role(parent_id);
CREATE INDEX role_code_index ON user_info.role(code);

-- 用户-角色关联表:记录每个用户拥有哪些角色信息
CREATE TABLE IF NOT EXISTS user_info.user_role(
    id serial PRIMARY KEY,
    user_id integer,
    role_id integer,
    created timestamptz,
    creator varchar(72),
    edited timestamptz,
    editor varchar(72),
    deleted bool DEFAULT FALSE
);
-- 注释
COMMENT ON TABLE user_info.user_role IS '用户-角色关联表:记录每个用户拥有哪些角色信息';

COMMENT ON COLUMN user_info.user_role.id IS 'ID';
COMMENT ON COLUMN user_info.user_role.user_id IS '用户ID';
COMMENT ON COLUMN user_info.user_role.role_id IS '角色ID';
COMMENT ON COLUMN user_info.user_role.created IS '创建时间';
COMMENT ON COLUMN user_info.user_role.creator IS '创建人';
COMMENT ON COLUMN user_info.user_role.edited IS '修改时间';
COMMENT ON COLUMN user_info.user_role.editor IS '修改人';
COMMENT ON COLUMN user_info.user_role.deleted IS '逻辑删除:true=删除,false=未删除';
-- 索引
CREATE INDEX user_role_id_index ON user_info.user_role(id);
CREATE INDEX user_role_user_id_index ON user_info.user_role(user_id);
CREATE INDEX user_role_role_id_index ON user_info.user_role(role_id);


-- 角色-权限关联表:记录每个角色拥有哪些权限信息
CREATE TABLE IF NOT EXISTS user_info.role_permission(
    id serial PRIMARY KEY,
    role_id integer,
    permission_id integer,
    created timestamptz,
    creator varchar(72),
    edited timestamptz,
    editor varchar(72),
    deleted bool DEFAULT FALSE
);
--注释
COMMENT ON TABLE user_info.role_permission IS '角色-权限关联表:记录每个角色拥有哪些权限信息';

COMMENT ON COLUMN user_info.role_permission.id IS 'ID';
COMMENT ON COLUMN user_info.role_permission.role_id IS '角色ID';
COMMENT ON COLUMN user_info.role_permission.permission_id IS '权限ID';
COMMENT ON COLUMN user_info.role_permission.created IS '创建时间';
COMMENT ON COLUMN user_info.role_permission.creator IS '创建人';
COMMENT ON COLUMN user_info.role_permission.edited IS '修改时间';
COMMENT ON COLUMN user_info.role_permission.editor IS '修改人';
COMMENT ON COLUMN user_info.role_permission.deleted IS '逻辑删除:0=未删除,1=已删除';
-- 索引
CREATE INDEX role_permission_id_index ON user_info.role_permission(id);
CREATE INDEX role_permission_role_id_index ON user_info.role_permission(role_id);
CREATE INDEX role_permission_permission_id_index ON user_info.role_permission(permission_id);