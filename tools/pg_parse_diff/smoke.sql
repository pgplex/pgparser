select 1;
select * from foo where a = 1 and b is not null;
select a, b from t order by a desc limit 5;
insert into t(a,b) values (1,2);
update t set a = a + 1 where id in (select id from t2);
delete from t where exists (select 1);
