CREATE TABLE layer.layerinfo(
		id  serial PRIMARY KEY,
		uid uuid NOT NULL UNIQUE,
		schema varchar(255) NOT NULL,
		name varchar(255) NOT NULL,
		attr json NOT NULL,
		layertype smallint NOT NULL default 1,
		description text,
		createat timestamptz,
		updateat timestamptz
)

INSERT INTO layer.layerinfo(
		uid,
		schema,
		name,
		attr,
		layertype
	) VALUES(
		'123e4567-e89b-12d3-a456-426655440000',
		'layer',
		'city',
		'{
			"name":"string"
		}',
		1
)
