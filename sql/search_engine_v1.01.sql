DROP TABLE IF EXISTS `data`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8 */;
CREATE TABLE `data` (
                          `data_id` int NOT NULL AUTO_INCREMENT,
                          `content` text NOT NULL,
                          `images_url` text NOT NULL,
                          PRIMARY KEY (`data_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `choice`
--

LOCK TABLES `data` WRITE;

DROP TABLE IF EXISTS `member`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8 */;
CREATE TABLE `member` (
                          `user_id` int NOT NULL AUTO_INCREMENT,
                          `nickname` varchar(255) DEFAULT NULL,
                          `username` varchar(255) DEFAULT NULL,
                          `password` varchar(255) DEFAULT NULL,
                          `user_type` int DEFAULT NULL,
                          `hobby` varchar(255) DEFAULT NULL,
                          `is_deleted` tinyint(2) NOT NULL DEFAULT 0,
                          PRIMARY KEY (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_general_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Dumping data for table `member`
--

LOCK TABLES `member` WRITE;