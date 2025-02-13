DO
$$ BEGIN
   IF NOT EXISTS (SELECT FROM pg_database WHERE datname = 'spy_cat_agency') THEN
      CREATE DATABASE spy_cat_agency;
   END IF;
END $$;

\c spy_cat_agency;

CREATE TABLE cats (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    years_of_experience INT NOT NULL,
    breed VARCHAR(255) NOT NULL,
    salary DECIMAL(10, 2) NOT NULL
);

CREATE TABLE missions (
    id SERIAL PRIMARY KEY,
    cat_id INT REFERENCES cats(id) ON DELETE SET NULL,
    complete BOOLEAN NOT NULL DEFAULT FALSE
);

CREATE TABLE targets (
    id SERIAL PRIMARY KEY,
    mission_id INT REFERENCES missions(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    country VARCHAR(255) NOT NULL,
    notes TEXT,
    complete BOOLEAN NOT NULL DEFAULT FALSE
);

INSERT INTO cats (name, years_of_experience, breed, salary) VALUES
('Whiskers', 5, 'Persian', 3000.00),
('Shadow', 3, 'Maine Coon', 2500.00),
('Luna', 4, 'Siamese', 2700.00);

INSERT INTO missions (cat_id, complete) VALUES
(1, FALSE),
(2, FALSE),
(3, TRUE);

INSERT INTO targets (mission_id, name, country, notes, complete) VALUES
(1, 'Target A', 'USA', 'Initial data collected', FALSE),
(1, 'Target B', 'Canada', 'Observing behavior', FALSE),
(2, 'Target C', 'Germany', 'Planning next steps', FALSE),
(3, 'Target D', 'France', 'Completed initial phase', TRUE);
