package patterns

import (
	_ "embed"

	"github.com/Red-Sock/rscli/internal/io/folder"
)

// Clients connection files
var (
	//go:embed pattern_c/internal/clients/redis/conn.go.pattern
	redisConnFile []byte
	RedisConnFile = &folder.Folder{
		Name:    ConnFileName,
		Content: redisConnFile,
	}

	//go:embed pattern_c/internal/clients/sqldb/conn.go.pattern
	sqlConnFile []byte
	SqlConnFile = &folder.Folder{
		Name:    ConnFileName,
		Content: sqlConnFile,
	}

	//go:embed pattern_c/internal/clients/sqldb/postgres.go.pattern
	postgresConnFile []byte
	PostgresConnFile = &folder.Folder{
		Name:    "postgres.go",
		Content: postgresConnFile,
	}
	//go:embed pattern_c/internal/clients/sqldb/sqlite.go.pattern
	sqliteConnFile []byte
	SqliteConnFile = &folder.Folder{
		Name:    "sqlite.go",
		Content: sqliteConnFile,
	}

	//go:embed pattern_c/internal/clients/telegram/conn.go.pattern
	tgConnFile []byte
	TgConnFile = &folder.Folder{
		Name:    ConnFileName,
		Content: tgConnFile,
	}

	//go:embed pattern_c/internal/clients/grpc/conn.go.pattern
	grpcClientConnFile []byte
	GrpcClientConnFile = &folder.Folder{
		Name:    ConnFileName,
		Content: grpcClientConnFile,
	}
)
