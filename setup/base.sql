CREATE TABLE users
(
  key bigint DEFAULT 1,
  id serial NOT NULL,
  username character varying,
  password character varying,
  role character varying,
  expire timestamp without time zone DEFAULT now(),
  sessionId bigint NOT NULL DEFAULT 0
)