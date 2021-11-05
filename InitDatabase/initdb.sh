#!/bin/bash
name=test_db
userdb=testuser
#текущий директорий
DB=`pwd`
echo " Имя БД: name = " $name
echo " Пользователь БД: userdb = " $userdb
#    Удаление и Создание БД
dropdb $name -U $userdb
createdb -U $userdb $name

echo "  ********  Создание таблицы  *********"
psql -d $name -f  $DB/employe.psql -U $userdb

