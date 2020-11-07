CREATE TABLE articles (
    url_key varchar(255) NOT NULL,
    title varchar(255) NOT NULL,
    content text,
    release_at timestamp with time zone,
    created_at timestamp with time zone,
    updated_at timestamp with time zone,
    PRIMARY KEY (url_key)
)