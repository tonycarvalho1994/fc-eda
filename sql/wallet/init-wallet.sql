USE wallet;

CREATE TABLE IF NOT EXISTS clients (id varchar(255), name varchar(255), email varchar(255), created_at date);
CREATE TABLE IF NOT EXISTS accounts (id varchar(255), client_id varchar(255), balance int, created_at date);
CREATE TABLE IF NOT EXISTS transactions (id varchar(255), account_id_from varchar(255), account_id_to varchar(255), amount int, created_at date);

INSERT INTO clients (id, name, email, created_at) VALUES ('737b3ad6-3a8d-4a6a-9eb6-4c250df8a02d', 'John Doe', '', '2024-11-12');
INSERT INTO clients (id, name, email, created_at) VALUES ('d3a7bb0b-03c0-4b66-9f66-39388fd0ad67', 'Jane Doe', '', '2024-11-12');

INSERT INTO accounts (id, client_id, balance, created_at) VALUES ('190d35b2-64c5-41e8-b36a-3e89a5945754', '737b3ad6-3a8d-4a6a-9eb6-4c250df8a02d', 1000, '2024-11-12');
INSERT INTO accounts (id, client_id, balance, created_at) VALUES ('ae3310e1-3036-4373-b902-b189f35415c0', 'd3a7bb0b-03c0-4b66-9f66-39388fd0ad67', 1000, '2024-11-12');
