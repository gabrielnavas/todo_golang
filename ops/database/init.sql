CREATE SCHEMA IF NOT EXISTS todos;
CREATE SCHEMA IF NOT EXISTS users;

CREATE TABLE IF NOT EXISTS todos.todo_status (
  id serial,
  name VARCHAR(255) NOT NULL,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW(),
  PRIMARY KEY (id)
);

CREATE TABLE IF NOT EXISTS todos.todo (
  id serial,
  tstts_id INT,
  title VARCHAR(255) NOT NULL,
  description VARCHAR(255) NOT NULL,
  image BYTEA DEFAULT null,
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW(),
  PRIMARY KEY (id),
  FOREIGN KEY (tstts_id) 
  	REFERENCES todos.todo_status(id)
    ON UPDATE CASCADE
    ON DELETE RESTRICT
);

CREATE TABLE IF NOT EXISTS users.user (
  id serial,
  name VARCHAR(255) NOT NULL,
  email VARCHAR(255) NOT NULL UNIQUE,
  username VARCHAR(255) NOT NULL UNIQUE,
  password VARCHAR(255) NOT NULL,
  level_access INT NOT NULL, 
  created_at TIMESTAMP DEFAULT NOW(),
  updated_at TIMESTAMP DEFAULT NOW(),
  photo BYTEA DEFAULT null,
  PRIMARY KEY (id)
);