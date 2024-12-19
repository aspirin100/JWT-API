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

Для создания JWT токенов была использована библиотека [golang-jwt](github.com/golang-jwt/jwt/v5).

При написании кода использовался линтер [golangci-lint](github.com/golangci/golangci-lint/cmd/golangci-lint)

## Тесты

Написаны интеграционные и unit-тесты. Для интеграционных тестов была использована библиотека [testcontainers-go](github.com/testcontainers/testcontainers-go). Unit-тесты использовали моковые данные, созданные при помощи библиотеки [gomock](go.uber.org/mock/gomock)

## Использование

Необходимо:
- Docker;
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


