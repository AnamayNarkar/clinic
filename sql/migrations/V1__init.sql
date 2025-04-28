CREATE TABLE doctors (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL UNIQUE,
    username TEXT NOT NULL UNIQUE,
    email TEXT NOT NULL UNIQUE,
    hased_password TEXT NOT NULL,
    specialization TEXT NOT NULL,
    salt TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE receptionists (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL UNIQUE,
    username TEXT NOT NULL UNIQUE,
    email TEXT NOT NULL UNIQUE,
    hased_password TEXT NOT NULL,
        salt TEXT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE patients (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid() NOT NULL UNIQUE,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE,
    hased_password TEXT NOT NULL,
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
    description TEXT,
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
    description TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    appointment_results BYTEA NOT NULL -- prescription, diagnosis, in a pdf or img or smth
);

-- insert initial data
INSERT INTO doctors (username, email, specialization, hased_password, salt) VALUES
('doctor1','doctor1@gmail.com','cardiologist','$2a$10$cEc62z1ZWy4cHXTPAGK7PONaCyvsDLjn5sKNmNT45IjtiC9glxlFq','8ddbc941-6958-4c10-aed3-913816d1486b'),
('doctor2','doctor2@gmail.com','neurologist','$2a$10$cEc62z1ZWy4cHXTPAGK7PONaCyvsDLjn5sKNmNT45IjtiC9glxlFq','8ddbc941-6958-4c10-aed3-913816d1486b'),
('doctor3','doctor3@gmail.com','dermatologist','$2a$10$cEc62z1ZWy4cHXTPAGK7PONaCyvsDLjn5sKNmNT45IjtiC9glxlFq','8ddbc941-6958-4c10-aed3-913816d1486b');

INSERT INTO receptionists (username, email, hased_password, salt) VALUES
('receptionist1','recep1@gmail.com','$2a$10$cEc62z1ZWy4cHXTPAGK7PONaCyvsDLjn5sKNmNT45IjtiC9glxlFq','8ddbc941-6958-4c10-aed3-913816d1486b'),
('receptionist2','recep2@gmail.com','$2a$10$0X1v5x4Q9Zc3Y6m7J8gkOe1j5z5b5f5f5f5f5f5f5f5f5f5f5f5f', '8ddbc941-6958-4c10-aed3-913816d1486b');

INSERT INTO patients (first_name, last_name, email, hased_password, salt, phone) VALUES
('patient','1','patient1@gmail.com','$2a$10$cEc62z1ZWy4cHXTPAGK7PONaCyvsDLjn5sKNmNT45IjtiC9glxlFq','8ddbc941-6958-4c10-aed3-913816d1486b','1234567890');