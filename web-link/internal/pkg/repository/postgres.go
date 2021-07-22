package repository

import (
	"context"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pehks1980/go_gb_be1_kurs/web-link/internal/pkg/model"
)

type tUserRole string

// Constants for type user rolr pg
const (
	SUPERUSER tUserRole = "SUPERUSER"
	CREATOR   tUserRole = "CREATOR"
	USER      tUserRole = "USER"
)

// User - go struct of pg db
type User struct {
	ID               int       `db:"id"`
	UID              string    `db:"uid"`
	Name             string    `db:"name"`
	Passwd           string    `db:"passwd"`
	Email            string    `db:"email"`
	CreatedOn        time.Time `db:"created_on"`
	LastLogin        time.Time `db:"last_login"`
	IsActive         bool      `db:"is_active"`
	UserRole         tUserRole `db:"user_role"`
	IsBalanceBlocked bool      `db:"is_balance_blocked"`
	Balance          string    `db:"balance"`
}

// UserData - go struct of pg db - related to user data contains all shortlink url counters
type UserData struct {
	ID       int       `db:"id"`
	UserID   int       `db:"user_id"`
	UID      string    `db:"uid"`
	URL      string    `db:"url"`
	ShortURL string    `db:"short_url"`
	DateTime time.Time `db:"date_time"`
	IsActive bool      `db:"is_active"`
	Redirs   int       `db:"redirs"`
}

// UsersTransactions - go struct of pg db - related to transactions b/w users
type UsersTransactions struct {
	ID          int       `db:"id"`
	DateTime    time.Time `db:"date_time"`
	UserIDFrom  int       `db:"user_id_from"`
	UserIDTo    int       `db:"user_id_to"`
	Amount      string    `db:"amount"`
	Description string    `db:"description"`
	Successful  bool      `db:"successful"`
}

// PgRepo init pg go struct holds connex to db
type PgRepo struct {
	URL    string
	CTX    context.Context
	DBPool *pgxpool.Pool
}

// WhoAmI - driver id 1-pg 0-file
func (pgr *PgRepo) WhoAmI() uint64 {
	return 1
}

// CloseConn - close db connex when server quit
func (pgr *PgRepo) CloseConn() {
	pgr.DBPool.Close()
}

// New Init of pg driver
func (pgr *PgRepo) New(filename string) RepoIf {
	ctx := context.Background()
	// Строка для подключения к базе данных
	url := filename //"postgres://postuser:postpassword@192.168.1.204:5432/a4"
	cfg, err := pgxpool.ParseConfig(url)
	if err != nil {
		log.Fatal(err)
	}
	// Pool соединений обязательно ограничивать сверху
	cfg.MaxConns = 8
	cfg.MinConns = 4
	// HealthCheckPeriod - частота пингования соединения с Postgres
	cfg.HealthCheckPeriod = 1 * time.Minute
	// MaxConnLifetime - сколько времени будет жить соединение.
	//можно устанавливать большие значения
	cfg.MaxConnLifetime = 24 * time.Hour
	// MaxConnIdleTime - время жизни неиспользуемого соединения,
	cfg.MaxConnIdleTime = 30 * time.Minute
	// ConnectTimeout устанавливает ограничение по времени
	// на весь процесс установки соединения и аутентификации.
	cfg.ConnConfig.ConnectTimeout = 1 * time.Second
	// Лимиты в net.Dialer позволяют достичь предсказуемого
	// поведения в случае обрыва сети.
	cfg.ConnConfig.DialFunc = (&net.Dialer{
		KeepAlive: cfg.HealthCheckPeriod,
		// Timeout на установку соединения гарантирует,
		// что не будет зависаний при попытке установить соединение.
		Timeout: cfg.ConnConfig.ConnectTimeout,
	}).DialContext
	// pgx предоставляет набор адаптеров для популярных логеров
	//это позволяет организовать сбор ошибок при работе с базой
	//@see: https://github.com/jackc/pgx/tree/master/log
	// cfg.ConnConfig = zerologadapter.NewLogger(logger)
	dbpool, err := pgxpool.ConnectConfig(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}

	return &PgRepo{CTX: ctx,
		URL:    url,
		DBPool: dbpool,
	}
}

// Get - get data string from pg repo
// uid - user uid, key - shortlink
// if uid == suid (SUPERUSER uid) - retreives information despite original uid
func (pgr *PgRepo) Get(uid, key string, su bool) (model.DataEl, error) {

	grGet := func(ctx context.Context, dbpool *pgxpool.Pool, uid, shorturl string, su bool) (UserData, error) {
		const sql = `
	SELECT id, user_id, url, redirs, is_active, short_url, date_time, uid FROM users_data
    	WHERE uid = $1 AND short_url = $2;
	`
		const sqlsu = `
	SELECT id, user_id, url, redirs, is_active, short_url, date_time, uid FROM users_data
    	WHERE short_url = $1;
	`
		var rows pgx.Rows
		var err error

		if su {
			rows, err = dbpool.Query(ctx, sqlsu, shorturl)
		} else {
			rows, err = dbpool.Query(ctx, sql, uid, shorturl)
		}

		var userdata UserData

		if err != nil {
			return UserData{}, fmt.Errorf("failed to query data: %w", err) // Вызов Close нужен, чтобы вернуть соединение в пул
		}
		defer rows.Close()

		for rows.Next() {
			err = rows.Scan(&userdata.ID,
				&userdata.UserID,
				&userdata.URL,
				&userdata.Redirs,
				&userdata.IsActive,
				&userdata.ShortURL,
				&userdata.DateTime,
				&userdata.UID,
			)

			if err != nil {
				return UserData{}, fmt.Errorf("failed to scan row: %w", err)
			}
		}

		if rows.Err() != nil {
			return UserData{}, fmt.Errorf("failed to read response: %w", rows.Err())
		}

		return userdata, nil
	}

	var userdata UserData
	var err error
	suid, _ := pgr.FindSuperUser()
	if suid == uid {
		userdata, err = grGet(pgr.CTX, pgr.DBPool, uid, key, true)
	} else {
		userdata, err = grGet(pgr.CTX, pgr.DBPool, uid, key, false)
	}

	if err != nil {
		return model.DataEl{}, err
	}

	//adjust field Active db - bool , api - int
	var activeInt = 0
	if userdata.IsActive {
		activeInt = 1
	}

	return model.DataEl{UID: userdata.UID,
		URL:      userdata.URL,
		Shorturl: userdata.ShortURL,
		Datetime: userdata.DateTime,
		Active:   activeInt,
		Redirs:   userdata.Redirs}, nil
}

// Put - store data string to pg repo
// uid - user uid, key - shortlink
// if uid == suid (SUPERUSER uid) - updates repo information despite original uid
func (pgr *PgRepo) Put(uid, key string, value model.DataEl, su bool) error {

	grPut := func(ctx context.Context, dbpool *pgxpool.Pool, uid, key string, userdata *UserData) error {
		const sql = `
	INSERT INTO users_data (user_id,url,short_url,redirs,date_time,uid)
    VALUES ((SELECT id FROM users WHERE uid = $1),$2,$3,$4,$5,$1)
        ON CONFLICT ON CONSTRAINT users_data_shorturl_user_id_keys
            DO UPDATE SET url = excluded.url,
                          redirs = excluded.redirs,
                          date_time = excluded.date_time,
                          uid = excluded.uid;
	`
		_, err := dbpool.Exec(ctx, sql,
			uid,
			userdata.URL,
			userdata.ShortURL,
			userdata.Redirs,
			userdata.DateTime,
		)
		if err != nil {
			return fmt.Errorf("failed to add/change userdata: %w", err)
		}

		return nil
	}

	//adjust type bool int
	var isActiveBool = false
	if value.Active == 1 {
		isActiveBool = true
	}

	userdata := UserData{UID: value.UID,
		URL:      value.URL,
		ShortURL: value.Shorturl,
		DateTime: value.Datetime,
		IsActive: isActiveBool,
		Redirs:   value.Redirs,
	}

	err := grPut(pgr.CTX, pgr.DBPool, uid, key, &userdata)
	if err != nil {
		return err
	}
	return nil
}

// Del - delete data entity from pg repo
// uid - user uid, key - shortlink
// if uid == suid (SUPERUSER uid) - updates repo information despite original uid
func (pgr *PgRepo) Del(uid, key string, su bool) error {

	grDel := func(ctx context.Context, dbpool *pgxpool.Pool, uid, shorturl string, su bool) error {
		const sql = `
	DELETE FROM users_data
    	WHERE uid = $1 and short_url = $2;
	`
		const sqlsu = `
	DELETE FROM users_data
    	WHERE short_url = $1;
	`
		var err error
		if su {
			_, err = dbpool.Exec(ctx, sqlsu,
				shorturl,
			)

		} else {
			_, err = dbpool.Exec(ctx, sql,
				uid,
				shorturl,
			)
		}

		if err != nil {
			return fmt.Errorf("failed to del userdata: %w", err)
		}

		return nil
	}

	suid, _ := pgr.FindSuperUser()
	var err error
	if suid == uid {
		err = grDel(pgr.CTX, pgr.DBPool, uid, key, true)
	} else {
		err = grDel(pgr.CTX, pgr.DBPool, uid, key, true)
	}

	if err != nil {
		return err
	}
	return nil
}

// List - list all keys for this user uid
func (pgr *PgRepo) List(uid string) ([]string, error) {

	grList := func(ctx context.Context, dbpool *pgxpool.Pool, uid string) ([]string, error) {
		const sql = `
	SELECT short_url FROM users_data
		WHERE uid = $1;
	`
		rows, err := dbpool.Query(ctx, sql, uid)

		var usersShortURL []string

		if err != nil {
			return nil, fmt.Errorf("failed to query data: %w", err) // Вызов Close нужен, чтобы вернуть соединение в пул
		}
		defer rows.Close()

		for rows.Next() {
			var userShortURL string

			err = rows.Scan(&userShortURL)

			if err != nil {
				return nil, fmt.Errorf("failed to scan row: %w", err)
			}

			usersShortURL = append(usersShortURL, userShortURL)
		}

		if rows.Err() != nil {
			return nil, fmt.Errorf("failed to read response: %w", rows.Err())
		}

		return usersShortURL, nil
	}

	links, err := grList(pgr.CTX, pgr.DBPool, uid)
	if err != nil {
		return nil, err
	}

	return links, nil
}

// TransactionFunc - inline func for a transasction выноска из inTx
type TransactionFunc func(context.Context, pgx.Tx) (string, error)

// inTx - transaction func
func inTx(ctx context.Context, dbpool *pgxpool.Pool, f TransactionFunc) (string, error) {
	//Begin;
	transaction, err := dbpool.Begin(ctx)
	if err != nil {
		return "", err
	}
	// launch function
	URL, err1 := f(ctx, transaction)
	if err1 != nil {
		rbErr := transaction.Rollback(ctx)
		if rbErr != nil {
			log.Print(rbErr)
		}
		return "", err1
	}

	//Commit;
	err = transaction.Commit(ctx)
	if err != nil {
		rbErr := transaction.Rollback(ctx)
		if rbErr != nil {
			log.Print(rbErr)
		}
		return "", err
	}
	return URL, nil
}

// GetUn - find unique shortlink in storage for shortopen api method
// + update redir count (protected by lock)
func (pgr *PgRepo) GetUn(shortlink string) (string, error) {

	grGetUn := func(ctx context.Context, dbpool *pgxpool.Pool, shorturl string) (string, error) {

		URL, err := inTx(ctx, dbpool, func(ctx context.Context, tx pgx.Tx) (string, error) {
			const sql1 = `SELECT url, user_id from users_data
    						WHERE short_url = $1;
			`
			rows, err := tx.Query(ctx, sql1, shorturl)
			if err != nil {
				return "", err
			}
			var URL string
			var userID int

			for rows.Next() {
				err = rows.Scan(&URL, &userID)
				if err != nil {
					return "", err
				}
			}

			const sql2 = `
			UPDATE users_data
    			SET redirs = redirs + 1
        			WHERE user_id = $1 AND short_url = $2;
			`
			_, err = tx.Exec(ctx, sql2, userID, shorturl)
			if err != nil {
				return "", err
			}

			return URL, nil
		})

		if err != nil {
			return "", err
		}

		return URL, nil
	}

	URL, err := grGetUn(pgr.CTX, pgr.DBPool, shortlink)

	if err != nil {
		return "", err
	}

	return URL, nil
}

// additional methods for 'improved' interface
// user crud

// PutUser new user add or update
func (pgr *PgRepo) PutUser(value model.User) (string, error) {

	grAddUser := func(ctx context.Context, dbpool *pgxpool.Pool, user *User) (int, error) {
		const sql = `
	INSERT INTO users (uid, name, passwd, email, user_role, created_on, last_login, balance)
    	VALUES ($1, $2, $3, $4, $5, current_timestamp, current_timestamp, $6::numeric)
		ON CONFLICT ON CONSTRAINT users_uid_key
		DO UPDATE SET passwd = excluded.passwd,
                      last_login = excluded.last_login,
                      balance = excluded.balance
		returning id;
	`
		var id int
		err := dbpool.QueryRow(ctx, sql,
			user.UID,
			user.Name,
			user.Passwd,
			user.Email,
			user.UserRole,
			user.Balance,
		).Scan(&id)
		if err != nil {
			return 0, fmt.Errorf("failed to add user: %w", err)
		}

		return id, nil
	}

	uid := value.Name + value.Email
	passwd := value.Passwd
	var role tUserRole
	switch value.Role {
	case "USER":
		role = USER
	case "CREATOR":
		role = CREATOR
	case "SUPERUSER":
		role = SUPERUSER
	}

	user := User{
		UID:              uid,
		Name:             value.Name,
		Passwd:           passwd,
		Email:            value.Email,
		IsActive:         true,
		UserRole:         role,
		IsBalanceBlocked: false,
		Balance:          "100.00",
	}

	_, err := grAddUser(pgr.CTX, pgr.DBPool, &user)
	if err != nil {
		return "", err
	}
	return uid, nil
}

// GetUser get user and put it in model.User struct
func (pgr *PgRepo) GetUser(uid string) (model.User, error) {

	// get user data
	grGetUser := func(ctx context.Context, dbpool *pgxpool.Pool, uid string) (User, error) {
		const sql = `
	SELECT id, uid, name, passwd, email, is_active, created_on, balance::varchar, last_login, is_balance_blocked, user_role FROM users
    	WHERE uid = $1;
	`
		rows, err := dbpool.Query(ctx, sql, uid)

		var user User

		if err != nil {
			return User{}, fmt.Errorf("failed to query user: %w", err)
		}
		defer rows.Close()

		for rows.Next() {

			err = rows.Scan(&user.ID,
				&user.UID,
				&user.Name,
				&user.Passwd,
				&user.Email,
				&user.IsActive,
				&user.CreatedOn,
				&user.Balance,
				&user.LastLogin,
				&user.IsBalanceBlocked,
				&user.UserRole,
			)

			if err != nil {
				return User{}, fmt.Errorf("failed to scan row: %w", err)
			}

		}

		if rows.Err() != nil {
			return User{}, fmt.Errorf("failed to read response: %w", rows.Err())
		}

		return user, nil
	}

	pguser, err := grGetUser(pgr.CTX, pgr.DBPool, uid)
	if err != nil {
		return model.User{}, err
	}

	var modelrole string
	switch pguser.UserRole {
	case SUPERUSER:
		modelrole = "SUPERUSER"
	case USER:
		modelrole = "USER"
	case CREATOR:
		modelrole = "CREATOR"
	}

	apiuser := model.User{Name: pguser.Name,
		Email:   pguser.Email,
		Role:    modelrole,
		Balance: pguser.Balance,
	}

	return apiuser, nil
}

// DelUser delete user
// name - user name, email - email  = unique combination for user
func (pgr *PgRepo) DelUser(name, email string) error {
	// delete user (and anything related to him)
	grDelUser := func(ctx context.Context, dbpool *pgxpool.Pool, name, email string) error {
		const sql = `
	DELETE FROM users
		WHERE name = $1 AND email = $2;
	`
		_, err := dbpool.Exec(ctx, sql, name, email)
		if err != nil {
			return fmt.Errorf("failed to del user: %w", err)
		}
		// todo when delete user we need to put all its transactions to log for an archive
		return nil
	}

	err := grDelUser(pgr.CTX, pgr.DBPool, name, email)
	if err != nil {
		return err
	}

	return nil
}

// payments

// FindSuperUser - gets suid of superuser
func (pgr *PgRepo) FindSuperUser() (string, error) {
	grFindSU := func(ctx context.Context, dbpool *pgxpool.Pool) (string, error) {
		// find superuser - must be one superuser in this game
		const sql1 = `SELECT uid from users
					WHERE user_role = 'SUPERUSER';
			`
		rows, err1 := dbpool.Query(ctx, sql1)

		if err1 != nil {
			return "", err1
		}
		var suid string

		for rows.Next() {
			err := rows.Scan(&suid)
			if err != nil {
				return "", err
			}
		}
		return suid, nil
	}
	suid, err := grFindSU(pgr.CTX, pgr.DBPool)

	if err != nil {
		return "", err
	}

	return suid, nil
}

// PayUser - pay amount for uidA to uidB as transaction
func (pgr *PgRepo) PayUser(uidA, uidB, amount string) error {
	// pay money transaction b/w users
	grPayUser := func(ctx context.Context, dbpool *pgxpool.Pool, uidA, uidB string, amount string) error {

		const sql = `
		INSERT INTO users_transactions (date_time,  user_id_from,  user_id_to, amount, description, successful )
        	VALUES (current_timestamp,
                (select id from users where uid = $1),
                (select id from users where uid = $2),
                $3::numeric,
                $4,
                FALSE)
		RETURNING id;
	`
		var transID int
		descrText := "Payment +" + amount + " from " + uidA + " for " + uidB
		err := dbpool.QueryRow(ctx, sql, uidA, uidB, amount, descrText).Scan(&transID)
		if err != nil {
			return err
		}

		_, err = inTx(ctx, dbpool, func(ctx context.Context, tx pgx.Tx) (string, error) {

			const sql1 = `SELECT balance::varchar, is_balance_blocked from users
    						WHERE uid = $1;
			`
			rows, err1 := tx.Query(ctx, sql1, uidA)
			if err1 != nil {
				return "", err1
			}
			var prebalanceA float64
			var preisBlockedA bool

			for rows.Next() {
				var balance string
				_ = rows.Scan(&balance, &preisBlockedA)
				prebalanceA, _ = strconv.ParseFloat(balance, 64)
			}

			if prebalanceA < 0 || preisBlockedA {
				err1 = fmt.Errorf("balance uidA is less 0 or blocked")
				return "", err1
			}

			// uidA pays uidB amount
			const sql2 = `
		UPDATE users SET balance = balance + ($1::numeric) where uid = $2;
		`
			_, err1 = tx.Exec(ctx, sql2, amount, uidB)
			if err1 != nil {
				return "", err1
			}
			_, err1 = tx.Exec(ctx, sql2, "-"+amount, uidA)
			if err1 != nil {
				return "", err1
			}

			const sql3 = `
		UPDATE users_transactions
			SET successful = TRUE
				WHERE id = $1;
		`
			_, err1 = tx.Exec(ctx, sql3, transID)
			if err1 != nil {
				return "", err1
			}

			//post update is_balance_blocked of these users
			const sql4 = `SELECT balance::varchar FROM users
    						WHERE uid = $1;
			`
			rows, err1 = tx.Query(ctx, sql4, uidA)
			if err1 != nil {
				return "", err1
			}
			var balanceA float64

			for rows.Next() {
				var balance string
				_ = rows.Scan(&balance)
				balanceA, _ = strconv.ParseFloat(balance, 64)
			}

			const sql5 = `
			UPDATE users
				SET is_balance_blocked = TRUE
					WHERE uid = $1;
			`
			if balanceA < 0 {
				// update user A
				_, err1 = tx.Exec(ctx, sql5, uidA)
				if err1 != nil {
					return "", err1
				}
			}

			rows, err1 = tx.Query(ctx, sql4, uidB)
			if err1 != nil {
				return "", err1
			}
			var balanceB float64

			for rows.Next() {
				var balance string
				_ = rows.Scan(&balance)
				balanceB, _ = strconv.ParseFloat(balance, 64)
			}

			if balanceB < 0 {
				// update user B
				_, err1 = tx.Exec(ctx, sql5, uidB)
				if err1 != nil {
					return "", err1
				}
			}

			return "", nil
		})

		if err != nil {
			return err
		}

		return nil
	}

	err := grPayUser(pgr.CTX, pgr.DBPool, uidA, uidB, amount)
	if err != nil {
		return err
	}
	return nil
}

// GetAll get all data items (with links) from pg db sorted by date
func (pgr *PgRepo) GetAll() (model.Data, error) {

	grGetAll := func(ctx context.Context, dbpool *pgxpool.Pool) ([]UserData, error) {
		const sql = `
	SELECT id, user_id, url, redirs, is_active, short_url, date_time, uid FROM users_data
    	ORDER BY date_time;
	`
		rows, err := dbpool.Query(ctx, sql)

		if err != nil {
			return nil, fmt.Errorf("failed to query data: %w", err) // Вызов Close нужен, чтобы вернуть соединение в пул
		}
		defer rows.Close()

		var usersdata []UserData

		for rows.Next() {
			var userdata UserData

			err = rows.Scan(&userdata.ID,
				&userdata.UserID,
				&userdata.URL,
				&userdata.Redirs,
				&userdata.IsActive,
				&userdata.ShortURL,
				&userdata.DateTime,
				&userdata.UID,
			)

			if err != nil {
				return nil, fmt.Errorf("failed to scan row: %w", err)
			}

			usersdata = append(usersdata, userdata)
		}

		if rows.Err() != nil {
			return nil, fmt.Errorf("failed to read response: %w", rows.Err())
		}

		return usersdata, nil
	}

	usersdata, err := grGetAll(pgr.CTX, pgr.DBPool)
	if err != nil {
		return model.Data{}, err
	}

	//reload pg usersdata to model datael
	var alldata model.Data
	for _, userdata := range usersdata {
		//adjust field Active db - bool , api - int
		var activeInt = 0
		if userdata.IsActive {
			activeInt = 1
		}

		modeldata := model.DataEl{UID: userdata.UID,
			URL:      userdata.URL,
			Shorturl: userdata.ShortURL,
			Datetime: userdata.DateTime,
			Active:   activeInt,
			Redirs:   userdata.Redirs,
		}

		alldata.Data = append(alldata.Data, modeldata)

	}
	return alldata, nil
}