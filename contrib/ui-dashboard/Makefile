all:
	python -m grpc_tools.protoc -I ../ -I../proto --python_out=. --grpc_python_out=. ../proto/*.proto
	python -m grpc_tools.protoc -I ../ -I../schema --python_out=. --grpc_python_out=. ../schema/*.proto
	python -m grpc_tools.protoc -I ../ -I../api --python_out=. --grpc_python_out=. ../api/*.proto
