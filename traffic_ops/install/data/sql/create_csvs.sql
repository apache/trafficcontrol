SELECT name, description
FROM profile
where name not like '4'
INTO OUTFILE '/opt/traffic_ops/install/data/csv/profile.csv'
FIELDS TERMINATED BY ','
ENCLOSED BY '"'
LINES TERMINATED BY '\n';


select profile.name, parameter.name, parameter.config_file, parameter.value 
from profile, parameter, profile_parameter
where profile.id = profile_parameter.profile
and parameter.id = profile_parameter.parameter
and profile.name not like '%4%'
and parameter.config_file not like '%url_sig%'
INTO OUTFILE '/opt/traffic_ops/install/data/csv/profile_parameter.csv'
FIELDS TERMINATED BY ','
ENCLOSED BY '"'
LINES TERMINATED BY '\n';


select name, config_file, value from parameter 
where id in (select pp.parameter from profile_parameter pp, profile p where p.name not like '%4%')
and config_file not like '%url_sig%'
INTO OUTFILE '/opt/traffic_ops/install/data/csv/parameter.csv'
FIELDS TERMINATED BY ','
ENCLOSED BY '"'
LINES TERMINATED BY '\n';

--need to replace \" with \"" for go