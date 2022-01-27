@echo off

set CUR=%~dp0

pushd

cd /D "%CUR%"

set CMD=protoc --proto_path=protobuf --go_opt=paths=source_relative --go_out=server/rpc --go-grpc_out=server/rpc --go-grpc_opt=paths=source_relative protobuf\*.proto
echo %CMD%
%CMD%

popd
