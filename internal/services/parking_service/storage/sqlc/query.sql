-- name: GetParkingByUUID :one
select * from cars where uuid = $1 LIMIT 1;