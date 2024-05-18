DO
$$
    DECLARE
        john_id    INT;
        richard_id INT;
    BEGIN
        SELECT id INTO john_id FROM "user" WHERE email = 'john.doe@example.com';
        SELECT id INTO richard_id FROM "user" WHERE email = 'rich@example.com';

        DELETE FROM account WHERE user_id IN (john_id, richard_id);

        DELETE FROM "user" WHERE id IN (john_id, richard_id);
    END
$$;