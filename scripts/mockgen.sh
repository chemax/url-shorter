mockgen -source=./util/util.go -destination=./mocks/storage/storage_mock.go StorageInterface
mockgen -source=./util/util.go -destination=./mocks/db/db_mock.go DBInterface
mockgen -source=./util/util.go -destination=./mocks/logger/log_mock.go LoggerInterface
mockgen -source=./util/util.go -destination=./mocks/config/config_mock.go ConfigInterface