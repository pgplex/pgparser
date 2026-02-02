#include "postgres.h"

#include <errno.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include <unistd.h>

#include "mb/pg_wchar.h"
#include "miscadmin.h"
#include "nodes/nodes.h"
#include "parser/parser.h"
#include "utils/elog.h"
#include "utils/errcodes.h"
#include "utils/guc.h"
#include "utils/memutils.h"
#include "utils/pg_locale.h"

static void
die(const char *msg)
{
	fprintf(stderr, "pg_parse_helper: %s\n", msg);
	exit(1);
}

static char *
read_all(FILE *fp)
{
	size_t cap = 8192;
	size_t len = 0;
	char *buf = (char *) malloc(cap);
	if (!buf)
		die("out of memory");

	for (;;)
	{
		size_t n = fread(buf + len, 1, cap - len, fp);
		len += n;
		if (n == 0)
		{
			if (ferror(fp))
				die("failed to read input");
			break;
		}
		if (len == cap)
		{
			cap *= 2;
			char *next = (char *) realloc(buf, cap);
			if (!next)
				die("out of memory");
			buf = next;
		}
	}
	buf = (char *) realloc(buf, len + 1);
	if (!buf)
		die("out of memory");
	buf[len] = '\0';
	return buf;
}

static void
pg_init(const char *argv0)
{
	MyProcPid = getpid();
	MemoryContextInit();
	set_pglocale_pgservice(argv0, PG_TEXTDOMAIN("postgres"));
	InitializeGUCOptions();
	SetDatabaseEncoding(PG_UTF8);
	SetClientEncoding(PG_UTF8);
}

static void
usage(const char *prog)
{
	fprintf(stderr,
			"Usage: %s [--mode=default|type_name|plpgsql_expr|plpgsql_assign1|plpgsql_assign2|plpgsql_assign3] [--file path]\n"
			"Reads SQL from --file or stdin and prints nodeToString(raw_parser(...)).\n",
			prog);
	exit(2);
}

static RawParseMode
parse_mode_from_arg(const char *arg)
{
	if (strcmp(arg, "default") == 0)
		return RAW_PARSE_DEFAULT;
	if (strcmp(arg, "type_name") == 0)
		return RAW_PARSE_TYPE_NAME;
	if (strcmp(arg, "plpgsql_expr") == 0)
		return RAW_PARSE_PLPGSQL_EXPR;
	if (strcmp(arg, "plpgsql_assign1") == 0)
		return RAW_PARSE_PLPGSQL_ASSIGN1;
	if (strcmp(arg, "plpgsql_assign2") == 0)
		return RAW_PARSE_PLPGSQL_ASSIGN2;
	if (strcmp(arg, "plpgsql_assign3") == 0)
		return RAW_PARSE_PLPGSQL_ASSIGN3;
	return RAW_PARSE_DEFAULT;
}

int
main(int argc, char **argv)
{
	const char *file_path = NULL;
	RawParseMode mode = RAW_PARSE_DEFAULT;
	FILE *fp = NULL;
	char *sql = NULL;

	for (int i = 1; i < argc; i++)
	{
		if (strcmp(argv[i], "--help") == 0 || strcmp(argv[i], "-h") == 0)
			usage(argv[0]);
		if (strncmp(argv[i], "--mode=", 7) == 0)
		{
			mode = parse_mode_from_arg(argv[i] + 7);
			continue;
		}
		if (strcmp(argv[i], "--file") == 0 && i + 1 < argc)
		{
			file_path = argv[++i];
			continue;
		}
		usage(argv[0]);
	}

	pg_init(argv[0]);

	if (file_path)
	{
		fp = fopen(file_path, "rb");
		if (!fp)
			die("failed to open file");
	}
	else
	{
		fp = stdin;
	}

	sql = read_all(fp);
	if (file_path)
		fclose(fp);

	if (sql[0] == '\0')
		die("empty input");

	PG_TRY();
	{
		List *parsetree_list = raw_parser(sql, mode);
		if (parsetree_list == NIL)
			die("parser returned empty tree");
		char *out = nodeToString((const void *) parsetree_list);
		if (!out)
			die("nodeToString returned NULL");
		printf("%s\n", out);
	}
	PG_CATCH();
	{
		ErrorData *edata = CopyErrorData();
		FlushErrorState();
		fprintf(stderr, "pg_parse_helper: %s (sqlstate %s)\n",
				edata->message ? edata->message : "parse error",
				edata->sqlerrcode ? unpack_sql_state(edata->sqlerrcode) : "XXXXX");
		FreeErrorData(edata);
		return 1;
	}
	PG_END_TRY();

	return 0;
}
