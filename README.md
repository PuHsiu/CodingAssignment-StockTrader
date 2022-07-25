A simple code for stock trade demo

# How To Used?

1. Start a http server without TLS by command `make run`, server will listen to 8080 port.
2. Send request to server, and monitor state changed on terminal. See API from swagger.yaml.
3. By default, AAPL, GOOG, MSFT are available to trade.

# Design
3 layer architecture is used as design, 
1. Application layer map to presentation/, http restful API is support now, registered by RestRouter() func
2. Service layer map to services/, stock service is support now, matching the buy/sell order
3. Model layer map to model/, define a rich model to control model, like stock order.
