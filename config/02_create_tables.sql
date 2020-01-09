-- DROP TABLE frequency_table;
CREATE TABLE frequency_table (
	id serial NOT NULL,
	"name" varchar(200) UNIQUE NOT NULL,
	date_created timestamp NOT NULL,
	last_updated timestamp NULL,
	CONSTRAINT frequency_table_pk PRIMARY KEY (id)
);

ALTER TABLE frequency_table OWNER TO postgres;
GRANT ALL ON TABLE frequency_table TO postgres;

-- DROP TABLE frequency_table_item;
CREATE TABLE frequency_table_item (
	frequency_table_id int4 NOT NULL,
	word varchar(50) NOT NULL,
	times int4 NOT NULL,
	CONSTRAINT frequency_table_item_un UNIQUE (frequency_table_id, word)
);

ALTER TABLE frequency_table_item OWNER TO postgres;
GRANT ALL ON TABLE frequency_table_item TO postgres;
