CREATE TABLE IF NOT EXISTS user (
  `id` INT NOT NULL AUTO_INCREMENT,
  `username` VARCHAR(100) NOT NULL,
  `email` VARCHAR(255) NOT NULL,
  `password` VARCHAR(255) NOT NULL,
  `phone` VARCHAR(15) NOT NULL,
  `full_name` VARCHAR(255) NOT NULL,
  `firebase_uid` VARCHAR(255) NULL,
  `address` VARCHAR(15) NOT NULL,
  `created_at` DATETIME NOT NULL,
  `updated_at` DATETIME NULL,
  `deleted_at` DATETIME NULL,
  PRIMARY KEY (id),
  CONSTRAINT `UQ_user_email` UNIQUE (`email`),
  CONSTRAINT `UQ_user_username` UNIQUE (`username`),
  CONSTRAINT `UQ_user_phone` UNIQUE (`phone`),
  CONSTRAINT `UQ_user_firebase_uid` UNIQUE (`firebase_uid`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4;