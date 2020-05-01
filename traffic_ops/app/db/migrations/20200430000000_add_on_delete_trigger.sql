-- SQL in section 'Up' is executed when this migration is applied
CREATE OR REPLACE FUNCTION on_delete_current_timestamp_last_updated() RETURNS trigger
    LANGUAGE plpgsql
    AS
    $$
BEGIN
  update last_deleted set last_updated = now() where tab_name = %(TD['args'][0]);
  RETURN NEW;
END;
$$;

-- +goose Up
-- SQL in section 'Up' is executed when this migration is applied
CREATE TRIGGER on_delete_current_timestamp
AFTER DELETE
ON api_capability
FOR EACH ROW EXECUTE PROCEDURE on_delete_current_timestamp_last_updated('api_capability');

CREATE TRIGGER on_delete_current_timestamp
AFTER DELETE
ON asn
FOR EACH ROW EXECUTE PROCEDURE on_delete_current_timestamp_last_updated('asn');

CREATE TRIGGER on_delete_current_timestamp
AFTER DELETE
ON cachegroup
FOR EACH ROW EXECUTE PROCEDURE on_delete_current_timestamp_last_updated('cachegroup');

CREATE TRIGGER on_delete_current_timestamp
AFTER DELETE
ON cachegroup_fallbacks
FOR EACH ROW EXECUTE PROCEDURE on_delete_current_timestamp_last_updated('cachegroup_fallbacks');

CREATE TRIGGER on_delete_current_timestamp
AFTER DELETE
ON cachegroup_localization_method
FOR EACH ROW EXECUTE PROCEDURE on_delete_current_timestamp_last_updated('cachegroup_localization_method');

CREATE TRIGGER on_delete_current_timestamp
AFTER DELETE
ON cachegroup_parameter
FOR EACH ROW EXECUTE PROCEDURE on_delete_current_timestamp_last_updated('cachegroup_parameter');

CREATE TRIGGER on_delete_current_timestamp
AFTER DELETE
ON capability
FOR EACH ROW EXECUTE PROCEDURE on_delete_current_timestamp_last_updated('capability');

CREATE TRIGGER on_delete_current_timestamp
AFTER DELETE
ON cdn
FOR EACH ROW EXECUTE PROCEDURE on_delete_current_timestamp_last_updated('cdn');

CREATE TRIGGER on_delete_current_timestamp
AFTER DELETE
ON coordinate
FOR EACH ROW EXECUTE PROCEDURE on_delete_current_timestamp_last_updated('coordinate');

CREATE TRIGGER on_delete_current_timestamp
AFTER DELETE
ON deliveryservice
FOR EACH ROW EXECUTE PROCEDURE on_delete_current_timestamp_last_updated('deliveryservice');

CREATE TRIGGER on_delete_current_timestamp
AFTER DELETE
ON deliveryservice_regex
FOR EACH ROW EXECUTE PROCEDURE on_delete_current_timestamp_last_updated('deliveryservice_regex');

CREATE TRIGGER on_delete_current_timestamp
AFTER DELETE
ON deliveryservice_request
FOR EACH ROW EXECUTE PROCEDURE on_delete_current_timestamp_last_updated('deliveryservice_request');

CREATE TRIGGER on_delete_current_timestamp
AFTER DELETE
ON deliveryservice_request_comment
FOR EACH ROW EXECUTE PROCEDURE on_delete_current_timestamp_last_updated('deliveryservice_request_comment');

CREATE TRIGGER on_delete_current_timestamp
AFTER DELETE
ON deliveryservice_server
FOR EACH ROW EXECUTE PROCEDURE on_delete_current_timestamp_last_updated('deliveryservice_server');

CREATE TRIGGER on_delete_current_timestamp
AFTER DELETE
ON deliveryservice_tmuser
FOR EACH ROW EXECUTE PROCEDURE on_delete_current_timestamp_last_updated('deliveryservice_tmuser');

CREATE TRIGGER on_delete_current_timestamp
AFTER DELETE
ON division
FOR EACH ROW EXECUTE PROCEDURE on_delete_current_timestamp_last_updated('division');

CREATE TRIGGER on_delete_current_timestamp
AFTER DELETE
ON federation
FOR EACH ROW EXECUTE PROCEDURE on_delete_current_timestamp_last_updated('federation');

CREATE TRIGGER on_delete_current_timestamp
AFTER DELETE
ON federation_deliveryservice
FOR EACH ROW EXECUTE PROCEDURE on_delete_current_timestamp_last_updated('federation_deliveryservice');

CREATE TRIGGER on_delete_current_timestamp
AFTER DELETE
ON federation_federation_resolver
FOR EACH ROW EXECUTE PROCEDURE on_delete_current_timestamp_last_updated('federation_federation_resolver');

CREATE TRIGGER on_delete_current_timestamp
AFTER DELETE
ON federation_resolver
FOR EACH ROW EXECUTE PROCEDURE on_delete_current_timestamp_last_updated('federation_resolver');

CREATE TRIGGER on_delete_current_timestamp
AFTER DELETE
ON federation_tmuser
FOR EACH ROW EXECUTE PROCEDURE on_delete_current_timestamp_last_updated('federation_tmuser');

CREATE TRIGGER on_delete_current_timestamp
AFTER DELETE
ON hwinfo
FOR EACH ROW EXECUTE PROCEDURE on_delete_current_timestamp_last_updated('hwinfo');

CREATE TRIGGER on_delete_current_timestamp
AFTER DELETE
ON job
FOR EACH ROW EXECUTE PROCEDURE on_delete_current_timestamp_last_updated('job');

CREATE TRIGGER on_delete_current_timestamp
AFTER DELETE
ON job_agent
FOR EACH ROW EXECUTE PROCEDURE on_delete_current_timestamp_last_updated('job_agent');

CREATE TRIGGER on_delete_current_timestamp
AFTER DELETE
ON job_status
FOR EACH ROW EXECUTE PROCEDURE on_delete_current_timestamp_last_updated('job_status');

CREATE TRIGGER on_delete_current_timestamp
AFTER DELETE
ON log
FOR EACH ROW EXECUTE PROCEDURE on_delete_current_timestamp_last_updated('log');

CREATE TRIGGER on_delete_current_timestamp
AFTER DELETE
ON origin
FOR EACH ROW EXECUTE PROCEDURE on_delete_current_timestamp_last_updated('origin');

CREATE TRIGGER on_delete_current_timestamp
AFTER DELETE
ON parameter
FOR EACH ROW EXECUTE PROCEDURE on_delete_current_timestamp_last_updated('parameter');

CREATE TRIGGER on_delete_current_timestamp
AFTER DELETE
ON phys_location
FOR EACH ROW EXECUTE PROCEDURE on_delete_current_timestamp_last_updated('phys_location');

CREATE TRIGGER on_delete_current_timestamp
AFTER DELETE
ON profile
FOR EACH ROW EXECUTE PROCEDURE on_delete_current_timestamp_last_updated('profile');

CREATE TRIGGER on_delete_current_timestamp
AFTER DELETE
ON profile_parameter
FOR EACH ROW EXECUTE PROCEDURE on_delete_current_timestamp_last_updated('profile_parameter');

CREATE TRIGGER on_delete_current_timestamp
AFTER DELETE
ON regex
FOR EACH ROW EXECUTE PROCEDURE on_delete_current_timestamp_last_updated('regex');

CREATE TRIGGER on_delete_current_timestamp
AFTER DELETE
ON region
FOR EACH ROW EXECUTE PROCEDURE on_delete_current_timestamp_last_updated('region');

CREATE TRIGGER on_delete_current_timestamp
AFTER DELETE
ON role
FOR EACH ROW EXECUTE PROCEDURE on_delete_current_timestamp_last_updated('role');

CREATE TRIGGER on_delete_current_timestamp
AFTER DELETE
ON role_capability
FOR EACH ROW EXECUTE PROCEDURE on_delete_current_timestamp_last_updated('role_capability');

CREATE TRIGGER on_delete_current_timestamp
AFTER DELETE
ON server
FOR EACH ROW EXECUTE PROCEDURE on_delete_current_timestamp_last_updated('server');

CREATE TRIGGER on_delete_current_timestamp
AFTER DELETE
ON servercheck
FOR EACH ROW EXECUTE PROCEDURE on_delete_current_timestamp_last_updated('servercheck');

CREATE TRIGGER on_delete_current_timestamp
AFTER DELETE
ON snapshot
FOR EACH ROW EXECUTE PROCEDURE on_delete_current_timestamp_last_updated('snapshot');

CREATE TRIGGER on_delete_current_timestamp
AFTER DELETE
ON staticdnsentry
FOR EACH ROW EXECUTE PROCEDURE on_delete_current_timestamp_last_updated('staticdnsentry');

CREATE TRIGGER on_delete_current_timestamp
AFTER DELETE
ON stats_summary
FOR EACH ROW EXECUTE PROCEDURE on_delete_current_timestamp_last_updated('stats_summary');

CREATE TRIGGER on_delete_current_timestamp
AFTER DELETE
ON status
FOR EACH ROW EXECUTE PROCEDURE on_delete_current_timestamp_last_updated('status');

CREATE TRIGGER on_delete_current_timestamp
AFTER DELETE
ON steering_target
FOR EACH ROW EXECUTE PROCEDURE on_delete_current_timestamp_last_updated('steering_target');

CREATE TRIGGER on_delete_current_timestamp
AFTER DELETE
ON tenant
FOR EACH ROW EXECUTE PROCEDURE on_delete_current_timestamp_last_updated('tenant');

CREATE TRIGGER on_delete_current_timestamp
AFTER DELETE
ON tm_user
FOR EACH ROW EXECUTE PROCEDURE on_delete_current_timestamp_last_updated('tm_user');

CREATE TRIGGER on_delete_current_timestamp
AFTER DELETE
ON to_extension
FOR EACH ROW EXECUTE PROCEDURE on_delete_current_timestamp_last_updated('to_extension');

CREATE TRIGGER on_delete_current_timestamp
AFTER DELETE
ON type
FOR EACH ROW EXECUTE PROCEDURE on_delete_current_timestamp_last_updated('type');

CREATE TRIGGER on_delete_current_timestamp
AFTER DELETE
ON user_role
FOR EACH ROW EXECUTE PROCEDURE on_delete_current_timestamp_last_updated('user_role');

CREATE TRIGGER on_delete_current_timestamp
AFTER DELETE
ON server_capability
FOR EACH ROW EXECUTE PROCEDURE on_delete_current_timestamp_last_updated('server_capability');

CREATE TRIGGER on_delete_current_timestamp
AFTER DELETE
ON server_server_capability
FOR EACH ROW EXECUTE PROCEDURE on_delete_current_timestamp_last_updated('server_server_capability');


CREATE TRIGGER on_delete_current_timestamp
AFTER DELETE
ON deliveryservices_required_capability
FOR EACH ROW EXECUTE PROCEDURE on_delete_current_timestamp_last_updated('deliveryservices_required_capability');

-- +goose Down
-- SQL section 'Down' is executed when this migration is rolled back
DROP TRIGGER IF EXISTS on_delete_current_timestamp on api_capability;
DROP TRIGGER IF EXISTS on_delete_current_timestamp on asn;
DROP TRIGGER IF EXISTS on_delete_current_timestamp on cachegroup;
DROP TRIGGER IF EXISTS on_delete_current_timestamp on cachegroup_fallbacks;
DROP TRIGGER IF EXISTS on_delete_current_timestamp on cachegroup_localization_method;
DROP TRIGGER IF EXISTS on_delete_current_timestamp on cachegroup_parameter;
DROP TRIGGER IF EXISTS on_delete_current_timestamp on capability;
DROP TRIGGER IF EXISTS on_delete_current_timestamp on cdn;
DROP TRIGGER IF EXISTS on_delete_current_timestamp on coordinate;
DROP TRIGGER IF EXISTS on_delete_current_timestamp on deliveryservice;
DROP TRIGGER IF EXISTS on_delete_current_timestamp on deliveryservice_regex;
DROP TRIGGER IF EXISTS on_delete_current_timestamp on deliveryservice_request;
DROP TRIGGER IF EXISTS on_delete_current_timestamp on deliveryservice_request_comment;
DROP TRIGGER IF EXISTS on_delete_current_timestamp on deliveryservice_server;
DROP TRIGGER IF EXISTS on_delete_current_timestamp on deliveryservice_tmuser;
DROP TRIGGER IF EXISTS on_delete_current_timestamp on division;
DROP TRIGGER IF EXISTS on_delete_current_timestamp on federation;
DROP TRIGGER IF EXISTS on_delete_current_timestamp on federation_deliveryservice;
DROP TRIGGER IF EXISTS on_delete_current_timestamp on federation_federation_resolver;
DROP TRIGGER IF EXISTS on_delete_current_timestamp on federation_resolver;
DROP TRIGGER IF EXISTS on_delete_current_timestamp on federation_tmuser;
DROP TRIGGER IF EXISTS on_delete_current_timestamp on hwinfo;
DROP TRIGGER IF EXISTS on_delete_current_timestamp on job;
DROP TRIGGER IF EXISTS on_delete_current_timestamp on job_agent;
DROP TRIGGER IF EXISTS on_delete_current_timestamp on job_status;
DROP TRIGGER IF EXISTS on_delete_current_timestamp on log;
DROP TRIGGER IF EXISTS on_delete_current_timestamp on origin;
DROP TRIGGER IF EXISTS on_delete_current_timestamp on parameter;
DROP TRIGGER IF EXISTS on_delete_current_timestamp on phys_location;
DROP TRIGGER IF EXISTS on_delete_current_timestamp on profile;
DROP TRIGGER IF EXISTS on_delete_current_timestamp on profile_parameter;
DROP TRIGGER IF EXISTS on_delete_current_timestamp on regex;
DROP TRIGGER IF EXISTS on_delete_current_timestamp on region;
DROP TRIGGER IF EXISTS on_delete_current_timestamp on role;
DROP TRIGGER IF EXISTS on_delete_current_timestamp on role_capability;
DROP TRIGGER IF EXISTS on_delete_current_timestamp on server;
DROP TRIGGER IF EXISTS on_delete_current_timestamp on servercheck;
DROP TRIGGER IF EXISTS on_delete_current_timestamp on snapshot;
DROP TRIGGER IF EXISTS on_delete_current_timestamp on staticdnsentry;
DROP TRIGGER IF EXISTS on_delete_current_timestamp on stats_summary;
DROP TRIGGER IF EXISTS on_delete_current_timestamp on status;
DROP TRIGGER IF EXISTS on_delete_current_timestamp on steering_target;
DROP TRIGGER IF EXISTS on_delete_current_timestamp on tenant;
DROP TRIGGER IF EXISTS on_delete_current_timestamp on tm_user;
DROP TRIGGER IF EXISTS on_delete_current_timestamp on to_extension;
DROP TRIGGER IF EXISTS on_delete_current_timestamp on type;
DROP TRIGGER IF EXISTS on_delete_current_timestamp on user_role;
DROP TRIGGER IF EXISTS on_delete_current_timestamp on server_capability;
DROP TRIGGER IF EXISTS on_delete_current_timestamp on server_server_capability;
DROP TRIGGER IF EXISTS on_delete_current_timestamp on deliveryservices_required_capability;