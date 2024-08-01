# Сервис обработки сообщений

API предоставляет возможность отправки смс сообщений с кодом подтверждения, сбор метрик, используя Prometheus, визуализацию метрик в Grafana

Архитектура сервиса разработана в соответствии с паттерном Transactional Outbox. Представляет из себя 3 сущности:
- Интерфейс взаимодействия с клиентом (API). Обрабатывает POST HTTP запросы на эндпоинт и сохраняет данные в Postgres. 
- Процессор. С определенным интервалом времени обращается к базе данных, резервирует и получает сообщения, после чего отправляет их в Kafka. Количество получаемых сообщений в рамках одной итерации и интервал конфигурируются в .env, параметры `PROCESSOR_TAKE_MESSAGE_LIMIT` и `PROCESSOR_RESERVATION_TIME` соответственно. 
- Консьюмер. Получает сообщения из Kafka и отправляет запрос с данными на сторонний сервис для дальнейшей обработки и отправки смс-сообщения. 

Для запуска требуется наличие установленных [git](https://git-scm.com), [docker](https://www.docker.com), Make.

Конфигурация по умолчанию:
- Postgres прослушивает порт `5432`
- Kafka прослушивает порт `29092`
- Prometheus прослушивает порт `9090`
- Grafana прослушивает порт `3000`
- API прослушивает порты `8080` и `8082` для метрик

При необходимости, конфигурацию можно изменить в `.env`, находящимся в корневой директории проекта:
- `APP_PORT` - Порт, которые прослушивает API
- `APP_MODE` - Debug мод (debug, release)
- `PROCESSOR_TAKE_MESSAGE_LIMIT` - Количество сообщений, получаемых процессором из базы данных в рамках одной итерации. 
- `PROCESSOR_RESERVATION_TIME` - Время, на которое сообщение будет заблокировано для получения другими процессорами из базы данных.
- `PROCESSOR_TOPIC` - Топик Kafka, в который процессор отправляет сообщения. 
- `POSTGRES_HOST` - Хост Postgres.
- `POSTGRES_PORT` - Порт Postgres.
- `POSTGRES_DB` - Название базы в Postgres.
- `POSTGRES_USER` - Имя пользователя в Postgres.
- `POSTGRES_PASSWORD` - Пароль пользователя в Postgres.
- `DB_SSL_MODE` - SSL мод Postgres.
- `PRODUCER_URL` - Строка в формате `host:port`. Host и Port для подключения Kafka.
- `CODE_LENGTH` - Длина кода, который будет отправлен в смс-сообщении.
- `CONSUMER_URL` - Строка в формате `host:port`. Host и Port для подключения Kafka.
- `CONSUMER_TOPIC` - Топик Kafka для подписки консьюмера.
- `CONSUMER_GROUP_ID` - Идентификатор группы консьюмеров.
- `CONSUMER_AUTO_OFFSET_RESET` - Конфигурация оффсета, по умолчанию `latest`
- `PROMETHEUS_SERVER_PORT` - Порт сервера метрик.
      
### Запуск

```

git clone https://github.com/denisovdev/go_kafka_sms_sender

cd go_kafka_sms_sender  

make up

```

После запуска сервисы доступны по следующим запросам:
- API - [localhost:8080](http://localhost:8080)
- Prometheus - [localhost:9090](http://localhost:9090)
- Grafana - [localhost:3000](http://localhost:3000) Логин и пароль по умолчанию: `admin` `admin`

### Методы API

##### POST `localhost:8080/api`
Запрос на отправку смс-сообщения. В параметрах необходимо передать номер телефона в международном формате без знака "+"
###### Пример запроса
``` bash
curl -X POST -H "Content-Type: application/json" -d '{"phone": "89999999999"}' http://localhost:8080/api/
```

###### Пример ответа со статусом 200
```
"message successfuly sended"
```

###### Пример ответа со статусом 400
```json
{
	"error": "invalid request body"
}
```
##### GET `localhost:8082/metrics`
Запрос на получение сырых метрик, собранных Prometheus
###### Пример запроса
``` bash
curl --location 'http://localhost:8082/metrics'
```
