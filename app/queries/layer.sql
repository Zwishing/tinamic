CREATE TABLE layers.table_layer(
		id  serial PRIMARY KEY,
		uid uuid NOT NULL UNIQUE,
		schema varchar(255) NOT NULL,
		name varchar(255) NOT NULL,
		attr json NOT NULL,
        geometry_type varchar(255),
        id_column varchar(255),
        geometry_column varchar(255),
        srid int,
        center float[],
        bounds float[],
        min_zoom int,
        max_zoom int,
        tile_url varchar(255),
		description text,
		create_at timestamptz,
		update_at timestamptz
)
