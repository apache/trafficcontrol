
-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
INSERT INTO
   PARAMETER (name, config_file, VALUE) 
   SELECT
      'scope',
      config_file,
      'server' 
   FROM
      PARAMETER 
   WHERE
      name = 'location' 
      AND config_file LIKE 'to_ext_%' 
      OR name = 'location' 
      AND config_file = '12M_facts' 
      OR name = 'location' 
      AND config_file = 'ip_allow.config' 
      OR name = 'location' 
      AND config_file = 'parent.config' 
      OR name = 'location' 
      AND config_file = 'records.config' 
      OR name = 'location' 
      AND config_file = 'remap.config' 
      OR name = 'location' 
      AND config_file = 'hosting.config' 
      OR name = 'location' 
      AND config_file = 'traffic_ops_ort_syncds.cron' 
      OR name = 'location' 
      AND config_file = 'cache.config';
INSERT INTO
   PARAMETER (name, config_file, VALUE) 
   SELECT
      'scope',
      config_file,
      'cdn' 
   FROM
      PARAMETER 
   WHERE
      name = 'location' 
      AND config_file LIKE 'cacheurl%' 
      OR name = 'location' 
      AND config_file LIKE 'hdr_rw_%' 
      OR name = 'location' 
      AND config_file LIKE 'regex_remap_%' 
      OR name = 'location' 
      AND config_file = 'regex_revalidate.config' 
      OR name = 'location' 
      AND config_file LIKE 'set_dscp_%' 
      OR name = 'location' 
      AND config_file = 'ssl_multicert.config' 
      OR name = 'location' 
      AND config_file = 'bg_fetch.config';
INSERT INTO
   PARAMETER (name, config_file, VALUE) 
   SELECT
      'scope',
      config_file,
      'profile' 
   FROM
      PARAMETER 
   WHERE
      name = 'location' 
      AND config_file = '50-ats.rules' 
      OR name = 'location' 
      AND config_file = 'astats.config' 
      OR name = 'location' 
      AND config_file = 'drop_qstring.config' 
      OR name = 'location' 
      AND config_file = 'logs_xml.config' 
      OR name = 'location' 
      AND config_file = 'plugin.config' 
      OR name = 'location' 
      AND config_file = 'storage.config' 
      OR name = 'location' 
      AND config_file = 'sysctl.conf' 
      OR name = 'location' 
      AND config_file LIKE 'url_sig_%' 
      OR name = 'location' 
      AND config_file = 'volume.config';
INSERT INTO
   profile_parameter (profile, PARAMETER) (
   SELECT DISTINCT
      p3.profile, p2.id AS scope_id 
   FROM
      (
         SELECT
            id,
            config_file,
            name 
         FROM
            PARAMETER 
         WHERE
            name = 'location'
      )
      p1 
      INNER JOIN
         (
            SELECT
               id,
               config_file,
               name 
            FROM
               PARAMETER 
            WHERE
               name = 'scope'
         )
         p2 
         ON (p1.config_file = p2.config_file) 
      INNER JOIN
         profile_parameter p3 
         ON (p1.id = p3.PARAMETER));


-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DELETE
FROM
   profile_parameter 
WHERE
   PARAMETER IN 
   (
      SELECT DISTINCT
         p1.id AS scope_id 
      FROM
         (
            SELECT
               id,
               config_file,
               name 
            FROM
               PARAMETER 
            WHERE
               name = 'scope'
         )
         p1 
         INNER JOIN
            profile_parameter p2 
            ON (p1.id = p2.PARAMETER)
   )
;
DELETE
FROM
   PARAMETER 
WHERE
   name = 'scope';
