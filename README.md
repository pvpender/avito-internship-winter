# avito-internship-winter

## Требования
1. Git
2. Docker версии 24 или выше
3. Веб-браузер или консольная утилита для создания http запросов (curl, например)
4. GNU Make версии 4 (опционально)

## Установка
1. Клонировать репозиторий
```bash
git clone https://github.com/pvpender/avito-internship-winter
```
2. Перейти в папку
```bash
cd avito-internship-winter
```
3. Запустить

Make
```bash
make up
```
docker
```bash
docker compose up --build -d
```

4. Выключить

Make
```bash
make down
```
docker
```bash
docker compose down
```


## Тестирование

Можно использовать swagger, который доступен по адресу `http://localhost:8080/swagger/index.html`

**Важно**

Так как версия сваггера 2.0, то после авторизации, её нужно добавлять в окно авторизации как 
`Bearer Токен`, а не просто токен. Пример того, как в заголовках передаётся токен в `curl`

```bash
curl -X 'GET' \
  'http://localhost:8080/api/info' \
  -H 'accept: application/json' \
  -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxfQ.jYyRJbb0WImFoUUdcslQQfwnXTHJzne-6tsPd8Hrw0I'
```

**Важно**

Не забывайте указывать `Content-Type: application/json` в заголовке, иначе CORS не пропустит!

Также можно использовать `test.http` как подсказку по запросам.

## Покрытие тестами

Доступно в файлах `cover.html` и `cover.out`

`config_test.yaml` нужен только для интеграционных тестов, которые запускаются с флагом `integration`

## Нагрузочное тестирование

Проведено с помощью [vegeta](https://github.com/tsenart/vegeta). Результат в файле plot.html и на скрине

![image](https://github.com/user-attachments/assets/42f5f7cd-f83f-4f19-b3e5-92426d0a749a)


Для RPS 1k условие задания выполняется, с увеличением его в 2 раза увеличивается задержка ответа.

## Пояснения

* Так, как swagger остался неизменным, то и ответы эндпоинтов не менялись, в связи с чем
не было предусмотрено пагинации для `GET /api/info`, хотя по-хорошему она нужна
* Автор знает что `config.yaml` и подобные не нужно коммитить в репозиторий, а добавлять в `.gitignore`, тут он закомичен для простоты запуска, в других проектах я так не делаю, правда
