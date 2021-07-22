BEGIN;
-- проверка на предмет добавления данных пользователя и его ссылки
INSERT INTO users (uid, name, passwd, email, created_on, user_role)
VALUES ('123','Moby','123','',current_timestamp,'SUPERUSER');

INSERT INTO users_data (user_id,url,short_url,redirs)
VALUES ((SELECT id FROM users WHERE name = 'Moby'),'www.mail.ru','asdfrg.dfg1',0);

INSERT INTO users (uid, name, passwd, email, created_on, user_role, balance)
VALUES ('321','Вован','123','',current_timestamp,'USER',99.00);

INSERT INTO users_data (user_id,url,short_url,redirs)
values ((select id from users where name = 'Вован'),'www.mail.ru','asdfrg.dfg2',0);

commit;
