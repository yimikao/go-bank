ALTER TABLE IF EXISTS "accounts" DROP CONSTRAINT IF EXISTS "owner_currency_key";
ALTER TABLE IF EXISTS "accounts" DROP CONSTRAINT IF EXISTS "accounts_owner_fkey";
-- cmd 2 gotten from Table Plus -> Info -> foreign key constraint name: accounts_owner_fkey
DROP TABLE IF EXISTS "users";