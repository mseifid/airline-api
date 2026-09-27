ALTER TABLE airport
    DROP CONSTRAINT fk_airport_city;

DROP INDEX IF EXISTS idx_airport_city_id;

ALTER TABLE airport
    DROP COLUMN city_id;

ALTER TABLE airport
    ADD COLUMN country VARCHAR(100) NOT NULL,
    ADD COLUMN province VARCHAR(100),
    ADD COLUMN city VARCHAR(100) NOT NULL;

CREATE INDEX idx_airport_city
    ON airport(city);

DROP TABLE IF EXISTS city;
DROP TABLE IF EXISTS province;
DROP TABLE IF EXISTS country;