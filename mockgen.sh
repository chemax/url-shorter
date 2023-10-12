mockgen -source=./internal/storage/storage.go -destination=./mocks/storage/storage_mock.go StorageInterface

mockgen -source=./internal/db/db.go -destination=./mocks/db/db_mock.go

mockgen -source=./internal/config/config.go -destination=./mocks/config/config_mock.go

mockgen -source=./internal/logger/logger.go -destination=./mocks/logger/logger_mock.go

mockgen -destination=mocks/mock_store.go -package=mocks project/store Store