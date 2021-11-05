-- dropping entries and transfers table before dropping the accounts table because there’s a foreign key constraint in entries and transfers that references accounts records.
DROP TABLE IF EXISTS entries;
DROP TABLE IF EXISTS transfers;
DROP TABLE IF EXISTS accounts;
