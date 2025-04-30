CREATE TABLE doctors (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL UNIQUE,
    username TEXT NOT NULL UNIQUE,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    specialization TEXT NOT NULL,
    salt TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE receptionists (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL UNIQUE,
    username TEXT NOT NULL UNIQUE,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    salt TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE patients (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL UNIQUE,
    username TEXT NOT NULL UNIQUE,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    salt TEXT NOT NULL,
    phone TEXT NOT NULL UNIQUE,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE applications ( -- basically a request for an appointment which is yer to be confirmbed by the receptionist for the doctor
    id UUID PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL UNIQUE,
    patient_id UUID NOT NULL REFERENCES patients(id) ON DELETE CASCADE,
    doctor_id UUID NOT NULL REFERENCES doctors(id) ON DELETE CASCADE,
    status TEXT NOT NULL,
    description TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE appointments (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL UNIQUE,
    patient_id UUID NOT NULL REFERENCES patients(id) ON DELETE CASCADE,
    doctor_id UUID NOT NULL REFERENCES doctors(id) ON DELETE CASCADE,
    application_id UUID NOT NULL REFERENCES applications(id) ON DELETE CASCADE,
    time_of_appointment TIMESTAMP NOT NULL,
    status TEXT NOT NULL,
    description TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    appointment_results TEXT NOT NULL
);

CREATE TABLE admins (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL UNIQUE,
    username TEXT NOT NULL UNIQUE,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    salt TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO admins (username, email, password_hash, salt) VALUES
('admin', 'admin@hospital.com', '$2a$10$cEc62z1ZWy4cHXTPAGK7PONaCyvsDLjn5sKNmNT45IjtiC9glxlFq', '8ddbc941-6958-4c10-aed3-913816d1486b');

-- insert initial data
INSERT INTO doctors (id,username, first_name, last_name, email, password_hash, specialization, salt) VALUES
('6e2d7b7e-1820-4fdf-81cc-68233e42539b','doctor1','john','doe','doctor1@gmail.com','$2a$10$cEc62z1ZWy4cHXTPAGK7PONaCyvsDLjn5sKNmNT45IjtiC9glxlFq','cardiologist','8ddbc941-6958-4c10-aed3-913816d1486b'),
('6e2d7b7e-1820-4fdf-81cc-68233e42539c','doctor2','jane','doe','doctor2@gmail.com','$2a$10$cEc62z1ZWy4cHXTPAGK7PONaCyvsDLjn5sKNmNT45IjtiC9glxlFq','neurologist','8ddbc941-6958-4c10-aed3-913816d1486b'),
('6e2d7b7e-1820-4fdf-81cc-68233e42539a','doctor3','manish','potey','doctor3@gmail.com','$2a$10$cEc62z1ZWy4cHXTPAGK7PONaCyvsDLjn5sKNmNT45IjtiC9glxlFq','orthopedic','8ddbc941-6958-4c10-aed3-913816d1486b');

INSERT INTO receptionists (username, first_name, last_name, email, password_hash, salt) VALUES
('receptionist1','recep','1','recep1@gmail.com','$2a$10$cEc62z1ZWy4cHXTPAGK7PONaCyvsDLjn5sKNmNT45IjtiC9glxlFq','8ddbc941-6958-4c10-aed3-913816d1486b'),
('receptionist2','recep','2','recep2@gmail.com','$2a$10$cEc62z1ZWy4cHXTPAGK7PONaCyvsDLjn5sKNmNT45IjtiC9glxlFq', '8ddbc941-6958-4c10-aed3-913816d1486b');

INSERT INTO patients (id,username,first_name, last_name, email, password_hash, salt, phone) VALUES
('5a91b64e-4fa2-43cf-be47-23fa7acfb4c0','pat1','patient','1','patient1@gmail.com','$2a$10$cEc62z1ZWy4cHXTPAGK7PONaCyvsDLjn5sKNmNT45IjtiC9glxlFq','8ddbc941-6958-4c10-aed3-913816d1486b','1234567890');