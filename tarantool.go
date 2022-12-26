package gormtnt

import (
	"context"
	"database/sql"
	"fmt"
	"regexp"

	_ "github.com/aeroideaservices/tnt"

	"gorm.io/gorm"
	"gorm.io/gorm/callbacks"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/migrator"
	"gorm.io/gorm/schema"
)

type Config struct {
	DriverName               string
	ServerVersion            string
	DSN                      string
	Conn                     gorm.ConnPool
	DefaultDatetimePrecision *int
}

type Dialector struct {
	*Config
}

var (
	CrateClauses             = []string{"INSERT", "VALUES"}
	QueryClauses             = []string{}
	UpdateClauses            = []string{"UPDATE", "SET", "WHERE", "ORDER BY", "LIMIT"}
	DeleteClauses            = []string{"DELETE", "FROM", "WHERE", "ORDER BY", "LIMIT"}
	defaultDatetimePrecision = 3
)

func Open(dsn string) gorm.Dialector {
	return &Dialector{Config: &Config{DSN: dsn}}
}

func New(config Config) gorm.Dialector {
	return &Dialector{Config: &config}
}
func (dialector Dialector) Name() string {
	return "tarantool"
}

func (dialector Dialector) Initialize(db *gorm.DB) (err error) {
	if dialector.DriverName == "" {
		dialector.DriverName = "tnt"
	}

	if dialector.DefaultDatetimePrecision == nil {
		dialector.DefaultDatetimePrecision = &defaultDatetimePrecision
	}

	if dialector.Conn != nil {
		db.ConnPool = dialector.Conn
	} else {
		db.ConnPool, err = sql.Open(dialector.DriverName, dialector.DSN)
		if err != nil {
			return err
		}
	}

	if dialector.ServerVersion == "" {
		err = db.ConnPool.QueryRowContext(context.Background(), "SELECT VERSION()").Scan(&dialector.ServerVersion)
		if err != nil {
			return err
		}
	}

	callbackConfig := &callbacks.Config{
		CreateClauses: CrateClauses,
		QueryClauses:  QueryClauses,
		UpdateClauses: UpdateClauses,
		DeleteClauses: DeleteClauses,
	}

	callbacks.RegisterDefaultCallbacks(db, callbackConfig)

	for k, v := range dialector.ClauseBuilders() {
		db.ClauseBuilders[k] = v
	}
	return
}

func (dialector Dialector) Migrator(db *gorm.DB) gorm.Migrator {
	return Migrator{migrator.Migrator{Config: migrator.Config{
		DB:                          db,
		Dialector:                   dialector,
		CreateIndexAfterCreateTable: true,
	}}}
}

func (dialector Dialector) DataTypeOf(field *schema.Field) string {
	if t, ok := field.TagSettings["type"]; ok {
		switch t {
		case "uuid":
			return "UUID"
		case "datetime":
			return "DATETIME"
		}
	}
	switch field.DataType {
	case schema.Bool:
		return "BOOLEAN"
	case schema.Int:
		return "INTEGER"
	case schema.Uint:
		return "UNSIGNED"
	case schema.Float:
		return "DECIMAL"
	case schema.String:
		if field.Size > 0 {
			return fmt.Sprintf("VARCHAR(%d)", field.Size)
		}
		return "TEXT"
	case schema.Time:
		return "DATETIME" // todo
	case schema.Bytes:
		return "VARBINARY" // todo
	default:
		return string(field.DataType)
	}
}

func (dialector Dialector) DefaultValueOf(_ *schema.Field) clause.Expression {
	return clause.Expr{SQL: "DEFAULT"}
}

func (dialector Dialector) BindVarTo(writer clause.Writer, stmt *gorm.Statement, v interface{}) {
	writer.WriteByte('?')
}

func (dialector Dialector) QuoteTo(writer clause.Writer, str string) {
	var (
		underQuoted, selfQuoted bool
		continuousBacktick      int8
		shiftDelimiter          int8
	)

	for _, v := range []byte(str) {
		switch v {
		case '"':
			continuousBacktick++
			if continuousBacktick == 2 {
				writer.WriteString(`""`)
				continuousBacktick = 0
			}
		case '.':
			if continuousBacktick > 0 || !selfQuoted {
				shiftDelimiter = 0
				underQuoted = false
				continuousBacktick = 0
				writer.WriteByte('"')
			}
			writer.WriteByte(v)
			continue
		default:
			if shiftDelimiter-continuousBacktick <= 0 && !underQuoted {
				writer.WriteByte('"')
				underQuoted = true
				if selfQuoted = continuousBacktick > 0; selfQuoted {
					continuousBacktick -= 1
				}
			}

			for ; continuousBacktick > 0; continuousBacktick -= 1 {
				writer.WriteString(`""`)
			}

			writer.WriteByte(v)
		}
		shiftDelimiter++
	}

	if continuousBacktick > 0 && !selfQuoted {
		writer.WriteString(`""`)
	}
	writer.WriteByte('"')
}

var numericPlaceholder = regexp.MustCompile(`\$(\d+)`)

func (dialector Dialector) Explain(sql string, vars ...interface{}) string {
	return logger.ExplainSQL(sql, numericPlaceholder, `'`, vars...)
}

const (
	ClauseValues = "VALUES"
)

func (dialector Dialector) ClauseBuilders() map[string]clause.ClauseBuilder {
	clauseBuilders := map[string]clause.ClauseBuilder{
		ClauseValues: func(c clause.Clause, builder clause.Builder) {
			if values, ok := c.Expression.(clause.Values); ok && len(values.Columns) == 0 {
				builder.WriteString("VALUES()")
				return
			}
			c.Build(builder)
		},
	}
	return clauseBuilders
}
