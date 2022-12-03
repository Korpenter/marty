package storage

const (
	createUsers = `CREATE TABLE IF NOT EXISTS users (
                login varchar(32) PRIMARY KEY,
    			password varchar(256),
    			balance numeric(12,4) DEFAULT 0
                )`
	createOrders = `CREATE TABLE IF NOT EXISTS orders (
    			id text PRIMARY KEY,
    			user_login varchar(32) NOT NULL,
        		uploaded_at timestamp DEFAULT CURRENT_TIMESTAMP,
    			accrual numeric(10,2) DEFAULT 0,
    			status varchar(10) DEFAULT 'NEW' CHECK (status in ('NEW', 'REGISTERED', 'PROCESSING', 'PROCESSED', 'INVALID')),
    			FOREIGN KEY(user_login) REFERENCES users(login)
                )`
	createWithdrawals = `CREATE TABLE IF NOT EXISTS withdrawals (
    			id SERIAL PRIMARY KEY,
    			order_id text,
				accrual numeric(10,2),
    			user_login varchar(32),
				processed_at timestamp DEFAULT CURRENT_TIMESTAMP,
    			FOREIGN KEY(user_login) REFERENCES users(login)
				)`
	createUser = `
	INSERT INTO users (login, password)
	VALUES ($1, $2)
	ON CONFLICT DO NOTHING
	RETURNING login`
	getHashByLogin = `SELECT password FROM users WHERE login = $1`
	addOrder       = `INSERT INTO orders (id, user_login)
				      VALUES($1, $2)
					  ON CONFLICT DO NOTHING`
	getOrderUserid  = `SELECT user_login FROM orders WHERE id = $1`
	getOrdersByUser = `SELECT id, status, accrual, uploaded_at FROM orders WHERE user_login=$1`
	getUserBalance  = `SELECT balance, sum(accrual)
			     	   FROM users u LEFT JOIN withdrawals w ON u.login = w.user_login
			     	   where u.login = $1 group by login, balance`
	getUserWithdrawals   = `SELECT order_id, accrual, processed_at FROM withdrawals WHERE user_login=$1`
	updateOrder          = `UPDATE orders SET status=$1 WHERE id=$2`
	updateProcessedOrder = `UPDATE orders SET status=$1, accrual=$2 WHERE id=$3 RETURNING user_login`
	updateBalance        = `UPDATE users SET balance=balance+$1 WHERE login = $2`
	userVerifyBalance    = `UPDATE users SET balance=balance-$1 WHERE balance>$1 AND login=$2`
	userWithdraw         = `INSERT INTO withdrawals (order_id, accrual, user_login)
						 VALUES ($1, $2, $3)`
	dropTables = `DROP TABLE withdrawals
				  DROP TABLE orders
				  DROP TABLE users`
)
