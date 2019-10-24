-- Here is a sample init file for your mysql database. This can be any mysql
-- exported file too. Just call it the name referenced in the docker-compose.yml
-- entry for 'mysql'. Anything in the container's docker-entrypoint-initdb.d
-- location will be automatically loaded on the container's creation
CREATE TABLE pet (id INT NOT NULL AUTO_INCREMENT PRIMARY KEY, name VARCHAR(20), owner VARCHAR(20), species VARCHAR(20), sex CHAR(1));

INSERT INTO pet (name, owner, species, sex) VALUES ('Doggo', 'Cam', 'Mut', 'M');
