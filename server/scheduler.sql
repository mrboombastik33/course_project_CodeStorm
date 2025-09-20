

USE `rfid_project`;

SET GLOBAL event_scheduler = ON;

CREATE EVENT update_keys_event
ON SCHEDULE EVERY 1 MINUTE
DO
  UPDATE key_inf k
  JOIN bookings b ON k.esp_id = b.esp_id
  JOIN access_inf a ON b.name = a.name
  SET k.uid = a.uid
  WHERE TIME(b.booking_time) = TIME(NOW())
    AND b.day_of_week = DAYNAME(NOW());
