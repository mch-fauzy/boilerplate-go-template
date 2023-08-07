DROP TABLE IF EXISTS `foo_item`;
DROP TABLE IF EXISTS `foo`;

CREATE TABLE IF NOT EXISTS `foo` (
  `entity_id` CHAR(36) NOT NULL,
  `name` VARCHAR(255) NOT NULL,
  `total_quantity` INT NOT NULL,
  `total_price` DECIMAL(14,2) NOT NULL,
  `total_discount` DECIMAL(14,2) NOT NULL,
  `shipping_fee` DECIMAL(14,2) NOT NULL,
  `grand_total` DECIMAL(14,2) NOT NULL,
  `status` ENUM('new', 'pending', 'verified', 'paid', 'inTransit', 'delivered', 'failedToDeliver') NOT NULL,
  `created` TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `created_by` CHAR(36) NOT NULL,
  `updated` TIMESTAMP NULL DEFAULT NULL,
  `updated_by` CHAR(36) NULL DEFAULT NULL,
  `deleted` TIMESTAMP NULL DEFAULT NULL,
  `deleted_by` CHAR(36) NULL DEFAULT NULL,
  PRIMARY KEY (`entity_id`),
  UNIQUE `idx_foo_1` (`name`),
  INDEX `idx_foo_2` (`status`),
  INDEX `idx_foo_3` (`created`),
  INDEX `idx_foo_4` (`created_by`),
  INDEX `idx_foo_5` (`updated`),
  INDEX `idx_foo_6` (`updated_by`),
  INDEX `idx_foo_7` (`deleted`),
  INDEX `idx_foo_8` (`deleted_by`)
) ENGINE=InnoDB
DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `foo_item` (
  `entity_id` CHAR(36) NOT NULL,
  `foo_id` CHAR(36) NOT NULL,
  `sku` VARCHAR(20) NOT NULL,
  `product_name` VARCHAR(255) NOT NULL,
  `quantity` INT NOT NULL,
  `unit_price` DECIMAL(14,2) NOT NULL,
  `total_price` DECIMAL(14,2) NOT NULL,
  `discount` DECIMAL(14,2) NOT NULL,
  `grand_total` DECIMAL(14,2) NOT NULL,
  PRIMARY KEY (`entity_id`),
  CONSTRAINT `fk_foo_item_foo_id` FOREIGN KEY (`foo_id`)
    REFERENCES `foo` (`entity_id`)
    ON UPDATE NO ACTION
    ON DELETE NO ACTION,
  UNIQUE `idx_foo_item_1` (`foo_id`, `sku`),
  INDEX `idx_foo_item_2` (`sku`),
  INDEX `idx_foo_item_3` (`product_name`)
) ENGINE=InnoDB
DEFAULT CHARSET=utf8;

INSERT INTO `foo`
(`entity_id`, `name`, `total_quantity`, `total_price`, `total_discount`, `shipping_fee`, `grand_total`, `status`, `created`, `created_by`, `updated`, `updated_by`)
VALUES
('4e80c5bf-b79b-4c90-8f91-82647f439e55', 'The First Foo', 5, 65000, 3900, 15000, 76100, 'new', NOW(), 'd2a9ca76-2468-40c0-87ee-477fcf0a73c3', NULL, NULL),
('3c28318f-35be-4323-94b7-69aeffe26be8', 'The Second Foo', 9, 205000, 31000, 17500, 191500,'inTransit', DATE_SUB(NOW(), INTERVAL 2 DAY), '4e7a814e-78f1-40f1-9a75-9d8ac25b3415', DATE_SUB(NOW(), INTERVAL 2 HOUR), 'd2a9ca76-2468-40c0-87ee-477fcf0a73c3');

INSERT INTO `foo_item`
(`entity_id`, `foo_id`, `sku`, `product_name`, `quantity`, `unit_price`, `total_price`, `discount`, `grand_total`)
VALUES
('7e94b76a-0fc7-4422-bdb0-0caa2f80e43f', '4e80c5bf-b79b-4c90-8f91-82647f439e55', 'SKU-00001', 'Product Name 1', 2, 10000, 20000, 1200, 18800),
('c43ce49f-c689-4f06-9f58-7dec2952beeb', '4e80c5bf-b79b-4c90-8f91-82647f439e55', 'SKU-00002', 'Product Name 2', 3, 15000, 45000, 2700, 42300),
('b0b586ba-c0e2-4f2d-bb0c-eb0a1657f87f', '3c28318f-35be-4323-94b7-69aeffe26be8', 'SKU-00003', 'Product Name 3', 4, 20000, 80000, 16000, 64000),
('c12815ba-23bb-4a2f-8786-10d85ac8b20d', '3c28318f-35be-4323-94b7-69aeffe26be8', 'SKU-00004', 'Product Name 4', 5, 25000, 125000, 15000, 110000);
