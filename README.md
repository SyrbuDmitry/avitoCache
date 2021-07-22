# Тестовое задание на позицию стажера backend в юнит Geo

Цель задания – разработать приложение имплементацию in-memory [Redis](https://redis.io/) кеша.

Детали реализации:
* Писать код можно на любом языке программирования
* Предоставить инструкцию по запуску приложения. В идеале (но не обязательно) – использовать контейнеризацию с возможностью запустить проект командой `docker-compose up`
* Финальную версию нужно выложить на github.com (просьба не делать форк этого репозитория, дабы не плодить плагиат)

Необходимы функционал:

* Клиент и сервер tcp(telnet)/REST API
* Key-value хранилище строк, списков, словарей
* Возможность установить TTL на каждый ключ. Время устанавливается в секундах
* Реализовать операторы: GET, SET, DEL, KEYS
* Реализовать покрытие несколькими тестами функционала

Дополнительно (необязательно):

* Реализовать операторы: HGET, HSET, LGET, LSET
* Реализовать сохранение на диск
* Масштабирование (на серверной или на клиентское стороне)
* Авторизация
* Нагрузочные тесты

Справка:

Описание Redis-команд можно найти [здесь](https://redis.io/commands)

## Реализованные методы
* GET - GET K
* SET - SET K V
* HGET - HGET K FIELD
* HSET - HSET K FIELD VALUE
* LGET - LGET K IND
* LSET - LSET K IND VAL
* TTL - TTL K
* EXPIRE - EXPIRE K RIME
* LPUSH/RPUSH - LPUSH/RPUSH K VAL
* LRANGE - LRANGE K START_IND END_IND
* DELETE - DELETE KEY
* KEYS - KEYS
* HGETALL - HGETALL

## Пример работы

```bash

$ SET key val ex 600
Ok!

$ GET key
val

$ TTL KEY
600

$ GET key1
No such key!

```

## Установка

```bash

$ git clone https://github.com/SyrbuDmitry/avitoCache.git

Запуск сервера
$ cd ~/avitoCache/server
$ go run redisServer.go

Запуск клиента
$ cd ~/avitoCache/client
$ go run redisClient.go
```


