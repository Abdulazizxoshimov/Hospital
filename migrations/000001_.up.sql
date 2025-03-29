CREATE TABLE users (
    id uuid PRIMARY KEY,
    full_name VARCHAR(255),
    username VARCHAR(100) UNIQUE NOT NULL,
    password TEXT NOT NULL,
    role VARCHAR(20) CHECK (role IN ('user', 'doctor', 'admin')) NOT NULL,
    email VARCHAR(255) UNIQUE,
    phone_number VARCHAR(15),
    refresh_token TEXT,
    created_at TIMESTAMP DEFAULT now()
);
INSERT INTO users (id, full_name, username, password, role, email, phone_number) 
VALUES ('6fbdd241-db3a-4d86-81ca-2e131b68515c', 'Abdulaziz Xoshimov', 'abdulaziz2004', '$2a$14$vFpAjqut.m4XLqx/BH914eXBoFSkWVXwJyvTPHfxCSEsTS8336mzK', 'admin', 'abdulazizxoshimov22@gmail.com', '+99890773142044');

CREATE TABLE doctors (
    id uuid PRIMARY KEY,
    user_id uuid UNIQUE NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    specialization VARCHAR(100) NOT NULL,
    working_hours VARCHAR(50) NOT NULL, 
    extra_info JSONB DEFAULT '{}'::jsonb,  
    created_at TIMESTAMP DEFAULT now()
);

CREATE TABLE doctor_availability (
    id SERIAL PRIMARY KEY,
    doctor_id uuid NOT NULL REFERENCES doctors(id) ON DELETE CASCADE,
    available_date DATE NOT NULL,
    start_time TIMESTAMP NOT NULL,
    end_time TIMESTAMP NOT NULL,
    is_booked BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT now(),
    UNIQUE (doctor_id, available_date, start_time)
);

CREATE TABLE appointments (
    id SERIAL PRIMARY KEY,
    patient_id uuid NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    doctor_id uuid NOT NULL REFERENCES doctors(id) ON DELETE CASCADE,
    appointment_time JSONB NOT NULL,
    start_time TIMESTAMP NOT NULL,
    status VARCHAR(20) CHECK (status IN ('scheduled', 'completed', 'cancelled')) DEFAULT 'scheduled',
    created_at TIMESTAMP DEFAULT now(),
    UNIQUE (doctor_id, start_time) 
);

CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_email ON users(email);
