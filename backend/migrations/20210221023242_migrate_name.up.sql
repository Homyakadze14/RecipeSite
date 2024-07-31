CREATE TABLE IF NOT EXISTS users(
    id INT PRIMARY KEY GENERATED BY DEFAULT AS IDENTITY,
    email TEXT NOT NULL UNIQUE,
    login VARCHAR(20) NOT NULL UNIQUE,
    password TEXT NOT NULL,
    about VARCHAR(1500),
    icon_url TEXT NOT NULL,
    created_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS recipes(
    id INT PRIMARY KEY GENERATED BY DEFAULT AS IDENTITY,
    user_id INT references users(id) ON DELETE CASCADE,
    title VARCHAR(50) NOT NULL,
    about TEXT NOT NULL,
    complexitiy int NOT NULL,
    need_time VARCHAR(20) NOT NULL,
    ingridients VARCHAR(1500) NOT NULL,
    photos_urls TEXT NOT NULL,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS comments(
    id INT PRIMARY KEY GENERATED BY DEFAULT AS IDENTITY,
    user_id INT references users(id) ON DELETE CASCADE,
    recipe_id INT references recipes(id) ON DELETE CASCADE,
    text VARCHAR(250) NOT NULL,
    created_at TIMESTAMP,
    updated_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS likes(
    id INT PRIMARY KEY GENERATED BY DEFAULT AS IDENTITY,
    user_id INT references users(id) ON DELETE CASCADE,
    recipe_id INT references recipes(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS bookmarks(
    id INT PRIMARY KEY GENERATED BY DEFAULT AS IDENTITY,
    user_id INT references users(id) ON DELETE CASCADE,
    recipe_id INT references recipes(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS sessions(
    id TEXT,
    user_id INT REFERENCES users(id)
);

CREATE TABLE IF NOT EXISTS subscriptions(
    id INT PRIMARY KEY GENERATED BY DEFAULT AS IDENTITY,
    creator_id INT references users(id) ON DELETE CASCADE,
    subscriber_id INT references users(id) ON DELETE CASCADE
)