CREATE TABLE IF NOT EXISTS product
(
    ean  TEXT PRIMARY KEY,
    name TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS nutrition
(
    ean TEXT REFERENCES product (ean) ON DELETE CASCADE ON UPDATE CASCADE,
    kcal INTEGER NOT NULL
);

CREATE TABLE IF NOT EXISTS unit
(
    id    INTEGER PRIMARY KEY,
    value TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS packaging
(
    ean     TEXT REFERENCES product (ean) ON DELETE CASCADE ON UPDATE CASCADE,
    value   REAL    NOT NULL,
    unit_id INTEGER NOT NULL REFERENCES unit (id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS nutrition_quantity
(
    ean     TEXT REFERENCES product (ean) ON DELETE CASCADE ON UPDATE CASCADE,
    value   INTEGER NOT NULL,
    unit_id INTEGER NOT NULL REFERENCES unit (id) ON DELETE CASCADE ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS nutrient_type
(
    id   INTEGER PRIMARY KEY,
    type TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS nutrient
(
    ean     TEXT REFERENCES product (ean) ON DELETE CASCADE ON UPDATE CASCADE,
    type_id INTEGER REFERENCES nutrient_type (id) ON DELETE CASCADE ON UPDATE CASCADE,
    value   REAL    NOT NULL,
    unit_id INTEGER NOT NULL REFERENCES unit (id) ON DELETE CASCADE ON UPDATE CASCADE,
    PRIMARY KEY (ean, type_id)
);


CREATE TABLE IF NOT EXISTS vitamin_type
(
    id   INTEGER PRIMARY KEY,
    type TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS vitamin
(
    ean     TEXT REFERENCES product (ean) ON DELETE CASCADE ON UPDATE CASCADE,
    type_id INTEGER REFERENCES vitamin_type (id) ON DELETE CASCADE ON UPDATE CASCADE,
    value   REAL    NOT NULL,
    unit_id INTEGER NOT NULL REFERENCES unit (id) ON DELETE CASCADE ON UPDATE CASCADE,
    PRIMARY KEY (ean, type_id)
);

CREATE TABLE IF NOT EXISTS mineral_type
(
    id   INTEGER PRIMARY KEY,
    type TEXT NOT NULL UNIQUE
);

CREATE TABLE IF NOT EXISTS mineral
(
    ean     TEXT REFERENCES product (ean) ON DELETE CASCADE ON UPDATE CASCADE,
    type_id INTEGER REFERENCES mineral_type (id) ON DELETE CASCADE ON UPDATE CASCADE,
    value   REAL    NOT NULL,
    unit_id INTEGER NOT NULL REFERENCES unit (id) ON DELETE CASCADE ON UPDATE CASCADE,
    PRIMARY KEY (ean, type_id)
);
