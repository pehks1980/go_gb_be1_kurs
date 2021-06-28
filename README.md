# gb_go_be1 beckend 1 final kurs work

Описание и принцип работы программы.

Программа использует файловое хранилище в формате json

fields and json record structure:

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
Реализован метод crud для хранения записей,<br>
методы для аутентификации: - выдача пары jwt token,<br> 
метод открытия короткой ссылки и редиректа на URL.
Подробнее можно ознакомиться с ним в Свагере.
(в папке swagger) – конфиг YAML с описанием всех апи методов 
и структур. Ошибки.

Принцип работы файло-хранилища.

// FileRepo - структура для файло-стораджа<br>
type FileRepo struct {<br>
sync.RWMutex<br>
fileName string<br>
fileData map[string]model.DataEl<br>
}<br>

Апи работает с этой структурой через интерфейс,
она хранит мапу ключем которой является пара “uid:shorturl”
а значением элемент вышепредставленнной json-структуры.
Сам физический файл (storage.json) создается на диске и апдейтится, при любом изменении этой мапы
(при удалении элемента, меняется флаг active=0 и на диск такие записи не скидываются)
Работа с хранилищем осуществляется через интерфейс:

type linkSvc interface {<br>
Get(uid, key string) (model.DataEl, error)<br>
Put(uid, key string, value model.DataEl) error<br>
Del(uid, key string) error<br>
List(uid string) ([]string, error)<br>
GetUn(shortlink string) (model.DataEl, error)<br>
}<br>

Первые 3 метода реализуют стандартный доступ к мапе (в случае файлового хранилища) типа crud, но, нужно задавать uid для указания пользователя.
Метод list служит для того, чтобы, где необходимо, получить список всех ссылок такого то пользователя.
Метод GetUn – выполняет поиск уникальной короткой ссылки и увеличивает (в режиме лока структуры) счетчик redir +1. Служит для открывания короткой ссылки и перехода redir по URL.

Для работы с API был написан клиент на node js. Файлы находятся в папке nodejs/proj1
Запуск его из этой папки командой: node server.js
Заход: http://127.0.0.1:8090/
