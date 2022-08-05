CREATE TABLE batch (
	id	                TEXT PRIMARY KEY,
	type    	        INTEGER,
	created_datetime    TEXT,
    updated_datetime    TEXT
);

create TABLE job (
    id                    INTEGER,
    batch_id              TEXT,
    payload               TEXT,
    result                TEXT,
    endpoint              INTEGER,  
    status                TEXT,
    error_msg             TEXT,
    max_response          INTEGER,
    retry_interval        INTEGER,
    retry_count           INTEGER,      
    created_datetime	  TEXT,
    updated_datetime      TEXT,
    UNIQUE(id,batch_id)
);

create table endpoint (
    batch_id TEXT,
    uri TEXT,
    method TEXT,
    auth_type TEXT,
    auth_srvcert TEXT,
    auth_uname TEXT,
    auth_pass TEXT,
    headers TEXT,
    UNIQUE(batch_id)
);

create table schedule (
    batch_id        TEXT,
    start_datetime      TEXT,
    end_datetime        TEXT,
    scheduled_count INTEGER
);



/*create TABLE batch_type (
    batch_id    TEXT,
    type        TEXT,
    datetime	TEXT
)*/

