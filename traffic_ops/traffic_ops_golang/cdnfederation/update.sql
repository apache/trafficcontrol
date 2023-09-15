UPDATE federation SET
	cname = $1,
	ttl = $2,
	"description" = $3
WHERE
  id = $4
RETURNING
	last_updated
