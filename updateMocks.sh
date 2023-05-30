#!/bin/bash

mockgen -source=models/apiModels/storage.go -destination=./handlers/mocks/IStorageMocks.go
