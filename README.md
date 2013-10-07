Listens to AMQP messages and sends them to websockets client.

Disclaimer: This is my first Go project and I might be doing it wrong.

```sh
go get code.google.com/p/go.net/websocket github.com/streadway/amqp

rabbitmqadmin declare exchange name=test type=direct

go run src/amqpcast.go \
--amqp-url amqp://localhost:5672/ \
--amqp-exchange=test \
--amqp-key=test

open http://localhost:12345/

rabbitmqadmin publish exchange=test routing_key=test payload="hello world"
```

Thanks to @garyburd for this helpful gist:

https://gist.github.com/garyburd/1316852
