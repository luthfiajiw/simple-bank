-- name: CreateTransfer :one
INSERT INTO transfer (
  from_id_accounts,
  to_id_accounts,
  amount
) VALUES (
  $1, $2, $3
) RETURNING *;

-- name: GetTransfer :one
SELECT * FROM transfer
WHERE id = $1 LIMIT 1;

-- name: ListTransfer :many
SELECT * FROM transfer
WHERE
  from_id_accounts = $1 OR
  to_id_accounts = $2
ORDER BY id
LIMIT $3
OFFSET $4;