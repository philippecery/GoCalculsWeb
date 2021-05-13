-- Removing root remote access
DROP USER root@'%';

-- Creating maths user
CREATE USER IF NOT EXISTS maths@'10.0.3.%' IDENTIFIED VIA mysql_native_password USING '*2A08194354781F68BFC5C54925EA31DC689459D6';
GRANT SELECT,INSERT,UPDATE,DELETE on TABLE maths.* TO maths@'10.0.3.%';