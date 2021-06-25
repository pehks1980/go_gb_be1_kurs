# gb_go_be1 beckend 1 final kurs work

model: файловое хранилище

format: json

fields  and json record structure:

 "data": [<br>
{<br>
"uid": “string” - идентификатор пользователя (уникальный номер)<br>
"url": "string", -  исходная ссылка<br>
"shorturl" : "string", -  короткая ссылка генерируемая<br>
"datetime": "2021-05-31T00:15:11.177Z", -  дата создания изменения<br>
"active": 1 -  активная ссылка (0- удалена)<br>
"redirs": 0 – сч-ик переходов по ссылке<br>
},<br>
]<br>


go structure which represents json format:

type Data []struct { <br>
UID string 			`json:"uid"`<br>
URL string 			`json:"url"`<br>
Shorturl string 		`json:"shorturl"`<br>
Datetime time.Time 		`json:"datetime"`<br>
Active int  			`json:"active"`<br>
Redirs int 			`json:"redirs"`<br>
}<br>

Принцип работы файло-хранилища.

Апи работает с структурой которая хранит мапу ключем которой является пара “uid:shorturl”
а значением элемент вышепредставленнной структуры. Сам физический файл на диске создается и апдейтится, при любом изменении этой мапы
(при удалении элемента, меняется флаг active=0 и на диск такие записи не скидываются)
Работа с хранилищем осуществляется через интерфейс:

type linkSvc interface {<br>
Get(uid, key string) (model.DataEl, error)<br>
Put(uid, key string, value model.DataEl) error<br>
Del(uid, key string) error<br>
List(uid string) ([]string, error)<br>
GetUn(shortlink string) (model.DataEl, error)<br>
}<br>

Первые 3 метода реализуют стандартный доступ к мапе типа crud, но, нужно задавать uid для указания пользователя.
Метод list служит для того, чтобы, где необходимо, получить список всех ссылок такого то пользователя.
Метод GetUn – выполняет поиск уникальной короткой ссылки и увеличивает счетчик redir +1. Служит для открывания короткой ссылки и перехода redir по URL.

Документация по API находится в папке swagger – конфиг YAML с описанием всех методов и структур.

Для работы с API был написан клиент на node js. Файлы находятся в папке nodejs/proj1
Запуск его из этой папки командой: node server.js
Заход: http://127.0.0.1:8090/
