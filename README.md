# gb_go_be1 beckend 1 final kurs work

Описание и принцип работы программы.

Программа использует файловое хранилище в формате json

fields and json record structure:

 "data": [<br>
{<br>
 "uid": “string” - идентификатор пользователя (уникальный ключ)<br>
"url": "string", -  исходная ссылка<br>
"shorturl" : "string", -  короткая ссылка, генерируемая (как уникальный ключ)<br>
"datetime": "2021-05-31T00:15:11.177Z", -  дата создания изменения<br>
"active": 1 -  активная ссылка (0- удалена)<br>
"redirs": 0 – сч-ик переходов по ссылке<br>
},<br>
]<br>


structure which represents json format:

type Data []struct { <br>
UID string 			`json:"uid"`<br>
URL string 			`json:"url"`<br>
Shorturl string 		`json:"shorturl"`<br>
Datetime time.Time 		`json:"datetime"`<br>
Active int  			`json:"active"`<br>
Redirs int 			`json:"redirs"`<br>
}<br>

Апи.<br>
Реализован метод crud для хранения записей,
методы для аутентификации: - выдача пары jwt token,<br>
метод открытия короткой ссылки и редиректа на URL.
Подробнее можно ознакомиться с ним в Свагере.
(в папке swagger) – конфиг YAML с описанием всех апи методов
и структур. Ошибки.<br>

Для работы с программой, первое что нужно сделать - пройти
аутентификацию и получить ключи для авторизации.<br>

curl -X POST http://127.0.0.1:8000/user/auth -H "Content-Type: application/json" -d '{"uid":"any user"}'<br>

в ответ придут пара ключей, <br>
{"accessToken":"eyJhbG...","refreshToken":"sdsdsd....."}

accessToken ключ нужно использовать для доступа к апи,
указывая его всякий раз при обращении к методам апи, например:<br>

curl -X GET http://127.0.0.1:8000/links/all -H "Authorization: Bearer eyJhbGciOiJIUzI1NiI...токен...."<br>

-должен выдать список линков этого пользователя "any user",<br>

но если происходит ошибка авторизации может прийти такое сообщение: <br>

{"errors":[{"code":1,"message":"Token has expired time"}]}




<br>

Принцип работы файло-хранилища.

// FileRepo - структура для файло-стораджа<br>
type FileRepo struct {<br>
sync.RWMutex<br>
fileName string<br>
fileData map[string]model.DataEl<br>
}<br>

Апи работает с этой структурой через интерфейс,
FileRepo хранит мапу fileData, ключем которой является пара “uid:shorturl”
а значением элемент файла-json-структуры (представленной в начале).
Сам физический файл (storage.json) создается на диске и апдейтится, при любом изменении этой мапы
(при удалении элемента, меняется флаг active=0 и на диск такие записи не скидываются).<br>

Работа с хранилищем осуществляется через интерфейс:

type linkSvc interface {<br>
Get(uid, key string) (model.DataEl, error)<br>
Put(uid, key string, value model.DataEl) error<br>
Del(uid, key string) error<br>
List(uid string) ([]string, error)<br>
GetUn(shortlink string) (model.DataEl, error)<br>
}<br>

Первые 3 метода реализуют стандартный доступ типа crud, но, нужно задавать uid для указания конкретного пользователя.
Метод list служит для того, чтобы, где необходимо, получить список всех ссылок конкретного пользователя.
Метод GetUn – выполняет поиск уникальной короткой ссылки во всем хранилище (без указания пользователя)
и увеличивает (в режиме лока структуры) счетчик redir +1.
Этот метод служит для открывания короткой ссылки и перехода redir по URL.

Для работы с API был также написан клиент на node js. Файлы находятся в папке nodejs/proj1
Запуск из этой папки командой: node server.js
Заход: http://127.0.0.1:8090/

![Иллюстрация к проекту](https://github.com/pehks1980/go_gb_be1_kurs/blob/main/image/image.png)

UPDATE: heroku api address: https://web-link19801.herokuapp.com <br>
when starting web-client:<br>

user@MacBook proj1 % node server.js<br>
server node.js started http://127.0.0.1:8090  (API URL: https://web-link19801.herokuapp.com ) <br>

#posgres migrations make: (dir called `migrations`)<br>
-- go get -tags 'postgres' -u github.com/golang-migrate/migrate/cmd/migrate or brew install golang-migrate (mac)<br>
-- migrate create -seq -ext sql -dir migrations init_schema<br>
-- clean db eg 'a5' must exist: psql -h localhost -U postgres -w -c "create database a5;"<br>
# use:
init: migrate -database "postgres://postuser:password@192.168.1.204:5432/a5?sslmode=disable" -path migrations up (put your creds)<br>
rollback: migrate -database "postgres://postuser:password@192.168.1.204:5432/a5?sslmode=disable" -path migrations down<br>
