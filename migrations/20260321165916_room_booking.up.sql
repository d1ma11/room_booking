CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users
(
    id         UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email      TEXT NOT NULL UNIQUE,
    password   TEXT,
    role       TEXT NOT NULL CHECK (role IN ('admin', 'user')),
    created_at TIMESTAMPTZ      DEFAULT now()
);
CREATE INDEX idx_users_role ON users (role);

CREATE TABLE rooms
(
    id          UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name        TEXT NOT NULL,
    description TEXT,
    capacity    INT,
    created_at  TIMESTAMPTZ      DEFAULT now()
);

CREATE TABLE schedules
(
    id           UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    room_id      UUID NOT NULL REFERENCES rooms (id) ON DELETE CASCADE,
    days_of_week INT[] NOT NULL,
    start_time   TIME NOT NULL,
    end_time     TIME NOT NULL,

    CONSTRAINT unique_schedule_per_room UNIQUE (room_id),
    CONSTRAINT valid_time CHECK (start_time < end_time)
);
CREATE INDEX idx_schedules_room_id ON schedules (room_id);

CREATE TABLE slots
(
    id       UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    room_id  UUID        NOT NULL REFERENCES rooms (id) ON DELETE CASCADE,
    start_at TIMESTAMPTZ NOT NULL,
    end_at   TIMESTAMPTZ NOT NULL,

    CONSTRAINT valid_slot_time CHECK (start_at < end_at)
);
CREATE INDEX idx_slots_room_date ON slots (room_id, start_at);
CREATE INDEX idx_slots_start_at ON slots (start_at);
CREATE UNIQUE INDEX idx_unique_slot ON slots (room_id, start_at, end_at);

CREATE TABLE bookings
(
    id              UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    slot_id         UUID NOT NULL REFERENCES slots (id) ON DELETE CASCADE,
    user_id         UUID NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    status          TEXT NOT NULL CHECK (status IN ('active', 'cancelled')),
    conference_link TEXT,
    created_at      TIMESTAMPTZ      DEFAULT now()
);
CREATE UNIQUE INDEX idx_unique_active_booking ON bookings (slot_id) WHERE status = 'active';
CREATE INDEX idx_bookings_user_id ON bookings (user_id);
CREATE INDEX idx_bookings_user_active ON bookings (user_id, status);
CREATE INDEX idx_bookings_slot_id ON bookings (slot_id);