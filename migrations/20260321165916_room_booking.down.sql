DROP INDEX IF EXISTS idx_bookings_created_at;
DROP INDEX IF EXISTS idx_bookings_slot_id;
DROP INDEX IF EXISTS idx_bookings_user_active;
DROP INDEX IF EXISTS idx_bookings_user_id;
DROP INDEX IF EXISTS idx_unique_active_booking;
DROP TABLE IF EXISTS bookings;

DROP INDEX IF EXISTS idx_unique_slot;
DROP INDEX IF EXISTS idx_slots_start_at;
DROP INDEX IF EXISTS idx_slots_room_date;
DROP TABLE IF EXISTS slots;

DROP INDEX IF EXISTS idx_schedules_room_id;
DROP TABLE IF EXISTS schedules;

DROP TABLE IF EXISTS rooms;

DROP INDEX IF EXISTS idx_users_role;
DROP TABLE IF EXISTS users;

DROP EXTENSION IF EXISTS "uuid-ossp";