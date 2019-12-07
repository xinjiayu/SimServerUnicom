CREATE TABLE `sim_unicom`
(
    `id`                int(11)      NOT NULL,
    `iccid`             varchar(100) NOT NULL,
    `imsi`              varchar(100)          DEFAULT NULL,
    `msisdn`            varchar(100)          DEFAULT NULL,
    `imei`              varchar(100)          DEFAULT NULL,
    `status`            varchar(100)          DEFAULT NULL,
    `rateplan`          varchar(100)          DEFAULT NULL,
    `communicationplan` varchar(100)          DEFAULT NULL,
    `customer`          varchar(100)          DEFAULT NULL,
    `endconsumerid`     varchar(100)          DEFAULT NULL,
    `dateactivated`     varchar(100)          DEFAULT NULL,
    `dateadded`         varchar(100)          DEFAULT NULL,
    `dateupdated`       varchar(100)          DEFAULT NULL,
    `dateshipped`       varchar(100)          DEFAULT NULL,
    `accountid`         varchar(100)          DEFAULT NULL,
    `fixedipaddress`    varchar(100)          DEFAULT NULL,
    `accountcustom1`    varchar(100)          DEFAULT NULL,
    `accountcustom2`    varchar(100)          DEFAULT NULL,
    `accountcustom3`    varchar(100)          DEFAULT NULL,
    `accountcustom4`    varchar(100)          DEFAULT NULL,
    `accountcustom5`    varchar(100)          DEFAULT NULL,
    `accountcustom6`    varchar(100)          DEFAULT NULL,
    `accountcustom7`    varchar(100)          DEFAULT NULL,
    `accountcustom8`    varchar(100)          DEFAULT NULL,
    `accountcustom9`    varchar(100)          DEFAULT NULL,
    `accountcustom10`   varchar(100)          DEFAULT NULL,
    `simnotes`          varchar(100)          DEFAULT NULL,
    `deviceid`          varchar(100)          DEFAULT NULL,
    `modemid`           varchar(100)          DEFAULT NULL,
    `globalsimtype`     varchar(100)          DEFAULT NULL,
    `ctddatausage`      bigint(19)   NOT NULL DEFAULT '0' COMMENT '本期使用流量'
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8 COMMENT ='联通sim卡信息表';

ALTER TABLE `sim_unicom`
    ADD PRIMARY KEY (`id`, `iccid`) USING BTREE;

ALTER TABLE `sim_unicom`
    MODIFY `id` int(11) NOT NULL AUTO_INCREMENT,
    AUTO_INCREMENT = 1;



CREATE TABLE `sim_flow`
(
    `id`                int(11)             NOT NULL COMMENT 'ID',
    `iccid`             varchar(30)         NOT NULL COMMENT 'sim卡号',
    `year`              int(4)              NOT NULL COMMENT '年份',
    `month`             tinyint(2) UNSIGNED NOT NULL DEFAULT '0' COMMENT '月份',
    `provider`          varchar(30)         NOT NULL DEFAULT '' COMMENT '提供商',
    `carrier`           varchar(30)         NOT NULL DEFAULT '' COMMENT '运营商',
    `d1`                bigint(19) UNSIGNED          DEFAULT '0' COMMENT '流量',
    `d2`                bigint(19) UNSIGNED          DEFAULT '0' COMMENT '流量',
    `d3`                bigint(19) UNSIGNED          DEFAULT '0' COMMENT '流量',
    `d4`                bigint(19) UNSIGNED          DEFAULT '0' COMMENT '流量',
    `d5`                bigint(19) UNSIGNED          DEFAULT '0' COMMENT '流量',
    `d6`                bigint(19) UNSIGNED          DEFAULT '0' COMMENT '流量',
    `d7`                bigint(19) UNSIGNED          DEFAULT '0' COMMENT '流量',
    `d8`                bigint(19) UNSIGNED          DEFAULT '0' COMMENT '流量',
    `d9`                bigint(19) UNSIGNED          DEFAULT '0' COMMENT '流量',
    `d10`               bigint(19) UNSIGNED          DEFAULT '0' COMMENT '流量',
    `d11`               bigint(19) UNSIGNED          DEFAULT '0' COMMENT '流量',
    `d12`               bigint(19) UNSIGNED          DEFAULT '0' COMMENT '流量',
    `d13`               bigint(19) UNSIGNED          DEFAULT '0' COMMENT '流量',
    `d14`               bigint(19) UNSIGNED          DEFAULT '0' COMMENT '流量',
    `d15`               bigint(19) UNSIGNED          DEFAULT '0' COMMENT '流量',
    `d16`               bigint(19) UNSIGNED          DEFAULT '0' COMMENT '流量',
    `d17`               bigint(19) UNSIGNED          DEFAULT '0' COMMENT '流量',
    `d18`               bigint(19) UNSIGNED          DEFAULT '0' COMMENT '流量',
    `d19`               bigint(19) UNSIGNED          DEFAULT '0' COMMENT '流量',
    `d20`               bigint(19) UNSIGNED          DEFAULT '0' COMMENT '流量',
    `d21`               bigint(19) UNSIGNED          DEFAULT '0' COMMENT '流量',
    `d22`               bigint(19) UNSIGNED          DEFAULT '0' COMMENT '流量',
    `d23`               bigint(19) UNSIGNED          DEFAULT '0' COMMENT '流量',
    `d24`               bigint(19) UNSIGNED          DEFAULT '0' COMMENT '流量',
    `d25`               bigint(19) UNSIGNED          DEFAULT '0' COMMENT '流量',
    `d26`               bigint(19) UNSIGNED          DEFAULT '0' COMMENT '流量',
    `d27`               bigint(19) UNSIGNED          DEFAULT '0' COMMENT '流量',
    `d28`               bigint(19) UNSIGNED          DEFAULT '0' COMMENT '流量',
    `d29`               bigint(19) UNSIGNED          DEFAULT '0' COMMENT '流量',
    `d30`               bigint(19) UNSIGNED          DEFAULT '0' COMMENT '流量',
    `d31`               bigint(19) UNSIGNED          DEFAULT '0' COMMENT '流量',
    `createtime`        int(10)             NOT NULL DEFAULT '0' COMMENT '创建时间',
    `updatetime`        int(10) UNSIGNED    NOT NULL DEFAULT '0' COMMENT '更新时间',
    `ctd_session_count` int(11)             NOT NULL COMMENT '数据会话次数',
    `rate_plan`         varchar(100)        NOT NULL COMMENT '套餐计划名称',
    `status`            varchar(50)         NOT NULL COMMENT '状态'
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8 COMMENT ='日流量表'
  ROW_FORMAT = COMPACT;


ALTER TABLE `sim_flow`
    ADD PRIMARY KEY (`id`) USING BTREE,
    ADD UNIQUE KEY `iccid` (`iccid`, `month`, `year`) USING BTREE;

ALTER TABLE `sim_flow`
    MODIFY `id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'ID',
    AUTO_INCREMENT = 1;


CREATE TABLE `unicom_notice`
(
    `id`         int(10)      NOT NULL COMMENT 'ID',
    `event_id`   varchar(100) NOT NULL DEFAULT '' COMMENT '标识',
    `event_type` varchar(50)  NOT NULL DEFAULT '' COMMENT '类型',
    `timestamp`  int(10)      NOT NULL DEFAULT '0' COMMENT '请求时间',
    `iccid`      varchar(50)  NOT NULL DEFAULT '' COMMENT 'ICCID',
    `data_usage` int(11)      NOT NULL DEFAULT '0' COMMENT '流量',
    `data1`      varchar(50)  NOT NULL DEFAULT '' COMMENT '数据一',
    `data2`      varchar(50)  NOT NULL DEFAULT '' COMMENT '数据二'
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8 COMMENT ='联通平台流量通知';

ALTER TABLE `unicom_notice`
    ADD PRIMARY KEY (`id`);

ALTER TABLE `unicom_notice`
    MODIFY `id` int(10) NOT NULL AUTO_INCREMENT COMMENT 'ID', AUTO_INCREMENT=1;