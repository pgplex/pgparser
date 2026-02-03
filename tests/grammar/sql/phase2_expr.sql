-- Phase 2: core expression grammar gaps
SELECT TREAT('1' AS int);
-- expect-error: UNIQUE predicate is not yet implemented
SELECT UNIQUE (SELECT 1);
SELECT ROW();
SELECT 1 OPERATOR(pg_catalog.=) ANY (SELECT 1);
SELECT (SELECT 1)[1];
SELECT (SELECT 1).foo;
SELECT CAST('a' AS NATIONAL CHAR(5));
SELECT CAST('a' AS NCHAR(5));
SELECT SUBSTRING('abcdef', 2, 3);
SELECT OVERLAY('abcdef', 'Z', 3);
SELECT 33 * ANY ('{1,2,3}');
SELECT 33 * ANY (44);
