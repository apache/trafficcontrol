DELETE FROM federation
WHERE id = $1
RETURNING
	cname,
	"description",
	id,
	last_updated,
	ttl
