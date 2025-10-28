USE deandb;
SET NAMES utf8mb4;

-- 1) Districts
CREATE TABLE IF NOT EXISTS districts (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  name VARCHAR(255) NOT NULL,
  state CHAR(2) NULL,
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id)
) ENGINE=InnoDB;

-- 2) Schools (belongs to a district)
CREATE TABLE IF NOT EXISTS schools (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  district_id BIGINT UNSIGNED NOT NULL,
  name VARCHAR(255) NOT NULL,
  school_type ENUM('ELEMENTARY', 'MIDDLE', 'HIGH', 'OTHER') NOT NULL DEFAULT 'OTHER',
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  CONSTRAINT fk_school_district
    FOREIGN KEY (district_id) REFERENCES districts(id)
    ON DELETE CASCADE ON UPDATE RESTRICT,
  KEY idx_school_districts (district_id)
) ENGINE=InnoDB;

-- 3) Students (belongs to a school)
CREATE TABLE IF NOT EXISTS students (
  id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
  school_id BIGINT UNSIGNED NOT NULL,
  first_name VARCHAR(100) NOT NULL,
  last_name VARCHAR(100) NOT NULL,
  dob DATE NULL,
  grade_level TINYINT UNSIGNED NUll, -- eg 0-12 or null if N/A
  parent_email VARCHAR(255) NULL,
  enrolled_on DATE NULL,
  status ENUM('ACTIVE', 'INACTIVE', 'GRADUATED', 'WITHDRAWN') NOT NULL DEFAULT 'ACTIVE',
  created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
  updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (id),
  CONSTRAINT fk_student_school
    FOREIGN KEY (school_id) REFERENCES schools(id)
    ON DELETE RESTRICT ON UPDATE RESTRICT,
  KEY idx_student_school (school_id, last_name, first_name),
  KEY idx_student_name (last_name, first_name)
) ENGINE=InnoDB;

-- ------------SEED DATA------------
INSERT INTO districts (name, state) VALUES
  ('North District', 'MI'),
  ('South District', 'MI');

-- Use inserted ids to seed schools
INSERT INTO schools (district_id, name, school_type) VALUES
  ((SELECT id FROM districts WHERE name='North District'), 'Maple Elementary', 'ELEMENTARY'),
  ((SELECT id FROM districts WHERE name='North District'), 'Pine Middle', 'MIDDLE'),
  ((SELECT id FROM districts WHERE name='South District'), 'Cedar High', 'HIGH');

-- Seed students
INSERT INTO students (school_id, first_name, last_name, dob, grade_level, parent_email, enrolled_on, status) VALUES
  ((SELECT s.id FROM schools s WHERE s.name='Maple Elementary'), 'Ava',   'Johnson',  '2014-05-12', 5, 'johnson@example.com',  '2024-09-03', 'ACTIVE'),
  ((SELECT s.id FROM schools s WHERE s.name='Pine Middle'), 'Liam',  'Martinez', '2011-10-02', 8, 'martinez@example.com','2024-09-03', 'ACTIVE'),
  ((SELECT s.id FROM schools s WHERE s.name='Cedar High'), 'Noah',  'Smith',    '2008-03-22',11, 'smith@example.com',   '2024-09-03', 'ACTIVE'),
  ((SELECT s.id FROM schools s WHERE s.name='Cedar High'), 'Emma',  'Brown',    '2007-12-01',12, NULL,                       '2023-09-03', 'GRADUATED');