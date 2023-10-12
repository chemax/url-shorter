mockgen -source=./util/util.go -destination=./mocks/storage/storage_mock.go StorageInterface
mockgen -source=./util/util.go -destination=./mocks/storage/storage_mock.go DBInterface
mockgen -source=./util/util.go -destination=./mocks/storage/storage_mock.go LoggerInterface
mockgen -source=./util/util.go -destination=./mocks/storage/storage_mock.go ConfigInterface