# ab_log_bmx280

##Create db and table

```
 CREATE DATABASE ab_log_db;

CREATE TABLE pressure (
    p_id serial not null primary key,
    p_val float  NOT NULL,
    p_date  timestamp default NULL
);
```
