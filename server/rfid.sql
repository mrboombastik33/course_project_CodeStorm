-- =================================================================
-- DATABASE SETUP
-- =================================================================

-- Drop the database if it already exists to ensure a clean start
DROP DATABASE IF EXISTS `rfid_project`;

-- Create the new database
CREATE DATABASE `rfid_project`;

-- Select the database to perform operations on
USE `rfid_project`;


-- =================================================================
-- TABLE CREATION (Using your latest schema)
-- =================================================================

CREATE TABLE access_inf (
    person_id INT AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(40) UNIQUE NOT NULL,   
    user_id VARCHAR(11) UNIQUE NOT NULL 
);

CREATE TABLE key_inf (
    esp_id INT UNIQUE NOT NULL,
    audience VARCHAR(30) PRIMARY KEY NOT NULL,
    uid VARCHAR(11) DEFAULT 'FF', -- Set a default value of 'FF' for availability
    FOREIGN KEY (uid) REFERENCES access_inf(user_id)
        ON DELETE SET NULL
        ON UPDATE CASCADE
);

CREATE TABLE bookings (
    id INT AUTO_INCREMENT PRIMARY KEY,
    group_name VARCHAR(10) NOT NULL,
    user_name VARCHAR(40) NOT NULL,  
    audience VARCHAR(30) UNIQUE NOT NULL, -- This audience is booked
    booking_time TIME NOT NULL,
    day_of_week CHAR(10) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_name) REFERENCES access_inf(name)
        ON DELETE CASCADE
        ON UPDATE CASCADE,
    FOREIGN KEY (audience) REFERENCES key_inf(audience)
        ON DELETE CASCADE
        ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS event_log (
    id INT AUTO_INCREMENT PRIMARY KEY,
    message VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);


-- =================================================================
-- INITIAL DATA INSERTION
-- =================================================================

INSERT INTO access_inf (name, user_id)
VALUES 
('user', '276FE546'),
('admin', 'FF');

INSERT INTO key_inf (esp_id, audience, uid)
VALUES
(312, '101A', 'FF'),
(313, '102B', 'FF'),
(314, '103C', 'FF'),
(315, '807', 'FF');


-- =================================================================
-- EVENT SCHEDULER SETUP AND THE CORRECTED EVENT LOGIC
-- =================================================================

-- Ensure the MySQL event scheduler is enabled globally
SET GLOBAL event_scheduler = ON;

-- Drop the event if it exists from a previous setup
DROP EVENT IF EXISTS manage_bookings_event;

CREATE EVENT manage_bookings_event
ON SCHEDULE EVERY 1 MINUTE
DO
BEGIN
  -- Set the language to Ukrainian for correct day name matching
  SET lc_time_names = 'uk_UA';

  -- *** THE FIX IS HERE ***
  -- STEP 1: ASSIGN KEYS FOR RECENTLY STARTED BOOKINGS
  -- We now check for bookings that started between 1 minute ago and now.
  -- This creates a reliable window to catch the booking, regardless of the
  -- exact second the event runs.
  UPDATE key_inf k
  JOIN bookings b ON k.audience = b.audience
  JOIN access_inf a ON b.user_name = a.name
  SET k.uid = a.user_id
  WHERE b.day_of_week = DAYNAME(NOW())
    AND b.booking_time <= TIME(NOW())
    AND b.booking_time > SUBTIME(TIME(NOW()), '00:01:00') -- Occurred in the last minute
    AND k.uid = 'FF';

  -- STEP 2: REVOKE KEYS FOR EXPIRED BOOKINGS
  -- This logic remains correct. It finds keys for bookings that ended
  -- more than 2 minutes ago and sets them back to 'FF'.
  UPDATE key_inf k
  JOIN bookings b ON k.audience = b.audience
  SET k.uid = 'FF'
  WHERE b.day_of_week = DAYNAME(NOW())
    AND TIME(NOW()) > ADDTIME(b.booking_time, '00:02:00');

  -- STEP 3: LOG THE DELETION of expired bookings
  -- This logic is also correct.
  INSERT INTO event_log (message)
  SELECT CONCAT('Deleting expired booking for user: ', user_name, ', Audience: ', audience, ', at: ', booking_time)
  FROM bookings
  WHERE day_of_week = DAYNAME(NOW())
    AND TIME(NOW()) > ADDTIME(booking_time, '00:02:00');

  -- STEP 4: DELETE the expired bookings from the schedule
  -- This runs last and is also correct.
  DELETE FROM bookings
  WHERE day_of_week = DAYNAME(NOW())
    AND TIME(NOW()) > ADDTIME(booking_time, '00:02:00');

END;
