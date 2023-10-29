# dbench

### source data
[price data](https://www.gov.uk/government/statistical-data-sets/price-paid-data-downloads)

download seed data in csv format
```
curl -L -O http://prod.publicdata.landregistry.gov.uk.s3-website-eu-west-1.amazonaws.com/pp-monthly-update-new-version.csv
```

create a database with table `prices` 
```
CREATE TABLE `prices` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT,
  `uuid` varchar(100) NOT NULL,
  `price` bigint(20) NOT NULL,
  `time_stamp` bigint(20) NOT NULL,
  `post_code` varchar(20) NOT NULL,
  `p_type` varchar(20) DEFAULT NULL,
  `is_new` tinyint(1) DEFAULT NULL,
  `duration` varchar(20) DEFAULT NULL,
  `addr1` varchar(255) DEFAULT NULL,
  `addr2` varchar(255) DEFAULT NULL,
  `street` varchar(100) DEFAULT NULL,
  `locality` varchar(100) DEFAULT NULL,
  `town` varchar(100) DEFAULT NULL,
  `district` varchar(100) DEFAULT NULL,
  `county` varchar(100) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `uuid` (`uuid`),
  KEY `price` (`price`),
  KEY `post_code` (`post_code`),
  KEY `time_stamp` (`time_stamp`),
  KEY `p_type` (`p_type`),
  KEY `duration` (`duration`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;
```
