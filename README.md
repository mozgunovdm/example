Представлен API сервис, который выводит информацию про руководителя и всех его сотрудников, соблюдая иерархию в древовидной форме.

Информация о каждом сотруднике храниться в базе данных и содержать следующие данные:
  1 Имя;
  2 Должность;
  3 Дата приема на работу;
У каждого сотрудника может быть не более 1 начальника. Иерархия ограничена 5 уровнями;

База данных и программа создается в docker запуском docker-compose.yml
#docker-compose up -d

#все примеры запросов через curl даны для windows
По запросу status в API выводится ответ о состоянии работы 
#curl -i -X GET http://localhost:8888/status
#result {"status":"ok"}

По запросу employes/{id} в API с указанием id сотрудника будет выводиться ответ
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

По запросу employes в API с указанием данных сотрудника, добавляется запись в базу данных. В ответ получаем id нового сотрудника:
#пример сотрудника без руководителя
#curl --header "Content-Type: application/json" --request POST --data "{\"name\":\"Poll\", \"job\":\"xyz\", \"employed_at\":\"2021-11-11\"}" http://localhost:8888/employes 
#ответ {"id":"1"}

Если неободимо задать руководителя сотруднику, то добовляем поле head_id:
#curl --header "Content-Type: application/json" --request POST --data "{\"name\":\"Bill\", \"job\":\"robot\", \"employed_at\":\"2021-11-12\", \"head_id\":\"1\"}" http://localhost:8888/employes


