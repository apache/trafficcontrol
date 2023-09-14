INSERT INTO federation (
	cname,
	ttl,
	"description"
) VALUES (
	$1,
	$2,
	$3
)
RETURNING
	id,
	last_updated
