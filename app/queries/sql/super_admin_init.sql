-- 创建一个用户
INSERT INTO user_info.user(id,name,cell_phone,salt,password,created,creator,edited,editor,deleted)
VALUES (0,'Alan','15600755813','','admin123',now(),'Alan',now(),'Alan',false);

INSERT INTO user_info.account(id,user_id,login_account,category,created,creator,edited,editor,deleted)
VALUES (0,0,'admin',1,now(),'Alan',now(),'Alan',false);

-- 模块权限
INSERT INTO user_info.permission(id,parent_id, code, name, category,created,creator,edited,editor,deleted)
VALUES (0,-1,'admin','用户管理',1,now(),'Alan',now(),'Alan',false);

INSERT INTO user_info.permission(id,parent_id, code, name, category,created,creator,edited,editor,deleted)
VALUES (1,-1,'sources','源数据',1,now(),'Alan',now(),'Alan',false);

INSERT INTO user_info.permission(id,parent_id, code, name, category,created,creator,edited,editor,deleted)
VALUES (2,-1,'services','服务管理',1,now(),'Alan',now(),'Alan',false);

-- 创建超级管理员、管理员、普通用户和游客的角色
INSERT INTO user_info.role(id,parent_id,code,name,created,creator,edited,editor,deleted)
VALUES (0,-1,'super_admin','超级管理员',now(),'Alan',now(),'Alan',false);

INSERT INTO user_info.role(id,parent_id,code,name,created,creator,edited,editor,deleted)
VALUES (1,0,'admin','管理员',now(),'Alan',now(),'Alan',false);

INSERT INTO user_info.role(id,parent_id,code,name,created,creator,edited,editor,deleted)
VALUES (2,1,'ordinary_user','普通用户',now(),'Alan',now(),'Alan',false);

INSERT INTO user_info.role(id,parent_id,code,name,created,creator,edited,editor,deleted)
VALUES (3,2,'tourist','游客',now(),'Alan',now(),'Alan',false);

-- 用户角色:设置一个超级管理员的角色的用户
INSERT INTO user_info.user_role(id,user_id,role_id,created,creator,edited,editor,deleted)
VALUES (0,0,0,now(),'Alan',now(),'Alan',false);

-- 角色权限设置
INSERT INTO user_info.role_permission(id,role_id,permission_id,created,creator,edited,editor,deleted)
VALUES (0,0,0,now(),'Alan',now(),'Alan',false);

INSERT INTO user_info.role_permission(id,role_id,permission_id,created,creator,edited,editor,deleted)
VALUES (1,0,1,now(),'Alan',now(),'Alan',false);

INSERT INTO user_info.role_permission(id,role_id,permission_id,created,creator,edited,editor,deleted)
VALUES (2,0,2,now(),'Alan',now(),'Alan',false);