~/go/bin/mockgen -source=./interfaces/interfaces.go -destination=./mocks/storage/storage_mock.go StorageInterface
~/go/bin/mockgen -source=./interfaces/interfaces.go -destination=./mocks/db/db_mock.go DBInterface
~/go/bin/mockgen -source=./interfaces/interfaces.go -destination=./mocks/logger/log_mock.go LoggerInterface
~/go/bin/mockgen -source=./interfaces/interfaces.go -destination=./mocks/config/config_mock.go ConfigInterface
~/go/bin/mockgen -source=./interfaces/interfaces.go -destination=./mocks/users/users_mock.go UsersInterface