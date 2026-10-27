CREATE TABLE accounts
(
    account_id      SERIAL PRIMARY KEY,
    document_number VARCHAR(100) NOT NULL
);

CREATE TABLE operation_types
(
    operation_type_id SERIAL PRIMARY KEY,
    description       VARCHAR(100) NOT NULL,
    direction VARCHAR(10) NOT NULL
);

CREATE TABLE transactions
(
    transaction_id    SERIAL PRIMARY KEY,
    account_id        INTEGER   NOT NULL,
    operation_type_id INTEGER   NOT NULL,
    amount            INTEGER   NOT NULL,
    event_date        TIMESTAMP NOT NULL,

    CONSTRAINT fk_accounts FOREIGN KEY (account_id) REFERENCES accounts (account_id),
    CONSTRAINT fk_operation_types FOREIGN KEY (operation_type_id) REFERENCES operation_types (operation_type_id)
);



insert into accounts (document_number)values ('123456');
INSERT INTO operation_types
(operation_type_id, description, direction)
VALUES
    (1, 'Normal Purchase', 'DEBIT'),
    (2, 'Purchase with installments', 'DEBIT'),
    (3, 'Withdrawal', 'DEBIT'),
    (4, 'Credit Voucher', 'CREDIT');


insert into transactions (account_id, operation_type_id, amount, event_date)
values (1, 1, -5000, NOW());
insert into transactions (account_id, operation_type_id, amount, event_date)
values (1, 1, -2350, NOW());
insert into transactions (account_id, operation_type_id, amount, event_date)
values (1, 1, -1870, NOW());
insert into transactions (account_id, operation_type_id, amount, event_date)
values (1, 4, 6000, NOW());


