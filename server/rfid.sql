DROP DATABASE IF EXISTS `rfid_project`;
CREATE DATABASE `rfid_project`;
USE `rfid_project`;

CREATE TABLE `access_inf` (
    `person_id` INT AUTO_INCREMENT PRIMARY KEY,
    `name` CHAR(40) UNIQUE NOT NULL,
    `user_id` CHAR(11) UNIQUE
);

CREATE TABLE `key_inf` (
    `esp_id` INT PRIMARY KEY NOT NULL,
    `audience` CHAR(30) NOT NULL,
    `uid` CHAR(11),
    FOREIGN KEY (uid) REFERENCES access_inf(user_id)
        ON DELETE SET NULL
        ON UPDATE CASCADE
);

CREATE TABLE `bookings` (
    `id` INT AUTO_INCREMENT PRIMARY KEY,
    `group_name` CHAR(10) NOT NULL,
    `user_name` CHAR(11) NOT NULL,
    `esp_id` INT NOT NULL,
    `booking_time` TIME NOT NULL,
    `day_of_week` CHAR(10) NOT NULL,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_name) REFERENCES access_inf(name)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    FOREIGN KEY (esp_id) REFERENCES key_inf(esp_id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);


-- Дані для тестів
INSERT INTO access_inf (name, user_id)
VALUES 
('Pavlo', '276FE546'),
('admin', 'FF');

INSERT INTO key_inf (esp_id, audience, uid)
VALUES
(312, '101A', 'FF'),
(313, '102B', 'FF'),
(314, '103C', '276FE546'),
(315, '807', 'FF');

-- INSERT INTO access_inf(person_id, person_name,uid)
-- VALUES (1, 'Pavlo', '276FE546'), (3, 'admin', 'FF');
-- INSERT INTO key_inf (esp_id, audience, uid)
-- VALUES(312, 'FF'), (313, 'FF'), (314, '276FE546');
