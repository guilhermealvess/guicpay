-- name: InsertNewAccount :exec
INSERT INTO accounts (id,account_type,customer_name,document_number,email,password_encoded,salt_hash_password,phone_number,status,created_at,updated_at) 
VALUES (?,?,?,?,?,?,?,?,?,?,?);

-- name: InsertNewTransaction :exec
INSERT INTO transactions (id,correlated_id,account_id,transaction_type,timestamp,amount)
VALUES (?,?,?,?,?,?);

-- name: FindAccountByID :one
SELECT sqlc.embed(ac), json_group_array(json_object('id', tr.id, 'amount', tr.amount, 'account_id', tr.account_id, 'correlated_id', tr.correlated_id, 'timestamp', tr."timestamp", 'transaction_type', tr.transaction_type)) as transactions
FROM accounts ac
LEFT JOIN transactions tr ON ac.id = tr.account_id
WHERE ac.id = ? GROUP BY ac.id;
