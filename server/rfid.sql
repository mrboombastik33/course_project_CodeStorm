-- Drop the database if it exists to start with a clean slate
DROP DATABASE IF EXISTS `rfid_project`;

-- Create the database
CREATE DATABASE `rfid_project`;

-- Select the database to use
USE `rfid_project`;

-- Create the table for access information
CREATE TABLE access_inf (
    person_id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(40) UNIQUE NOT NULL,   
    user_id VARCHAR(11) UNIQUE NOT NULL 
);

-- Create the table for key (ESP) information
CREATE TABLE key_inf (
    esp_id INT PRIMARY KEY NOT NULL,
    audience VARCHAR(30) NOT NULL,
    uid VARCHAR(11),
    FOREIGN KEY (uid) REFERENCES access_inf(user_id)
        ON DELETE SET NULL
        ON UPDATE CASCADE
);

-- Create the table for bookings
CREATE TABLE bookings (
    id INT AUTO_INCREMENT PRIMARY KEY,
    group_name VARCHAR(10) NOT NULL,
    user_name VARCHAR(40) NOT NULL,  
    esp_id INT NOT NULL,
    booking_time TIME NOT NULL,
    day_of_week CHAR(10) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_name) REFERENCES access_inf(name)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    FOREIGN KEY (esp_id) REFERENCES key_inf(esp_id)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);

-- Create the table for logging event actions
CREATE TABLE IF NOT EXISTS event_log (
    id INT AUTO_INCREMENT PRIMARY KEY,
    message VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Insert sample data for testing purposes
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

-- Ensure the MySQL event scheduler is turned on
SET GLOBAL event_scheduler = ON;

    
-- First, drop the old event
DROP EVENT IF EXISTS manage_bookings_event;

-- Create the corrected event
CREATE EVENT manage_bookings_event
ON SCHEDULE EVERY 1 MINUTE
DO
BEGIN
  -- SET THE LANGUAGE for this session to Ukrainian
  SET lc_time_names = 'uk_UA';

  -- Part 1: Update key information (this will now work correctly)
  UPDATE key_inf k
  JOIN bookings b ON k.esp_id = b.esp_id
  JOIN access_inf a ON b.user_name = a.name
  SET k.uid = a.user_id
  WHERE b.booking_time = TIME(NOW())
    AND b.day_of_week = DAYNAME(NOW());

  -- Part 2: Log the bookings about to be deleted
  INSERT INTO event_log (message)
  SELECT CONCAT('Deleting expired booking for user: ', user_name, ', ESP ID: ', esp_id, ', scheduled at: ', booking_time)
  FROM bookings
  WHERE day_of_week = DAYNAME(NOW())
    AND TIME(NOW()) > ADDTIME(booking_time, '00:02:00');

  -- Part 3: Delete old bookings (this will now work correctly)
  DELETE FROM bookings
  WHERE day_of_week = DAYNAME(NOW())
    AND TIME(NOW()) > ADDTIME(booking_time, '00:02:00');
END;

  