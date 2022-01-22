# Пример "Employes"

API сервис, для добавления и выводит информацию о сотрудниках и всех его подчиненных, соблюдая иерархию.

Информация о каждом сотруднике храниться в базе данных и содержать следующие данные:
 - 1 Имя;
 - 2 Должность;
 - 3 Дата приема на работу;

У каждого сотрудника может быть только один непосредственный руководитель. Иерархия ограничена 5 уровнями

## Как использовать

База данных и программа создается в docker запуском docker-compose.yml
```
docker-compose up -d
```

По запросу status в API выводится ответ о состоянии работы

Пример:

 ```
 curl -i -X GET http://localhost:8888/status
 ```

Ответ:

 ```
 result {"status":"ok"}
 ```

По запросу employes/{id} в API с указанием id сотрудника

Пример:

```
 curl -i -X GET http://localhost:8888/employes/1
```

Ответ:

```
{"employe":
  {"id":"1",
  "name":"testName",
  "job":"testJob",
  "employe_at":"2020-01-01",
  "employes":[
    {"id":"12",
    "name":"testName2",
    "job":"testJob2",
    "employe_at":"2000-01-01"}]
 }}
 ```

По запросу employes с указанием данных сотрудника, добавляется запись в базу данных. В ответ получаем id нового сотрудника.

Пример добавления сотрудника без руководителя:
```
curl --header "Content-Type: application/json" --request POST --data "{\"name\":\"Poll\", \"job\":\"xyz\", \"employed_at\":\"2021-11-11\"}" http://localhost:8888/employes
```
Ответ:
```
{"id":"1"}
```
Если неободимо задать руководителя, то добовляем поле head_id и id существующего сотрудника.

Пример добавления сотрудника с руководителем:

```
curl --header "Content-Type: application/json" --request POST --data "{\"name\":\"Bill\", \"job\":\"robot\", \"employed_at\":\"2021-11-12\", \"head_id\":\"1\"}" http://localhost:8888/employes
```

Ответ:
```
{"id":"1"}
```

## Технологии
* Go 
* PostgreSQL
* Go kit

