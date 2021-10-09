# Open-Platform

the api is in [yapi](https://yapi.hustunique.com/project/11/interface/api)

The Open-platform is the toolbox of the UniqueStudio. 

There are 2 ways to use it:

- HTTP (old way)
- gRPC (recommended)

The HTTP APIs are in yapi below. For gRPC, the proto IDL is in [UniqueIDL.SMS](https://github.com/UniqueStudio/UniqueIDL/blob/master/sms.proto). 

In deployment environment, all services are in a subnet, therefore, using gRPC without TLS is sufficient. For debug use, TLS is required. The server cert should ask DevOps to acquire.

## Ability

1. Push SMS