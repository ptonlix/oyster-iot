/*
 Navicat Premium Data Transfer

 Source Server         : oyster
 Source Server Type    : MySQL
 Source Server Version : 80028
 Source Host           : localhost:3306
 Source Schema         : oyster

 Target Server Type    : MySQL
 Target Server Version : 80028
 File Encoding         : 65001

 Date: 11/05/2022 10:29:15
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for business
-- ----------------------------
DROP TABLE IF EXISTS `business`;
CREATE TABLE `business` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '业务ID',
  `name` varchar(255) NOT NULL DEFAULT '' COMMENT '业务名称',
  `remark` varchar(255) NOT NULL DEFAULT '' COMMENT '业务信息说明',
  `createdat` datetime NOT NULL COMMENT '创建时间',
  `updatedat` datetime NOT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Table structure for device
-- ----------------------------
DROP TABLE IF EXISTS `device`;
CREATE TABLE `device` (
  `id` int NOT NULL AUTO_INCREMENT COMMENT '设备ID',
  `assets_num` varchar(255) NOT NULL DEFAULT '' COMMENT '设备资产编码',
  `device_name` varchar(255) NOT NULL DEFAULT '' COMMENT '设备名称',
  `token` varchar(255) NOT NULL DEFAULT '' COMMENT '设备验证token',
  `protocol` varchar(255) NOT NULL DEFAULT '' COMMENT '设备采用的协议 MQTT TCP ...',
  `publish` varchar(255) NOT NULL DEFAULT '' COMMENT '设备发布消息的主题',
  `subscribe` varchar(255) NOT NULL DEFAULT '' COMMENT '设备订阅的主题',
  `type` varchar(64) NOT NULL DEFAULT '' COMMENT '设备类型',
  `business_id` int DEFAULT NULL COMMENT '设备关联的业务ID',
  PRIMARY KEY (`id`),
  UNIQUE KEY `assets_num` (`assets_num`),
  KEY `device_business_id` (`business_id`)
) ENGINE=InnoDB AUTO_INCREMENT=13 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Table structure for device_data
-- ----------------------------
DROP TABLE IF EXISTS `device_data`;
CREATE TABLE `device_data` (
  `id` int NOT NULL AUTO_INCREMENT,
  `dev_assets_num` varchar(255) NOT NULL DEFAULT '' COMMENT '设备资产编码',
  `dev_type` varchar(64) NOT NULL DEFAULT '',
  `msg` longtext NOT NULL COMMENT '设备上报的信息',
  `ts` datetime NOT NULL COMMENT '保存信息时间戳',
  PRIMARY KEY (`id`),
  KEY `device_data_dev_assets_num` (`dev_assets_num`)
) ENGINE=InnoDB AUTO_INCREMENT=879 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Table structure for users
-- ----------------------------
DROP TABLE IF EXISTS `users`;
CREATE TABLE `users` (
  `id` int NOT NULL AUTO_INCREMENT,
  `enabled` tinyint(1) NOT NULL DEFAULT '0',
  `email` varchar(255) NOT NULL DEFAULT '',
  `username` varchar(255) NOT NULL DEFAULT '',
  `password` varchar(255) NOT NULL DEFAULT '',
  `firstname` varchar(255) NOT NULL DEFAULT '',
  `lastname` varchar(255) NOT NULL DEFAULT '',
  `mobile` varchar(255) NOT NULL DEFAULT '',
  `remark` varchar(255) NOT NULL DEFAULT '',
  `is_admin` tinyint(1) NOT NULL DEFAULT '0',
  `wxopenid` varchar(255) NOT NULL DEFAULT '',
  `wxunionid` varchar(255) NOT NULL DEFAULT '',
  `createdat` datetime NOT NULL COMMENT '创建时间',
  `updatedat` datetime NOT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `username` (`username`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

SET FOREIGN_KEY_CHECKS = 1;
