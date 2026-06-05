# Spot instrument serice

### Описание

Spot instrument service предназначен для управления работой рынков и доступа к нимб предоставления списка рынков, доступных пользователю по его запросу.
Сервис работает на фрейвоке `gRPC`.

### Контракты

Так как сервис работает на фреймворке `gRPC`, ему требуется контракт. Контракт для сервис `UserService` создан с помощью языка `protobuf` и включает в себя следующие виды *rpc*-запросов:
 - ViewMarket.

#### Предоставления списка доступных рынков

Для предоставления списка дотупных рынков используется метод `ViewMarket`.
Сигнатура метода:
```protobuf
rpc ViewMarket(ViewMarketReq) returns (ViewMarketResp);
```

В качестве запроса используется сообщение `ViewMarketReq`.
```protobuf
message ViewMarketReq {

common.UserRole user_roles = 1;

string user_id = 2 [(validate.rules).string.uuid = true];

int32 page_size = 3 [(validate.rules).int32 = {

gte: 0

lte: 50

}];

string page_token = 4;

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
В запросе используется информация о роли пользователя `user_roles`, идентификатор пользователя `user_id`, количество рынков, возвращаемое за один запрос `page_size`, токен рынка, полученный в предыдущем запросе, если запрос не первичный `page_token`. `Page_size` и `page_token` используются для пагинации.

Помимо валидации в самом контракте, происходит и валидация полученных данных на сервисе. Если переданные данные не совпадают по типу поля или отсутствуют, то метод вернет ошибку. 

В качестве ответа используется сообщение `ViewMarketResp`. 
```protobuf
message ViewMarketResp {

repeated Market enable_markets = 1;

string next_page_token = 2;

}

message Market {

string market_id = 1 [(validate.rules).string.uuid = true];

string market_name = 2 [(validate.rules).string.min_len = 1];

}
```
В ответе содержится информация о рынках, доступных пользователю на момент запроса `market_id` и `market_name`, токен последнего рынка в предоставленном списке `next_page_token`.
Если на момент запроса не один рынок для пользователя недоступен, метод вернет ошибку.

### Запуск сервиса

Перед запуском сервиса требуется указать `host` и `port` на котором будет запущен сервис и opentelemetry collector. Указываются они в файле `.env`, расположенному по пути:
```shell
./SpotInstrumentService/config/.env
```

Помимо этого, сервис использует *Kafka* и *Redis*. Для этих компонентов так же нужно указать параметры подключения. Сервис использует *Kafka reader*, поэтому к нему так же нужно указать `host` и `port`, а к *Redis* еще и `password`. 

Пример:
```protobuf
REDIS_HOST = "localhost"

REDIS_PORT = "6379"

REDIS_PASSWORD = ""

  

SERVER_HOST = "localhost"

SERVER_PORT = "8080"

  

COLLECTOR_HOST = "localhost"

COLLECTOR_PORT = "4317"

  

KAFKA_HOST = "localhost"

KAFKA_PORT = "29092"
```

После указания адресов подключения, можно настроить отдельные компоненты работы сервиса в файле конфигурации расположенному по пути:
```shell
./SpotInstrumentService/config/config.yaml
```

Пример `config.yaml`

```yaml
server:
	port:
	host:
	network: "tcp"
	request_per_second_limit: 100

storage:
	timeout: 10s

redis:
	host:
	port:
	password:
	db: 0
	pool_size: 8
	min_id_le_conns: 5
	dial_timeout: 5s
	read_timeout: 3s
	write_timeout: 3s
	timeout: 10s

collector:
	host:
	port:
	trace_percentage: 100
	metric_interval: 1s
  
kafka:
	host:
	port:
	topic: "markets"
	max_attempts: 10
	write_backoff_min: 100ms
	write_backoff_max: 2s
	batch_size: 1
	batch_bytes: 1024
	get_markets_interval: 10s
	relay_interval: 3s
```

Как можно заметить, в данном `yaml` файле поля с названием `port` и `server` пустые. Информация для их заполения берется из переменных окружения обозначенных в файле `.env`.

`storage` - это конфигуратор для управления рынками. Поле `timeout` отвечает за интервал изменения состояния рынка.

Так же к файле конфигурации присутствую и другие. Данные поля испольуются внутри бизнес логики сервиса.

После настройки, сервис готов к работе. Запуск из корня происходит с помощью команды:
```shell
go run SpotInstrumentService/cmd/main.go
```

При успешном запуске сервиса в консоле выведется сообщение:
```shell
2026-06-05T22:59:33.796+0300    INFO   apprunner/fx_apprunner.go:383    Server start on port:8080
```
Сервис готов к работе.