SELECT
	ds.tenant_id,
	federation.id AS id,
	federation.cname,
	federation.ttl,
	federation.description,
	federation.last_updated,
	ds.id AS ds_id,
	ds.xml_id
FROM federation
JOIN
	federation_deliveryservice AS fd
	ON federation.id = fd.federation
JOIN
	deliveryservice AS ds ON
	ds.id = fd.deliveryservice
JOIN
	cdn AS c
	ON c.id = ds.cdn_id
