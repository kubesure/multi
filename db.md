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
    endpoint              TEXT,  
    status                TEXT,
    error_msg             TEXT,
    max_response          INTEGER,
    retry_interval        INTEGER,
    retry_count           INTEGER,      
    created_datetime	  TEXT,
    updated_datetime      TEXT,
    UNIQUE(id,batch_id)
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

