# JWT-API

API было написано в рамках [тестового задания](https://medods.notion.site/Test-task-BackDev-623508ed85474f48a721e43ab00e9916)

## Подход
Использовался подход **Specification-first** или **API-first** с ипользованием инструментов [Open API](https://en.wikipedia.org/wiki/OpenAPI_Specification).

Коротко о подходе:

**API-First** - это подход к разработке серверных приложений, при котором API является наиболее важной 
частью продукта. При применении такого подхода у API продукта есть собственный цикл разработки, 
соответственно создается его артефакт(документ). Очевидно, что в сервисе присутствует 
зависимость от артефакта API, чем обеспечивается актуальность спецификации в этом 
артефакте. **Необходимо сначала разработать API, а потом уже работать над 
реализацией этого API в вашем сервисе**. Указанная очередность разработки 
является необходимостью при приверженности к подходу API-First, но не 
определением самого подхода. First в наименовании  означает, что 
API - это первый по важности продукт.

**Преимущества подхода**:

- Позволяет разрабатывать параллельно тесты, приложение-потребитель API и имплементацию API;
- Не нужно предполагать, на каком языке будет реализовано приложение-клиент;
- Появляется обязательный этап разработки API;
- Спецификация API всегда актуальна;
- Уменьшается риск возникновения различных ошибок при разработке;
- Уменьшаются трудозатраты на разработку.


## Имплементация

Для реализации API , используя подход API-first, были использованы инструменты для генерации кода на 
основе спецификации. Спецификация находится в корне проекта(openapi_v1.yml).
Для генерации кода существует несколько инструментов, выбран был [ogen](https://github.com/ogen-go/ogen). 
Сгенерированный код находится в папке [internal/oas/generated](./internal/oas/generated). Там же 
находится файл конфигурации инструмента [internal/oas/.ogen.yml](./internal/oas/.ogen.yml). 

Для создания JWT токенов была использована библиотека [golang-jwt](https://github.com/golang-jwt/jwt).

При написании кода использовался линтер [golangci-lint](https://github.com/golangci/golangci-lint/).

Взаимодействие с базой данных происходит при помощи миграций. Библиотека для этого - [goose](https://github.com/pressly/goose/).
Файлы миграции(структура БД):
[Инициализация](./migrations/20241217184221_init.sql):
```sql
-- +goose Up
-- +goose StatementBegin
create table if not exists refresh_tokens
(
    pair_id       uuid primary key,
    user_id       uuid        not null,
    refresh_token bytea       not null,
    created_at    timestamptz not null default now(),
    used          bool        not null default false
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists refresh_tokens;
-- +goose StatementEnd
```

[Установление отношений с таблицей юзеров](migrations/20241218183811_relation.sql):
```sql
-- +goose Up
-- +goose StatementBegin
create table if not exists users (
    id uuid primary key,
    email text not null
);

insert into users (id, email) values ('e05fa11d-eec3-4fba-b223-d6516800a047', 'test1@exampl.com');
insert into users (id, email) values ('3966749e-45d4-460d-8e59-34235672f03b', 'test2@exampl.com');
insert into users (id, email) values ('b3f7c269-1e35-4139-b882-2ec0b6629f7e', 'test3@exampl.com');
insert into users (id, email) values ('70d77738-2f5f-447e-8fa3-c36b238d9301', 'test4@exampl.com');

alter table refresh_tokens
    add constraint fk_user_id foreign key (user_id) references users(id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
drop table if exists users;
-- +goose StatementEnd
```

Миграции позволяют откатывать изменения, произведенные в файлах миграций. Инициализируются при запуске сервера:
```golang
db, err := sql.Open("postgres", config.PostgresDSN)
	if err != nil {
		logger.Fatal(err.Error())
}
...

err = goose.Up(db, ".")
	if err != nil {
		logger.Fatal(err.Error())
}

```

## Тесты

Написаны интеграционные и unit-тесты. Для интеграционных тестов была использована библиотека [testcontainers-go](https://github.com/testcontainers/testcontainers-go).
Unit-тесты использовали моковые данные, созданные при помощи библиотеки [gomock](https://github.com/uber-go/mock)

## Использование

Необходимо:
- запущенный Docker Desktop;
- curl.

Поднятие докер-контейнеров:

```shell
make docker-compose-up
```


Создание токенов:

```shell
curl -X 'POST' \
  'http://localhost:8000/users/3966749e-45d4-460d-8e59-34235672f03b/tokens' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -H 'X-Forwarded-For: 192.168.0.10' \
  -d '{
  "additionalInfo": {
    "additionalProp1": "string",
    "additionalProp2": "string",
    "additionalProp3": "string"
  }
}'
```

Рефреш токенов:

```shell
curl -X 'PUT' \
  'http://localhost:8000/users/3966749e-45d4-460d-8e59-34235672f03b/tokens' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -H 'X-Forwarded-For: 192.168.0.10' \
  -d '{"accessToken":"...","refreshToken":"..."}'
```

ВНИМАНИE! Вместо **...** в запросе должны быть созданные ранее токены!

Запуск тестов:

```shell
make test
```


