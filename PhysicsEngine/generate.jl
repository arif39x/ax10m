using ProtoBuf

protoc(`-I=../ --julia_out=src ../math.proto`)
