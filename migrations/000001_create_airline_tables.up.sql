CREATE TABLE agency (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    wallet_balance BIGINT NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),

    CONSTRAINT agency_wallet_balance_non_negative
        CHECK (wallet_balance >= 0)
);

CREATE TABLE api_key (
    id BIGSERIAL PRIMARY KEY,
    agency_id INTEGER NOT NULL,
    key_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    revoked_at TIMESTAMP,

    CONSTRAINT fk_api_key_agency
        FOREIGN KEY (agency_id)
        REFERENCES agency(id)
        ON DELETE CASCADE,

    CONSTRAINT api_key_key_hash_unique
        UNIQUE (key_hash)
);

CREATE INDEX idx_api_key_agency_id
    ON api_key(agency_id);

CREATE TABLE airport (
    id SERIAL PRIMARY KEY,
    code VARCHAR(10) NOT NULL,
    name VARCHAR(255) NOT NULL,
    country VARCHAR(100) NOT NULL,
    province VARCHAR(100),
    city VARCHAR(100) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),

    CONSTRAINT airport_code_unique
        UNIQUE (code)
);

CREATE INDEX idx_airport_city
ON airport (city);

CREATE TABLE airplane (
    id SMALLSERIAL PRIMARY KEY,
    type VARCHAR(100) NOT NULL,
    capacity SMALLINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),

    CONSTRAINT airplane_capacity_positive
        CHECK (capacity > 0)
);

CREATE TABLE flight (
    id BIGSERIAL PRIMARY KEY,

    airplane_id SMALLINT NOT NULL,
    departure_airport_id INTEGER NOT NULL,
    arrival_airport_id INTEGER NOT NULL,

    departure_at TIMESTAMP NOT NULL,
    arrival_at TIMESTAMP NOT NULL,

    price INTEGER NOT NULL,

    capacity SMALLINT NOT NULL,
    available_seats SMALLINT NOT NULL,

    status VARCHAR(100) NOT NULL,

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_flight_airplane
        FOREIGN KEY (airplane_id)
        REFERENCES airplane(id)
        ON DELETE RESTRICT,

    CONSTRAINT fk_flight_departure_airport
        FOREIGN KEY (departure_airport_id)
        REFERENCES airport(id)
        ON DELETE RESTRICT,

    CONSTRAINT fk_flight_arrival_airport
        FOREIGN KEY (arrival_airport_id)
        REFERENCES airport(id)
        ON DELETE RESTRICT,

    CONSTRAINT flight_price_non_negative
        CHECK (price >= 0),

    CONSTRAINT flight_capacity_positive
        CHECK (capacity > 0),

    CONSTRAINT flight_available_seats_valid
        CHECK (
            available_seats >= 0
            AND available_seats <= capacity
        ),

    CONSTRAINT flight_different_airports
        CHECK (departure_airport_id <> arrival_airport_id),

    CONSTRAINT flight_arrival_after_departure
        CHECK (arrival_at > departure_at)
);

CREATE INDEX idx_flight_departure_airport_id
    ON flight(departure_airport_id);

CREATE INDEX idx_flight_arrival_airport_id
    ON flight(arrival_airport_id);

CREATE INDEX idx_flight_active_departure_time
    ON flight (departure_at)
    WHERE status = 'active';

CREATE TABLE ticket (
    id BIGSERIAL PRIMARY KEY,

    agency_id INTEGER NOT NULL,
    flight_id BIGINT NOT NULL,

    seat_count SMALLINT NOT NULL,

    unit_price INTEGER NOT NULL,
    total_price BIGINT NOT NULL,

    status VARCHAR(20) NOT NULL DEFAULT 'active',

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    cancelled_at TIMESTAMP,

    CONSTRAINT fk_ticket_agency
        FOREIGN KEY (agency_id)
        REFERENCES agency(id)
        ON DELETE RESTRICT,

    CONSTRAINT fk_ticket_flight
        FOREIGN KEY (flight_id)
        REFERENCES flight(id)
        ON DELETE RESTRICT,

    CONSTRAINT ticket_seat_count_positive
        CHECK (seat_count > 0),

    CONSTRAINT ticket_unit_price_non_negative
        CHECK (unit_price >= 0),

    CONSTRAINT ticket_total_price_non_negative
        CHECK (total_price >= 0)
);

CREATE INDEX idx_ticket_agency_id
    ON ticket(agency_id);

CREATE INDEX idx_ticket_flight_id
    ON ticket(flight_id);

CREATE INDEX idx_ticket_status
    ON ticket(status);

CREATE TABLE passenger (
    id BIGSERIAL PRIMARY KEY,

    ticket_id BIGINT NOT NULL,

    name VARCHAR(255) NOT NULL,
    mobile VARCHAR(50) NOT NULL,
    national_code VARCHAR(50) NOT NULL,

    created_at TIMESTAMP NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_passenger_ticket
        FOREIGN KEY (ticket_id)
        REFERENCES ticket(id)
        ON DELETE CASCADE
);

CREATE INDEX idx_passenger_ticket_id
    ON passenger(ticket_id);