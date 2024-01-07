# Project Change Log
## v1.3.0 - (3 Changes)
- Plugged EC2 into internal service layers
- Updated github/workflows to be one for each lambda and then the ec2 version of project
- Added better logging to user repository layer

## v1.2.0 - (4 Changes)
- Added user repository layer
- Added auth service layer
- Plugged public lambda into auth service layer for registering and logging in
- Added intergration test xxx_it.go for running lambda tests locally 

## v1.1.0 - (2 Changes)
- Added basic template structure for ec2 version of project using github.com/gofiber/fiber (v1.14.6) with routes
- Added basic template strucutre for lamnbda version of project using github.com/aws/aws-lambda-go (v1.43.0) with no switching

## v1.0.0 - (5 Changes)
- Added database Init Script for starting up the shop database from scratch
- Added a basic "Hello World" APIGatewayWebsocketProxyRequest template for the auth lambda.
- Introduced a util package containing basic utility methods and functions.
- Introduced a constants package for HTTP error codes and descriptions.
- Added a types package for SocketErrors.
