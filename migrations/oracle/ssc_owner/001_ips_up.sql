
-- DROP SEQUENCE ips_group_id_seq;
-- CREATE SEQUENCE ips_group_id_seq;

DROP TABLE my_owner.ips;
CREATE TABLE my_owner.ips (
  id NUMBER(11) GENERATED AS IDENTITY NOT NULL PRIMARY KEY,
  ip NUMBER(11) NOT NULL
);

INSERT INTO my_owner.ips (ip) VALUES (0);
INSERT INTO my_owner.ips (ip) VALUES (1);
INSERT INTO my_owner.ips (ip) VALUES (23456);
INSERT INTO my_owner.ips (ip) VALUES (-9876543210);
INSERT INTO my_owner.ips (ip) VALUES (6789012345);

-- INSERT ALL
--   INTO ips (ip) VALUES (0)
--   INTO ips (ip) VALUES (1)
--   INTO ips (ip) VALUES (23456)
--   INTO ips (ip) VALUES (-9876543210)
--   INTO ips (ip) VALUES (6789012345)
-- SELECT 1 from dual;

SELECT * FROM my_owner.ips;
