CREATE TABLE users (
                       id SERIAL PRIMARY KEY ,
                       username TEXT NOT NULL UNIQUE ,
                       password TEXT NOT NULL ,
                       coins INTEGER NOT NULL
);

CREATE TABLE merch (
                       id SERIAL PRIMARY KEY ,
                       item_type TEXT NOT NULL UNIQUE ,
                       price INTEGER NOT NULL
);


CREATE TABLE purchases_history(
                                  id SERIAL PRIMARY KEY ,
                                  user_id INTEGER REFERENCES users (id),
                                  item_id INTEGER REFERENCES merch (id)
);

CREATE TABLE send_history(
                             id SERIAL PRIMARY KEY ,
                             from_user INTEGER references users (id),
                             to_user INTEGER references users (id) ,
                             amount INTEGER
);

INSERT INTO merch(item_type, price) VALUES ('t-shirt', 80),
                                           ('cup', 20),
                                           ('book', 50),
                                           ('pen', 10),
                                           ('powerbank', 200),
                                           ('hoody', 300),
                                           ('umbrella', 200),
                                           ('socks', 10),
                                           ('wallet', 50),
                                           ('pink-hoody', 500);

INSERT INTO users (username, password, coins)  VALUES ('Nic', '1', 10000), ('Jo', '2', 10000);