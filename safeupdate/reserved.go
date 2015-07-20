package safeupdate

// DO NOT EDIT THIS FILE, it's been auto-generated
//
//go:generate go run update-reserved/main.go

var reservedWords = map[string]bool{
	"ABORT":          true,
	"ABSOLUTE":       true,
	"ACTION":         true,
	"ADD":            true,
	"AFTER":          true,
	"AGENT":          true,
	"AGGREGATE":      true,
	"ALL":            true,
	"ALLOCATE":       true,
	"ALTER":          true,
	"ANALYZE":        true,
	"AND":            true,
	"ANY":            true,
	"ARCHIVE":        true,
	"ARE":            true,
	"ARRAY":          true,
	"AS":             true,
	"ASC":            true,
	"ASCII":          true,
	"ASENSITIVE":     true,
	"ASSERTION":      true,
	"ASYMMETRIC":     true,
	"AT":             true,
	"ATOMIC":         true,
	"ATTACH":         true,
	"ATTRIBUTE":      true,
	"AUTH":           true,
	"AUTHORIZATION":  true,
	"AUTHORIZE":      true,
	"AUTO":           true,
	"AVG":            true,
	"BACK":           true,
	"BACKUP":         true,
	"BASE":           true,
	"BATCH":          true,
	"BEFORE":         true,
	"BEGIN":          true,
	"BETWEEN":        true,
	"BIGINT":         true,
	"BINARY":         true,
	"BIT":            true,
	"BLOB":           true,
	"BLOCK":          true,
	"BOOLEAN":        true,
	"BOTH":           true,
	"BREADTH":        true,
	"BUCKET":         true,
	"BULK":           true,
	"BY":             true,
	"BYTE":           true,
	"CALL":           true,
	"CALLED":         true,
	"CALLING":        true,
	"CAPACITY":       true,
	"CASCADE":        true,
	"CASCADED":       true,
	"CASE":           true,
	"CAST":           true,
	"CATALOG":        true,
	"CHAR":           true,
	"CHARACTER":      true,
	"CHECK":          true,
	"CLASS":          true,
	"CLOB":           true,
	"CLOSE":          true,
	"CLUSTER":        true,
	"CLUSTERED":      true,
	"CLUSTERING":     true,
	"CLUSTERS":       true,
	"COALESCE":       true,
	"COLLATE":        true,
	"COLLATION":      true,
	"COLLECTION":     true,
	"COLUMN":         true,
	"COLUMNS":        true,
	"COMBINE":        true,
	"COMMENT":        true,
	"COMMIT":         true,
	"COMPACT":        true,
	"COMPILE":        true,
	"COMPRESS":       true,
	"CONDITION":      true,
	"CONFLICT":       true,
	"CONNECT":        true,
	"CONNECTION":     true,
	"CONSISTENCY":    true,
	"CONSISTENT":     true,
	"CONSTRAINT":     true,
	"CONSTRAINTS":    true,
	"CONSTRUCTOR":    true,
	"CONSUMED":       true,
	"CONTINUE":       true,
	"CONVERT":        true,
	"COPY":           true,
	"CORRESPONDING":  true,
	"COUNT":          true,
	"COUNTER":        true,
	"CREATE":         true,
	"CROSS":          true,
	"CUBE":           true,
	"CURRENT":        true,
	"CURSOR":         true,
	"CYCLE":          true,
	"DATA":           true,
	"DATABASE":       true,
	"DATE":           true,
	"DATETIME":       true,
	"DAY":            true,
	"DEALLOCATE":     true,
	"DEC":            true,
	"DECIMAL":        true,
	"DECLARE":        true,
	"DEFAULT":        true,
	"DEFERRABLE":     true,
	"DEFERRED":       true,
	"DEFINE":         true,
	"DEFINED":        true,
	"DEFINITION":     true,
	"DELETE":         true,
	"DELIMITED":      true,
	"DEPTH":          true,
	"DEREF":          true,
	"DESC":           true,
	"DESCRIBE":       true,
	"DESCRIPTOR":     true,
	"DETACH":         true,
	"DETERMINISTIC":  true,
	"DIAGNOSTICS":    true,
	"DIRECTORIES":    true,
	"DISABLE":        true,
	"DISCONNECT":     true,
	"DISTINCT":       true,
	"DISTRIBUTE":     true,
	"DO":             true,
	"DOMAIN":         true,
	"DOUBLE":         true,
	"DROP":           true,
	"DUMP":           true,
	"DURATION":       true,
	"DYNAMIC":        true,
	"EACH":           true,
	"ELEMENT":        true,
	"ELSE":           true,
	"ELSEIF":         true,
	"EMPTY":          true,
	"ENABLE":         true,
	"END":            true,
	"EQUAL":          true,
	"EQUALS":         true,
	"ERROR":          true,
	"ESCAPE":         true,
	"ESCAPED":        true,
	"EVAL":           true,
	"EVALUATE":       true,
	"EXCEEDED":       true,
	"EXCEPT":         true,
	"EXCEPTION":      true,
	"EXCEPTIONS":     true,
	"EXCLUSIVE":      true,
	"EXEC":           true,
	"EXECUTE":        true,
	"EXISTS":         true,
	"EXIT":           true,
	"EXPLAIN":        true,
	"EXPLODE":        true,
	"EXPORT":         true,
	"EXPRESSION":     true,
	"EXTENDED":       true,
	"EXTERNAL":       true,
	"EXTRACT":        true,
	"FAIL":           true,
	"FALSE":          true,
	"FAMILY":         true,
	"FETCH":          true,
	"FIELDS":         true,
	"FILE":           true,
	"FILTER":         true,
	"FILTERING":      true,
	"FINAL":          true,
	"FINISH":         true,
	"FIRST":          true,
	"FIXED":          true,
	"FLATTERN":       true,
	"FLOAT":          true,
	"FOR":            true,
	"FORCE":          true,
	"FOREIGN":        true,
	"FORMAT":         true,
	"FORWARD":        true,
	"FOUND":          true,
	"FREE":           true,
	"FROM":           true,
	"FULL":           true,
	"FUNCTION":       true,
	"FUNCTIONS":      true,
	"GENERAL":        true,
	"GENERATE":       true,
	"GET":            true,
	"GLOB":           true,
	"GLOBAL":         true,
	"GO":             true,
	"GOTO":           true,
	"GRANT":          true,
	"GREATER":        true,
	"GROUP":          true,
	"GROUPING":       true,
	"HANDLER":        true,
	"HASH":           true,
	"HAVE":           true,
	"HAVING":         true,
	"HEAP":           true,
	"HIDDEN":         true,
	"HOLD":           true,
	"HOUR":           true,
	"IDENTIFIED":     true,
	"IDENTITY":       true,
	"IF":             true,
	"IGNORE":         true,
	"IMMEDIATE":      true,
	"IMPORT":         true,
	"IN":             true,
	"INCLUDING":      true,
	"INCLUSIVE":      true,
	"INCREMENT":      true,
	"INCREMENTAL":    true,
	"INDEX":          true,
	"INDEXED":        true,
	"INDEXES":        true,
	"INDICATOR":      true,
	"INFINITE":       true,
	"INITIALLY":      true,
	"INLINE":         true,
	"INNER":          true,
	"INNTER":         true,
	"INOUT":          true,
	"INPUT":          true,
	"INSENSITIVE":    true,
	"INSERT":         true,
	"INSTEAD":        true,
	"INT":            true,
	"INTEGER":        true,
	"INTERSECT":      true,
	"INTERVAL":       true,
	"INTO":           true,
	"INVALIDATE":     true,
	"IS":             true,
	"ISOLATION":      true,
	"ITEM":           true,
	"ITEMS":          true,
	"ITERATE":        true,
	"JOIN":           true,
	"KEY":            true,
	"KEYS":           true,
	"LAG":            true,
	"LANGUAGE":       true,
	"LARGE":          true,
	"LAST":           true,
	"LATERAL":        true,
	"LEAD":           true,
	"LEADING":        true,
	"LEAVE":          true,
	"LEFT":           true,
	"LENGTH":         true,
	"LESS":           true,
	"LEVEL":          true,
	"LIKE":           true,
	"LIMIT":          true,
	"LIMITED":        true,
	"LINES":          true,
	"LIST":           true,
	"LOAD":           true,
	"LOCAL":          true,
	"LOCALTIME":      true,
	"LOCALTIMESTAMP": true,
	"LOCATION":       true,
	"LOCATOR":        true,
	"LOCK":           true,
	"LOCKS":          true,
	"LOG":            true,
	"LOGED":          true,
	"LONG":           true,
	"LOOP":           true,
	"LOWER":          true,
	"MAP":            true,
	"MATCH":          true,
	"MATERIALIZED":   true,
	"MAX":            true,
	"MAXLEN":         true,
	"MEMBER":         true,
	"MERGE":          true,
	"METHOD":         true,
	"METRICS":        true,
	"MIN":            true,
	"MINUS":          true,
	"MINUTE":         true,
	"MISSING":        true,
	"MOD":            true,
	"MODE":           true,
	"MODIFIES":       true,
	"MODIFY":         true,
	"MODULE":         true,
	"MONTH":          true,
	"MULTI":          true,
	"MULTISET":       true,
	"NAME":           true,
	"NAMES":          true,
	"NATIONAL":       true,
	"NATURAL":        true,
	"NCHAR":          true,
	"NCLOB":          true,
	"NEW":            true,
	"NEXT":           true,
	"NO":             true,
	"NONE":           true,
	"NOT":            true,
	"NULL":           true,
	"NULLIF":         true,
	"NUMBER":         true,
	"NUMERIC":        true,
	"OBJECT":         true,
	"OF":             true,
	"OFFLINE":        true,
	"OFFSET":         true,
	"OLD":            true,
	"ON":             true,
	"ONLINE":         true,
	"ONLY":           true,
	"OPAQUE":         true,
	"OPEN":           true,
	"OPERATOR":       true,
	"OPTION":         true,
	"OR":             true,
	"ORDER":          true,
	"ORDINALITY":     true,
	"OTHER":          true,
	"OTHERS":         true,
	"OUT":            true,
	"OUTER":          true,
	"OUTPUT":         true,
	"OVER":           true,
	"OVERLAPS":       true,
	"OVERRIDE":       true,
	"OWNER":          true,
	"PAD":            true,
	"PARALLEL":       true,
	"PARAMETER":      true,
	"PARAMETERS":     true,
	"PARTIAL":        true,
	"PARTITION":      true,
	"PARTITIONED":    true,
	"PARTITIONS":     true,
	"PATH":           true,
	"PERCENT":        true,
	"PERCENTILE":     true,
	"PERMISSION":     true,
	"PERMISSIONS":    true,
	"PIPE":           true,
	"PIPELINED":      true,
	"PLAN":           true,
	"POOL":           true,
	"POSITION":       true,
	"PRECISION":      true,
	"PREPARE":        true,
	"PRESERVE":       true,
	"PRIMARY":        true,
	"PRIOR":          true,
	"PRIVATE":        true,
	"PRIVILEGES":     true,
	"PROCEDURE":      true,
	"PROCESSED":      true,
	"PROJECT":        true,
	"PROJECTION":     true,
	"PROPERTY":       true,
	"PROVISIONING":   true,
	"PUBLIC":         true,
	"PUT":            true,
	"QUERY":          true,
	"QUIT":           true,
	"QUORUM":         true,
	"RAISE":          true,
	"RANDOM":         true,
	"RANGE":          true,
	"RANK":           true,
	"RAW":            true,
	"READ":           true,
	"READS":          true,
	"REAL":           true,
	"REBUILD":        true,
	"RECORD":         true,
	"RECURSIVE":      true,
	"REDUCE":         true,
	"REF":            true,
	"REFERENCE":      true,
	"REFERENCES":     true,
	"REFERENCING":    true,
	"REGEXP":         true,
	"REGION":         true,
	"REINDEX":        true,
	"RELATIVE":       true,
	"RELEASE":        true,
	"REMAINDER":      true,
	"RENAME":         true,
	"REPEAT":         true,
	"REPLACE":        true,
	"REQUEST":        true,
	"RESET":          true,
	"RESIGNAL":       true,
	"RESOURCE":       true,
	"RESPONSE":       true,
	"RESTORE":        true,
	"RESTRICT":       true,
	"RESULT":         true,
	"RETURN":         true,
	"RETURNING":      true,
	"RETURNS":        true,
	"REVERSE":        true,
	"REVOKE":         true,
	"RIGHT":          true,
	"ROLE":           true,
	"ROLES":          true,
	"ROLLBACK":       true,
	"ROLLUP":         true,
	"ROUTINE":        true,
	"ROW":            true,
	"ROWS":           true,
	"RULE":           true,
	"RULES":          true,
	"SAMPLE":         true,
	"SATISFIES":      true,
	"SAVE":           true,
	"SAVEPOINT":      true,
	"SCAN":           true,
	"SCHEMA":         true,
	"SCOPE":          true,
	"SCROLL":         true,
	"SEARCH":         true,
	"SECOND":         true,
	"SECTION":        true,
	"SEGMENT":        true,
	"SEGMENTS":       true,
	"SELECT":         true,
	"SELF":           true,
	"SEMI":           true,
	"SENSITIVE":      true,
	"SEPARATE":       true,
	"SEQUENCE":       true,
	"SERIALIZABLE":   true,
	"SESSION":        true,
	"SET":            true,
	"SETS":           true,
	"SHARD":          true,
	"SHARE":          true,
	"SHARED":         true,
	"SHORT":          true,
	"SHOW":           true,
	"SIGNAL":         true,
	"SIMILAR":        true,
	"SIZE":           true,
	"SKEWED":         true,
	"SMALLINT":       true,
	"SNAPSHOT":       true,
	"SOME":           true,
	"SOURCE":         true,
	"SPACE":          true,
	"SPACES":         true,
	"SPARSE":         true,
	"SPECIFIC":       true,
	"SPECIFICTYPE":   true,
	"SPLIT":          true,
	"SQL":            true,
	"SQLCODE":        true,
	"SQLERROR":       true,
	"SQLEXCEPTION":   true,
	"SQLSTATE":       true,
	"SQLWARNING":     true,
	"START":          true,
	"STATE":          true,
	"STATIC":         true,
	"STATUS":         true,
	"STORAGE":        true,
	"STORE":          true,
	"STORED":         true,
	"STREAM":         true,
	"STRING":         true,
	"STRUCT":         true,
	"STYLE":          true,
	"SUB":            true,
	"SUBMULTISET":    true,
	"SUBPARTITION":   true,
	"SUBSTRING":      true,
	"SUBTYPE":        true,
	"SUM":            true,
	"SUPER":          true,
	"SYMMETRIC":      true,
	"SYNONYM":        true,
	"SYSTEM":         true,
	"TABLE":          true,
	"TABLESAMPLE":    true,
	"TEMP":           true,
	"TEMPORARY":      true,
	"TERMINATED":     true,
	"TEXT":           true,
	"THAN":           true,
	"THEN":           true,
	"THROUGHPUT":     true,
	"TIME":           true,
	"TIMESTAMP":      true,
	"TIMEZONE":       true,
	"TINYINT":        true,
	"TO":             true,
	"TOKEN":          true,
	"TOTAL":          true,
	"TOUCH":          true,
	"TRAILING":       true,
	"TRANSACTION":    true,
	"TRANSFORM":      true,
	"TRANSLATE":      true,
	"TRANSLATION":    true,
	"TREAT":          true,
	"TRIGGER":        true,
	"TRIM":           true,
	"TRUE":           true,
	"TRUNCATE":       true,
	"TTL":            true,
	"TUPLE":          true,
	"TYPE":           true,
	"UNDER":          true,
	"UNDO":           true,
	"UNION":          true,
	"UNIQUE":         true,
	"UNIT":           true,
	"UNKNOWN":        true,
	"UNLOGGED":       true,
	"UNNEST":         true,
	"UNPROCESSED":    true,
	"UNSIGNED":       true,
	"UNTIL":          true,
	"UPDATE":         true,
	"UPPER":          true,
	"URL":            true,
	"USAGE":          true,
	"USE":            true,
	"USER":           true,
	"USERS":          true,
	"USING":          true,
	"UUID":           true,
	"VACUUM":         true,
	"VALUE":          true,
	"VALUED":         true,
	"VALUES":         true,
	"VARCHAR":        true,
	"VARIABLE":       true,
	"VARIANCE":       true,
	"VARINT":         true,
	"VARYING":        true,
	"VIEW":           true,
	"VIEWS":          true,
	"VIRTUAL":        true,
	"VOID":           true,
	"WAIT":           true,
	"WHEN":           true,
	"WHENEVER":       true,
	"WHERE":          true,
	"WHILE":          true,
	"WINDOW":         true,
	"WITH":           true,
	"WITHIN":         true,
	"WITHOUT":        true,
	"WORK":           true,
	"WRAPPED":        true,
	"WRITE":          true,
	"YEAR":           true,
	"ZONE":           true,
}
