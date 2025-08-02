/*
 Navicat Premium Dump SQL

 Source Server         : 233
 Source Server Type    : MySQL
 Source Server Version : 80035 (8.0.35)
 Source Host           : 192.168.60.233:3306
 Source Schema         : xxl_job

 Target Server Type    : MySQL
 Target Server Version : 80035 (8.0.35)
 File Encoding         : 65001

 Date: 02/08/2025 13:03:02
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for xxl_job_group
-- ----------------------------
DROP TABLE IF EXISTS `xxl_job_group`;
CREATE TABLE `xxl_job_group` (
  `id` int NOT NULL AUTO_INCREMENT,
  `app_name` varchar(64) NOT NULL COMMENT '执行器AppName',
  `title` varchar(12) NOT NULL COMMENT '执行器名称',
  `address_type` tinyint NOT NULL DEFAULT '0' COMMENT '执行器地址类型：0=自动注册、1=手动录入',
  `address_list` text COMMENT '执行器地址列表，多地址逗号分隔',
  `update_time` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Records of xxl_job_group
-- ----------------------------
BEGIN;
INSERT INTO `xxl_job_group` (`id`, `app_name`, `title`, `address_type`, `address_list`, `update_time`) VALUES (1, 'xxl-job-executor-sample', '通用执行器Sample', 0, NULL, '2025-08-02 14:02:57');
INSERT INTO `xxl_job_group` (`id`, `app_name`, `title`, `address_type`, `address_list`, `update_time`) VALUES (2, 'xxl-job-executor-sample-ai', 'AI执行器Sample', 0, NULL, '2025-08-02 14:02:57');
INSERT INTO `xxl_job_group` (`id`, `app_name`, `title`, `address_type`, `address_list`, `update_time`) VALUES (3, 'worldExecutor', 'coder任务执行器', 1, 'http://127.0.0.1:9999', '2025-08-02 11:32:25');
COMMIT;

-- ----------------------------
-- Table structure for xxl_job_info
-- ----------------------------
DROP TABLE IF EXISTS `xxl_job_info`;
CREATE TABLE `xxl_job_info` (
  `id` int NOT NULL AUTO_INCREMENT,
  `job_group` int NOT NULL COMMENT '执行器主键ID',
  `job_desc` varchar(255) NOT NULL,
  `add_time` datetime DEFAULT NULL,
  `update_time` datetime DEFAULT NULL,
  `author` varchar(64) DEFAULT NULL COMMENT '作者',
  `alarm_email` varchar(255) DEFAULT NULL COMMENT '报警邮件',
  `schedule_type` varchar(50) NOT NULL DEFAULT 'NONE' COMMENT '调度类型',
  `schedule_conf` varchar(128) DEFAULT NULL COMMENT '调度配置，值含义取决于调度类型',
  `misfire_strategy` varchar(50) NOT NULL DEFAULT 'DO_NOTHING' COMMENT '调度过期策略',
  `executor_route_strategy` varchar(50) DEFAULT NULL COMMENT '执行器路由策略',
  `executor_handler` varchar(255) DEFAULT NULL COMMENT '执行器任务handler',
  `executor_param` varchar(512) DEFAULT NULL COMMENT '执行器任务参数',
  `executor_block_strategy` varchar(50) DEFAULT NULL COMMENT '阻塞处理策略',
  `executor_timeout` int NOT NULL DEFAULT '0' COMMENT '任务执行超时时间，单位秒',
  `executor_fail_retry_count` int NOT NULL DEFAULT '0' COMMENT '失败重试次数',
  `glue_type` varchar(50) NOT NULL COMMENT 'GLUE类型',
  `glue_source` mediumtext COMMENT 'GLUE源代码',
  `glue_remark` varchar(128) DEFAULT NULL COMMENT 'GLUE备注',
  `glue_updatetime` datetime DEFAULT NULL COMMENT 'GLUE更新时间',
  `child_jobid` varchar(255) DEFAULT NULL COMMENT '子任务ID，多个逗号分隔',
  `trigger_status` tinyint NOT NULL DEFAULT '0' COMMENT '调度状态：0-停止，1-运行',
  `trigger_last_time` bigint NOT NULL DEFAULT '0' COMMENT '上次调度时间',
  `trigger_next_time` bigint NOT NULL DEFAULT '0' COMMENT '下次调度时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Records of xxl_job_info
-- ----------------------------
BEGIN;
INSERT INTO `xxl_job_info` (`id`, `job_group`, `job_desc`, `add_time`, `update_time`, `author`, `alarm_email`, `schedule_type`, `schedule_conf`, `misfire_strategy`, `executor_route_strategy`, `executor_handler`, `executor_param`, `executor_block_strategy`, `executor_timeout`, `executor_fail_retry_count`, `glue_type`, `glue_source`, `glue_remark`, `glue_updatetime`, `child_jobid`, `trigger_status`, `trigger_last_time`, `trigger_next_time`) VALUES (1, 1, '示例任务01', '2025-08-02 10:25:40', '2025-08-02 10:25:40', 'XXL', '', 'CRON', '0 0 0 * * ? *', 'DO_NOTHING', 'FIRST', 'demoJobHandler', '', 'SERIAL_EXECUTION', 0, 0, 'BEAN', '', 'GLUE代码初始化', '2025-08-02 10:25:40', '', 0, 0, 0);
INSERT INTO `xxl_job_info` (`id`, `job_group`, `job_desc`, `add_time`, `update_time`, `author`, `alarm_email`, `schedule_type`, `schedule_conf`, `misfire_strategy`, `executor_route_strategy`, `executor_handler`, `executor_param`, `executor_block_strategy`, `executor_timeout`, `executor_fail_retry_count`, `glue_type`, `glue_source`, `glue_remark`, `glue_updatetime`, `child_jobid`, `trigger_status`, `trigger_last_time`, `trigger_next_time`) VALUES (2, 2, 'Ollama示例任务01', '2025-08-02 10:25:40', '2025-08-02 10:25:40', 'XXL', '', 'NONE', '', 'DO_NOTHING', 'FIRST', 'ollamaJobHandler', '{\n    \"input\": \"慢SQL问题分析思路\",\n    \"prompt\": \"你是一个研发工程师，擅长解决技术类问题。\"\n}', 'SERIAL_EXECUTION', 0, 0, 'BEAN', '', 'GLUE代码初始化', '2025-08-02 10:25:40', '', 0, 0, 0);
INSERT INTO `xxl_job_info` (`id`, `job_group`, `job_desc`, `add_time`, `update_time`, `author`, `alarm_email`, `schedule_type`, `schedule_conf`, `misfire_strategy`, `executor_route_strategy`, `executor_handler`, `executor_param`, `executor_block_strategy`, `executor_timeout`, `executor_fail_retry_count`, `glue_type`, `glue_source`, `glue_remark`, `glue_updatetime`, `child_jobid`, `trigger_status`, `trigger_last_time`, `trigger_next_time`) VALUES (3, 2, 'Dify示例任务', '2025-08-02 10:25:40', '2025-08-02 10:25:40', 'XXL', '', 'NONE', '', 'DO_NOTHING', 'FIRST', 'difyWorkflowJobHandler', '{\n    \"inputs\":{\n        \"input\":\"查询班级各学科前三名\"\n    },\n    \"user\": \"xxl-job\",\n    \"baseUrl\": \"http://localhost/v1\",\n    \"apiKey\": \"app-OUVgNUOQRIMokfmuJvBJoUTN\"\n}', 'SERIAL_EXECUTION', 0, 0, 'BEAN', '', 'GLUE代码初始化', '2025-08-02 10:25:40', '', 0, 0, 0);
INSERT INTO `xxl_job_info` (`id`, `job_group`, `job_desc`, `add_time`, `update_time`, `author`, `alarm_email`, `schedule_type`, `schedule_conf`, `misfire_strategy`, `executor_route_strategy`, `executor_handler`, `executor_param`, `executor_block_strategy`, `executor_timeout`, `executor_fail_retry_count`, `glue_type`, `glue_source`, `glue_remark`, `glue_updatetime`, `child_jobid`, `trigger_status`, `trigger_last_time`, `trigger_next_time`) VALUES (4, 3, '状态检查任务', '2025-08-02 11:09:22', '2025-08-02 11:39:50', 'coder', '', 'CRON', '0 0/5 * * * ?', 'DO_NOTHING', 'FIRST', 'jobGenerateIndex', '', 'SERIAL_EXECUTION', 0, 0, 'BEAN', '', 'GLUE代码初始化', '2025-08-02 11:09:22', '', 1, 1754114400000, 1754114700000);
COMMIT;

-- ----------------------------
-- Table structure for xxl_job_lock
-- ----------------------------
DROP TABLE IF EXISTS `xxl_job_lock`;
CREATE TABLE `xxl_job_lock` (
  `lock_name` varchar(50) NOT NULL COMMENT '锁名称',
  PRIMARY KEY (`lock_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Records of xxl_job_lock
-- ----------------------------
BEGIN;
INSERT INTO `xxl_job_lock` (`lock_name`) VALUES ('schedule_lock');
COMMIT;

-- ----------------------------
-- Table structure for xxl_job_log
-- ----------------------------
DROP TABLE IF EXISTS `xxl_job_log`;
CREATE TABLE `xxl_job_log` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `job_group` int NOT NULL COMMENT '执行器主键ID',
  `job_id` int NOT NULL COMMENT '任务，主键ID',
  `executor_address` varchar(255) DEFAULT NULL COMMENT '执行器地址，本次执行的地址',
  `executor_handler` varchar(255) DEFAULT NULL COMMENT '执行器任务handler',
  `executor_param` varchar(512) DEFAULT NULL COMMENT '执行器任务参数',
  `executor_sharding_param` varchar(20) DEFAULT NULL COMMENT '执行器任务分片参数，格式如 1/2',
  `executor_fail_retry_count` int NOT NULL DEFAULT '0' COMMENT '失败重试次数',
  `trigger_time` datetime DEFAULT NULL COMMENT '调度-时间',
  `trigger_code` int NOT NULL COMMENT '调度-结果',
  `trigger_msg` text COMMENT '调度-日志',
  `handle_time` datetime DEFAULT NULL COMMENT '执行-时间',
  `handle_code` int NOT NULL COMMENT '执行-状态',
  `handle_msg` text COMMENT '执行-日志',
  `alarm_status` tinyint NOT NULL DEFAULT '0' COMMENT '告警状态：0-默认、1-无需告警、2-告警成功、3-告警失败',
  PRIMARY KEY (`id`),
  KEY `I_trigger_time` (`trigger_time`),
  KEY `I_handle_code` (`handle_code`),
  KEY `I_jobid_jobgroup` (`job_id`,`job_group`),
  KEY `I_job_id` (`job_id`)
) ENGINE=InnoDB AUTO_INCREMENT=41 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Records of xxl_job_log
-- ----------------------------
BEGIN;
INSERT INTO `xxl_job_log` (`id`, `job_group`, `job_id`, `executor_address`, `executor_handler`, `executor_param`, `executor_sharding_param`, `executor_fail_retry_count`, `trigger_time`, `trigger_code`, `trigger_msg`, `handle_time`, `handle_code`, `handle_msg`, `alarm_status`) VALUES (1, 1, 4, NULL, 'jobGenerateIndex', '', NULL, 0, '2025-08-02 11:12:00', 500, '任务触发类型：Cron触发<br>调度机器：192.168.60.233<br>执行器-注册方式：自动注册<br>执行器-地址列表：null<br>路由策略：第一个<br>阻塞处理策略：单机串行<br>任务超时时间：0<br>失败重试次数：0<br><br><span style=\"color:#00c0ef;\" > >>>>>>>>>>>触发调度<<<<<<<<<<< </span><br>调度失败：执行器地址为空<br><br>', NULL, 0, NULL, 2);
INSERT INTO `xxl_job_log` (`id`, `job_group`, `job_id`, `executor_address`, `executor_handler`, `executor_param`, `executor_sharding_param`, `executor_fail_retry_count`, `trigger_time`, `trigger_code`, `trigger_msg`, `handle_time`, `handle_code`, `handle_msg`, `alarm_status`) VALUES (2, 1, 4, NULL, 'jobGenerateIndex', '', NULL, 0, '2025-08-02 11:15:00', 500, '任务触发类型：Cron触发<br>调度机器：192.168.60.233<br>执行器-注册方式：自动注册<br>执行器-地址列表：null<br>路由策略：第一个<br>阻塞处理策略：单机串行<br>任务超时时间：0<br>失败重试次数：0<br><br><span style=\"color:#00c0ef;\" > >>>>>>>>>>>触发调度<<<<<<<<<<< </span><br>调度失败：执行器地址为空<br><br>', NULL, 0, NULL, 2);
INSERT INTO `xxl_job_log` (`id`, `job_group`, `job_id`, `executor_address`, `executor_handler`, `executor_param`, `executor_sharding_param`, `executor_fail_retry_count`, `trigger_time`, `trigger_code`, `trigger_msg`, `handle_time`, `handle_code`, `handle_msg`, `alarm_status`) VALUES (3, 1, 4, NULL, 'jobGenerateIndex', '', NULL, 0, '2025-08-02 11:20:00', 500, '任务触发类型：Cron触发<br>调度机器：192.168.60.233<br>执行器-注册方式：自动注册<br>执行器-地址列表：null<br>路由策略：第一个<br>阻塞处理策略：单机串行<br>任务超时时间：0<br>失败重试次数：0<br><br><span style=\"color:#00c0ef;\" > >>>>>>>>>>>触发调度<<<<<<<<<<< </span><br>调度失败：执行器地址为空<br><br>', NULL, 0, NULL, 2);
INSERT INTO `xxl_job_log` (`id`, `job_group`, `job_id`, `executor_address`, `executor_handler`, `executor_param`, `executor_sharding_param`, `executor_fail_retry_count`, `trigger_time`, `trigger_code`, `trigger_msg`, `handle_time`, `handle_code`, `handle_msg`, `alarm_status`) VALUES (4, 1, 4, NULL, 'jobGenerateIndex', '', NULL, 0, '2025-08-02 11:25:00', 500, '任务触发类型：Cron触发<br>调度机器：192.168.60.233<br>执行器-注册方式：自动注册<br>执行器-地址列表：null<br>路由策略：第一个<br>阻塞处理策略：单机串行<br>任务超时时间：0<br>失败重试次数：0<br><br><span style=\"color:#00c0ef;\" > >>>>>>>>>>>触发调度<<<<<<<<<<< </span><br>调度失败：执行器地址为空<br><br>', NULL, 0, NULL, 2);
INSERT INTO `xxl_job_log` (`id`, `job_group`, `job_id`, `executor_address`, `executor_handler`, `executor_param`, `executor_sharding_param`, `executor_fail_retry_count`, `trigger_time`, `trigger_code`, `trigger_msg`, `handle_time`, `handle_code`, `handle_msg`, `alarm_status`) VALUES (5, 1, 4, NULL, 'jobGenerateIndex', '', NULL, 0, '2025-08-02 11:30:00', 500, '任务触发类型：Cron触发<br>调度机器：192.168.60.233<br>执行器-注册方式：自动注册<br>执行器-地址列表：null<br>路由策略：第一个<br>阻塞处理策略：单机串行<br>任务超时时间：0<br>失败重试次数：0<br><br><span style=\"color:#00c0ef;\" > >>>>>>>>>>>触发调度<<<<<<<<<<< </span><br>调度失败：执行器地址为空<br><br>', NULL, 0, NULL, 2);
INSERT INTO `xxl_job_log` (`id`, `job_group`, `job_id`, `executor_address`, `executor_handler`, `executor_param`, `executor_sharding_param`, `executor_fail_retry_count`, `trigger_time`, `trigger_code`, `trigger_msg`, `handle_time`, `handle_code`, `handle_msg`, `alarm_status`) VALUES (6, 3, 4, 'http://127.0.0.1:9999', 'jobGenerateIndex', '', NULL, 0, '2025-08-02 11:34:00', 500, '任务触发类型：Cron触发<br>调度机器：192.168.60.233<br>执行器-注册方式：手动录入<br>执行器-地址列表：[http://127.0.0.1:9999]<br>路由策略：第一个<br>阻塞处理策略：单机串行<br>任务超时时间：0<br>失败重试次数：0<br><br><span style=\"color:#00c0ef;\" > >>>>>>>>>>>触发调度<<<<<<<<<<< </span><br>触发调度：<br>address：http://127.0.0.1:9999<br>code：500<br>msg：xxl-job remoting error(拒绝连接), for url : http://127.0.0.1:9999/run', NULL, 0, NULL, 2);
INSERT INTO `xxl_job_log` (`id`, `job_group`, `job_id`, `executor_address`, `executor_handler`, `executor_param`, `executor_sharding_param`, `executor_fail_retry_count`, `trigger_time`, `trigger_code`, `trigger_msg`, `handle_time`, `handle_code`, `handle_msg`, `alarm_status`) VALUES (7, 3, 4, 'http://127.0.0.1:9999', 'jobGenerateIndex', '', NULL, 0, '2025-08-02 11:35:00', 500, '任务触发类型：Cron触发<br>调度机器：192.168.60.233<br>执行器-注册方式：手动录入<br>执行器-地址列表：[http://127.0.0.1:9999]<br>路由策略：第一个<br>阻塞处理策略：单机串行<br>任务超时时间：0<br>失败重试次数：0<br><br><span style=\"color:#00c0ef;\" > >>>>>>>>>>>触发调度<<<<<<<<<<< </span><br>触发调度：<br>address：http://127.0.0.1:9999<br>code：500<br>msg：xxl-job remoting error(拒绝连接), for url : http://127.0.0.1:9999/run', NULL, 0, NULL, 2);
INSERT INTO `xxl_job_log` (`id`, `job_group`, `job_id`, `executor_address`, `executor_handler`, `executor_param`, `executor_sharding_param`, `executor_fail_retry_count`, `trigger_time`, `trigger_code`, `trigger_msg`, `handle_time`, `handle_code`, `handle_msg`, `alarm_status`) VALUES (8, 3, 4, 'http://127.0.0.1:9999', 'jobGenerateIndex', '', NULL, 0, '2025-08-02 11:36:00', 500, '任务触发类型：Cron触发<br>调度机器：192.168.60.233<br>执行器-注册方式：手动录入<br>执行器-地址列表：[http://127.0.0.1:9999]<br>路由策略：第一个<br>阻塞处理策略：单机串行<br>任务超时时间：0<br>失败重试次数：0<br><br><span style=\"color:#00c0ef;\" > >>>>>>>>>>>触发调度<<<<<<<<<<< </span><br>触发调度：<br>address：http://127.0.0.1:9999<br>code：500<br>msg：xxl-job remoting error(拒绝连接), for url : http://127.0.0.1:9999/run', NULL, 0, NULL, 2);
INSERT INTO `xxl_job_log` (`id`, `job_group`, `job_id`, `executor_address`, `executor_handler`, `executor_param`, `executor_sharding_param`, `executor_fail_retry_count`, `trigger_time`, `trigger_code`, `trigger_msg`, `handle_time`, `handle_code`, `handle_msg`, `alarm_status`) VALUES (9, 3, 4, 'http://127.0.0.1:9999', 'jobGenerateIndex', '', NULL, 0, '2025-08-02 11:37:00', 500, '任务触发类型：Cron触发<br>调度机器：192.168.60.233<br>执行器-注册方式：手动录入<br>执行器-地址列表：[http://127.0.0.1:9999]<br>路由策略：第一个<br>阻塞处理策略：单机串行<br>任务超时时间：0<br>失败重试次数：0<br><br><span style=\"color:#00c0ef;\" > >>>>>>>>>>>触发调度<<<<<<<<<<< </span><br>触发调度：<br>address：http://127.0.0.1:9999<br>code：500<br>msg：xxl-job remoting error(拒绝连接), for url : http://127.0.0.1:9999/run', NULL, 0, NULL, 2);
INSERT INTO `xxl_job_log` (`id`, `job_group`, `job_id`, `executor_address`, `executor_handler`, `executor_param`, `executor_sharding_param`, `executor_fail_retry_count`, `trigger_time`, `trigger_code`, `trigger_msg`, `handle_time`, `handle_code`, `handle_msg`, `alarm_status`) VALUES (10, 3, 4, 'http://127.0.0.1:9999', 'jobGenerateIndex', '', NULL, 0, '2025-08-02 11:38:00', 500, '任务触发类型：Cron触发<br>调度机器：192.168.60.233<br>执行器-注册方式：手动录入<br>执行器-地址列表：[http://127.0.0.1:9999]<br>路由策略：第一个<br>阻塞处理策略：单机串行<br>任务超时时间：0<br>失败重试次数：0<br><br><span style=\"color:#00c0ef;\" > >>>>>>>>>>>触发调度<<<<<<<<<<< </span><br>触发调度：<br>address：http://127.0.0.1:9999<br>code：500<br>msg：xxl-job remoting error(拒绝连接), for url : http://127.0.0.1:9999/run', NULL, 0, NULL, 2);
INSERT INTO `xxl_job_log` (`id`, `job_group`, `job_id`, `executor_address`, `executor_handler`, `executor_param`, `executor_sharding_param`, `executor_fail_retry_count`, `trigger_time`, `trigger_code`, `trigger_msg`, `handle_time`, `handle_code`, `handle_msg`, `alarm_status`) VALUES (11, 3, 4, 'http://127.0.0.1:9999', 'jobGenerateIndex', '', NULL, 0, '2025-08-02 11:39:00', 500, '任务触发类型：Cron触发<br>调度机器：192.168.60.233<br>执行器-注册方式：手动录入<br>执行器-地址列表：[http://127.0.0.1:9999]<br>路由策略：第一个<br>阻塞处理策略：单机串行<br>任务超时时间：0<br>失败重试次数：0<br><br><span style=\"color:#00c0ef;\" > >>>>>>>>>>>触发调度<<<<<<<<<<< </span><br>触发调度：<br>address：http://127.0.0.1:9999<br>code：500<br>msg：xxl-job remoting error(拒绝连接), for url : http://127.0.0.1:9999/run', NULL, 0, NULL, 2);
INSERT INTO `xxl_job_log` (`id`, `job_group`, `job_id`, `executor_address`, `executor_handler`, `executor_param`, `executor_sharding_param`, `executor_fail_retry_count`, `trigger_time`, `trigger_code`, `trigger_msg`, `handle_time`, `handle_code`, `handle_msg`, `alarm_status`) VALUES (12, 3, 4, 'http://127.0.0.1:9999', 'jobGenerateIndex', '', NULL, 0, '2025-08-02 11:40:00', 500, '任务触发类型：Cron触发<br>调度机器：192.168.60.233<br>执行器-注册方式：手动录入<br>执行器-地址列表：[http://127.0.0.1:9999]<br>路由策略：第一个<br>阻塞处理策略：单机串行<br>任务超时时间：0<br>失败重试次数：0<br><br><span style=\"color:#00c0ef;\" > >>>>>>>>>>>触发调度<<<<<<<<<<< </span><br>触发调度：<br>address：http://127.0.0.1:9999<br>code：500<br>msg：xxl-job remoting error(拒绝连接), for url : http://127.0.0.1:9999/run', NULL, 0, NULL, 2);
INSERT INTO `xxl_job_log` (`id`, `job_group`, `job_id`, `executor_address`, `executor_handler`, `executor_param`, `executor_sharding_param`, `executor_fail_retry_count`, `trigger_time`, `trigger_code`, `trigger_msg`, `handle_time`, `handle_code`, `handle_msg`, `alarm_status`) VALUES (13, 3, 4, 'http://127.0.0.1:9999', 'jobGenerateIndex', '', NULL, 0, '2025-08-02 11:45:00', 200, '任务触发类型：Cron触发<br>调度机器：192.168.60.233<br>执行器-注册方式：手动录入<br>执行器-地址列表：[http://127.0.0.1:9999]<br>路由策略：第一个<br>阻塞处理策略：单机串行<br>任务超时时间：0<br>失败重试次数：0<br><br><span style=\"color:#00c0ef;\" > >>>>>>>>>>>触发调度<<<<<<<<<<< </span><br>触发调度：<br>address：http://127.0.0.1:9999<br>code：200<br>msg：invoked successfully!', '2025-08-02 11:45:00', 500, '<text style=\'color:red\'>任务ID: [4]<br>任务名称: [jobGenerateIndex]<br>参数: <br> panic: runtime error: invalid memory address or nil pointer dereference</text><br>', 2);
INSERT INTO `xxl_job_log` (`id`, `job_group`, `job_id`, `executor_address`, `executor_handler`, `executor_param`, `executor_sharding_param`, `executor_fail_retry_count`, `trigger_time`, `trigger_code`, `trigger_msg`, `handle_time`, `handle_code`, `handle_msg`, `alarm_status`) VALUES (14, 3, 4, 'http://127.0.0.1:9999', 'jobGenerateIndex', '', NULL, 0, '2025-08-02 11:50:00', 200, '任务触发类型：Cron触发<br>调度机器：192.168.60.233<br>执行器-注册方式：手动录入<br>执行器-地址列表：[http://127.0.0.1:9999]<br>路由策略：第一个<br>阻塞处理策略：单机串行<br>任务超时时间：0<br>失败重试次数：0<br><br><span style=\"color:#00c0ef;\" > >>>>>>>>>>>触发调度<<<<<<<<<<< </span><br>触发调度：<br>address：http://127.0.0.1:9999<br>code：200<br>msg：invoked successfully!', '2025-08-02 11:50:00', 500, '<text style=\'color:red\'>任务ID: [4]<br>任务名称: [jobGenerateIndex]<br>参数: <br> panic: runtime error: invalid memory address or nil pointer dereference</text><br>', 2);
INSERT INTO `xxl_job_log` (`id`, `job_group`, `job_id`, `executor_address`, `executor_handler`, `executor_param`, `executor_sharding_param`, `executor_fail_retry_count`, `trigger_time`, `trigger_code`, `trigger_msg`, `handle_time`, `handle_code`, `handle_msg`, `alarm_status`) VALUES (15, 3, 4, 'http://127.0.0.1:9999', 'jobGenerateIndex', '', NULL, 0, '2025-08-02 11:55:00', 200, '任务触发类型：Cron触发<br>调度机器：192.168.60.233<br>执行器-注册方式：手动录入<br>执行器-地址列表：[http://127.0.0.1:9999]<br>路由策略：第一个<br>阻塞处理策略：单机串行<br>任务超时时间：0<br>失败重试次数：0<br><br><span style=\"color:#00c0ef;\" > >>>>>>>>>>>触发调度<<<<<<<<<<< </span><br>触发调度：<br>address：http://127.0.0.1:9999<br>code：200<br>msg：invoked successfully!', '2025-08-02 11:55:00', 500, '<text style=\'color:red\'>任务ID: [4]<br>任务名称: [jobGenerateIndex]<br>参数: <br> panic: runtime error: invalid memory address or nil pointer dereference</text><br>', 2);
INSERT INTO `xxl_job_log` (`id`, `job_group`, `job_id`, `executor_address`, `executor_handler`, `executor_param`, `executor_sharding_param`, `executor_fail_retry_count`, `trigger_time`, `trigger_code`, `trigger_msg`, `handle_time`, `handle_code`, `handle_msg`, `alarm_status`) VALUES (16, 3, 4, 'http://127.0.0.1:9999', 'jobGenerateIndex', '', NULL, 0, '2025-08-02 12:00:00', 200, '任务触发类型：Cron触发<br>调度机器：192.168.60.233<br>执行器-注册方式：手动录入<br>执行器-地址列表：[http://127.0.0.1:9999]<br>路由策略：第一个<br>阻塞处理策略：单机串行<br>任务超时时间：0<br>失败重试次数：0<br><br><span style=\"color:#00c0ef;\" > >>>>>>>>>>>触发调度<<<<<<<<<<< </span><br>触发调度：<br>address：http://127.0.0.1:9999<br>code：200<br>msg：invoked successfully!', '2025-08-02 12:00:00', 500, '<text style=\'color:red\'>任务ID: [4]<br>任务名称: [jobGenerateIndex]<br>参数: <br> panic: runtime error: invalid memory address or nil pointer dereference</text><br>', 2);
INSERT INTO `xxl_job_log` (`id`, `job_group`, `job_id`, `executor_address`, `executor_handler`, `executor_param`, `executor_sharding_param`, `executor_fail_retry_count`, `trigger_time`, `trigger_code`, `trigger_msg`, `handle_time`, `handle_code`, `handle_msg`, `alarm_status`) VALUES (17, 3, 4, 'http://127.0.0.1:9999', 'jobGenerateIndex', '', NULL, 0, '2025-08-02 12:05:00', 200, '任务触发类型：Cron触发<br>调度机器：192.168.60.233<br>执行器-注册方式：手动录入<br>执行器-地址列表：[http://127.0.0.1:9999]<br>路由策略：第一个<br>阻塞处理策略：单机串行<br>任务超时时间：0<br>失败重试次数：0<br><br><span style=\"color:#00c0ef;\" > >>>>>>>>>>>触发调度<<<<<<<<<<< </span><br>触发调度：<br>address：http://127.0.0.1:9999<br>code：200<br>msg：invoked successfully!', '2025-08-02 12:05:00', 500, '<text style=\'color:red\'>任务ID: [4]<br>任务名称: [jobGenerateIndex]<br>参数: <br> panic: runtime error: invalid memory address or nil pointer dereference</text><br>', 2);
INSERT INTO `xxl_job_log` (`id`, `job_group`, `job_id`, `executor_address`, `executor_handler`, `executor_param`, `executor_sharding_param`, `executor_fail_retry_count`, `trigger_time`, `trigger_code`, `trigger_msg`, `handle_time`, `handle_code`, `handle_msg`, `alarm_status`) VALUES (18, 3, 4, 'http://127.0.0.1:9999', 'jobGenerateIndex', '', NULL, 0, '2025-08-02 12:10:00', 200, '任务触发类型：Cron触发<br>调度机器：192.168.60.233<br>执行器-注册方式：手动录入<br>执行器-地址列表：[http://127.0.0.1:9999]<br>路由策略：第一个<br>阻塞处理策略：单机串行<br>任务超时时间：0<br>失败重试次数：0<br><br><span style=\"color:#00c0ef;\" > >>>>>>>>>>>触发调度<<<<<<<<<<< </span><br>触发调度：<br>address：http://127.0.0.1:9999<br>code：200<br>msg：invoked successfully!', '2025-08-02 12:10:00', 500, '<text style=\'color:red\'>任务ID: [4]<br>任务名称: [jobGenerateIndex]<br>参数: <br> panic: runtime error: invalid memory address or nil pointer dereference</text><br>', 2);
INSERT INTO `xxl_job_log` (`id`, `job_group`, `job_id`, `executor_address`, `executor_handler`, `executor_param`, `executor_sharding_param`, `executor_fail_retry_count`, `trigger_time`, `trigger_code`, `trigger_msg`, `handle_time`, `handle_code`, `handle_msg`, `alarm_status`) VALUES (19, 3, 4, 'http://127.0.0.1:9999', 'jobGenerateIndex', '', NULL, 0, '2025-08-02 12:15:00', 200, '任务触发类型：Cron触发<br>调度机器：192.168.60.233<br>执行器-注册方式：手动录入<br>执行器-地址列表：[http://127.0.0.1:9999]<br>路由策略：第一个<br>阻塞处理策略：单机串行<br>任务超时时间：0<br>失败重试次数：0<br><br><span style=\"color:#00c0ef;\" > >>>>>>>>>>>触发调度<<<<<<<<<<< </span><br>触发调度：<br>address：http://127.0.0.1:9999<br>code：200<br>msg：invoked successfully!', '2025-08-02 12:15:00', 500, '<text style=\'color:red\'>任务ID: [4]<br>任务名称: [jobGenerateIndex]<br>参数: <br> panic: runtime error: invalid memory address or nil pointer dereference</text><br>', 2);
INSERT INTO `xxl_job_log` (`id`, `job_group`, `job_id`, `executor_address`, `executor_handler`, `executor_param`, `executor_sharding_param`, `executor_fail_retry_count`, `trigger_time`, `trigger_code`, `trigger_msg`, `handle_time`, `handle_code`, `handle_msg`, `alarm_status`) VALUES (20, 3, 4, 'http://127.0.0.1:9999', 'jobGenerateIndex', '', NULL, 0, '2025-08-02 12:20:00', 500, '任务触发类型：Cron触发<br>调度机器：192.168.60.233<br>执行器-注册方式：手动录入<br>执行器-地址列表：[http://127.0.0.1:9999]<br>路由策略：第一个<br>阻塞处理策略：单机串行<br>任务超时时间：0<br>失败重试次数：0<br><br><span style=\"color:#00c0ef;\" > >>>>>>>>>>>触发调度<<<<<<<<<<< </span><br>触发调度：<br>address：http://127.0.0.1:9999<br>code：500<br>msg：xxl-job remoting error(拒绝连接), for url : http://127.0.0.1:9999/run', NULL, 0, NULL, 2);
INSERT INTO `xxl_job_log` (`id`, `job_group`, `job_id`, `executor_address`, `executor_handler`, `executor_param`, `executor_sharding_param`, `executor_fail_retry_count`, `trigger_time`, `trigger_code`, `trigger_msg`, `handle_time`, `handle_code`, `handle_msg`, `alarm_status`) VALUES (21, 3, 4, 'http://127.0.0.1:9999', 'jobGenerateIndex', '', NULL, 0, '2025-08-02 12:25:00', 500, '任务触发类型：Cron触发<br>调度机器：192.168.60.233<br>执行器-注册方式：手动录入<br>执行器-地址列表：[http://127.0.0.1:9999]<br>路由策略：第一个<br>阻塞处理策略：单机串行<br>任务超时时间：0<br>失败重试次数：0<br><br><span style=\"color:#00c0ef;\" > >>>>>>>>>>>触发调度<<<<<<<<<<< </span><br>触发调度：<br>address：http://127.0.0.1:9999<br>code：500<br>msg：xxl-job remoting error(拒绝连接), for url : http://127.0.0.1:9999/run', NULL, 0, NULL, 2);
INSERT INTO `xxl_job_log` (`id`, `job_group`, `job_id`, `executor_address`, `executor_handler`, `executor_param`, `executor_sharding_param`, `executor_fail_retry_count`, `trigger_time`, `trigger_code`, `trigger_msg`, `handle_time`, `handle_code`, `handle_msg`, `alarm_status`) VALUES (22, 3, 4, 'http://127.0.0.1:9999', 'jobGenerateIndex', '', NULL, 0, '2025-08-02 12:30:00', 500, '任务触发类型：Cron触发<br>调度机器：192.168.60.233<br>执行器-注册方式：手动录入<br>执行器-地址列表：[http://127.0.0.1:9999]<br>路由策略：第一个<br>阻塞处理策略：单机串行<br>任务超时时间：0<br>失败重试次数：0<br><br><span style=\"color:#00c0ef;\" > >>>>>>>>>>>触发调度<<<<<<<<<<< </span><br>触发调度：<br>address：http://127.0.0.1:9999<br>code：500<br>msg：xxl-job remoting error(拒绝连接), for url : http://127.0.0.1:9999/run', NULL, 0, NULL, 2);
INSERT INTO `xxl_job_log` (`id`, `job_group`, `job_id`, `executor_address`, `executor_handler`, `executor_param`, `executor_sharding_param`, `executor_fail_retry_count`, `trigger_time`, `trigger_code`, `trigger_msg`, `handle_time`, `handle_code`, `handle_msg`, `alarm_status`) VALUES (23, 3, 4, 'http://127.0.0.1:9999', 'jobGenerateIndex', '', NULL, 0, '2025-08-02 12:35:00', 500, '任务触发类型：Cron触发<br>调度机器：192.168.60.233<br>执行器-注册方式：手动录入<br>执行器-地址列表：[http://127.0.0.1:9999]<br>路由策略：第一个<br>阻塞处理策略：单机串行<br>任务超时时间：0<br>失败重试次数：0<br><br><span style=\"color:#00c0ef;\" > >>>>>>>>>>>触发调度<<<<<<<<<<< </span><br>触发调度：<br>address：http://127.0.0.1:9999<br>code：500<br>msg：xxl-job remoting error(拒绝连接), for url : http://127.0.0.1:9999/run', NULL, 0, NULL, 2);
INSERT INTO `xxl_job_log` (`id`, `job_group`, `job_id`, `executor_address`, `executor_handler`, `executor_param`, `executor_sharding_param`, `executor_fail_retry_count`, `trigger_time`, `trigger_code`, `trigger_msg`, `handle_time`, `handle_code`, `handle_msg`, `alarm_status`) VALUES (24, 3, 4, 'http://127.0.0.1:9999', 'jobGenerateIndex', '', NULL, 0, '2025-08-02 12:40:00', 500, '任务触发类型：Cron触发<br>调度机器：192.168.60.233<br>执行器-注册方式：手动录入<br>执行器-地址列表：[http://127.0.0.1:9999]<br>路由策略：第一个<br>阻塞处理策略：单机串行<br>任务超时时间：0<br>失败重试次数：0<br><br><span style=\"color:#00c0ef;\" > >>>>>>>>>>>触发调度<<<<<<<<<<< </span><br>触发调度：<br>address：http://127.0.0.1:9999<br>code：500<br>msg：xxl-job remoting error(拒绝连接), for url : http://127.0.0.1:9999/run', NULL, 0, NULL, 2);
INSERT INTO `xxl_job_log` (`id`, `job_group`, `job_id`, `executor_address`, `executor_handler`, `executor_param`, `executor_sharding_param`, `executor_fail_retry_count`, `trigger_time`, `trigger_code`, `trigger_msg`, `handle_time`, `handle_code`, `handle_msg`, `alarm_status`) VALUES (25, 3, 4, 'http://127.0.0.1:9999', 'jobGenerateIndex', '', NULL, 0, '2025-08-02 12:45:00', 500, '任务触发类型：Cron触发<br>调度机器：192.168.60.233<br>执行器-注册方式：手动录入<br>执行器-地址列表：[http://127.0.0.1:9999]<br>路由策略：第一个<br>阻塞处理策略：单机串行<br>任务超时时间：0<br>失败重试次数：0<br><br><span style=\"color:#00c0ef;\" > >>>>>>>>>>>触发调度<<<<<<<<<<< </span><br>触发调度：<br>address：http://127.0.0.1:9999<br>code：500<br>msg：xxl-job remoting error(拒绝连接), for url : http://127.0.0.1:9999/run', NULL, 0, NULL, 2);
INSERT INTO `xxl_job_log` (`id`, `job_group`, `job_id`, `executor_address`, `executor_handler`, `executor_param`, `executor_sharding_param`, `executor_fail_retry_count`, `trigger_time`, `trigger_code`, `trigger_msg`, `handle_time`, `handle_code`, `handle_msg`, `alarm_status`) VALUES (26, 3, 4, 'http://127.0.0.1:9999', 'jobGenerateIndex', '', NULL, 0, '2025-08-02 12:50:00', 500, '任务触发类型：Cron触发<br>调度机器：192.168.60.233<br>执行器-注册方式：手动录入<br>执行器-地址列表：[http://127.0.0.1:9999]<br>路由策略：第一个<br>阻塞处理策略：单机串行<br>任务超时时间：0<br>失败重试次数：0<br><br><span style=\"color:#00c0ef;\" > >>>>>>>>>>>触发调度<<<<<<<<<<< </span><br>触发调度：<br>address：http://127.0.0.1:9999<br>code：500<br>msg：xxl-job remoting error(拒绝连接), for url : http://127.0.0.1:9999/run', NULL, 0, NULL, 2);
INSERT INTO `xxl_job_log` (`id`, `job_group`, `job_id`, `executor_address`, `executor_handler`, `executor_param`, `executor_sharding_param`, `executor_fail_retry_count`, `trigger_time`, `trigger_code`, `trigger_msg`, `handle_time`, `handle_code`, `handle_msg`, `alarm_status`) VALUES (27, 3, 4, 'http://127.0.0.1:9999', 'jobGenerateIndex', '', NULL, 0, '2025-08-02 12:55:00', 500, '任务触发类型：Cron触发<br>调度机器：192.168.60.233<br>执行器-注册方式：手动录入<br>执行器-地址列表：[http://127.0.0.1:9999]<br>路由策略：第一个<br>阻塞处理策略：单机串行<br>任务超时时间：0<br>失败重试次数：0<br><br><span style=\"color:#00c0ef;\" > >>>>>>>>>>>触发调度<<<<<<<<<<< </span><br>触发调度：<br>address：http://127.0.0.1:9999<br>code：500<br>msg：xxl-job remoting error(拒绝连接), for url : http://127.0.0.1:9999/run', NULL, 0, NULL, 2);
INSERT INTO `xxl_job_log` (`id`, `job_group`, `job_id`, `executor_address`, `executor_handler`, `executor_param`, `executor_sharding_param`, `executor_fail_retry_count`, `trigger_time`, `trigger_code`, `trigger_msg`, `handle_time`, `handle_code`, `handle_msg`, `alarm_status`) VALUES (28, 3, 4, 'http://127.0.0.1:9999', 'jobGenerateIndex', '', NULL, 0, '2025-08-02 13:00:00', 500, '任务触发类型：Cron触发<br>调度机器：192.168.60.233<br>执行器-注册方式：手动录入<br>执行器-地址列表：[http://127.0.0.1:9999]<br>路由策略：第一个<br>阻塞处理策略：单机串行<br>任务超时时间：0<br>失败重试次数：0<br><br><span style=\"color:#00c0ef;\" > >>>>>>>>>>>触发调度<<<<<<<<<<< </span><br>触发调度：<br>address：http://127.0.0.1:9999<br>code：500<br>msg：xxl-job remoting error(拒绝连接), for url : http://127.0.0.1:9999/run', NULL, 0, NULL, 2);
INSERT INTO `xxl_job_log` (`id`, `job_group`, `job_id`, `executor_address`, `executor_handler`, `executor_param`, `executor_sharding_param`, `executor_fail_retry_count`, `trigger_time`, `trigger_code`, `trigger_msg`, `handle_time`, `handle_code`, `handle_msg`, `alarm_status`) VALUES (29, 3, 4, 'http://127.0.0.1:9999', 'jobGenerateIndex', '', NULL, 0, '2025-08-02 13:05:00', 500, '任务触发类型：Cron触发<br>调度机器：192.168.60.233<br>执行器-注册方式：手动录入<br>执行器-地址列表：[http://127.0.0.1:9999]<br>路由策略：第一个<br>阻塞处理策略：单机串行<br>任务超时时间：0<br>失败重试次数：0<br><br><span style=\"color:#00c0ef;\" > >>>>>>>>>>>触发调度<<<<<<<<<<< </span><br>触发调度：<br>address：http://127.0.0.1:9999<br>code：500<br>msg：xxl-job remoting error(拒绝连接), for url : http://127.0.0.1:9999/run', NULL, 0, NULL, 2);
INSERT INTO `xxl_job_log` (`id`, `job_group`, `job_id`, `executor_address`, `executor_handler`, `executor_param`, `executor_sharding_param`, `executor_fail_retry_count`, `trigger_time`, `trigger_code`, `trigger_msg`, `handle_time`, `handle_code`, `handle_msg`, `alarm_status`) VALUES (30, 3, 4, 'http://127.0.0.1:9999', 'jobGenerateIndex', '', NULL, 0, '2025-08-02 13:10:00', 500, '任务触发类型：Cron触发<br>调度机器：192.168.60.233<br>执行器-注册方式：手动录入<br>执行器-地址列表：[http://127.0.0.1:9999]<br>路由策略：第一个<br>阻塞处理策略：单机串行<br>任务超时时间：0<br>失败重试次数：0<br><br><span style=\"color:#00c0ef;\" > >>>>>>>>>>>触发调度<<<<<<<<<<< </span><br>触发调度：<br>address：http://127.0.0.1:9999<br>code：500<br>msg：xxl-job remoting error(拒绝连接), for url : http://127.0.0.1:9999/run', NULL, 0, NULL, 2);
INSERT INTO `xxl_job_log` (`id`, `job_group`, `job_id`, `executor_address`, `executor_handler`, `executor_param`, `executor_sharding_param`, `executor_fail_retry_count`, `trigger_time`, `trigger_code`, `trigger_msg`, `handle_time`, `handle_code`, `handle_msg`, `alarm_status`) VALUES (31, 3, 4, 'http://127.0.0.1:9999', 'jobGenerateIndex', '', NULL, 0, '2025-08-02 13:15:00', 500, '任务触发类型：Cron触发<br>调度机器：192.168.60.233<br>执行器-注册方式：手动录入<br>执行器-地址列表：[http://127.0.0.1:9999]<br>路由策略：第一个<br>阻塞处理策略：单机串行<br>任务超时时间：0<br>失败重试次数：0<br><br><span style=\"color:#00c0ef;\" > >>>>>>>>>>>触发调度<<<<<<<<<<< </span><br>触发调度：<br>address：http://127.0.0.1:9999<br>code：500<br>msg：xxl-job remoting error(拒绝连接), for url : http://127.0.0.1:9999/run', NULL, 0, NULL, 2);
INSERT INTO `xxl_job_log` (`id`, `job_group`, `job_id`, `executor_address`, `executor_handler`, `executor_param`, `executor_sharding_param`, `executor_fail_retry_count`, `trigger_time`, `trigger_code`, `trigger_msg`, `handle_time`, `handle_code`, `handle_msg`, `alarm_status`) VALUES (32, 3, 4, 'http://127.0.0.1:9999', 'jobGenerateIndex', '', NULL, 0, '2025-08-02 13:20:00', 500, '任务触发类型：Cron触发<br>调度机器：192.168.60.233<br>执行器-注册方式：手动录入<br>执行器-地址列表：[http://127.0.0.1:9999]<br>路由策略：第一个<br>阻塞处理策略：单机串行<br>任务超时时间：0<br>失败重试次数：0<br><br><span style=\"color:#00c0ef;\" > >>>>>>>>>>>触发调度<<<<<<<<<<< </span><br>触发调度：<br>address：http://127.0.0.1:9999<br>code：500<br>msg：xxl-job remoting error(拒绝连接), for url : http://127.0.0.1:9999/run', NULL, 0, NULL, 2);
INSERT INTO `xxl_job_log` (`id`, `job_group`, `job_id`, `executor_address`, `executor_handler`, `executor_param`, `executor_sharding_param`, `executor_fail_retry_count`, `trigger_time`, `trigger_code`, `trigger_msg`, `handle_time`, `handle_code`, `handle_msg`, `alarm_status`) VALUES (33, 3, 4, 'http://127.0.0.1:9999', 'jobGenerateIndex', '', NULL, 0, '2025-08-02 13:25:00', 500, '任务触发类型：Cron触发<br>调度机器：192.168.60.233<br>执行器-注册方式：手动录入<br>执行器-地址列表：[http://127.0.0.1:9999]<br>路由策略：第一个<br>阻塞处理策略：单机串行<br>任务超时时间：0<br>失败重试次数：0<br><br><span style=\"color:#00c0ef;\" > >>>>>>>>>>>触发调度<<<<<<<<<<< </span><br>触发调度：<br>address：http://127.0.0.1:9999<br>code：500<br>msg：xxl-job remoting error(拒绝连接), for url : http://127.0.0.1:9999/run', NULL, 0, NULL, 2);
INSERT INTO `xxl_job_log` (`id`, `job_group`, `job_id`, `executor_address`, `executor_handler`, `executor_param`, `executor_sharding_param`, `executor_fail_retry_count`, `trigger_time`, `trigger_code`, `trigger_msg`, `handle_time`, `handle_code`, `handle_msg`, `alarm_status`) VALUES (34, 3, 4, 'http://127.0.0.1:9999', 'jobGenerateIndex', '', NULL, 0, '2025-08-02 13:30:00', 500, '任务触发类型：Cron触发<br>调度机器：192.168.60.233<br>执行器-注册方式：手动录入<br>执行器-地址列表：[http://127.0.0.1:9999]<br>路由策略：第一个<br>阻塞处理策略：单机串行<br>任务超时时间：0<br>失败重试次数：0<br><br><span style=\"color:#00c0ef;\" > >>>>>>>>>>>触发调度<<<<<<<<<<< </span><br>触发调度：<br>address：http://127.0.0.1:9999<br>code：500<br>msg：xxl-job remoting error(拒绝连接), for url : http://127.0.0.1:9999/run', NULL, 0, NULL, 2);
INSERT INTO `xxl_job_log` (`id`, `job_group`, `job_id`, `executor_address`, `executor_handler`, `executor_param`, `executor_sharding_param`, `executor_fail_retry_count`, `trigger_time`, `trigger_code`, `trigger_msg`, `handle_time`, `handle_code`, `handle_msg`, `alarm_status`) VALUES (35, 3, 4, 'http://127.0.0.1:9999', 'jobGenerateIndex', '', NULL, 0, '2025-08-02 13:35:00', 500, '任务触发类型：Cron触发<br>调度机器：192.168.60.233<br>执行器-注册方式：手动录入<br>执行器-地址列表：[http://127.0.0.1:9999]<br>路由策略：第一个<br>阻塞处理策略：单机串行<br>任务超时时间：0<br>失败重试次数：0<br><br><span style=\"color:#00c0ef;\" > >>>>>>>>>>>触发调度<<<<<<<<<<< </span><br>触发调度：<br>address：http://127.0.0.1:9999<br>code：500<br>msg：xxl-job remoting error(拒绝连接), for url : http://127.0.0.1:9999/run', NULL, 0, NULL, 2);
INSERT INTO `xxl_job_log` (`id`, `job_group`, `job_id`, `executor_address`, `executor_handler`, `executor_param`, `executor_sharding_param`, `executor_fail_retry_count`, `trigger_time`, `trigger_code`, `trigger_msg`, `handle_time`, `handle_code`, `handle_msg`, `alarm_status`) VALUES (36, 3, 4, 'http://127.0.0.1:9999', 'jobGenerateIndex', '', NULL, 0, '2025-08-02 13:40:00', 500, '任务触发类型：Cron触发<br>调度机器：192.168.60.233<br>执行器-注册方式：手动录入<br>执行器-地址列表：[http://127.0.0.1:9999]<br>路由策略：第一个<br>阻塞处理策略：单机串行<br>任务超时时间：0<br>失败重试次数：0<br><br><span style=\"color:#00c0ef;\" > >>>>>>>>>>>触发调度<<<<<<<<<<< </span><br>触发调度：<br>address：http://127.0.0.1:9999<br>code：500<br>msg：xxl-job remoting error(拒绝连接), for url : http://127.0.0.1:9999/run', NULL, 0, NULL, 2);
INSERT INTO `xxl_job_log` (`id`, `job_group`, `job_id`, `executor_address`, `executor_handler`, `executor_param`, `executor_sharding_param`, `executor_fail_retry_count`, `trigger_time`, `trigger_code`, `trigger_msg`, `handle_time`, `handle_code`, `handle_msg`, `alarm_status`) VALUES (37, 3, 4, 'http://127.0.0.1:9999', 'jobGenerateIndex', '', NULL, 0, '2025-08-02 13:45:00', 500, '任务触发类型：Cron触发<br>调度机器：192.168.60.233<br>执行器-注册方式：手动录入<br>执行器-地址列表：[http://127.0.0.1:9999]<br>路由策略：第一个<br>阻塞处理策略：单机串行<br>任务超时时间：0<br>失败重试次数：0<br><br><span style=\"color:#00c0ef;\" > >>>>>>>>>>>触发调度<<<<<<<<<<< </span><br>触发调度：<br>address：http://127.0.0.1:9999<br>code：500<br>msg：xxl-job remoting error(拒绝连接), for url : http://127.0.0.1:9999/run', NULL, 0, NULL, 2);
INSERT INTO `xxl_job_log` (`id`, `job_group`, `job_id`, `executor_address`, `executor_handler`, `executor_param`, `executor_sharding_param`, `executor_fail_retry_count`, `trigger_time`, `trigger_code`, `trigger_msg`, `handle_time`, `handle_code`, `handle_msg`, `alarm_status`) VALUES (38, 3, 4, 'http://127.0.0.1:9999', 'jobGenerateIndex', '', NULL, 0, '2025-08-02 13:50:00', 200, '任务触发类型：Cron触发<br>调度机器：192.168.60.233<br>执行器-注册方式：手动录入<br>执行器-地址列表：[http://127.0.0.1:9999]<br>路由策略：第一个<br>阻塞处理策略：单机串行<br>任务超时时间：0<br>失败重试次数：0<br><br><span style=\"color:#00c0ef;\" > >>>>>>>>>>>触发调度<<<<<<<<<<< </span><br>触发调度：<br>address：http://127.0.0.1:9999<br>code：200<br>msg：invoked successfully!', '2025-08-02 13:50:00', 200, 'success', 0);
INSERT INTO `xxl_job_log` (`id`, `job_group`, `job_id`, `executor_address`, `executor_handler`, `executor_param`, `executor_sharding_param`, `executor_fail_retry_count`, `trigger_time`, `trigger_code`, `trigger_msg`, `handle_time`, `handle_code`, `handle_msg`, `alarm_status`) VALUES (39, 3, 4, 'http://127.0.0.1:9999', 'jobGenerateIndex', '', NULL, 0, '2025-08-02 13:55:00', 200, '任务触发类型：Cron触发<br>调度机器：192.168.60.233<br>执行器-注册方式：手动录入<br>执行器-地址列表：[http://127.0.0.1:9999]<br>路由策略：第一个<br>阻塞处理策略：单机串行<br>任务超时时间：0<br>失败重试次数：0<br><br><span style=\"color:#00c0ef;\" > >>>>>>>>>>>触发调度<<<<<<<<<<< </span><br>触发调度：<br>address：http://127.0.0.1:9999<br>code：200<br>msg：invoked successfully!', '2025-08-02 13:55:00', 200, 'success', 0);
INSERT INTO `xxl_job_log` (`id`, `job_group`, `job_id`, `executor_address`, `executor_handler`, `executor_param`, `executor_sharding_param`, `executor_fail_retry_count`, `trigger_time`, `trigger_code`, `trigger_msg`, `handle_time`, `handle_code`, `handle_msg`, `alarm_status`) VALUES (40, 3, 4, 'http://127.0.0.1:9999', 'jobGenerateIndex', '', NULL, 0, '2025-08-02 14:00:00', 200, '任务触发类型：Cron触发<br>调度机器：192.168.60.233<br>执行器-注册方式：手动录入<br>执行器-地址列表：[http://127.0.0.1:9999]<br>路由策略：第一个<br>阻塞处理策略：单机串行<br>任务超时时间：0<br>失败重试次数：0<br><br><span style=\"color:#00c0ef;\" > >>>>>>>>>>>触发调度<<<<<<<<<<< </span><br>触发调度：<br>address：http://127.0.0.1:9999<br>code：200<br>msg：invoked successfully!', '2025-08-02 14:00:00', 200, 'success', 0);
COMMIT;

-- ----------------------------
-- Table structure for xxl_job_log_report
-- ----------------------------
DROP TABLE IF EXISTS `xxl_job_log_report`;
CREATE TABLE `xxl_job_log_report` (
  `id` int NOT NULL AUTO_INCREMENT,
  `trigger_day` datetime DEFAULT NULL COMMENT '调度-时间',
  `running_count` int NOT NULL DEFAULT '0' COMMENT '运行中-日志数量',
  `suc_count` int NOT NULL DEFAULT '0' COMMENT '执行成功-日志数量',
  `fail_count` int NOT NULL DEFAULT '0' COMMENT '执行失败-日志数量',
  `update_time` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `i_trigger_day` (`trigger_day`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Records of xxl_job_log_report
-- ----------------------------
BEGIN;
INSERT INTO `xxl_job_log_report` (`id`, `trigger_day`, `running_count`, `suc_count`, `fail_count`, `update_time`) VALUES (1, '2025-08-02 00:00:00', 0, 3, 37, NULL);
INSERT INTO `xxl_job_log_report` (`id`, `trigger_day`, `running_count`, `suc_count`, `fail_count`, `update_time`) VALUES (2, '2025-08-01 00:00:00', 0, 0, 0, NULL);
INSERT INTO `xxl_job_log_report` (`id`, `trigger_day`, `running_count`, `suc_count`, `fail_count`, `update_time`) VALUES (3, '2025-07-31 00:00:00', 0, 0, 0, NULL);
COMMIT;

-- ----------------------------
-- Table structure for xxl_job_logglue
-- ----------------------------
DROP TABLE IF EXISTS `xxl_job_logglue`;
CREATE TABLE `xxl_job_logglue` (
  `id` int NOT NULL AUTO_INCREMENT,
  `job_id` int NOT NULL COMMENT '任务，主键ID',
  `glue_type` varchar(50) DEFAULT NULL COMMENT 'GLUE类型',
  `glue_source` mediumtext COMMENT 'GLUE源代码',
  `glue_remark` varchar(128) NOT NULL COMMENT 'GLUE备注',
  `add_time` datetime DEFAULT NULL,
  `update_time` datetime DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Records of xxl_job_logglue
-- ----------------------------
BEGIN;
COMMIT;

-- ----------------------------
-- Table structure for xxl_job_registry
-- ----------------------------
DROP TABLE IF EXISTS `xxl_job_registry`;
CREATE TABLE `xxl_job_registry` (
  `id` int NOT NULL AUTO_INCREMENT,
  `registry_group` varchar(50) NOT NULL,
  `registry_key` varchar(255) NOT NULL,
  `registry_value` varchar(255) NOT NULL,
  `update_time` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `i_g_k_v` (`registry_group`,`registry_key`,`registry_value`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=173 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Records of xxl_job_registry
-- ----------------------------
BEGIN;
INSERT INTO `xxl_job_registry` (`id`, `registry_group`, `registry_key`, `registry_value`, `update_time`) VALUES (158, 'EXECUTOR', 'worldExecutor', 'http://127.0.0.1:9999', '2025-08-02 14:03:01');
COMMIT;

-- ----------------------------
-- Table structure for xxl_job_user
-- ----------------------------
DROP TABLE IF EXISTS `xxl_job_user`;
CREATE TABLE `xxl_job_user` (
  `id` int NOT NULL AUTO_INCREMENT,
  `username` varchar(50) NOT NULL COMMENT '账号',
  `password` varchar(50) NOT NULL COMMENT '密码',
  `role` tinyint NOT NULL COMMENT '角色：0-普通用户、1-管理员',
  `permission` varchar(255) DEFAULT NULL COMMENT '权限：执行器ID列表，多个逗号分割',
  PRIMARY KEY (`id`),
  UNIQUE KEY `i_username` (`username`) USING BTREE
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Records of xxl_job_user
-- ----------------------------
BEGIN;
INSERT INTO `xxl_job_user` (`id`, `username`, `password`, `role`, `permission`) VALUES (1, 'admin', 'e10adc3949ba59abbe56e057f20f883e', 1, NULL);
COMMIT;

SET FOREIGN_KEY_CHECKS = 1;
