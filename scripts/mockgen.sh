mockgen -source=./interfaces/interfaces.go -destination=./mocks/storage/storage_mock.go StorageInterface
mockgen -source=./interfaces/interfaces.go -destination=./mocks/db/db_mock.go DBInterface
mockgen -source=./interfaces/interfaces.go -destination=./mocks/logger/log_mock.go LoggerInterface
mockgen -source=./interfaces/interfaces.go -destination=./mocks/config/config_mock.go ConfigInterface
mockgen -source=./interfaces/interfaces.go -destination=./mocks/users/users_mock.go UsersInterface