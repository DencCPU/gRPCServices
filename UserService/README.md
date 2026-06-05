# User service

### Описание 

UserService предназначен для регистрации и авторизации пользователей и управления их правами доступа.  Сервис работает на фрейвоке `gRPC`.

### Контракты

Так как сервис работает на фреймворке `gRPC`, ему требуется контракт. Контракт для сервис `UserService` создан с помощью языка `protobuf` и включает в себя следующие виды *rpc*-запросов:
- RegistrationNewUser;
- UpdateTokens;
- ValidationTokens;
- Authentication.

#### Регистрация пользователя

Регистрация нового пользователя в системе происходит при помощи метода `RegistrationNewUser`. Метод вызывается из сервиса `APIGetway`.
Сигнатура метода:
```protobuf
rpc RegistrationNewUser(RegistrationUserReq) returns (RegistrationUserResp)
```
В качестве запроса используется сообщение `RegistrationUserReq`.
```protobuf
message RegistrationUserReq {

string name = 1 [(validate.rules).string.min_len = 2];

string email = 2 [(validate.rules).string.email = true];

string password = 3 [(validate.rules).string.min_len = 8];

}
```
В запросе содержится информация о имени пользователя `name`, его `email` и пароль `password`.

 Помимо валидации в самом контракте, происходит и валидация полученных данных на сервисе. Если переданные данные не совпадают по типу поля или отсутствуют, то метод вернет ошибку. 

Ответом на запрос является сообщение `RegistrationUserResp`.
```protobuf
message RegistrationUserResp {

string access_token = 1;

string refresh_token = 2;

google.protobuf.Timestamp expire_at = 3;

}
```
В ответе метода содержится информация о короткоживущем токене, в котором содержится необходимая информация о пользователе `access_token`, долгоживущем токене `refresh_token` и время жизни токена `access._token` `expire_at`.
Если пользователь с данным `email` уже существует, то метод вернет ошибку.

#### Обновление токенов доступа

Обновление токенов доступа происходит, когда истекает время токенов доступа. За обновление токенов отвечает метод `UpdateTokens`.
Сигнатура метода:
```protobuf
rpc UpdateTokens(UpdateTokensReq) returns (UpdateTokensResp);
```

В качестве запроса используется сообщение `UpdateTokensReq`.
```protobuf
message UpdateTokensReq {

string access_token = 1 [(validate.rules).string.min_len = 1];

string refresh_token = 2 [(validate.rules).string.uuid = true];

}
```
В запросе передается короткоживущий токен `access_token` и долгоживущий токен `refresh_token`. `Refresh_token` нужен для проверки блокировки пользователя в системе или его наличия в целом.

Помимо валидации в самом контракте, происходит и валидация полученных данных на сервисе. Если переданные данные не совпадают по типу поля или отсутствуют, то метод вернет ошибку. 

В качестве ответа метод вернет сообщение `UpdateTokenResp`.
```protobuf
message UpdateTokensResp {

string access_token = 1;

string refresh_token = 2;

google.protobuf.Timestamp expire_at = 3;

}
```

В ответе метода содержится информация об обновленных короткоживущем токене, в котором содержится необходимая информация о пользователе `access_token`, долгоживущем токене `refresh_token` и время жизни токена `access._token` `expire_at`.

#### Валидация токена

Валидация токена доступа происходит с помощью метода `ValidationTokens`. Данный метод проверяет не истекло время жизни токена `access_token` и валиден ли он. Метод вызывается из сервиса `APIGetway`.
Сигнатура метода:
```protobuf
rpc ValidationTokens(ValidationReq) returns (ValidationResp);
```

В качестве запроса используется сообщение `ValidationReq`.

```protobuf
message ValidationReq {

string access_token = 1 [(validate.rules).string.min_len = 1];

}
```

В запросе передается короткоживущий токен `access_token`.

Помимо валидации в самом контракте, происходит и валидация полученных данных на сервисе. Если переданные данные не совпадают по типу поля или отсутствуют, то метод вернет ошибку. 

В качестве ответа метода передается сообщение `ValidationResp`.
```protobuf
message ValidationResp {

string user_id = 1;

common.UserRole role = 2;

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

В ответе содержится информация о идентификаторе пользователя `user_id` и роли пользователя `user_role`, необходимые для работы с другими сервисами.

#### Аутентификация пользователя

Аутентификациф пользователя происходит при помощи метода `Authentication`. Она требуется, если истекло время действия токена `refresh_token`.
Сигнатура метода:
```protobuf
rpc Authentication(AuthReq) returns (AuthResp);
```

В качестве запроса используется сообщение `AuthReq`. 
```protobuf
message AuthReq {

string email = 1 [(validate.rules).string.email = true];

string password = 2 [(validate.rules).string.min_len = 8];

}
```

В запросе передается информация о `email` пользователя и его пароль `password`. 
В качестве ответа метода возращает сообщение `AuthResp`.
```protobuf
message AuthResp {

string access_token = 1;

string refresh_token = 2;

google.protobuf.Timestamp expire_at = 3;

}
```
В ответе метода содержится информация об обновленных короткоживущем токене, в котором содержится необходимая информация о пользователе `access_token`, долгоживущем токене `refresh_token` и время жизни токена `access._token` `expire_at`.

Если `email` или `password`, переданные в запросе неверны или пользователь заблокирован, то метод вернет ошибку.

### Запуск сервиса

Перед запуском сервиса требуется указать `host` и `port` на котором будет запущен сервис и opentelemetry collector. Указываются они в файле `.env`, расположенному по пути:
```shell
./UserService/config/.env
```

Помимо этого, сервис использует *PostgresSQL* . Для этого компонента также нужно указать параметры подключения. 

Перед запуском *PostgreSQL* требуется создать базу данных с названием, указанным в переменной окружения `POSTGRES_NAME`.

Так как сервис использует *JWT* токены, то требуется и секретный ключ (подпись). Она так же указывается в переменной окуружения.

Пример:
```env
SERVER_HOST = "localhost"

SERVER_PORT = "8083"

  

POSTGRES_HOST = "localhost"

POSTGRES_PORT = "5432"

POSTGRES_USER = "postgres"

POSTGRES_PASSWORD = "12345"

POSTGRES_NAME = "user_service"

  

JWT_SECRET = "d354b87c-d70c-498c-99d3-631d93956033"

  

COLLECTOR_HOST = "localhost"

COLLECTOR_PORT = "4317"
```

После указания адресов подключения, можно настроить отдельные компоненты работы сервиса в файле конфигурации расположенному по пути:
```shell
./UserService/config/config.yaml
```

Пример `config.yaml`

```yaml
server:
	port:
	host:
	network: "tcp"
	request_per_second_limit: 100

postgres:
	host:
	port:
	user:
	password:
	name:
	sslmode: "disable"

jwt:
	secret:
	ttl: 60m

collector:
	host:
	port:
	trace_percentage: 100
	metric_interval: 1s
```

Как можно заметить, в данном `yaml` файле поля с названием `port` и `server` пустые. Информация для их заполения берется из переменных окружения обозначенных в файле `.env`.

`breakersetting` - это конфигуратор для настройки *circuit breaker*.

Так же к файле конфигурации присутствую и другие. Данные поля испольуются внутри бизнес логики сервиса.

После настройки, сервис готов к работе. Запуск из корня происходит с помощью команды:
```shell
go run UserService/cmd/main.go
```

При успешном запуске сервиса в консоле выведется сообщение:
```shell
2026-06-04T15:10:11.828+0300    INFOapprunner/fx_apprunner.go:273   Server start on port:8083
```
Сервис запущен и готов к работе.