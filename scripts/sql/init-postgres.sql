CREATE SCHEMA IF NOT EXISTS mockserver;

CREATE TABLE IF NOT EXISTS mockserver.rules
(
    "key"                 varchar(255) NOT NULL,
    "group"               varchar(255) NOT NULL,
    name                  varchar(255) NOT NULL,
    path                  varchar(255) NOT NULL,
    strategy              varchar(255) NOT NULL,
    method                varchar(45)  NOT NULL,
    status                varchar(255) NOT NULL,
    pattern               varchar(255) NOT NULL,
    next_response_index   int          NOT NULL DEFAULT 0,
    PRIMARY KEY ("key"),
    UNIQUE ("key")
);

CREATE TABLE IF NOT EXISTS mockserver.responses
(
    id           bigserial    NOT NULL,
    body         text         NOT NULL,
    content_type varchar(255) NOT NULL,
    http_status  int          NOT NULL,
    delay        int          DEFAULT 0,
    scene        varchar(255) DEFAULT NULL,
    rule_key     varchar(255) NOT NULL,
    description  varchar(255) DEFAULT NULL,
    PRIMARY KEY (id),
    CONSTRAINT fk_rules FOREIGN KEY (rule_key) REFERENCES mockserver.rules ("key")
);

CREATE INDEX IF NOT EXISTS idx_responses_rule_key ON mockserver.responses (rule_key);

CREATE TABLE IF NOT EXISTS mockserver.variables
(
    id         bigserial    NOT NULL,
    type       varchar(255) NOT NULL,
    name       varchar(255) NOT NULL,
    "key"      varchar(255) NOT NULL,
    rule_key   varchar(255) NOT NULL,
    assertions json         DEFAULT NULL,
    PRIMARY KEY (id),
    CONSTRAINT fk_rule FOREIGN KEY (rule_key) REFERENCES mockserver.rules ("key")
);

CREATE INDEX IF NOT EXISTS idx_variables_rule_key ON mockserver.variables (rule_key);
