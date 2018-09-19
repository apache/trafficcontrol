source /to-access.sh

$RIAK_ADMIN security grant search.admin on schema to admin
$RIAK_ADMIN security grant search.admin on index to admin
$RIAK_ADMIN security grant search.query on index to admin
$RIAK_ADMIN security grant search.query on index sslkeys to admin
$RIAK_ADMIN security grant search.query on index to riakuser
$RIAK_ADMIN security grant search.query on index sslkeys to riakuser
$RIAK_ADMIN security grant riak_core.set_bucket on any to admin
