-----------------------------------------------------------------------------------------------------
--                  СОЗДАНИЕ ТАБЛИЦЫ "СОТРУДНИКИ"
-----------------------------------------------------------------------------------------------------
/*---------------------------------------------------------------------------------------------------
//    id             1 id 
//    name           2 Имя
//    job            3 Должность
//    employed_at    4 Дата према на работу
//    head_id        5 Ссылка на идентификатор начальника
//---------------------------------------------------------------------------------------------------*/
/*---------------------------------------------------------------------------------------------------

//---------------------------------------------------------------------------------------------------*/
drop table if exists "employe" cascade;

create table "employe"
(
   "id"           serial,
   "name"         character varying,
   "job"          character varying,
   "employed_at"  character varying,
   "head_id"      integer,

   constraint "C_employe_PK" primary key ("id")
);

COMMENT ON TABLE employe IS 'Cотрудники';

COMMENT ON COLUMN employe.id IS 'Номер по порядку';
COMMENT ON COLUMN employe.name IS 'Имя';
COMMENT ON COLUMN employe.job IS 'Должность';
COMMENT ON COLUMN employe.employed_at IS 'Дата приема на работу';
COMMENT ON COLUMN employe.head_id IS 'Ссылка на идентификатор начальника';

------------------------------------------------------------------------------
--              Представление таблицы "СОТРУДНИКИ"
------------------------------------------------------------------------------
create or replace view "v_employe"
(
   "vid",	    		--  1  Идентификатор
   "vname",      		--  2  Имя
   "vjob",    	      --  3  Должность
   "vemployed_at",   --  4  Дата приема на работу
   "vhead_id"   		--  5  Ссылка на идентификатор начальника
)
as SELECT
   "id", 	  		   --  1  Идентификатор			 			      
   "name",	  		   --  2  Имя
   "job",    		   --  3  Должность
   "employed_at",    --  4  Дата приема на работу
   "head_id"   		--  5  Ссылка на идентификатор начальника		 			     

FROM "employe";

------------------------------------------------------------------------------
--              Проверка наличия идентификатора в таблице "СОТРУДНИКИ"
------------------------------------------------------------------------------
create or replace function "p_employe_id_exists"
(
   integer              -- 1 id
)
returns integer as'
declare
      ret_id bigint default -5;
begin
      select into ret_id "id" from "employe" where (id = $1);
      if ret_id is null then 
         return -2;  
      else 
         return ret_id;
      end if;
end;
'language plpgsql security definer; 

------------------------------------------------------------------------------
--              Добавление записи в таблицу "СОТРУДНИКИ"
------------------------------------------------------------------------------
create or replace function "p_employe_add"
(
     character varying,	      --  1  Имя
     character varying,       --  2  Должность
     character varying,       --  3  Дата приема на работу
     integer                  --  4  Ссылка на начальника
)  
returns integer as'
declare
      check_id integer default 1;
      ret_id bigint default -5;
begin
      if $4 is not null then
         select into check_id p_employe_id_exists($4);
         if check_id < 0 then
            return -1;
         end if;
      end if;

	   insert into "employe"
	   (  "name",
         "job",
         "employed_at",
	      "head_id")
	   values
	   (  trim($1),
         trim($2),
         trim($3),
         $4
      ) returning into ret_id "id";
      return ret_id;
end;
'language plpgsql security definer;

create or replace function "p_employe_add"
(
     character varying,	      --  1  Имя
     character varying,       --  2  Должность
     character varying        --  3  Дата приема на работу
)  
returns integer as'
declare
      ret_id bigint default -5;
begin
	   insert into "employe"
	   (  "name",
         "job",
         "employed_at")
	   values
	   (  trim($1),
         trim($2),
         trim($3)
      ) returning into ret_id "id";
      return ret_id;
end;
'language plpgsql security definer;

------------------------------------------------------------------------------
--              Удаление записи из таблицы "ТП"
------------------------------------------------------------------------------
create or replace function "p_employe_del"
(
   integer     -- 1 Идентификатор
)
returns integer as'
declare
      ret_id bigint default -5;
begin
      delete from "employe" where ("id" = $1);
      select into ret_id "id" from "employe" where "id" = $1;
      if ret_id is null then
         return 1;
      else 
         return -3;
      end if;
end;
'language plpgsql security definer;
