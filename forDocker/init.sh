#!/bin/bash
# name=test_db
# userdb=testuser
# #текущий директорий
# DB=`pwd`
# echo " Имя БД: name = " $name
# echo " Пользователь БД: userdb = " $userdb
# #    Удаление и Создание БД
# dropdb $name -U $userdb
# createdb -U $userdb $name

# echo "  ********  Создание таблицы  *********"
# psql -d $name -f  $DB/employe.psql -U $userdb

set -e

echo " Имя БД: name = " $POSTGRES_USER
echo " Пользователь БД: userdb = " $POSTGRES_DB
psql -v ON_ERROR_STOP=1 -U $POSTGRES_USER -d $POSTGRES_DB <<-EOSQL

    ----------------------------------------------------------------------------------------------------
    --                  Создание таблицы "Сотрудники"
    -----------------------------------------------------------------------------------------------------
    /*---------------------------------------------------------------------------------------------------
    //	TABLE employe                    Таблица 
    //	id                               Идентификатор
    //	name                             Имя
    //	job                              Должность
    //	employed_at                      Дата приема на работу
    //---------------------------------------------------------------------------------------------------
    */
    drop table if exists employe cascade;

    CREATE TABLE employe (
        id serial PRIMARY KEY,
        name character varying NOT NULL,		      
        job character varying NOT NULL,		      
        employed_at date NOT NULL
    );	
        

    ----------------------------------------------------------------------------------------------------
    --                  Создание таблицы "Связь сотрудников"
    -----------------------------------------------------------------------------------------------------
    /*---------------------------------------------------------------------------------------------------
    //	TABLE employe                    Таблица Связь сотрудников
    //	stuff_id                         Идентификатор сотрудника
    //	head_id                          Идентификатор начальника
    //---------------------------------------------------------------------------------------------------
    */
    drop table if exists relation_employes cascade;

    CREATE TABLE relation_employes (
        stuff_id integer NOT NULL,
        head_id integer NOT NULL CHECK (head_id <> stuff_id),
        
        CONSTRAINT stuff_unique UNIQUE (stuff_id),
        
        CONSTRAINT stuff_id_key FOREIGN KEY (stuff_id)
        REFERENCES public.employe (id) MATCH SIMPLE,
        
        CONSTRAINT head_id_key FOREIGN KEY (head_id)
        REFERENCES public.employe (id) MATCH SIMPLE
    );

    create unique index stuff_head_unique on relation_employes (stuff_id, head_id);


EOSQL