DROP DATABASE IF EXISTS `rfid_project`;

CREATE DATABASE `rfid_project`;
USE `rfid_project`;

CREATE TABLE `access_inf` (
`person_id` INT NOT NULL,
`person_name` CHAR(40) NOT NULL UNIQUE,
`phone` CHAR(19) UNIQUE,
`role_of` BOOLEAN,
`uid` CHAR(11) UNIQUE
);

CREATE TABLE `key_inf` (
`esp_id` INT PRIMARY KEY NOT NULL, 
`audience` CHAR(30) , 
`uid` CHAR(11),
FOREIGN KEY (uid) REFERENCES access_inf(uid)
ON DELETE SET NULL
ON UPDATE CASCADE
);

CREATE TABLE `bookings` (
    `id` INT AUTO_INCREMENT PRIMARY KEY,
    `group_name` CHAR(10) NOT NULL,
    `student_name` CHAR(40) NOT NULL,
    `audience` CHAR(30) NOT NULL,
    `booking_time` TIME NOT NULL,
    `day_of_week` CHAR(10) NOT NULL,
    `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- INSERT INTO access_inf(person_id, person_name, phone, role_of, uid)
-- VALUES (1, 'Pavlo', '380507777777', false, '276FE546'), (3, 'admin', '380677777777', true, 'FF');
-- INSERT INTO key_inf (esp_id, audience, uid)
-- VALUES(312, null, 'FF'), (313, null, 'FF'), (314, null, '276FE546');