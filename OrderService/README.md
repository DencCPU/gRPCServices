# Order service

### Описание

**OrderService** позволяет клиентам создавать заказы и следить за их статусом. Статус можно получить как единоразовым запросом, так и в реальном времени — через открытое соединение, по которому сервис сам присылает обновления. Сервис работает на фрейвоке `gRPC`.

### Контракты

Так как сервис работает на фреймворке `gRPC`, ему требуется контракт. Контракт создан с помощью языка `protobuf`. Контракт для сервиса `OrderService` включает в себя следующие виды *rpc*-запросов:
- CreateOrder;
- GetOrderStatus;
- StreamGetOrder.
#### Создание заказа

Запрос на создание заказа приходит из сервиса `APIGetway`. Для этого используется метод `CreateOrder` из контракта *protobuf*. Сигнатура метода:
```protobuf
rpc CreateOrder(CreateOrderReq) returns (CreateOrderResp);
```
В качестве запроса принимается сообщение `CreateOrderReq`.
```protobuf
//Request for CreateOrder
message CreateOrderReq {
string user_id = 1 [(validate.rules).string.uuid = true];

string market_id = 2 [(validate.rules).string.uuid = true];

OrderType order_type = 3;

google.type.Money price = 4;

int64 quantity = 5 [(validate.rules).int64.gte = 1];

common.UserRole user_role = 6;

string indempotency_key = 7 [(validate.rules).string.min_len = 1];
}

enum OrderType {
ORDER_TYPE_UNSPECIFIED = 0;
ORDER_TYPE_NORMAL = 1;
ORDER_TYPE_EXPRESS = 2;
}

//Common package
enum UserRole {
USER_ROLE_UNSPECIFIED = 0;
USER_ROLE_BASIC_USER = 1;
USER_ROLE_PREMIUM_USER = 2;
USER_ROLE_ADMIN_USER = 3;
USER_ROLE_MODERATOR_USER = 4;
}
```
В запросе содержится идентификатор пользователя `user_id`, идентификатор рынка, на который делается заказ `market_id`, тип заказа `order_type`, цена заказа `price`,  количество единиц заказа `quantity`, роль пользователя `user_role`, и ключ идемпотентности, для предотвращения дублирования заказа `idempotency_key`.

Помимо валидации в самом контракте, происходит и валидация полученных данных на сервисе. Если переданные данные не совпадают по типу поля или отсутствуют, от метод вернет ошибку. 

Ответом на данный запрос является сообщение `CreateOrderResp`. 

```protobuf
//Response for CreateOrder
message CreateOrderResp {
string order_id = 1;
string order_status = 2;
}
```
Ответ метода, проверяется на вызывающем этот метод сервисе. В ответе указываются идентификатор заказа `order_id`, и статус заказа на момент созадания `order_status`.

#### Получение статуса заказа
Для получения статуса заказа используется *rpc*-метод контракта `GetOrderStatus`. Данный метод запрашивает сервис `APIGetway`.
Сигнатура метода:
```protobuf
rpc GetOrderStatus(GetOrderReq) returns (GetOrderResp);
```
В качестве запроса принимается сообщение `GetOrderReq`.
```protobuf
//Request for GetOrderStatus
message GetOrderReq {

string order_id = 1 [(validate.rules).string.uuid = true];

string user_id = 2 [(validate.rules).string.uuid = true];

}
```
В запросе передается идентификатор заказа `order_id` и идентификатор пользователя `user_id`.

Помимо валидации в самом контракте, происходит и валидация полученных данных на сервисе. Если переданные данные не совпадают по типу поля или отсутствуют, от метод вернет ошибку. 

В качестве ответа возвращается сообщение `GetOrderResp`.
```protobuf
//Response for GetOrderStatus
message GetOrderResp {
string order_status = 1;
string order_id = 2;
google.type.Money price = 3;
int64 quantity = 4;
string market_name = 5;
}
```
В ответе указывается статус заказа `order_status`, идентификатор заказа `order_id`, цена заказа `price`, количество единиц заказа `quantity` и наименование рынка, в котором сделан заказ `market_name`.

#### Получние статуса заказа в реальнос времени

Для получения статуса заказа в реальном времени используется *rpc* метод `StreamOrderUpdate`. 
Сигнатура метода:
```protobuf
rpc StreamOrderUpdate(StreamOrderUpdateReq) returns (stream StreamOrderUpdateResp);
```

В качестве запроса принимается сообщение `StreamOrderUpdateReq`.
```protobuf
message StreamOrderUpdateReq {

string order_id = 1 [(validate.rules).string.uuid = true];

string user_id = 2 [(validate.rules).string.uuid = true];

}
```
В запросе передается идентификатор заказа `order_id` и идентификатор пользователя `user_id`.

В качестве ответа возвращается стриминг сообщения `StreamOrderUpdateResp`.
```protobuf
message StreamOrderUpdateResp {

string order_status = 1;

google.protobuf.Timestamp update_status_time = 2;

}
```
В ответе указывается статус заказа `order_status` и время обновления статуса зазказа `update_status_time`.

### Запуск сервиса

Перед запуском сервиса требуется указать `host` и `port` на котором будет запущен сервис и opentelemetry collector. Указываются они в файле `.env`, расположенному по пути:
```shell
./OrderService/config/.env
```

Помимо этого, сервис использует *PostgresSQL* и *Kafka*. Для этих компонетнтов так же нужно указать параметры подключения. 

На сервисе искользуется только *Kafka consumer*, поэтому нужно указать список адесов брокеров, из которых *Kafka consumer* будет получать сообщения.

Перед запуском *PostgreSQL* требуется создать базу данных с названием, указанным в переменной окружения `POSTGRES_NAME`.

Пример:

```env
SERVER_HOST = "localhost"
SERVER_PORT = "8081"

POSTGRES_HOST = "localhost"
POSTGRES_PORT = "5432"
POSTGRES_USER = "postgres"
POSTGRES_PASSWORD = "12345"
POSTGRES_NAME = "order_service"

COLLECTOR_HOST = "localhost"
COLLECTOR_PORT = "4317"

KAFKA_BROKERS = "localhost:29092"
```

После указания адресов подключения, можно настроить отдельные компоненты работы сервиса в файле конфигурации расположенному по пути:

```shell
./OrderService/config/config.yaml
```

Пример `config.yaml`

```yaml
server:
	port:
	host:
	network: "tcp"
	request_per_second_limit: 100
	client_connection_timeout: 5s

postgres:
	host:
	port:
	user:
	password:
	name:
	sslmode: "disable"
	
	chan_size: 10
	idepmpotency_cache_ttl: 5s
	market_cache_ttl: 10s

  

breakersetting:
	name: "OrderService"
	max_request: 5
	interval: 60s
	timeout: 10s
	max_fail_request: 10

notify:
	ticker_interval: 1ms

collector:
	host:
	port:
	trace_percentage: 100
	metric_interval: 1s

kafka:
	brokers:
	topic: "markets"
	group_id: "order-service-group"
	min_bytes: 1
	max_bytes: 1048576
	max_wait: 2s
```
Как можно заметить, в данном `yaml` файле поля с названием `port` и `server` пустые. Информация для их заполения берется из переменных окружения обозначенных в файле `.env`.

`breakersetting` - это конфигуратор для настройки *circuit breaker*.

Так же к файле конфигурации присутствую и другие. Данные поля испольуются внутри бизнес логики сервиса.

После настройки, сервис готов к работе. Запуск из корня происходит с помощью команды:
```
go run OrderService/cmd/main.go
```

При успешном запуске сервиса в консоле выведется сообщение 
```shell
2026-06-04T19:41:15.707+0300    INFO  apprunner/fxappruner.go:418       The server is running on port:8081
```

Сервис запущен и готов к работе.