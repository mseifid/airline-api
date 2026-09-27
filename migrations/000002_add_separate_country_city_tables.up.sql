CREATE TABLE country (
    id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),

    CONSTRAINT country_name_unique
        UNIQUE (name)
);

CREATE TABLE province (
    id SERIAL PRIMARY KEY,
    country_id INTEGER NOT NULL,
    name VARCHAR(100) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_province_country
        FOREIGN KEY (country_id)
        REFERENCES country(id)
        ON DELETE RESTRICT,

    CONSTRAINT province_country_name_unique
        UNIQUE (country_id, name)
);

CREATE INDEX idx_province_country_id
    ON province(country_id);

CREATE TABLE city (
    id SERIAL PRIMARY KEY,
    province_id INTEGER NOT NULL,
    name VARCHAR(100) NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),

    CONSTRAINT fk_city_province
        FOREIGN KEY (province_id)
        REFERENCES province(id)
        ON DELETE RESTRICT,

    CONSTRAINT city_province_name_unique
        UNIQUE (province_id, name)
);

CREATE INDEX idx_city_province_id
    ON city(province_id);

ALTER TABLE airport
    DROP COLUMN country,
    DROP COLUMN province,
    DROP COLUMN city;

ALTER TABLE airport
    ADD COLUMN city_id INTEGER NOT NULL;

ALTER TABLE airport
    ADD CONSTRAINT fk_airport_city
        FOREIGN KEY (city_id)
        REFERENCES city(id)
        ON DELETE RESTRICT;

CREATE INDEX idx_airport_city_id
    ON airport(city_id);