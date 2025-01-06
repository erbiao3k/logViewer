-- MySQL dump 10.13  Distrib 5.7.36, for Linux (x86_64)
--
-- Host: 43.128.13.79    Database: log_project
-- ------------------------------------------------------
-- Server version	5.5.68-MariaDB

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!40101 SET NAMES utf8 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `pm_infos`
--

DROP TABLE IF EXISTS `pm_infos`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `pm_infos` (
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL,
  `pm_area` varchar(255) DEFAULT NULL,
  `pm_name` varchar(255) DEFAULT NULL,
  `pm_email` varchar(100) DEFAULT NULL,
  `pm_phone` varchar(255) DEFAULT NULL,
  `pm_passwd` varchar(255) DEFAULT NULL,
  `enable` tinyint(1) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uix_pm_infos_pm_email` (`pm_email`),
  KEY `idx_pm_infos_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=11 DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `pm_infos`
--

LOCK TABLES `pm_infos` WRITE;
/*!40000 ALTER TABLE `pm_infos` DISABLE KEYS */;
INSERT INTO `pm_infos` VALUES (5,'2022-01-14 15:18:17','2022-01-14 15:18:17',NULL,'admin','','adminEngineer@myemal.com','12345678901','icf9rTM77887Q44F!nOW20',1);
/*!40000 ALTER TABLE `pm_infos` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `project_infos`
--

DROP TABLE IF EXISTS `project_infos`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `project_infos` (
  `md5_data` varchar(255) DEFAULT NULL,
  `project_name` varchar(255) DEFAULT NULL,
  `project_env` varchar(255) DEFAULT NULL,
  `project_area` varchar(255) DEFAULT NULL,
  `svc_admin` varchar(255) DEFAULT NULL,
  `svc_portal` varchar(255) DEFAULT NULL,
  `svc_api` varchar(255) DEFAULT NULL,
  `svc_schedule` varchar(255) DEFAULT NULL,
  `svc_fos` varchar(255) DEFAULT NULL,
  `svc_convert` varchar(255) DEFAULT NULL,
  `svc_sign` varchar(255) DEFAULT NULL,
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_project_infos_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `project_infos`
--

LOCK TABLES `project_infos` WRITE;
/*!40000 ALTER TABLE `project_infos` DISABLE KEYS */;
/*!40000 ALTER TABLE `project_infos` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `project_log_commits`
--

DROP TABLE IF EXISTS `project_log_commits`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `project_log_commits` (
  `log_commit_md5` varchar(255) DEFAULT NULL,
  `project_name` varchar(255) DEFAULT NULL,
  `svc_name` varchar(255) DEFAULT NULL,
  `project_env` varchar(255) DEFAULT NULL,
  `log_date` varchar(255) DEFAULT NULL,
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL,
  `pm_name` varchar(255) DEFAULT NULL,
  `log_status` varchar(255) DEFAULT NULL,
  `pm_email` varchar(255) DEFAULT NULL,
  `svc_addr` varchar(255) DEFAULT NULL,
  `create_time` varchar(255) DEFAULT NULL,
  `pm_phone` varchar(255) DEFAULT NULL,
  `log_download_addr` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_project_log_commits_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `project_log_commits`
--

LOCK TABLES `project_log_commits` WRITE;
/*!40000 ALTER TABLE `project_log_commits` DISABLE KEYS */;
/*!40000 ALTER TABLE `project_log_commits` ENABLE KEYS */;
UNLOCK TABLES;

--
-- Table structure for table `white_lists`
--

DROP TABLE IF EXISTS `white_lists`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!40101 SET character_set_client = utf8 */;
CREATE TABLE `white_lists` (
  `project_name` varchar(255) DEFAULT NULL,
  `project_env` varchar(255) DEFAULT NULL,
  `public_ip` varchar(255) DEFAULT NULL,
  `id` int(10) unsigned NOT NULL AUTO_INCREMENT,
  `created_at` datetime DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_white_lists_deleted_at` (`deleted_at`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `white_lists`
--

LOCK TABLES `white_lists` WRITE;
/*!40000 ALTER TABLE `white_lists` DISABLE KEYS */;
INSERT INTO `white_lists` VALUES ('testProject','uat','127.0.0.1',1,'2022-01-16 18:48:42','2022-01-16 18:48:42',NULL);
/*!40000 ALTER TABLE `white_lists` ENABLE KEYS */;
UNLOCK TABLES;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2022-01-17 16:18:37
