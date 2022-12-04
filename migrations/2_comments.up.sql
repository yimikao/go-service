CREATE TABLE IF NOT EXISTS comments (
    id SERIAL PRIMARY KEY,
    comment VARCHAR NOT NULL,
    comment_date DATE DEFAULT CURRENT_DATE,
    user_id BIGINT REFERENCES users(id)
);

INSERT INTO users(name, id) VALUES('dev_test_user', 1);
INSERT INTO comments (comment, user_id) VALUES('first test comment', 1);
