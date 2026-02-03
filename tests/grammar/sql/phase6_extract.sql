-- Phase 6: extractor behavior coverage

-- psql terminators should end statements mid-line
SELECT 1 \gdesc
SELECT 2;

SELECT 3 \gexec
SELECT 4;

SELECT 5 \crosstabview
SELECT 6;

-- BEGIN ATOMIC block should not split on internal semicolons
CREATE FUNCTION atomic_fun(a int) RETURNS int LANGUAGE SQL
BEGIN ATOMIC
  SELECT a + 1;
  SELECT a + 2;
END;

-- CREATE RULE DO (...) should keep semicolons inside parens
CREATE RULE r1 AS ON INSERT TO foo DO (
  INSERT INTO foo VALUES (1);
  UPDATE foo SET a = 2;
);
