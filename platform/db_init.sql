-- Create Users Table
CREATE TABLE users (
                       id SERIAL PRIMARY KEY,
                       username TEXT NOT NULL UNIQUE,
                       email TEXT NOT NULL UNIQUE,
                       password_hash TEXT NOT NULL,
                       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create Roles Table
CREATE TABLE roles (
                       id SERIAL PRIMARY KEY,
                       name TEXT NOT NULL UNIQUE
);

-- Insert Default Roles (Admin, View)
INSERT INTO roles (name) VALUES ('Admin'), ('View');

-- Create Homes Table
CREATE TABLE homes (
                       id SERIAL PRIMARY KEY,
                       home_name TEXT NOT NULL,
                       user_id INTEGER NOT NULL REFERENCES users(id),
                       created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create Home-User Mapping Table
CREATE TABLE home_users (
                            id SERIAL PRIMARY KEY,
                            home_id INTEGER NOT NULL REFERENCES homes(id),
                            user_id INTEGER NOT NULL REFERENCES users(id),
                            role_id INTEGER NOT NULL REFERENCES roles(id),
                            created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create Devices Table
CREATE TABLE devices (
                         id SERIAL PRIMARY KEY,
                         device_id TEXT NOT NULL UNIQUE,
                         channel_id TEXT NOT NULL UNIQUE,
                         production_date DATE,
                         warranty INTEGER,
                         location TEXT,
                         is_active BOOLEAN DEFAULT TRUE,
                         user_id INTEGER NOT NULL REFERENCES users(id),
                         home_id INTEGER REFERENCES homes(id),
                         created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create Device Data Table
CREATE TABLE device_data (
                             id SERIAL PRIMARY KEY,
                             device_id TEXT NOT NULL,
                             home_id INTEGER REFERENCES homes(id),
                             data JSONB NOT NULL,
                             created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create Device Analytics Table
CREATE TABLE device_analytics (
                                  id SERIAL PRIMARY KEY,
                                  device_id TEXT NOT NULL,
                                  home_id INTEGER REFERENCES homes(id),
                                  metric TEXT NOT NULL,
                                  value NUMERIC,
                                  count INTEGER,
                                  min_value NUMERIC,
                                  max_value NUMERIC,
                                  avg_value NUMERIC,
                                  sum_value NUMERIC,
                                  aggregation_period TIMESTAMP,
                                  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
