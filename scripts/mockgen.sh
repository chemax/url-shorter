~/go/bin/mockgen -source=./httpserver/server.go -destination=./mocks/httpserver/logger

~/go/bin/mockgen -source=./internal/db/db.go -destination=./mocks/db/db_mock.go

~/go/bin/mockgen -source=./internal/handlers/handlers.go -destination=./mocks/handlers/handlers_mock.go

~/go/bin/mockgen -source=./internal/storage/storage.go -destination=./mocks/storage/storage_mock.go

~/go/bin/mockgen -source=./users/users.go -destination=./mocks/users/users_mock.go
