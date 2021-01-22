/**
 *
 */
export const enum ProfileType {
	/**
	 *
	 */
	ATS_PROFILE = "ATS_PROFILE",
	/**
	 *
	 */
	TR_PROFILE = "TR_PROFILE",
	/**
	 *
	 */
	TM_PROFILE = "TM_PROFILE",
	/**
	 *
	 */
	TS_PROFILE = "TS_PROFILE",
	/**
	 *
	 */
	TP_PROFILE = "TP_PROFILE",
	/**
	 *
	 */
	INFLUXDB_PROFILE = "INFLUXDB_PROFILE",
	/**
	 *
	 */
	RIAK_PROFILE = "RIAK_PROFILE",
	/**
	 *
	 */
	SPLUNK_PROFILE = "SPLUNK_PROFILE",
	/**
	 *
	 */
	DS_PROFILE = "DS_PROFILE",
	/**
	 *
	 */
	ORG_PROFILE = "ORG_PROFILE",
	/**
	 *
	 */
	KAFKA_PROFILE = "KAFKA_PROFILE",
	/**
	 *
	 */
	LOGSTASH_PROFILE = "LOGSTASH_PROFILE",
	/**
	 *
	 */
	ES_PROFILE = "ES_PROFILE",
	/**
	 *
	 */
	UNK_PROFILE = "UNK_PROFILE",
	/**
	 *
	 */
	GROVE_PROFILE = "GROVE_PROFILE"
}

/**
 *
 */
export interface Parameter {
	configFile: string;
	id?: number;
	lastUpdated?: Date | null;
	name: string;
	profiles: null | Array<string>;
	secure: boolean;
	value: string;
}

/**
 *
 */
export interface Profile {
	id?: number;
	lastUpdated?: Date;
	name: string;
	cdnName: string;
	cdn: number;
	routingDisabled: boolean;
	type: ProfileType;
	params?: Array<Parameter>;
}
