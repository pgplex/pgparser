-- Phase 1: infrastructure coverage for core expression nodes
SELECT CURRENT_DATE;
SELECT CURRENT_TIME(3);
SELECT COALESCE(NULL, 1, 2);
SELECT GREATEST(1, 2), LEAST(2, 3);
