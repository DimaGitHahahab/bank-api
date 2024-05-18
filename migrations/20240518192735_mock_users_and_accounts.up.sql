INSERT INTO "user" (name, email, password)
VALUES ('John Doe', 'john.doe@example.com', '<hashed_password>'),
       ('Richard Smith', 'rich@example.com', '<other_hashed_password>');

DO
$$
    DECLARE
        john_id INT;
        richard_id INT;
    BEGIN
        SELECT id INTO john_id FROM "user" WHERE email = 'john.doe@example.com';
        SELECT id INTO richard_id FROM "user" WHERE email = 'rich@example.com';

        INSERT INTO account (user_id, currency_id, amount)
        VALUES (john_id, (SELECT id FROM currency WHERE symbol = 'USD'), 1000.00),
               (john_id, (SELECT id FROM currency WHERE symbol = 'EUR'), 2000.00);

        INSERT INTO account (user_id, currency_id, amount)
        VALUES (richard_id, (SELECT id FROM currency WHERE symbol = 'GBP'), 3000.00),
               (richard_id, (SELECT id FROM currency WHERE symbol = 'JPY'), 4000.00);
    END
$$;