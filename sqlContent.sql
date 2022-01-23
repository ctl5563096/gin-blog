// 新博客新建表的sql语句
CREATE TABLE `t_go_article` (
    `id`          int(11)                           NOT NULL            AUTO_INCREMENT,
    `title`       varchar(30)                       NOT NULL DEFAULT '' COMMENT '文章标题',
    `summary`     varchar(1024)                     NOT NULL DEFAULT '' COMMENT '文章简介',
    `cover`       varchar(1024)                     NOT NULL DEFAULT '' COMMENT '封面图',
    `is_show`     tinyint(1)                        NOT NULL DEFAULT 1  COMMENT '是否展示,0是不展示,1是展示，默认为展示1',
    `is_delete`   tinyint(1)                        NOT NULL DEFAULT 0  COMMENT '是否删除文章,0是未删除,1是删除，默认为未删除0',
    `content`     text COLLATE utf8mb4_unicode_ci   NOT NULL DEFAULT '' COMMENT '文章内容',
    `view`        int(11)                           NOT NULL DEFAULT 0  COMMENT '观看文章人数',
    `hate`        int(11)                           NOT NULL DEFAULT 0  COMMENT '点踩数',
    `like`        int(11)                           NOT NULL DEFAULT 0  COMMENT '喜欢文章人数',
    `author`      varchar(30)                       NOT NULL DEFAULT '' COMMENT '作者',
    `created_at`  timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at`  timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `ext`         text COMMENT '扩展字段，json结构的字符串，一些纯展示的字段，可以作为此字段的内部字段，如 {"property":"xxx", ...}',
     PRIMARY KEY (`id`),
     UNIQUE KEY `titleIndex` (`title`),
     KEY `createdAt` (`created_at`)
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8mb4 COMMENT='go 博客重写文章表';


// 权限菜单表
CREATE TABLE `t_go_rule` (
     `id`              int(11)       unsigned NOT NULL AUTO_INCREMENT,
     `rule_name`       varchar(1024)          NOT NULL DEFAULT ''                COMMENT '权限名字',
     `pid`             int(11)       unsigned NOT NULL DEFAULT '0'               COMMENT '父权限Id',
     `is_menu`         tinyint(1)    unsigned NOT NULL DEFAULT '0'               COMMENT '是否为菜单',
     `url`             varchar(255)           NOT NULL DEFAULT ''                COMMENT '菜单路径',
     `icon`            varchar(255)           NOT NULL DEFAULT ''                COMMENT '菜单图标',
     `controller`      varchar(255)           NOT NULL DEFAULT ''                COMMENT '控制器',
     `action`          varchar(255)           NOT NULL DEFAULT ''                COMMENT '方法',
     `sort`            int(11)                NOT NULL DEFAULT '0'               COMMENT '排序',
     `created_at`      timestamp              NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
     `updated_at`      timestamp              NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
     PRIMARY KEY (`id`)
     KEY `sort` (`sort`)
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8mb4 COMMENT='go权限菜单表';

// 角色表
CREATE TABLE `t_go_role` (
     `id`              int(11)     unsigned   NOT NULL AUTO_INCREMENT,
     `role_name`       varchar(256)           NOT NULL DEFAULT '' COMMENT '角色名',
     `is_enabled`      tinyint(1)  unsigned   NOT NULL DEFAULT '0' COMMENT '是否启用',
     `created_at`      timestamp              NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
     `updated_at`      timestamp              NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
     PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8 COMMENT='角色表';


// 角色权限表
CREATE TABLE `t_go_role_rule` (
     `id`              int(11)       unsigned NOT NULL AUTO_INCREMENT,
     `role`            int(11)       unsigned NOT NULL DEFAULT '0' COMMENT '角色ID',
     `rule`            int(11)       unsigned NOT NULL DEFAULT '0' COMMENT '权限ID',
     `created_at`      timestamp              NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
     `updated_at`      timestamp              NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
     PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8 COMMENT='角色权限对应表';

// 本地资源和阿里云oss
CREATE TABLE `t_go_local_oss` (
      `id`              int(11)       unsigned NOT NULL AUTO_INCREMENT,
      `local`           varchar(256)           NOT NULL DEFAULT ''  COMMENT '本地地址',
      `oss`             varchar(256)           NOT NULL DEFAULT '' COMMENT 'oss地址',
      `created_at`      timestamp              NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
      `updated_at`      timestamp              NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
      PRIMARY KEY (`id`)
      KEY `local` (`local`)
      KEY `oss` (`oss`)
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8 COMMENT='阿里云与本地文件的对应关系,防止费用太贵导致用不起';

CREATE TABLE `t_go_parameter` (
    `id`            int(11) unsigned NOT NULL AUTO_INCREMENT,
    `name`          varchar(128)     NOT NULL DEFAULT '' COMMENT '名称',
    `code`          int(11)          NOT NULL COMMENT '参数编码',
    `param_name`    varchar(128)     NOT NULL DEFAULT '' COMMENT '参数名称',
    `param_value`   varchar(128)     NOT NULL DEFAULT '' COMMENT '参数值',
    `weight`        int(11)          NOT NULL DEFAULT '0' COMMENT '参数值',
    `is_enabled`    tinyint(1)       NOT NULL DEFAULT '1' COMMENT '是否启用参数 0为否 1为是',
    `created_at`    timestamp        NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at`    timestamp        NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '修改时间',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8 COMMENT='系统参数表';

// 资源和标签关系表
CREATE TABLE `t_go_resource_tags_relation` (
    `id`            int(11) unsigned NOT NULL AUTO_INCREMENT,
    `resource_id`   int(11)          NOT NULL DEFAULT 0 COMMENT '资源id',
    `resource_type` varchar(256)     NOT NULL DEFAULT '' COMMENT '相关的资源类型',
    `code`          int(11)          NOT NULL COMMENT '参数编码',
    `param_value`   int(11)          NOT NULL DEFAULT 0 COMMENT '参数值',
    `tag_style`     varchar(256)     NOT NULL DEFAULT ''  COMMENT 'tag样式',
    `created_at`    timestamp        NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at`    timestamp        NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '修改时间',
    PRIMARY KEY (`id`),
    KEY `t_resource_tag` (`resource_id`,`resource_type`,`code`)
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8 COMMENT='系统资源参数关系表';

CREATE TABLE `t_go_resource_tags_relation` (
   `id`             int(11)     unsigned NOT NULL AUTO_INCREMENT,
   `resource_id`    int(11)              NOT NULL DEFAULT '0' COMMENT '资源id',
   `resource_type`  varchar(256)         NOT NULL DEFAULT '' COMMENT '相关的资源类型',
   `code`           varchar(11)          NOT NULL DEFAULT '' COMMENT '参数编码',
   `param_value`    int(11)              NOT NULL DEFAULT '0' COMMENT '参数值',
   `tag_style`      varchar(30)          NOT NULL DEFAULT 'plain' COMMENT 'tag样式',
   `created_at`     timestamp            NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
   `updated_at`     timestamp            NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '修改时间',
   PRIMARY KEY (`id`),
   KEY `t_resource_tag` (`resource_id`,`resource_type`,`code`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8 COMMENT='系统资源参数关系表';

// 热门资源表
CREATE TABLE `t_go_hot_content` (
   `id`             int(11)     unsigned NOT NULL AUTO_INCREMENT,
   `resource_id`    int(11)              NOT NULL DEFAULT '0' COMMENT '资源id',
   `resource_type`  int(11)              NOT NULL COMMENT '资源类型',
   `title`          varchar(256)         NOT NULL DEFAULT ''  COMMENT '资源标题',
   `is_top`         tinyint(1)           NOT NULL DEFAULT 1   COMMENT '是否为置顶资源 1为否 2为是',
   `cover`          varchar(256)         NOT NULL DEFAULT ''  COMMENT '资源封面图',
   `is_delete`      tinyint(1)           NOT NULL DEFAULT 1   COMMENT '热门资源是否被删除 1为否 2为是',
   `created_at`     timestamp            NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
   `updated_at`     timestamp            NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '修改时间',
   PRIMARY KEY (`id`),
   KEY `resource_id` (`resource_id`)
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8 COMMENT='热门资源';

// 音乐资源表
CREATE TABLE `t_go_music` (
    `id`             int(11)     unsigned NOT NULL AUTO_INCREMENT,
    `title`          varchar(256)         NOT NULL DEFAULT ''  COMMENT '音乐标题',
    `summary`        varchar(1024)        NOT NULL DEFAULT ''  COMMENT '音乐简介',
    `lyric`          text COLLATE utf8mb4_unicode_ci   NOT NULL DEFAULT '' COMMENT '文章内容',
    `is_top`         tinyint(1)           NOT NULL DEFAULT 1   COMMENT '是否为置顶资源 1为否 2为是',
    `thumb`          varchar(256)         NOT NULL DEFAULT ''  COMMENT '资源缩略图',
    `cover`          varchar(256)         NOT NULL DEFAULT ''  COMMENT '资源封面图',
    `is_delete`      tinyint(1)           NOT NULL DEFAULT 1   COMMENT '是否被删除 1为否 2为是',
    `created_at`     timestamp            NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at`     timestamp            NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '修改时间',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=0 DEFAULT CHARSET=utf8mb4 COMMENT='收藏音乐表';
