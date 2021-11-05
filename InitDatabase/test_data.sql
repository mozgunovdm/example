select "p_employe_add" ('Lutor', 'Robot', '1991-01-01', null);
select "p_employe_add" ('Vagner', 'Jobless', '2012-03-15', null);
select "p_employe_add" ('Petrov', 'Trainee', '2020-09-02', 1);
select "p_employe_add" ('Ivanov', 'Trainee', '2020-08-02', 1);
select "p_employe_add" ('Semenov', 'Engineer', '2018-04-05', 1);
select "p_employe_add" ('Fomin', 'Engineer', '2018-04-05', 5);
select "p_employe_add" ('Kolov', 'Trainee', '2020-05-05', 5);
select "p_employe_add" ('Frolov', 'Trainee', '2020-04-05', 5);
select "p_employe_add" ('Holohol', 'Manager', '2019-05-05', 6);
select "p_employe_add" ('Fotov', 'Tester', '2019-02-02', 6);
select "p_employe_add" ('Gotov', 'Tester', '2000-02-02', 6);

INSERT INTO employe
    (id, name, job, employed_at)
SELECT
    i, 'name' || i, 'job' || i, 1980 + i % 40 || '-01-01'
FROM
    generate_series(1, 50001) AS s(i);
	
update employe set head_id = (id - id % 100) where (id % 10 = 0) and (id % 100 <> 0);
update employe set head_id = (id - id % 10) where id % 10 = 1 or id % 10 = 2;	
update employe set head_id = (id - id % 10 + 1) where id % 10 = 3 or id % 10 = 4;	
update employe set head_id = (id - id % 10 + 2) where id % 10 = 5 or id % 10 = 6;
update employe set head_id = (id - id % 10 + 5) where id % 10 = 7 or id % 10 = 8;
update employe set head_id = (id - id % 10 + 8) where id % 10 = 9;
ALTER SEQUENCE employe_id_seq RESTART WITH 50002