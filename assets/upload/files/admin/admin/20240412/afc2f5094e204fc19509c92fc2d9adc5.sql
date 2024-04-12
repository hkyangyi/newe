/*
 Navicat Premium Data Transfer

 Source Server         : 本地
 Source Server Type    : MySQL
 Source Server Version : 50644
 Source Host           : localhost:3306
 Source Schema         : newe

 Target Server Type    : MySQL
 Target Server Version : 50644
 File Encoding         : 65001

 Date: 11/04/2024 19:32:42
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for sys_api_list
-- ----------------------------
DROP TABLE IF EXISTS `sys_api_list`;
CREATE TABLE `sys_api_list`  (
  `id` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `menu_id` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '菜单id',
  `path` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT 'api 路径',
  `method` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '请求方式',
  `name` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT 'api名称',
  `create_time` bigint(20) NULL DEFAULT NULL,
  `update_time` bigint(20) NULL DEFAULT NULL,
  `status` int(1) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Compact;

-- ----------------------------
-- Records of sys_api_list
-- ----------------------------
INSERT INTO `sys_api_list` VALUES ('0163cd4dff5c44ee93aeb6e0a3f68ada', '9b144acec61e4886bdcbf0e6ad333481', '/newesys/admin/system/deptgetlist', 'get', '获取组织结构列表', 1706173243, 1706173243, 1);
INSERT INTO `sys_api_list` VALUES ('0461798d7f7245da8c6f3c5bab5f165d', '702d3b0676524eaa97de61db2aeb3616', '/newesys/admin/system/dictadd', 'post', '添加字典', 1706173489, 1706173489, 1);
INSERT INTO `sys_api_list` VALUES ('070914e683474faeaf4a2e9202b8df72', '44030ff97bf7446693ace817f71dfb7a', '/newesys/admin/system/ApiListSetStatus', 'put', '启停api', 1706173643, 1706173643, 1);
INSERT INTO `sys_api_list` VALUES ('0d3c9db1272243ebad7c45d44f9dd9b8', '6fb3210467194b64a6c60a90f0a6e3c5', '/newesys/admin/system/roledel', 'delete', '删除角色', 1706173444, 1706173444, 1);
INSERT INTO `sys_api_list` VALUES ('18ac3a56882a4248a592c0cd12568a81', 'd436aeea798d42a4bea423164f214623', '/newesys/admin/system/memberdel', 'delete', '删除员工', 1706172274, 1706172274, 1);
INSERT INTO `sys_api_list` VALUES ('1a1c2fb5fc474c7d84674c36dc8076d9', 'd436aeea798d42a4bea423164f214623', '/newesys/admin/system/memberedit', 'put', '编辑员工', 1706172232, 1706172232, 1);
INSERT INTO `sys_api_list` VALUES ('1b347deaa11247d6be007f76276ef676', '6fb3210467194b64a6c60a90f0a6e3c5', '/newesys/admin/system/roleedit', 'put', '编辑角色', 1706173424, 1706173424, 1);
INSERT INTO `sys_api_list` VALUES ('26f3660954794e67b85b7d12d68e6b10', '44030ff97bf7446693ace817f71dfb7a', '/newesys/admin/system/ApiListAdd', 'post', '添加api', 1706173585, 1706173585, 1);
INSERT INTO `sys_api_list` VALUES ('2d0256f898f144258ea944861f92cdc1', '6fb3210467194b64a6c60a90f0a6e3c5', '/newesys/admin/system/roleRulesEdit', 'put', '修改角色数据权限', 1709808309, 1709808647, 1);
INSERT INTO `sys_api_list` VALUES ('34f04b50dedd48e5a8dd97b2040764f1', '9b144acec61e4886bdcbf0e6ad333481', '/newesys/admin/system/deptrulessave', 'post', '编辑组织结构菜单权限', 1706173358, 1706173358, 1);
INSERT INTO `sys_api_list` VALUES ('41e1843c425145789e3bf5ebd45a2939', '44030ff97bf7446693ace817f71dfb7a', '/newesys/admin/system/ApiListEdit', 'put', '编辑api', 1706173599, 1706173599, 1);
INSERT INTO `sys_api_list` VALUES ('43eb887c22ee427eafe16e995a56d57f', '6fb3210467194b64a6c60a90f0a6e3c5', '/newesys/admin/system/roleadd', 'post', '添加角色', 1706173396, 1706173983, 1);
INSERT INTO `sys_api_list` VALUES ('50dfd071e389411ba5b4e9090c8efdf9', 'd436aeea798d42a4bea423164f214623', '/newesys/admin/system/membergetlist', 'get', '获取员工列表', 1705919667, 1705919667, 1);
INSERT INTO `sys_api_list` VALUES ('51211bd4a7344054a2bed3bafa94de3b', '702d3b0676524eaa97de61db2aeb3616', '/newesys/admin/system/dictdel', 'delete', '删除字典', 1706173526, 1706173526, 1);
INSERT INTO `sys_api_list` VALUES ('5b2d83c3b47f40d99020fa0d37759a75', '702d3b0676524eaa97de61db2aeb3616', '/newesys/admin/system/dictedit', 'get', '编辑字典', 1706173505, 1706173505, 1);
INSERT INTO `sys_api_list` VALUES ('61f1549a67d240feb5a0b62a57abe328', '6fb3210467194b64a6c60a90f0a6e3c5', '/newesys/admin/system/rolegetlist', 'get', '获取角色列表', 1706173383, 1706173383, 1);
INSERT INTO `sys_api_list` VALUES ('7451cf94e359481f8caae27aa486ee0b', '9b144acec61e4886bdcbf0e6ad333481', '/newesys/admin/system/deptdel', 'delete', '删除组织结构', 1706173294, 1706173294, 1);
INSERT INTO `sys_api_list` VALUES ('7a983997340340f6bd3ee169f4eb29d4', '9b144acec61e4886bdcbf0e6ad333481', '/newesys/admin/system/deptgetrules', 'get', '获取组织结构菜单权限列表', 1706173324, 1706173324, 1);
INSERT INTO `sys_api_list` VALUES ('7bb1756fe75948809d75b529bc3f1f9a', '84cb2653afa7499e829f00bc0c8367f7', '/newesys/admin/system/delMenuList', 'delete', '删除菜单', 1706173183, 1706173183, 1);
INSERT INTO `sys_api_list` VALUES ('8204b2c4f45c43f5bb03149de0f3f2df', '84cb2653afa7499e829f00bc0c8367f7', '/newesys/admin/system/addMenuList', 'post', '添加菜单', 1706173153, 1706173153, 1);
INSERT INTO `sys_api_list` VALUES ('84af270e7bdf4063815c1b6a0d5d9584', '44030ff97bf7446693ace817f71dfb7a', '/newesys/admin/system/ApiListDel', 'delete', '删除api', 1706173627, 1706173627, 1);
INSERT INTO `sys_api_list` VALUES ('a509cd28d6574e458bbd8f7572e0bcd1', '84cb2653afa7499e829f00bc0c8367f7', '/newesys/admin/system/editMenuList', 'put', '编辑菜单', 1706173166, 1706173166, 1);
INSERT INTO `sys_api_list` VALUES ('b0123da2ea964f2db0fd59084b9f217c', '6fb3210467194b64a6c60a90f0a6e3c5', '/newesys/admin/system/roleRulesGetList', 'get', '获取角色数据权限列表', 1709808279, 1709808279, 1);
INSERT INTO `sys_api_list` VALUES ('b2d7a48899114b6aacee94f22baf5c48', '84cb2653afa7499e829f00bc0c8367f7', '/newesys/admin/system/getMenuList', 'get', '获取列表', 1706173136, 1706173136, 1);
INSERT INTO `sys_api_list` VALUES ('c2ff8e889ec84da5985f14155b06cad3', '702d3b0676524eaa97de61db2aeb3616', '/newesys/admin/system/dictgetlist', 'get', '获取字典列表', 1706173474, 1706173474, 1);
INSERT INTO `sys_api_list` VALUES ('c7e112ebb1f143309405ea8e725fc241', '9b144acec61e4886bdcbf0e6ad333481', '/newesys/admin/system/deptedit', 'put', '编辑组织结构', 1706173278, 1706173278, 1);
INSERT INTO `sys_api_list` VALUES ('d3af6f4386a3417e94cb0713d3c77888', '702d3b0676524eaa97de61db2aeb3616', '/newesys/admin/system/dictgetbycode', 'get', '根据字典编码查询字典内容', 1706173563, 1706173563, 1);
INSERT INTO `sys_api_list` VALUES ('d893f3174fe14e90bd00768a8b1618ea', '44030ff97bf7446693ace817f71dfb7a', '/newesys/admin/system/ApiListGetList', 'get', '获取API列表', 1706173613, 1706173613, 1);
INSERT INTO `sys_api_list` VALUES ('e62e164ec7664389938825b48c5c1445', 'd436aeea798d42a4bea423164f214623', '/newesys/admin/system/memberadd', 'post', '添加员工', 1705919755, 1705919755, 1);
INSERT INTO `sys_api_list` VALUES ('ea7aae887e5a42fd9a5d069706d55e43', '9b144acec61e4886bdcbf0e6ad333481', '/newesys/admin/system/deptadd', 'post', '添加组织结构', 1706173264, 1706173264, 1);

-- ----------------------------
-- Table structure for sys_depart
-- ----------------------------
DROP TABLE IF EXISTS `sys_depart`;
CREATE TABLE `sys_depart`  (
  `id` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT 'uuid',
  `pid` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '父级ID',
  `name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '分组名称（机构名称）',
  `code` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '分组编码',
  `type` int(11) NULL DEFAULT NULL COMMENT '类型（1集团，2公司，3部门，4服务门店）',
  `telephone` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '联系电话',
  `phone` varchar(30) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '联系手机',
  `address` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '地址',
  `sort_no` int(11) NULL DEFAULT NULL COMMENT '排序',
  `create_time` bigint(20) NULL DEFAULT NULL COMMENT '创建时间',
  `update_time` bigint(20) NULL DEFAULT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Compact;

-- ----------------------------
-- Records of sys_depart
-- ----------------------------
INSERT INTO `sys_depart` VALUES ('1d0b1d775cda4244abb4c4c7f7152771', '', 'newe', 'A8', 1, '', '', '', 0, 1703582271, 1703582271);
INSERT INTO `sys_depart` VALUES ('8031a1e774e54152a3952895b0471ce3', '', '测试区', 'A10', 3, '', '', '', 1, 1709794790, 1709794790);
INSERT INTO `sys_depart` VALUES ('c1a4ad97b18a404cbbe83a31d1ee02d6', '', '系统管理', 'A01', 1, '', '', '', 0, 1651113401, 1651113401);

-- ----------------------------
-- Table structure for sys_depart_dict
-- ----------------------------
DROP TABLE IF EXISTS `sys_depart_dict`;
CREATE TABLE `sys_depart_dict`  (
  `id` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `depart_id` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '结构ID',
  `lv` int(255) NULL DEFAULT NULL COMMENT '等级',
  `code` bigint(20) NULL DEFAULT NULL COMMENT '编号',
  `serve_code` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '服务器编号',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Compact;

-- ----------------------------
-- Records of sys_depart_dict
-- ----------------------------
INSERT INTO `sys_depart_dict` VALUES ('1cd119d6dd534e4ca5847b659d0bfc27', 'c3e2171228874dc6848c856ee2551df3', 0, 5, 'A');
INSERT INTO `sys_depart_dict` VALUES ('5212e873d8ff446f8c8863c74d680977', 'cb17314537ee4492a3d9a6c5034dfd9f', 0, 1, 'A');
INSERT INTO `sys_depart_dict` VALUES ('5db1b676ca9645a181a10c113f05ebc6', '1d0b1d775cda4244abb4c4c7f7152771', 0, 8, 'A');
INSERT INTO `sys_depart_dict` VALUES ('82bd8d48cd124186b5ade22fb196832a', 'd4e4b0d6101a44559ef0cd8799c5d1a0', 0, 9, 'A');
INSERT INTO `sys_depart_dict` VALUES ('b2bf8876015845cb9345dd097942b7e6', 'b171d9ee59024300a8f8e18facadc4f7', 0, 0, 'A');
INSERT INTO `sys_depart_dict` VALUES ('d10e1dd3e1714f6bb7e5f4b22ef082a1', '419177857cd142ef9692d9c0255a25a2', 0, 6, 'A');
INSERT INTO `sys_depart_dict` VALUES ('e3adb8b1d23f452dafa95745da068761', '1f7eb5b61dbd420c891cf9301bb9239d', 0, 4, 'A');
INSERT INTO `sys_depart_dict` VALUES ('e9158e66e0d84bd1bebec378dae44ac8', '8031a1e774e54152a3952895b0471ce3', 0, 10, 'A');
INSERT INTO `sys_depart_dict` VALUES ('f48f921fd24b440091985d90c327b373', '193d63dbd70d4b36a9456cab3c3ce358', 0, 3, 'A');
INSERT INTO `sys_depart_dict` VALUES ('f4be0354737943a6853e012892d09ab8', 'c4de30b343784b6394a0fa1ff924ecc3', 0, 7, 'A');
INSERT INTO `sys_depart_dict` VALUES ('fdbd1041c4e443c38e167fc4615324ac', '9edf9176834f4c9daa1b4ac78e53a782', 0, 2, 'A');

-- ----------------------------
-- Table structure for sys_depart_rules
-- ----------------------------
DROP TABLE IF EXISTS `sys_depart_rules`;
CREATE TABLE `sys_depart_rules`  (
  `id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `depart_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '组织结构ID',
  `menu_id` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `org_code` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '组织结构编码',
  `create_time` bigint(20) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Compact;

-- ----------------------------
-- Records of sys_depart_rules
-- ----------------------------
INSERT INTO `sys_depart_rules` VALUES ('0409bc6cd6254e6180120650a2d9a549', 'c1a4ad97b18a404cbbe83a31d1ee02d6', 'a5c819b08e50487a88f031c4a308e084', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('07d1519f4d2c4095b7da267c0e2e666d', '9edf9176834f4c9daa1b4ac78e53a782', 'd436aeea798d42a4bea423164f214623', '', 1683893170);
INSERT INTO `sys_depart_rules` VALUES ('0835e687211a4c63a4898e87d0f8a89b', 'c1a4ad97b18a404cbbe83a31d1ee02d6', '8a6758c7e3b34053b394896c6ad7e1ca', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('0a2cc42fddb94c4b9d2df3537203b3fd', 'c1a4ad97b18a404cbbe83a31d1ee02d6', '5be3f23196934fb0a12abb130a591017', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('0c3e41478e614875816cfa570bcb54f8', '1d0b1d775cda4244abb4c4c7f7152771', '165dc6313c8b4bc1ab6662bbfb6bcd90', '', 1709794852);
INSERT INTO `sys_depart_rules` VALUES ('0d7c27a0857a45f2a77d73db98ea1be8', 'c1a4ad97b18a404cbbe83a31d1ee02d6', '4ef5cf77dd8d4528a735fd15ecd5d50d', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('0dbbe86d0dc94125a58feb89d80e5d07', 'c1a4ad97b18a404cbbe83a31d1ee02d6', 'e62df83e3c8a4c279f2e84def1184ac0', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('0f6a9d9fb0d54ba19c91d0fa4c62b2d3', 'c1a4ad97b18a404cbbe83a31d1ee02d6', 'cabf695c50ab4929a95c1595cc3250a8', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('0f9806506c26435db4d4a19821a0e4fc', 'c1a4ad97b18a404cbbe83a31d1ee02d6', '9a4ff26d59f746cfaa67dab0a85bf2ab', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('119d3f0aa8cc4c93a6b16309fddf0727', 'c1a4ad97b18a404cbbe83a31d1ee02d6', '18a9c42ba84846c3a860df02b7c33e2c', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('127548e7f51b472abaa0b4cc5110fdf4', '9edf9176834f4c9daa1b4ac78e53a782', 'b6a35f9b920448c7a58d8d9bdbdb10a4', '', 1683887853);
INSERT INTO `sys_depart_rules` VALUES ('12b509aedd154d9c90bb2ad936514f34', 'c1a4ad97b18a404cbbe83a31d1ee02d6', '19092509384c483494921ab287b38f09', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('133a068c0aab47a3b1495ffe3d9b3828', 'c1a4ad97b18a404cbbe83a31d1ee02d6', 'b44ea7fb5e254fe7aefe201fabe1a272', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('151001a2b4dd49de88236a03b3c63f11', 'c1a4ad97b18a404cbbe83a31d1ee02d6', '855b1407736b4f0b993c2d0e06bcbc8b', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('165032547ee84cc1b7e1c917529fd9b9', 'c1a4ad97b18a404cbbe83a31d1ee02d6', 'd3d90dd60b4e41a796ab64a8b193421f', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('1695d4f66cb24cb89aa9a1134b7be8d3', 'c1a4ad97b18a404cbbe83a31d1ee02d6', '156e4bdb2a6c4be7b7f5c611969ce869', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('17c71d4a2ef64c07933324e5bd4f0d5a', 'c1a4ad97b18a404cbbe83a31d1ee02d6', 'c9849a79a23342a4b3addd4ad8780ec8', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('18ec27652dd946fe834822cc865d5a7b', 'c1a4ad97b18a404cbbe83a31d1ee02d6', 'f3f494263d84485fb2f0893de8808237', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('19634e1376944102ae3baf3abcb00741', 'c1a4ad97b18a404cbbe83a31d1ee02d6', '13bec0ea69d4406ebb8241df60595730', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('19ee0439b11649438e0e9877f65c1cb1', 'c1a4ad97b18a404cbbe83a31d1ee02d6', '6eb4f618a3f94de5aed8685d0b13e460', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('1c09bd928bfd4450bceaf8aebd484c32', 'c1a4ad97b18a404cbbe83a31d1ee02d6', 'c38345c5be0a498a9693eab94652fa03', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('1cda687678c74e078a38d04864e589b7', 'c1a4ad97b18a404cbbe83a31d1ee02d6', '20097cd989e943389b2c249be139f4c7', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('1d68c9356d8948e28fb19076f3b116f9', '1d0b1d775cda4244abb4c4c7f7152771', 'd436aeea798d42a4bea423164f214623', '', 1709794852);
INSERT INTO `sys_depart_rules` VALUES ('1e2d9f06f7454de8a63f92082403d796', 'c1a4ad97b18a404cbbe83a31d1ee02d6', '9c284933cd8c453f84ff2df138604238', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('1e7933d945db48f59a8f88feb43dc731', 'c1a4ad97b18a404cbbe83a31d1ee02d6', 'd603303c35374254a16c3b7f56f5030b', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('1f0a28facbfd4809b8da41dfb9db6ac3', 'c1a4ad97b18a404cbbe83a31d1ee02d6', '5abc7be997bc49d18ac0d1d2b185ee6f', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('1f57fe0d4f174093b0f5b53fecee0da0', 'c1a4ad97b18a404cbbe83a31d1ee02d6', '62a9aeb26bfd4ffc99272baa7fb5b523', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('20f4078b19874b0dbceecb525ea3bd2d', 'c1a4ad97b18a404cbbe83a31d1ee02d6', '622d88add55146cf919212fcb519b688', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('22692ead22a04106be09dcfa07c56650', 'c1a4ad97b18a404cbbe83a31d1ee02d6', '5f7c4cacc2d74c1986d0b1518424a5eb', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('23d92997440a4026b51e71c7b8065d79', '1f7eb5b61dbd420c891cf9301bb9239d', '19092509384c483494921ab287b38f09', '', 1683891482);
INSERT INTO `sys_depart_rules` VALUES ('2582498786a84e0482b541ea404549d0', 'c1a4ad97b18a404cbbe83a31d1ee02d6', '35ab55381d944d71a31e1dcc6865588a', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('277521abeaa5454f8ae3249e7534083f', 'c1a4ad97b18a404cbbe83a31d1ee02d6', '83093300c2fb40f2929802abe4d3d796', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('28a85a44b9d048039acb9e0e17c37014', 'c1a4ad97b18a404cbbe83a31d1ee02d6', '832a5eb4dc014f07b298001c067d1fbe', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('299d5082779c4a3d9b3f8b57441ca95d', '8031a1e774e54152a3952895b0471ce3', '19092509384c483494921ab287b38f09', '', 1709795472);
INSERT INTO `sys_depart_rules` VALUES ('29baedbf44184210a4a07f36f7ce03b8', 'c1a4ad97b18a404cbbe83a31d1ee02d6', 'e91ab6c14efb425596bdca621cdf314c', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('2bca3aaf5886467c9f04189ad0ddb112', 'c1a4ad97b18a404cbbe83a31d1ee02d6', 'f93f4dfc390741dda7fe9c4ba549ff00', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('2dea6f837c6140cfb343546328c32644', '8031a1e774e54152a3952895b0471ce3', '65ac433016984495ba9c62a0d381ed0b', '', 1709795472);
INSERT INTO `sys_depart_rules` VALUES ('2dec00f2d69d4889b757236c50b5ada4', 'c1a4ad97b18a404cbbe83a31d1ee02d6', '954f3644dc984b1a8c971ccf2b3ca7f5', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('2edc1f8eb3a44c498055ed42d466cdd8', 'c1a4ad97b18a404cbbe83a31d1ee02d6', '18e01cceb53a4dc88a766b47c7172e1e', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('3065b0b507db4d8fac1c7d08c88f9ffd', '8031a1e774e54152a3952895b0471ce3', 'd436aeea798d42a4bea423164f214623', '', 1709795472);
INSERT INTO `sys_depart_rules` VALUES ('3082114aa32340fd8cdad3a6aea64788', 'c1a4ad97b18a404cbbe83a31d1ee02d6', 'd95331a2c13e41f391f57890abd8fe0f', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('3091cbd0420a430d970c84f6c2de597f', 'c1a4ad97b18a404cbbe83a31d1ee02d6', '4afaa4b996424125a5913ed897d22d92', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('30f1533f99ea451db454d187faf483b6', 'c1a4ad97b18a404cbbe83a31d1ee02d6', 'b307984ef7c347d3b97f13faa70c76d7', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('31deb95ef02b457ba821befa6a54e55a', '8031a1e774e54152a3952895b0471ce3', '6fb3210467194b64a6c60a90f0a6e3c5', '', 1709795472);
INSERT INTO `sys_depart_rules` VALUES ('35936462a3f1401f974b916022e996c4', '1d0b1d775cda4244abb4c4c7f7152771', '44030ff97bf7446693ace817f71dfb7a', '', 1709794852);
INSERT INTO `sys_depart_rules` VALUES ('382c84f19055422aba77ddae29aca3b4', 'c1a4ad97b18a404cbbe83a31d1ee02d6', '3bf0ddd871aa4d1ca0df67af445b2f28', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('3a47c33126954802ac9d9e6d49eab15d', 'c1a4ad97b18a404cbbe83a31d1ee02d6', 'cc3752c9425a4b04a122eafc4a925fb3', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('3ab2305b47ee41cb86df24a1e468eae8', '9edf9176834f4c9daa1b4ac78e53a782', '243ba9a650284136b4b4a4ed8c5d9979', '', 1683893170);
INSERT INTO `sys_depart_rules` VALUES ('3b5c392752bd45e99698ee9cd19ba194', 'c1a4ad97b18a404cbbe83a31d1ee02d6', 'd022725a9885478e87cbaa5201245934', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('3bd1806b93494623872cf46509e0ce7b', 'c1a4ad97b18a404cbbe83a31d1ee02d6', '1cea1b2bf8624ff9bb62eb449899018f', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('3cf4691913f6429d92183a250a3ed27c', 'c1a4ad97b18a404cbbe83a31d1ee02d6', 'b005005a431f48bf89254bf7405d4cc3', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('3e1af39c7d6d415baf57d4ceacfab7e3', 'c1a4ad97b18a404cbbe83a31d1ee02d6', '486fcfa720a547c582f00fc9bdaba89f', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('3e224cc7f48744ae935f44d36fdc4321', '1d0b1d775cda4244abb4c4c7f7152771', '65ac433016984495ba9c62a0d381ed0b', '', 1709794852);
INSERT INTO `sys_depart_rules` VALUES ('3ed2566ed45e4e3c8ba31b690eb27ab9', 'c1a4ad97b18a404cbbe83a31d1ee02d6', 'bad57780a5454bf7a903033490794e59', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('3f7d11176a3e4f51816c61dec4ba5a84', 'c1a4ad97b18a404cbbe83a31d1ee02d6', 'a832aea449084fb1a991209a5383acdb', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('407d166f88734a8583cfb7e1b56a655d', 'c1a4ad97b18a404cbbe83a31d1ee02d6', 'd57ee3f99c5745e69efdf4fe9fe3cfd0', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('40f8393d6014472a8f6fe94add781a2e', 'c1a4ad97b18a404cbbe83a31d1ee02d6', '4474facaa77c4d99baee20136b6eda98', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('41cd8a562b434f5985195ee4f88c29bb', 'c1a4ad97b18a404cbbe83a31d1ee02d6', '09e2c124ff06405db60e621249734c87', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('421ebf7764fd47deaa5eccb9c15d8bf0', 'c1a4ad97b18a404cbbe83a31d1ee02d6', '9e3c7ec8fce940fb8914b88211c91613', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('433b4df4516347e7815084a24ab95b55', 'c1a4ad97b18a404cbbe83a31d1ee02d6', '3656d30ad4c64671bed5728ba4b622bf', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('44c239d256074e85acda8e2008ed28a7', 'c1a4ad97b18a404cbbe83a31d1ee02d6', '66e34fb2873f4596af2f5cf4f0fd1db3', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('44e055d504594fbfb1095df9c1be7d5a', 'c1a4ad97b18a404cbbe83a31d1ee02d6', 'eb12ba7d1fe946c485f3169ba5f26fae', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('45c24c4ce44249ad8cb24f3ce4614a81', '1f7eb5b61dbd420c891cf9301bb9239d', 'd436aeea798d42a4bea423164f214623', '', 1683891482);
INSERT INTO `sys_depart_rules` VALUES ('46ce7208a9a646808ad5e9efaf1fd919', 'c1a4ad97b18a404cbbe83a31d1ee02d6', '6b441948f2f640deacb9cc75e97d3d33', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('48fdeabc9b9146dbb4d4b7a8232a9b10', 'c1a4ad97b18a404cbbe83a31d1ee02d6', '7f0da824fb5746d28ab38220c0ca854c', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('4a1f52e0194b4609b3994422cadc9b73', 'c1a4ad97b18a404cbbe83a31d1ee02d6', 'f84f2fa3cc134066a4f031b0f4a1b7bf', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('4b494ea05b0642d5ace55311f91d4222', '1f7eb5b61dbd420c891cf9301bb9239d', 'b6a35f9b920448c7a58d8d9bdbdb10a4', '', 1683891482);
INSERT INTO `sys_depart_rules` VALUES ('4c0790e30b57482bbd7221a61caa8545', 'c1a4ad97b18a404cbbe83a31d1ee02d6', '5328353627c64f5da3d0ad97664e9094', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('4d322f0b28994b3bb30a6446272644b8', 'c1a4ad97b18a404cbbe83a31d1ee02d6', 'e2edce39bbab4ff7b477d8bc1ff46355', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('4d5e6c447c1442fb82153c91f1794500', 'c1a4ad97b18a404cbbe83a31d1ee02d6', '6e7229f6397747858b0821f6be9e05b0', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('4e15460b626740729d0e854b0658a8dd', 'c1a4ad97b18a404cbbe83a31d1ee02d6', 'e438ca5d5b7246f28e094a6968aaaed0', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('4ec83dd10fd743f7944ba0f609144ca2', 'c1a4ad97b18a404cbbe83a31d1ee02d6', 'fdf83b8f163444eeae67c8687dd000d7', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('52928f9891b142c3920609e51ebab568', 'c1a4ad97b18a404cbbe83a31d1ee02d6', 'ff96ba3fd99b44cd9e98be2f9ea731b8', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('5482bc8f2b8c4594a9338374affaf862', 'c1a4ad97b18a404cbbe83a31d1ee02d6', 'a0dc332d2a4744af98e902eadd547836', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('5893984b6f824b3ab1378af4bff6a968', 'c1a4ad97b18a404cbbe83a31d1ee02d6', 'd44350eb2d654815859f7b28f5115562', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('5a6f3715805b4b29b7072dd3e74bee35', 'c1a4ad97b18a404cbbe83a31d1ee02d6', '1ac1e90c0135411dac881063de5a299c', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('5b3287e211f6454f8ed5f2f7fdc3c174', 'c1a4ad97b18a404cbbe83a31d1ee02d6', '20047af77f794102beefd019f5214f03', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('5f5d8c287a3741c2a56b318e3e323151', 'c1a4ad97b18a404cbbe83a31d1ee02d6', '911304671b514f57bdc849b29cf7b2a9', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('5f91c0c36ea64e64ab910e368ffcfa2b', 'c1a4ad97b18a404cbbe83a31d1ee02d6', '32410d40574b49b2b799c0ebfc65b6f3', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('6238eb0341c0432187685c9511a86c5f', 'c1a4ad97b18a404cbbe83a31d1ee02d6', '6d263b4c58314e5b9517bb8838f5bac3', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('63575fe147334e30bd496f8015e93c65', 'c1a4ad97b18a404cbbe83a31d1ee02d6', 'c7663bdc393b48ac8d4c02970e95bab4', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('64e69300dba74b1a80e5bb628c3dbe02', 'c1a4ad97b18a404cbbe83a31d1ee02d6', '89b237003b1b4bc08fd2060eda7d6f48', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('650f002d906148929ea60cd0e238fc3c', 'c1a4ad97b18a404cbbe83a31d1ee02d6', '72fce948f50d4e2284f5ec2f260a2284', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('68caa01ef3bf46b785b405593e1f74f4', 'c1a4ad97b18a404cbbe83a31d1ee02d6', 'fc6e241275ab4b12b5e1d72835c5bbd7', '', 1683887833);
INSERT INTO `sys_depart_rules` VALUES ('697ef123357f478790f67260fe67f73f', 'c1a4ad97b18a404cbbe83a31d1ee02d6', 'de20c758be4047439fc96897907b2b80', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('6c6bfc76c60349bd82cb56438b033840', 'c1a4ad97b18a404cbbe83a31d1ee02d6', 'edb5b31a0f194824aebbc5b68435ad6c', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('6db86b012ac04d39a5ae62877c81d4ad', 'c1a4ad97b18a404cbbe83a31d1ee02d6', '01403f88847747b3bee767c7cc8edf7f', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('6f4ea443f020446b8eac831cd23febd4', 'c1a4ad97b18a404cbbe83a31d1ee02d6', '7fe5f9f9843145099761ac50aee0a3ac', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('6fbe58a32c324a01998168b50e328766', 'c1a4ad97b18a404cbbe83a31d1ee02d6', '84cb2653afa7499e829f00bc0c8367f7', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('710552239b4a4dd4afac380dbd7b8417', 'c1a4ad97b18a404cbbe83a31d1ee02d6', '45bab7125b454c23882c83d1f2103437', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('72198c6e194f4a7b977e2b18fdfd9aa9', '1d0b1d775cda4244abb4c4c7f7152771', '84cb2653afa7499e829f00bc0c8367f7', '', 1709794852);
INSERT INTO `sys_depart_rules` VALUES ('72276f3fe7324208bdefae0c6bfdc851', 'c1a4ad97b18a404cbbe83a31d1ee02d6', 'a6f32880e71540d7aa6b7a24a27c1b6d', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('74cd21060cda4800839e40f224395bd6', '9edf9176834f4c9daa1b4ac78e53a782', '3a09d8911f1f401898a3d4bd3018e217', '', 1683887853);
INSERT INTO `sys_depart_rules` VALUES ('75f2f4e0d8cf4df797106bbcbc25251a', 'c1a4ad97b18a404cbbe83a31d1ee02d6', 'f184f4d0ac3b45ee9996dd530abbff1e', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('77678c3ab3484daa9a957b19ae2ec69f', '8031a1e774e54152a3952895b0471ce3', 'dc3339db28c44e0fbc602cd44086094d', '', 1709795472);
INSERT INTO `sys_depart_rules` VALUES ('782ea27e33c540ffaee8ff4e80c2e00b', '9edf9176834f4c9daa1b4ac78e53a782', '702d3b0676524eaa97de61db2aeb3616', '', 1683887853);
INSERT INTO `sys_depart_rules` VALUES ('7bd73f3dc70c4f0488ea4beb4f007dc5', '1d0b1d775cda4244abb4c4c7f7152771', '19092509384c483494921ab287b38f09', '', 1709794852);
INSERT INTO `sys_depart_rules` VALUES ('7e651fcc813e4558b16fbe7e3136434d', 'c1a4ad97b18a404cbbe83a31d1ee02d6', 'f86479e189a64811987f6610d96fb96f', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('7f3ebd016935497c87ea37cfb1ea7fb7', 'c1a4ad97b18a404cbbe83a31d1ee02d6', '38b36df92d3343e4834316411997ee8b', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('8066a69ab535442eb1f5e9bf3cbfd79a', '1d0b1d775cda4244abb4c4c7f7152771', '702d3b0676524eaa97de61db2aeb3616', '', 1709794852);
INSERT INTO `sys_depart_rules` VALUES ('8072ad474acd4d33b8bc80d40472e15a', 'c1a4ad97b18a404cbbe83a31d1ee02d6', 'f2b3a035b9e147efbd899701d2582885', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('80aaea094d364aeaa231167612d65cc6', 'c1a4ad97b18a404cbbe83a31d1ee02d6', '9312704b69b64c43803847abaf791c44', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('8296948a4d8946f69912766c171c6b9f', '1f7eb5b61dbd420c891cf9301bb9239d', '934819acada046a5a29577df449aec68', '', 1683894737);
INSERT INTO `sys_depart_rules` VALUES ('82b64b3e9c804fdd83e58fbb45d0db60', 'c1a4ad97b18a404cbbe83a31d1ee02d6', 'ad525262e44146ef85e31392dceb35fb', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('848b100a38c0478987e68163cf4b4af7', '1f7eb5b61dbd420c891cf9301bb9239d', '84cb2653afa7499e829f00bc0c8367f7', '', 1683891482);
INSERT INTO `sys_depart_rules` VALUES ('84cf98126bd2491c89ab4ddea90acf63', 'c1a4ad97b18a404cbbe83a31d1ee02d6', 'f532191778fd4c47a616f9eda594bf9a', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('872a15f6cced42beaaa6637aaa2278d5', 'c1a4ad97b18a404cbbe83a31d1ee02d6', '3f8d8d059c5945aa9b6bb01a7c811abb', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('878921a2b7cd420b9dfdfd0c6f1cf749', 'c1a4ad97b18a404cbbe83a31d1ee02d6', '3a09d8911f1f401898a3d4bd3018e217', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('885d27dc9bd540008d8fdb1e67cd8687', 'c1a4ad97b18a404cbbe83a31d1ee02d6', 'cd9879099e7d4dc194d2c20fa26e8e4b', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('8916a984916c479fbf1942b6357c57d4', 'c1a4ad97b18a404cbbe83a31d1ee02d6', '6f7c1d1071a3402a82d91ce6795eec15', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('895c950a758f423394b852664dcabdff', 'c1a4ad97b18a404cbbe83a31d1ee02d6', '2dd7748ba05643e8ba8790b7f884ad31', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('8a4b662356804abc9ef7262489b09a8a', 'c1a4ad97b18a404cbbe83a31d1ee02d6', '2fc98de902cb400cbb54689fa8d40e87', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('8ab579a439a545b78d8a0235b395fd71', 'c1a4ad97b18a404cbbe83a31d1ee02d6', '1d7e814f789e4454ba5b3454f91a3d0f', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('8cb6d7ff1f594027ba86aff2ce4b1c10', 'c1a4ad97b18a404cbbe83a31d1ee02d6', '057965ca4b704d86b0a5dd57ea29e535', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('8cdf17ee3dda4c1685097c0db6f01359', 'c1a4ad97b18a404cbbe83a31d1ee02d6', '01c411fba224457b80e7d112024c2057', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('8e8263b64d1740ef9556fa197c394435', 'c1a4ad97b18a404cbbe83a31d1ee02d6', '5d377664250d4733997f13d971fe8b26', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('8fcc307f0693422eb0737c7d10944a25', 'c1a4ad97b18a404cbbe83a31d1ee02d6', '09b64c379dcb4d59aecb8f8e9deead72', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('903e97b541314b4eb761db8459d659bf', 'c1a4ad97b18a404cbbe83a31d1ee02d6', '7e36c639a2de4a1b83b3f45e81fe80bb', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('906d098856b344348f4c2680e2101aa8', 'c1a4ad97b18a404cbbe83a31d1ee02d6', '82990002f388438eb4c59d82b88ff3f6', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('912755d9344e42259948e188d3dc14d6', 'c1a4ad97b18a404cbbe83a31d1ee02d6', 'f8700a89ab52418ab45d3754d3c2adcb', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('91fca086d94147f28183eba5a0019f01', 'c1a4ad97b18a404cbbe83a31d1ee02d6', 'c9850f0b1908434e86c5b6ea2da9e5e5', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('92bad2ec60a5434c8da6394c2af661df', 'c1a4ad97b18a404cbbe83a31d1ee02d6', '6709571ddf0548c6b77c8ee94cf2069e', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('93111dbb66fb4ca5953ec4da77ee9ee1', 'c1a4ad97b18a404cbbe83a31d1ee02d6', '880ef7398ed644faa71242c0189be1e4', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('93bf6a04568c498688666395a6bd3fcb', 'c1a4ad97b18a404cbbe83a31d1ee02d6', '737252cc760e46eb9b2182765f1691a5', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('969d0d029ec34062be52a02670dc251a', 'c1a4ad97b18a404cbbe83a31d1ee02d6', '19085cde4f90404680a5ce14e274cdca', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('97c9e13946834c91bf821520adb7df00', 'c1a4ad97b18a404cbbe83a31d1ee02d6', '8835797e8ec44bb490ac36d6d4d48af1', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('97d7bb34cee548959440076d75f69557', 'c1a4ad97b18a404cbbe83a31d1ee02d6', '56f80b704bf5446880e8ca9d5620e86e', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('98b12a3a6a66470397feed3d56e21770', 'c1a4ad97b18a404cbbe83a31d1ee02d6', 'a875d9d7e3534895b8c26da32bb94923', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('98b39978e39d480e9717f8e2d5835847', 'c1a4ad97b18a404cbbe83a31d1ee02d6', '2115daaa2d294fc2b9be5e85d645a3db', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('991b3ba5eb7c410d8caea2305e0727d2', '8031a1e774e54152a3952895b0471ce3', '165dc6313c8b4bc1ab6662bbfb6bcd90', '', 1709795472);
INSERT INTO `sys_depart_rules` VALUES ('9a490840419e483c9812be78faaa0c6c', 'c1a4ad97b18a404cbbe83a31d1ee02d6', '819032e614c84a83aa64737f0ff43863', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('9b5534095d334dee8ccc002d57be2a22', 'c1a4ad97b18a404cbbe83a31d1ee02d6', '253d65d769374c0f9f7d0f4e2c741e3f', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('9c1c0ff27a074d8a907acb30dabf8da5', 'c1a4ad97b18a404cbbe83a31d1ee02d6', 'cdc8ee1068914109890ae9724930bd52', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('9c82146604c0416a92fb8c944696b234', 'c1a4ad97b18a404cbbe83a31d1ee02d6', 'fb7982029cef479d85bf7d7b3088aa3b', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('9cc921fa9e134f73989146bd3e22cb2f', 'c1a4ad97b18a404cbbe83a31d1ee02d6', '7a79d04c3707456893123fb565f37819', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('9db273ec3e6e4c1ca3bd971bbd01d53f', '1d0b1d775cda4244abb4c4c7f7152771', '9b144acec61e4886bdcbf0e6ad333481', '', 1709794852);
INSERT INTO `sys_depart_rules` VALUES ('a13345f558764003af7d01e2c43b03be', 'c1a4ad97b18a404cbbe83a31d1ee02d6', '67772187adf54df2ab9869176a4b5851', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('a1efaefd60b346ffa2b57838ca1b240a', 'c1a4ad97b18a404cbbe83a31d1ee02d6', '6ca835c0146d4e489d6befb6b8605a8f', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('a39a2077621949a0a83b5b4c3e271d6f', 'c1a4ad97b18a404cbbe83a31d1ee02d6', '0e1ae963f43a429a934caba77f0055ad', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('a69e15f352874ab79e2b7118cc433472', 'c1a4ad97b18a404cbbe83a31d1ee02d6', '215e7a8f3ea3408186113f9a839bae04', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('a75013ebdbb54b418cec76e78809b839', 'c1a4ad97b18a404cbbe83a31d1ee02d6', '2523741989d84b0ebcd8b025fd6a4462', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('a825913749ff4dd0add84f38807f4653', 'c1a4ad97b18a404cbbe83a31d1ee02d6', '96f87edd59ba42389a5aac2d9881f68d', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('a964f2ea36044c55aee80d2956202486', 'c1a4ad97b18a404cbbe83a31d1ee02d6', '8bbd3cc24fab4566afc0fb2142765e92', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('aa581428cad24f8fa3889d56c17a70ce', 'c1a4ad97b18a404cbbe83a31d1ee02d6', '702d3b0676524eaa97de61db2aeb3616', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('ad1c6e8ccfab4c9d9056e97f7697c632', 'c1a4ad97b18a404cbbe83a31d1ee02d6', '115d3121a44c4eabb87829a619bb1307', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('adf42656724b42f0bdb87d57384c506a', 'c1a4ad97b18a404cbbe83a31d1ee02d6', '63b4ed9e2ebc475787811df1ab12b2fb', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('aed521ad6bdd48c4800ca86945b7c261', 'c1a4ad97b18a404cbbe83a31d1ee02d6', 'd72806f73c5d48cd8c0b4cd3fbdb6e1c', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('af1c7fa2edb746aea5c25af4f11b8b36', '8031a1e774e54152a3952895b0471ce3', '84cb2653afa7499e829f00bc0c8367f7', '', 1709795472);
INSERT INTO `sys_depart_rules` VALUES ('b119c66ffe1d485bb70e7b5329e279cd', 'c1a4ad97b18a404cbbe83a31d1ee02d6', '3550d6cb7501456fb38db084adf32829', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('b3870d11593e40b4815ade782c9271a9', 'c1a4ad97b18a404cbbe83a31d1ee02d6', 'f52e7b02ae814df0889ec1436f8beac0', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('b553e0ef1c734f718bd1a206f791f2c1', 'c1a4ad97b18a404cbbe83a31d1ee02d6', 'db69a091af6e45f5af902c45c8b70bf2', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('b6b3ca14bfe943048615d0015e109fdc', '1f7eb5b61dbd420c891cf9301bb9239d', '3a09d8911f1f401898a3d4bd3018e217', '', 1683891482);
INSERT INTO `sys_depart_rules` VALUES ('b7d5277b925e40cdbd25f23691b30032', 'c1a4ad97b18a404cbbe83a31d1ee02d6', 'c86e0d8a96b24acc9883919c88b66932', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('b8f74214d9e045bcbe7f7adbb5f14fa5', '9edf9176834f4c9daa1b4ac78e53a782', '19092509384c483494921ab287b38f09', '', 1683887853);
INSERT INTO `sys_depart_rules` VALUES ('bb9a064c314f4ebb8b719156de7a331f', '1f7eb5b61dbd420c891cf9301bb9239d', '5661a26a8174420e8fa0a8e05fce29bd', '', 1683894737);
INSERT INTO `sys_depart_rules` VALUES ('be1d7ac99df94b7f9da46b29df99eecb', 'c1a4ad97b18a404cbbe83a31d1ee02d6', '0d57ad05447d49b281638270742a0fb0', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('be243a99c0a64ef7b3a38b84e8b127d2', 'c1a4ad97b18a404cbbe83a31d1ee02d6', '032fac52543747faaea6aea276abdc44', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('bf2e27af48a945e781c90a024d61f967', 'c1a4ad97b18a404cbbe83a31d1ee02d6', '337af61f06394a2bbc120854628f16d6', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('c1dea2ed2adc4e7586f8439f21c63a10', 'c1a4ad97b18a404cbbe83a31d1ee02d6', 'd053fc0d87f048219f66190ecc2e709c', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('c1f6653710064be1a4e110e1bbaadaa5', 'c1a4ad97b18a404cbbe83a31d1ee02d6', '9b144acec61e4886bdcbf0e6ad333481', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('c39c06411cee4e5abe75a73e9ad33185', 'c1a4ad97b18a404cbbe83a31d1ee02d6', '41812fcfce4e4cbf9ac8214825dd8fd9', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('c3c9183173114b738194e9b4bd8943e2', 'c1a4ad97b18a404cbbe83a31d1ee02d6', 'bd553cd0fe7441259505b36607b4e99b', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('c42f2fe1a131427b88cef2d0cab1a8e9', 'c1a4ad97b18a404cbbe83a31d1ee02d6', '4e33f9ea715243dd9e08fe4fe5478eff', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('c432e4255b014b35b51ecdf418c35ce5', 'c1a4ad97b18a404cbbe83a31d1ee02d6', '7b464ec84d4b46bf9fe9b38f6ed59559', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('c51d950406b042c897cb34a50031e05a', '9edf9176834f4c9daa1b4ac78e53a782', '84cb2653afa7499e829f00bc0c8367f7', '', 1683887853);
INSERT INTO `sys_depart_rules` VALUES ('c5f11bc75f7440cea71b5915beebf2b7', 'c1a4ad97b18a404cbbe83a31d1ee02d6', 'd436aeea798d42a4bea423164f214623', '', 1698992959);
INSERT INTO `sys_depart_rules` VALUES ('c81a13f3c1cd4a9fb039819d0c59a443', '8031a1e774e54152a3952895b0471ce3', '702d3b0676524eaa97de61db2aeb3616', '', 1709795472);
INSERT INTO `sys_depart_rules` VALUES ('ca488e5ae6984bdab0d3ec9ac2965697', '1f7eb5b61dbd420c891cf9301bb9239d', '243ba9a650284136b4b4a4ed8c5d9979', '', 1683894737);
INSERT INTO `sys_depart_rules` VALUES ('ca90492d5e284da1ae03b47c4923f719', '9edf9176834f4c9daa1b4ac78e53a782', '9b144acec61e4886bdcbf0e6ad333481', '', 1683887853);
INSERT INTO `sys_depart_rules` VALUES ('cc5db6838b444284a7e5115847d97d84', '1f7eb5b61dbd420c891cf9301bb9239d', '702d3b0676524eaa97de61db2aeb3616', '', 1683891482);
INSERT INTO `sys_depart_rules` VALUES ('cc72cbb884ac4ab1a303d4484db701ca', 'c1a4ad97b18a404cbbe83a31d1ee02d6', 'f0074822910b4226a0f0a741d42ba3fa', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('ccc925aef8434137ae3216a9a8f44ffb', 'c1a4ad97b18a404cbbe83a31d1ee02d6', '4c9d910eee54499faaf6e65008b9f844', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('ccd323dec8fa4772b713663f05a2c8e1', 'c1a4ad97b18a404cbbe83a31d1ee02d6', '71906de352334e7eb5c6cf9c033c8694', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('ceb6b60e75164d91a68ba4dac9a47fde', 'c1a4ad97b18a404cbbe83a31d1ee02d6', '32f7c11b48104cef9f027829f78a4fb3', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('d026e904ad85403c99741807837b3cfd', 'c1a4ad97b18a404cbbe83a31d1ee02d6', '7271242be4eb4e1893c176b28e63a317', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('d16a62e5910446d1935ea5837cc6bde5', 'c1a4ad97b18a404cbbe83a31d1ee02d6', '9eec65bd7d2d41a8a5680d0870939ba2', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('d1d9745a25054e07a90b5d0a3dd5af43', '1f7eb5b61dbd420c891cf9301bb9239d', 'de20c758be4047439fc96897907b2b80', '', 1683891482);
INSERT INTO `sys_depart_rules` VALUES ('d3e708017fdf44cb8ae328c7b1834fc0', 'c1a4ad97b18a404cbbe83a31d1ee02d6', 'e3d2983be6c74d45b6b59cfcf416ca23', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('d4d04830f4134d05aab69bba432b868f', 'c1a4ad97b18a404cbbe83a31d1ee02d6', 'de418390eab9451087118bf0f35614fe', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('d611ba0e98f6451580f20ea2b907f344', 'c1a4ad97b18a404cbbe83a31d1ee02d6', 'e49bc585e15142139933c5ed455d5896', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('d6c12ebbce4a47d0944a58817719d688', 'c1a4ad97b18a404cbbe83a31d1ee02d6', 'e6a24e81354245d48e7e7bf49741f54d', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('d796875407c24247ad37878798dcca83', 'c1a4ad97b18a404cbbe83a31d1ee02d6', '0538b5fa5f6c4ec3b9ed40fe84765a12', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('d89da08623ad4f07b3923537143bdfcb', '8031a1e774e54152a3952895b0471ce3', '44030ff97bf7446693ace817f71dfb7a', '', 1709795472);
INSERT INTO `sys_depart_rules` VALUES ('d8cfdd35a9044a35b83abc78a2e67742', 'c1a4ad97b18a404cbbe83a31d1ee02d6', 'c24a50e8069946d6997b48d73681a2d9', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('daf7f5a81122454098b100c0e70328f3', 'c1a4ad97b18a404cbbe83a31d1ee02d6', 'b6a35f9b920448c7a58d8d9bdbdb10a4', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('db40080761e94bf085099147dc3eefe3', 'c1a4ad97b18a404cbbe83a31d1ee02d6', '5512c3b60840403cb076d328cdc9d427', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('dc6c099fdc9c464d8429a34ae97e9a85', 'c1a4ad97b18a404cbbe83a31d1ee02d6', '32f14c776a75403680a580681244c6ce', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('df9bb9977130442bb5caf6efb487aed5', 'c1a4ad97b18a404cbbe83a31d1ee02d6', '75b0714720be444ca6b78d5dda6d6bb8', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('e1f06a9de76743f5a9c3150615ec2602', 'c1a4ad97b18a404cbbe83a31d1ee02d6', 'b97751d4c7cb447cb6ec1149f5b2c7c2', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('e204c7595ae949958a26b34308978459', 'c1a4ad97b18a404cbbe83a31d1ee02d6', '344be9b8528a4095bad35757ec356171', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('e21223e9878a489394f4452f7a920839', 'c1a4ad97b18a404cbbe83a31d1ee02d6', '85a0db15af4d4d6fa8b88f77e3157d7d', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('e454b8e84b854c01a8df3e72730cf4a4', '8031a1e774e54152a3952895b0471ce3', '9b144acec61e4886bdcbf0e6ad333481', '', 1709795472);
INSERT INTO `sys_depart_rules` VALUES ('e65317737c574813827d6e7b6f97da6e', 'c1a4ad97b18a404cbbe83a31d1ee02d6', '1f1d165cde3448d8a0840ad71feeb6fc', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('ee1c9c7b04e04129839a5d59046221ee', 'c1a4ad97b18a404cbbe83a31d1ee02d6', '0cb43590702e487ab34f6ad37f566dc7', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('f26305726f644b778d912774c3c461e3', 'c1a4ad97b18a404cbbe83a31d1ee02d6', '1d1c2716bedb44f59945e46375f874ee', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('f2c1adeb50934059a4d2a39ba0d88a0a', 'c1a4ad97b18a404cbbe83a31d1ee02d6', '7f0ea4df11e9419992d691dde2502e91', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('f33dc6ed9ff044d9a01a046b38e6a62d', 'c1a4ad97b18a404cbbe83a31d1ee02d6', '6c0415d763e94c6c8bc9394d1b554a81', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('f4564f15543a4683b9f433a735dcc2c3', 'c1a4ad97b18a404cbbe83a31d1ee02d6', 'd5bbea09765f4c5893ed025f17b1310b', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('f5fa084bac9349ffa68e67aa03a120de', 'c1a4ad97b18a404cbbe83a31d1ee02d6', '508b8b96ff244df9b5e8ec79d3d2be30', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('f6b8aa688ea64b0e80923cc0faf71c74', '1f7eb5b61dbd420c891cf9301bb9239d', '9b144acec61e4886bdcbf0e6ad333481', '', 1683891482);
INSERT INTO `sys_depart_rules` VALUES ('f79de1b3c5354a21a5059f0f6a8ab648', 'c1a4ad97b18a404cbbe83a31d1ee02d6', 'c668dc2debea4d82bd307594dab1579c', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('faa4746a68324cdd937104c221753f6e', 'c1a4ad97b18a404cbbe83a31d1ee02d6', 'e412b048248f455b8a12a50619679b0a', '', 1683887250);
INSERT INTO `sys_depart_rules` VALUES ('fb38eb3c74f246b78928ae68f61f8a9c', '1d0b1d775cda4244abb4c4c7f7152771', '6fb3210467194b64a6c60a90f0a6e3c5', '', 1709794852);
INSERT INTO `sys_depart_rules` VALUES ('fd974b269df54ae79179c939a3a0f5a9', '9edf9176834f4c9daa1b4ac78e53a782', 'de20c758be4047439fc96897907b2b80', '', 1683887853);
INSERT INTO `sys_depart_rules` VALUES ('ff36407f57c0406bb002910976db8887', 'c1a4ad97b18a404cbbe83a31d1ee02d6', '77df1345d01f42f4a7f07a84f9f9d77f', '', 1683887250);

-- ----------------------------
-- Table structure for sys_dict_list
-- ----------------------------
DROP TABLE IF EXISTS `sys_dict_list`;
CREATE TABLE `sys_dict_list`  (
  `id` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `parent_id` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `parent_name` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '字典名称',
  `value` int(10) NULL DEFAULT NULL COMMENT '值',
  `name` varchar(10) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '名称',
  `status` int(1) NULL DEFAULT NULL COMMENT '1正常，2停用，10已删除',
  `sort` int(10) NULL DEFAULT NULL COMMENT '排序',
  `type` int(1) NULL DEFAULT NULL COMMENT '1自定义 2 数据表',
  `table_key` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '数据表主键',
  `table_val` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '数据包值',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Compact;

-- ----------------------------
-- Records of sys_dict_list
-- ----------------------------
INSERT INTO `sys_dict_list` VALUES ('3f20d7bb21af45f4877fc4a10dea73fa', '', 'DepartType', 0, '组织结构类型', 1, 1, 1, NULL, NULL);
INSERT INTO `sys_dict_list` VALUES ('41d6d569364f4da1a9d5384c8aa21bff', '', 'sysMenus', 0, '菜单列表', 1, 1, 2, 'id', 'name');
INSERT INTO `sys_dict_list` VALUES ('59f52e645cb24125b095e74e4668236f', '3f20d7bb21af45f4877fc4a10dea73fa', '', 2, '公司', 1, 2, NULL, NULL, NULL);
INSERT INTO `sys_dict_list` VALUES ('6dcd9df894aa4dbea5c46736e74bf300', '3f20d7bb21af45f4877fc4a10dea73fa', '', 4, '门店', 1, 4, NULL, NULL, NULL);
INSERT INTO `sys_dict_list` VALUES ('7010f3e446a243eab7fab5f393f6dccb', '', 'sysRole', 0, '系统角色', 1, 3, 2, 'id', 'role_name');
INSERT INTO `sys_dict_list` VALUES ('89919246cc434591b4de187e5eaa84db', '3f20d7bb21af45f4877fc4a10dea73fa', '', 3, '部门', 1, 3, NULL, NULL, NULL);
INSERT INTO `sys_dict_list` VALUES ('9d08a40cbf6a4d7c8aaca56a9191dc41', '3f20d7bb21af45f4877fc4a10dea73fa', '', 1, '集团', 1, 1, NULL, NULL, NULL);
INSERT INTO `sys_dict_list` VALUES ('b9dcda059190495d9be6b659424b3950', '', 'sysDepart', 0, '组织结构', 1, 2, 2, 'id', 'name');

-- ----------------------------
-- Table structure for sys_member
-- ----------------------------
DROP TABLE IF EXISTS `sys_member`;
CREATE TABLE `sys_member`  (
  `id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `depart_id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '组织结构ID',
  `role_id` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '角色ID',
  `uid` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '会员ID',
  `username` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '登陆账号',
  `password` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '密码',
  `nickname` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '昵称',
  `realname` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '真实姓名',
  `headimgurl` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '头像',
  `mp` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '手机号',
  `idcard` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '身份证号码',
  `sex` int(1) NULL DEFAULT NULL COMMENT '性别 1男2女',
  `status` int(1) NULL DEFAULT NULL COMMENT '1正常，2禁用',
  `org_code` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '组织结构编码',
  `create_time` bigint(20) NULL DEFAULT NULL COMMENT '创建时间',
  `update_time` bigint(20) NULL DEFAULT NULL COMMENT '更新时间',
  `files` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '附件',
  `remark` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '备注',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Compact;

-- ----------------------------
-- Records of sys_member
-- ----------------------------
INSERT INTO `sys_member` VALUES ('72bf2cad7f6f4a3f8acd8d4808885293', 'c1a4ad97b18a404cbbe83a31d1ee02d6', '3aa86b48d87d4f118546906769b51df5', '', 'admin', '21232f297a57a5a743894a0e4a801fc3', '管理员', '管理员', '', '13333333333', '360425198806214612', 1, 1, 'A01', 1651113987, 1709724875, '', NULL);
INSERT INTO `sys_member` VALUES ('ddd0ef7c9ea643499b444e21bbf002e9', '8031a1e774e54152a3952895b0471ce3', 'a1ff8273cdf243f6972070a51323e0b4', '', 'demo', '21232f297a57a5a743894a0e4a801fc3', '', '测试', '', '13333333333', '3600000000000000', 1, 1, 'A10', 1709794908, 1709795330, '', '测试专用');

-- ----------------------------
-- Table structure for sys_menus
-- ----------------------------
DROP TABLE IF EXISTS `sys_menus`;
CREATE TABLE `sys_menus`  (
  `id` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `pid` varchar(64) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `name` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `component` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '组件地址',
  `icon` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '图片',
  `is_ext` int(1) NULL DEFAULT NULL COMMENT '是否外链',
  `keepalive` int(1) NULL DEFAULT NULL COMMENT '是否缓存',
  `show` int(1) NULL DEFAULT NULL COMMENT '是否显示',
  `type` int(1) NULL DEFAULT NULL COMMENT '类型 1目录2菜单 3按钮',
  `sort_no` int(11) NULL DEFAULT NULL COMMENT '排序',
  `route_path` varchar(60) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '路由地址',
  `permission` varchar(60) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '权限编码',
  `status` int(1) NULL DEFAULT NULL COMMENT '是否启用 0启用1禁用',
  `create_time` bigint(20) NULL DEFAULT NULL,
  `iframe_src` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '嵌套连接',
  `route_name` varchar(60) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '路由名称',
  `redirect` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL,
  `props` text CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL,
  `bindapi` int(1) NULL DEFAULT NULL COMMENT '绑定数据权限 1是 -1否',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Compact;

-- ----------------------------
-- Records of sys_menus
-- ----------------------------
INSERT INTO `sys_menus` VALUES ('165dc6313c8b4bc1ab6662bbfb6bcd90', 'd436aeea798d42a4bea423164f214623', '添加', '', '', -1, 1, 1, 3, 1, '', 'memberAdd', 1, 1698995482, '', '', '', '', NULL);
INSERT INTO `sys_menus` VALUES ('19092509384c483494921ab287b38f09', '', '系统管理', 'LAYOUT', 'ant-design:setting-outlined', 0, 0, 1, 1, 10, '/system', '', 1, 1652078930, '0', 'newesys', NULL, NULL, NULL);
INSERT INTO `sys_menus` VALUES ('44030ff97bf7446693ace817f71dfb7a', '19092509384c483494921ab287b38f09', 'api管理', '/system/api/index', 'ant-design:api-filled', -1, 1, 1, 2, 5, 'api', '', 1, 1705911802, '', 'SystemApi', '', '', 1);
INSERT INTO `sys_menus` VALUES ('65ac433016984495ba9c62a0d381ed0b', 'd436aeea798d42a4bea423164f214623', '编辑', '', '', -1, 1, 1, 3, 2, '', 'memberEdit', 1, 1698995504, '', '', '', '', NULL);
INSERT INTO `sys_menus` VALUES ('6fb3210467194b64a6c60a90f0a6e3c5', '19092509384c483494921ab287b38f09', '角色管理', '/system/role/index', 'ant-design:trademark-outlined', -1, 1, 1, 2, 3, 'role', '', 1, 1703641600, '', 'SystemRole', '', '', NULL);
INSERT INTO `sys_menus` VALUES ('702d3b0676524eaa97de61db2aeb3616', '19092509384c483494921ab287b38f09', '字典管理', '/system/dict/index', 'ant-design:appstore-filled', -1, 1, 1, 2, 3, 'dict', '', 1, 1683277674, '', 'SysDict', '', '', NULL);
INSERT INTO `sys_menus` VALUES ('84cb2653afa7499e829f00bc0c8367f7', '19092509384c483494921ab287b38f09', '菜单管理', '/system/menus/index', 'ant-design:menu-unfold-outlined', -1, 1, 1, 2, 1, 'SysMenus', '', 1, 1650025127, '', 'NeweMenus', NULL, NULL, 1);
INSERT INTO `sys_menus` VALUES ('9b144acec61e4886bdcbf0e6ad333481', '19092509384c483494921ab287b38f09', '组织结构', '/system/depart/index', 'ant-design:apartment-outlined', -1, 1, 1, 2, 2, 'depart', '', 1, 1682653838, '', 'SystemDepart', '', '', 1);
INSERT INTO `sys_menus` VALUES ('d436aeea798d42a4bea423164f214623', '19092509384c483494921ab287b38f09', '员工管理', '/system/member/index', 'ant-design:user-outlined', -1, 1, 1, 2, 0, 'member', '', 1, 1683887982, '', 'member', '', '', 1);
INSERT INTO `sys_menus` VALUES ('dc3339db28c44e0fbc602cd44086094d', 'd436aeea798d42a4bea423164f214623', '删除', '', '', -1, 1, 1, 3, 3, '', 'memberDel', 1, 1698995526, '', '', '', '', NULL);

-- ----------------------------
-- Table structure for sys_records
-- ----------------------------
DROP TABLE IF EXISTS `sys_records`;
CREATE TABLE `sys_records`  (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `state` bigint(20) NULL DEFAULT NULL COMMENT '请求状态',
  `api` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '请求地址',
  `create_time` bigint(20) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1 CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Compact;

-- ----------------------------
-- Table structure for sys_role
-- ----------------------------
DROP TABLE IF EXISTS `sys_role`;
CREATE TABLE `sys_role`  (
  `id` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL COMMENT '角色ID',
  `depart_id` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '组织结构ID',
  `org_code` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '组织结构代码',
  `role_name` varchar(20) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '角色名称',
  `status` int(1) NULL DEFAULT NULL COMMENT '角色状态(-1禁用，1启用)',
  `create_time` bigint(20) NULL DEFAULT NULL COMMENT '创建时间',
  `update_time` bigint(20) NULL DEFAULT NULL COMMENT '修改时间',
  `remark` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '备注',
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Compact;

-- ----------------------------
-- Records of sys_role
-- ----------------------------
INSERT INTO `sys_role` VALUES ('3aa86b48d87d4f118546906769b51df5', 'c1a4ad97b18a404cbbe83a31d1ee02d6', 'A01', '超级管理员', 1, 1709724807, 1709724807, '');
INSERT INTO `sys_role` VALUES ('a1ff8273cdf243f6972070a51323e0b4', '8031a1e774e54152a3952895b0471ce3', 'A10', '测试号', 1, 1709794822, 1709794822, '测试账号');

-- ----------------------------
-- Table structure for sys_role_rules
-- ----------------------------
DROP TABLE IF EXISTS `sys_role_rules`;
CREATE TABLE `sys_role_rules`  (
  `id` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NOT NULL,
  `api_id` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT 'apiID',
  `role_id` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '角色ID',
  `menu_id` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '所属菜单ID',
  `org_code` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '部门数据',
  `where_code` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci NULL DEFAULT NULL COMMENT '其他条件',
  `status` int(1) NULL DEFAULT NULL COMMENT '是否允许',
  `create_time` bigint(20) NULL DEFAULT NULL,
  `update_time` bigint(20) NULL DEFAULT NULL,
  PRIMARY KEY (`id`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8mb4 COLLATE = utf8mb4_general_ci ROW_FORMAT = Compact;

-- ----------------------------
-- Records of sys_role_rules
-- ----------------------------
INSERT INTO `sys_role_rules` VALUES ('03816e67c8e545ac912de0fdcbd678b2', '43eb887c22ee427eafe16e995a56d57f', '3aa86b48d87d4f118546906769b51df5', '6fb3210467194b64a6c60a90f0a6e3c5', '', '', 1, 1709809454, 1709809454);
INSERT INTO `sys_role_rules` VALUES ('044df45584864705ba770e1b45d62a66', 'd893f3174fe14e90bd00768a8b1618ea', 'a1ff8273cdf243f6972070a51323e0b4', '44030ff97bf7446693ace817f71dfb7a', '', '', 1, 1709794997, 1709794997);
INSERT INTO `sys_role_rules` VALUES ('0a819ecefbce4745867262ab6e2979f8', 'b2d7a48899114b6aacee94f22baf5c48', '3aa86b48d87d4f118546906769b51df5', '84cb2653afa7499e829f00bc0c8367f7', '', '', 1, 1709724819, 1709724819);
INSERT INTO `sys_role_rules` VALUES ('0c34db8919ea430fb5cb5186acd414a7', '51211bd4a7344054a2bed3bafa94de3b', 'ee795d17571d4187a2b667a2d2a8245a', '702d3b0676524eaa97de61db2aeb3616', '', '', 1, 1709191384, 1709191384);
INSERT INTO `sys_role_rules` VALUES ('0d01b8db097b4c938280c5e0eb7186a4', '50dfd071e389411ba5b4e9090c8efdf9', '3aa86b48d87d4f118546906769b51df5', 'd436aeea798d42a4bea423164f214623', '', '', 1, 1709724815, 1709807864);
INSERT INTO `sys_role_rules` VALUES ('113d4a0503d84111815cef3754d61e37', 'd3af6f4386a3417e94cb0713d3c77888', 'ee795d17571d4187a2b667a2d2a8245a', '702d3b0676524eaa97de61db2aeb3616', '', '', 1, 1709191385, 1709192728);
INSERT INTO `sys_role_rules` VALUES ('138b18d14ddd45fba5835f131466310b', 'b2d7a48899114b6aacee94f22baf5c48', 'ee795d17571d4187a2b667a2d2a8245a', '84cb2653afa7499e829f00bc0c8367f7', '', '123123', 1, 1709102838, 1709192996);
INSERT INTO `sys_role_rules` VALUES ('1bf381c5990d47e1ba000f32e6da0cbb', '61f1549a67d240feb5a0b62a57abe328', 'ee795d17571d4187a2b667a2d2a8245a', '6fb3210467194b64a6c60a90f0a6e3c5', 'A8', '', 1, 1709191381, 1709192470);
INSERT INTO `sys_role_rules` VALUES ('1e10f069f8864b6eb065599f8bafb2c4', '7a983997340340f6bd3ee169f4eb29d4', 'ee795d17571d4187a2b667a2d2a8245a', '9b144acec61e4886bdcbf0e6ad333481', 'A01', '', 1, 1709106009, 1709191375);
INSERT INTO `sys_role_rules` VALUES ('1ec50f6801df4ef8a6267e782565d464', 'c2ff8e889ec84da5985f14155b06cad3', 'a1ff8273cdf243f6972070a51323e0b4', '702d3b0676524eaa97de61db2aeb3616', '', '', 1, 1709794989, 1709794989);
INSERT INTO `sys_role_rules` VALUES ('1f79a166e9a8455793bd3e411dc2e4ca', '8204b2c4f45c43f5bb03149de0f3f2df', '3aa86b48d87d4f118546906769b51df5', '84cb2653afa7499e829f00bc0c8367f7', '', '', 1, 1709724818, 1709724818);
INSERT INTO `sys_role_rules` VALUES ('205c52542d3548acbb1d7c14af35e8cb', 'd893f3174fe14e90bd00768a8b1618ea', 'ee795d17571d4187a2b667a2d2a8245a', '44030ff97bf7446693ace817f71dfb7a', '', '', 1, 1709191390, 1709191390);
INSERT INTO `sys_role_rules` VALUES ('283ea31d7b5f4c578871a93002750b43', '0d3c9db1272243ebad7c45d44f9dd9b8', '3aa86b48d87d4f118546906769b51df5', '6fb3210467194b64a6c60a90f0a6e3c5', '', '', 1, 1709809238, 1709809444);
INSERT INTO `sys_role_rules` VALUES ('342d96ad21c844c28e8b72cef5e929ad', 'c2ff8e889ec84da5985f14155b06cad3', '3aa86b48d87d4f118546906769b51df5', '702d3b0676524eaa97de61db2aeb3616', '', '', 1, 1709724829, 1709724829);
INSERT INTO `sys_role_rules` VALUES ('3750d64a847f4088bcbba45b1027697e', '0461798d7f7245da8c6f3c5bab5f165d', 'ee795d17571d4187a2b667a2d2a8245a', '702d3b0676524eaa97de61db2aeb3616', '', '', 1, 1709191383, 1709191383);
INSERT INTO `sys_role_rules` VALUES ('3770ffd599c84927aaf2606dbf7c2653', '50dfd071e389411ba5b4e9090c8efdf9', 'a1ff8273cdf243f6972070a51323e0b4', 'd436aeea798d42a4bea423164f214623', 'A10', '', 1, 1709794955, 1709794955);
INSERT INTO `sys_role_rules` VALUES ('3dcf805af9bd4dacbb31ccd27469e12b', '1a1c2fb5fc474c7d84674c36dc8076d9', 'ee795d17571d4187a2b667a2d2a8245a', 'd436aeea798d42a4bea423164f214623', '', '', 1, 1709102622, 1709102622);
INSERT INTO `sys_role_rules` VALUES ('3ec20da8c88044d3a02456e5b6afb034', '61f1549a67d240feb5a0b62a57abe328', '3aa86b48d87d4f118546906769b51df5', '6fb3210467194b64a6c60a90f0a6e3c5', '', '', 1, 1709809451, 1709809451);
INSERT INTO `sys_role_rules` VALUES ('421ada07b55e47fd94be3113785ef962', '26f3660954794e67b85b7d12d68e6b10', 'ee795d17571d4187a2b667a2d2a8245a', '44030ff97bf7446693ace817f71dfb7a', '', '', 1, 1709191388, 1709191388);
INSERT INTO `sys_role_rules` VALUES ('43206fb3e353409489c85dc0f4025232', '0163cd4dff5c44ee93aeb6e0a3f68ada', 'ee795d17571d4187a2b667a2d2a8245a', '9b144acec61e4886bdcbf0e6ad333481', '', '', 1, 1709106007, 1709193001);
INSERT INTO `sys_role_rules` VALUES ('47106c6cf9e748baba06de7af2228bd1', 'ea7aae887e5a42fd9a5d069706d55e43', '3aa86b48d87d4f118546906769b51df5', '9b144acec61e4886bdcbf0e6ad333481', '', '', 1, 1709724823, 1709724823);
INSERT INTO `sys_role_rules` VALUES ('48c0649d2edd4dd28c2a1406ccc1371f', '50dfd071e389411ba5b4e9090c8efdf9', 'a1ff8273cdf243f6972070a51323e0b4', 'd436aeea798d42a4bea423164f214623', '', '', 1, 1709794939, 1709794939);
INSERT INTO `sys_role_rules` VALUES ('4c1cf5937af14559b42fc4c89314fe85', '1a1c2fb5fc474c7d84674c36dc8076d9', 'a1ff8273cdf243f6972070a51323e0b4', 'd436aeea798d42a4bea423164f214623', '', '', 1, 1709794947, 1709794947);
INSERT INTO `sys_role_rules` VALUES ('4ff13fcf003446a6b626c26f54388998', 'ea7aae887e5a42fd9a5d069706d55e43', 'ee795d17571d4187a2b667a2d2a8245a', '9b144acec61e4886bdcbf0e6ad333481', '', '', 1, 1709191377, 1709191377);
INSERT INTO `sys_role_rules` VALUES ('52e34127969545c6ad56d2a5429dafac', '18ac3a56882a4248a592c0cd12568a81', 'ee795d17571d4187a2b667a2d2a8245a', 'd436aeea798d42a4bea423164f214623', '', '', 1, 1709102621, 1709102621);
INSERT INTO `sys_role_rules` VALUES ('59ba000368fa4c6f8067143bb08191fc', '7bb1756fe75948809d75b529bc3f1f9a', 'a1ff8273cdf243f6972070a51323e0b4', '84cb2653afa7499e829f00bc0c8367f7', '', '', 1, 1709808664, 1709867612);
INSERT INTO `sys_role_rules` VALUES ('5aa4b5949a1a4026a6feaece934e9eb7', '7bb1756fe75948809d75b529bc3f1f9a', 'ee795d17571d4187a2b667a2d2a8245a', '84cb2653afa7499e829f00bc0c8367f7', '', '', 1, 1709102836, 1709102836);
INSERT INTO `sys_role_rules` VALUES ('5c8be18a20ab412c925b596d8f4c01fc', '5b2d83c3b47f40d99020fa0d37759a75', 'ee795d17571d4187a2b667a2d2a8245a', '702d3b0676524eaa97de61db2aeb3616', '', '', 1, 1709191384, 1709192912);
INSERT INTO `sys_role_rules` VALUES ('5d25690173b7452e891ffdfe815f7321', 'b0123da2ea964f2db0fd59084b9f217c', '3aa86b48d87d4f118546906769b51df5', '6fb3210467194b64a6c60a90f0a6e3c5', '', '', 1, 1709809476, 1709809822);
INSERT INTO `sys_role_rules` VALUES ('6220eab9b4ac41eaabc379dd4b2ec3e4', 'b2d7a48899114b6aacee94f22baf5c48', 'a1ff8273cdf243f6972070a51323e0b4', '84cb2653afa7499e829f00bc0c8367f7', 'A10', '', 1, 1709794964, 1709794964);
INSERT INTO `sys_role_rules` VALUES ('62ad9d510b5a48a5a3ea312a1b66e1be', '18ac3a56882a4248a592c0cd12568a81', '3aa86b48d87d4f118546906769b51df5', 'd436aeea798d42a4bea423164f214623', '', '', 1, 1709724814, 1709807867);
INSERT INTO `sys_role_rules` VALUES ('66bd29ec1103475bb4d9e2fddc9f1509', '41e1843c425145789e3bf5ebd45a2939', '3aa86b48d87d4f118546906769b51df5', '44030ff97bf7446693ace817f71dfb7a', '', '', 1, 1709724832, 1709724832);
INSERT INTO `sys_role_rules` VALUES ('673f0c5fce7140499558646ca4d018a5', '0d3c9db1272243ebad7c45d44f9dd9b8', 'a1ff8273cdf243f6972070a51323e0b4', '6fb3210467194b64a6c60a90f0a6e3c5', '', '', 1, 1709867625, 1709867625);
INSERT INTO `sys_role_rules` VALUES ('769d605b11c6434ca803d2691dbf2472', '7a983997340340f6bd3ee169f4eb29d4', 'a1ff8273cdf243f6972070a51323e0b4', '9b144acec61e4886bdcbf0e6ad333481', '', '', 1, 1709794972, 1709794972);
INSERT INTO `sys_role_rules` VALUES ('795bd08890c2457888446a88c2a5fb4c', '50dfd071e389411ba5b4e9090c8efdf9', 'ee795d17571d4187a2b667a2d2a8245a', 'd436aeea798d42a4bea423164f214623', '', '123123111', 1, 1709102622, 1709192712);
INSERT INTO `sys_role_rules` VALUES ('7d8f995bfdf640e8976f40dbd1b508b9', 'c7e112ebb1f143309405ea8e725fc241', 'ee795d17571d4187a2b667a2d2a8245a', '9b144acec61e4886bdcbf0e6ad333481', '', '', 1, 1709191377, 1709191377);
INSERT INTO `sys_role_rules` VALUES ('83a8fb355b434054b0d3e888d7e31e55', '1a1c2fb5fc474c7d84674c36dc8076d9', '3aa86b48d87d4f118546906769b51df5', 'd436aeea798d42a4bea423164f214623', '', '', 1, 1709724815, 1709807868);
INSERT INTO `sys_role_rules` VALUES ('865e6432529b400192fc6e57d2c79a85', '8204b2c4f45c43f5bb03149de0f3f2df', 'a1ff8273cdf243f6972070a51323e0b4', '84cb2653afa7499e829f00bc0c8367f7', '', '', 1, 1709867618, 1709867618);
INSERT INTO `sys_role_rules` VALUES ('86ec63b6199d4797827d1e2a6fb49359', 'b2d7a48899114b6aacee94f22baf5c48', 'a1ff8273cdf243f6972070a51323e0b4', '84cb2653afa7499e829f00bc0c8367f7', '', '', 1, 1709794962, 1709794962);
INSERT INTO `sys_role_rules` VALUES ('8963534661e842b5a72d6decc1bce493', 'c7e112ebb1f143309405ea8e725fc241', '3aa86b48d87d4f118546906769b51df5', '9b144acec61e4886bdcbf0e6ad333481', '', '', 1, 1709724822, 1709724822);
INSERT INTO `sys_role_rules` VALUES ('8ae287dac66a4e85b16527c345e7dd00', 'd3af6f4386a3417e94cb0713d3c77888', 'a1ff8273cdf243f6972070a51323e0b4', '702d3b0676524eaa97de61db2aeb3616', '', '', 1, 1709794993, 1709795008);
INSERT INTO `sys_role_rules` VALUES ('8ba8754763334ade9873a6658d9abd7c', 'e62e164ec7664389938825b48c5c1445', 'a1ff8273cdf243f6972070a51323e0b4', 'd436aeea798d42a4bea423164f214623', '', '', 1, 1709794944, 1709794944);
INSERT INTO `sys_role_rules` VALUES ('8d11de6dfc814a3a8c50c6d862f36e30', 'a509cd28d6574e458bbd8f7572e0bcd1', 'a1ff8273cdf243f6972070a51323e0b4', '84cb2653afa7499e829f00bc0c8367f7', '', '', 1, 1709867619, 1709867619);
INSERT INTO `sys_role_rules` VALUES ('8db40daf4ad5412baad85d59038c9d5a', '84af270e7bdf4063815c1b6a0d5d9584', '3aa86b48d87d4f118546906769b51df5', '44030ff97bf7446693ace817f71dfb7a', '', '', 1, 1709724833, 1709724833);
INSERT INTO `sys_role_rules` VALUES ('912f8f9523ce40c4974afb18401cfe56', '0163cd4dff5c44ee93aeb6e0a3f68ada', 'a1ff8273cdf243f6972070a51323e0b4', '9b144acec61e4886bdcbf0e6ad333481', '', '', 1, 1709794970, 1709794970);
INSERT INTO `sys_role_rules` VALUES ('935423de5ac4485cad00ba08451ef53b', '61f1549a67d240feb5a0b62a57abe328', 'a1ff8273cdf243f6972070a51323e0b4', '6fb3210467194b64a6c60a90f0a6e3c5', '', '', 1, 1709794980, 1709794980);
INSERT INTO `sys_role_rules` VALUES ('93f1a311d6154fa0a1d6cf69ae317a3c', '18ac3a56882a4248a592c0cd12568a81', 'a1ff8273cdf243f6972070a51323e0b4', 'd436aeea798d42a4bea423164f214623', '', '', 1, 1709867609, 1709867609);
INSERT INTO `sys_role_rules` VALUES ('98af9109f2e748f1b47bb64ca51930fc', 'a509cd28d6574e458bbd8f7572e0bcd1', '3aa86b48d87d4f118546906769b51df5', '84cb2653afa7499e829f00bc0c8367f7', '', '', 1, 1709724818, 1709724818);
INSERT INTO `sys_role_rules` VALUES ('98b6daca7faf43d196278d0a1e1031d0', 'c2ff8e889ec84da5985f14155b06cad3', 'a1ff8273cdf243f6972070a51323e0b4', '702d3b0676524eaa97de61db2aeb3616', '', '', 1, 1709794991, 1709795007);
INSERT INTO `sys_role_rules` VALUES ('9992e2ddfaec4f0d819df0a9df82bb12', '26f3660954794e67b85b7d12d68e6b10', '3aa86b48d87d4f118546906769b51df5', '44030ff97bf7446693ace817f71dfb7a', '', '', 1, 1709724832, 1709724832);
INSERT INTO `sys_role_rules` VALUES ('9fed99022e164cd1b65776b1dff01a46', '7451cf94e359481f8caae27aa486ee0b', 'ee795d17571d4187a2b667a2d2a8245a', '9b144acec61e4886bdcbf0e6ad333481', '', '', 1, 1709191377, 1709191377);
INSERT INTO `sys_role_rules` VALUES ('a02c87ce51774045914c6901b689cee0', '34f04b50dedd48e5a8dd97b2040764f1', '3aa86b48d87d4f118546906769b51df5', '9b144acec61e4886bdcbf0e6ad333481', '', '', 1, 1709724821, 1709809080);
INSERT INTO `sys_role_rules` VALUES ('a08e98d318a74a389f9351a7e3685531', 'e62e164ec7664389938825b48c5c1445', 'ee795d17571d4187a2b667a2d2a8245a', 'd436aeea798d42a4bea423164f214623', '', '', 1, 1709102623, 1709102623);
INSERT INTO `sys_role_rules` VALUES ('b0239f15cd5d4f778341214188890879', '7a983997340340f6bd3ee169f4eb29d4', 'a1ff8273cdf243f6972070a51323e0b4', '9b144acec61e4886bdcbf0e6ad333481', 'A10', '', 1, 1709794976, 1709794976);
INSERT INTO `sys_role_rules` VALUES ('ba1db756b4bf409eac672661cb81829a', '7bb1756fe75948809d75b529bc3f1f9a', '3aa86b48d87d4f118546906769b51df5', '84cb2653afa7499e829f00bc0c8367f7', '', '', 1, 1709724818, 1709724818);
INSERT INTO `sys_role_rules` VALUES ('be6ce15c31a14f5fb6ee1f6c5e29c0b3', 'd893f3174fe14e90bd00768a8b1618ea', '3aa86b48d87d4f118546906769b51df5', '44030ff97bf7446693ace817f71dfb7a', '', '', 1, 1709724833, 1709724833);
INSERT INTO `sys_role_rules` VALUES ('c137e7952c694f848fcdc4426f4e3a60', '070914e683474faeaf4a2e9202b8df72', 'ee795d17571d4187a2b667a2d2a8245a', '44030ff97bf7446693ace817f71dfb7a', '', '', 1, 1709191387, 1709191387);
INSERT INTO `sys_role_rules` VALUES ('c40ede3eb0df40f38994de68cdfcbe7d', '070914e683474faeaf4a2e9202b8df72', '3aa86b48d87d4f118546906769b51df5', '44030ff97bf7446693ace817f71dfb7a', '', '', 1, 1709724832, 1709724832);
INSERT INTO `sys_role_rules` VALUES ('c41477bb3a514725af8588b795f6a756', '1b347deaa11247d6be007f76276ef676', '3aa86b48d87d4f118546906769b51df5', '6fb3210467194b64a6c60a90f0a6e3c5', '', '', 1, 1709809446, 1709809695);
INSERT INTO `sys_role_rules` VALUES ('c65a04941b8143419b58e2cb9715e4e6', '41e1843c425145789e3bf5ebd45a2939', 'ee795d17571d4187a2b667a2d2a8245a', '44030ff97bf7446693ace817f71dfb7a', '', '', 1, 1709191388, 1709191388);
INSERT INTO `sys_role_rules` VALUES ('c6c846bbe9644dce9749f57e844c6374', 'a509cd28d6574e458bbd8f7572e0bcd1', 'ee795d17571d4187a2b667a2d2a8245a', '84cb2653afa7499e829f00bc0c8367f7', '', '', 1, 1709102838, 1709102838);
INSERT INTO `sys_role_rules` VALUES ('c842b2b2f16440d59db2c4bd9a9b16ae', '5b2d83c3b47f40d99020fa0d37759a75', '3aa86b48d87d4f118546906769b51df5', '702d3b0676524eaa97de61db2aeb3616', '', '', 1, 1709724828, 1709724828);
INSERT INTO `sys_role_rules` VALUES ('c8dc33957e134bcda4708172952e7709', '84af270e7bdf4063815c1b6a0d5d9584', 'ee795d17571d4187a2b667a2d2a8245a', '44030ff97bf7446693ace817f71dfb7a', '', '', 1, 1709191388, 1709191388);
INSERT INTO `sys_role_rules` VALUES ('cb38a9b1f6b343cd87777e6a87fb2bd3', '2d0256f898f144258ea944861f92cdc1', '3aa86b48d87d4f118546906769b51df5', '6fb3210467194b64a6c60a90f0a6e3c5', '', '', 1, 1709809470, 1709868117);
INSERT INTO `sys_role_rules` VALUES ('cb8af0db432f44d4bf2b81c7e8520e31', '8204b2c4f45c43f5bb03149de0f3f2df', 'ee795d17571d4187a2b667a2d2a8245a', '84cb2653afa7499e829f00bc0c8367f7', '', '', 1, 1709102837, 1709102837);
INSERT INTO `sys_role_rules` VALUES ('cf290c7798c54a90866bab40f51b7a3b', 'd3af6f4386a3417e94cb0713d3c77888', '3aa86b48d87d4f118546906769b51df5', '702d3b0676524eaa97de61db2aeb3616', '', '', 1, 1709724830, 1709724830);
INSERT INTO `sys_role_rules` VALUES ('d1595c981025485ea2b97896b4456d26', '43eb887c22ee427eafe16e995a56d57f', 'ee795d17571d4187a2b667a2d2a8245a', '6fb3210467194b64a6c60a90f0a6e3c5', '', '', 1, 1709191380, 1709191380);
INSERT INTO `sys_role_rules` VALUES ('d3d5b197b01d49aeb93d9146f97c622f', '0461798d7f7245da8c6f3c5bab5f165d', '3aa86b48d87d4f118546906769b51df5', '702d3b0676524eaa97de61db2aeb3616', '', '', 1, 1709724828, 1709724828);
INSERT INTO `sys_role_rules` VALUES ('d66828178d6b4fe8ba224ccb51f88e9a', 'c2ff8e889ec84da5985f14155b06cad3', 'ee795d17571d4187a2b667a2d2a8245a', '702d3b0676524eaa97de61db2aeb3616', '', '', 1, 1709191384, 1709192743);
INSERT INTO `sys_role_rules` VALUES ('d9278eee48204f34827ea8a5841fbbdd', '7bb1756fe75948809d75b529bc3f1f9a', 'a1ff8273cdf243f6972070a51323e0b4', '84cb2653afa7499e829f00bc0c8367f7', '', '', 1, 1709808664, 1709808664);
INSERT INTO `sys_role_rules` VALUES ('d92c33df31f841d39345d56d7d4ebdde', '0163cd4dff5c44ee93aeb6e0a3f68ada', '3aa86b48d87d4f118546906769b51df5', '9b144acec61e4886bdcbf0e6ad333481', '', '', 1, 1709724821, 1709724821);
INSERT INTO `sys_role_rules` VALUES ('dc6512cec3914c969027c0737bfcd14f', '7a983997340340f6bd3ee169f4eb29d4', '3aa86b48d87d4f118546906769b51df5', '9b144acec61e4886bdcbf0e6ad333481', '', '', 1, 1709724822, 1709724822);
INSERT INTO `sys_role_rules` VALUES ('e232d9b9b2224535b78e3154b0b57c90', '51211bd4a7344054a2bed3bafa94de3b', '3aa86b48d87d4f118546906769b51df5', '702d3b0676524eaa97de61db2aeb3616', '', '', 1, 1709724828, 1709724828);
INSERT INTO `sys_role_rules` VALUES ('e698b568176e4ea6821dbac2448fd9af', '34f04b50dedd48e5a8dd97b2040764f1', 'ee795d17571d4187a2b667a2d2a8245a', '9b144acec61e4886bdcbf0e6ad333481', '', '', 1, 1709191376, 1709191376);
INSERT INTO `sys_role_rules` VALUES ('ec4f959f732846949d9539ac826c305c', '0163cd4dff5c44ee93aeb6e0a3f68ada', 'a1ff8273cdf243f6972070a51323e0b4', '9b144acec61e4886bdcbf0e6ad333481', 'A10', '', 1, 1709794974, 1709794974);
INSERT INTO `sys_role_rules` VALUES ('ecc05bd0bc7c495bae13b2aefdfc1cfb', '7451cf94e359481f8caae27aa486ee0b', '3aa86b48d87d4f118546906769b51df5', '9b144acec61e4886bdcbf0e6ad333481', '', '', 1, 1709724822, 1709724822);
INSERT INTO `sys_role_rules` VALUES ('f024b050fa184d43821565fb4e875b3e', 'd3af6f4386a3417e94cb0713d3c77888', 'a1ff8273cdf243f6972070a51323e0b4', '702d3b0676524eaa97de61db2aeb3616', '', '', 1, 1709794990, 1709794990);
INSERT INTO `sys_role_rules` VALUES ('f36040ec2401422480d001abfb4768a1', '61f1549a67d240feb5a0b62a57abe328', 'a1ff8273cdf243f6972070a51323e0b4', '6fb3210467194b64a6c60a90f0a6e3c5', 'A10', '', 1, 1709794982, 1709794982);
INSERT INTO `sys_role_rules` VALUES ('fd0ba1107f3445b1999cf95798a2010e', '1b347deaa11247d6be007f76276ef676', 'ee795d17571d4187a2b667a2d2a8245a', '6fb3210467194b64a6c60a90f0a6e3c5', '', '', 1, 1709191380, 1709191380);
INSERT INTO `sys_role_rules` VALUES ('fddd4221178c484e9c8c2aa5f01db63c', '0d3c9db1272243ebad7c45d44f9dd9b8', 'ee795d17571d4187a2b667a2d2a8245a', '6fb3210467194b64a6c60a90f0a6e3c5', '', '', 1, 1709191379, 1709191379);
INSERT INTO `sys_role_rules` VALUES ('fef467d371d447a29b89be8f5185ae60', 'e62e164ec7664389938825b48c5c1445', '3aa86b48d87d4f118546906769b51df5', 'd436aeea798d42a4bea423164f214623', '', '', 1, 1709724816, 1709868112);

SET FOREIGN_KEY_CHECKS = 1;
