-- Grammar compatibility fixtures for postgres features missing in pgparser.

-- expect-error: syntax error
ALTER DATABASE mydb REFRESH COLLATION VERSION;

-- expect-error: syntax error
ALTER TYPE mood DROP VALUE 'sad';

-- expect-error: syntax error
ALTER EXTENSION hstore ADD CAST (text AS hstore);

-- expect-error: syntax error
ALTER EXTENSION hstore ADD TRANSFORM FOR jsonb LANGUAGE plpgsql;

-- expect-error: syntax error
CREATE ASSERTION positive CHECK (1 = 1);
