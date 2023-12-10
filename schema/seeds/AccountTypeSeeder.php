<?php


use Phinx\Seed\AbstractSeed;

class AccountTypeSeeder extends AbstractSeed
{
    /**
     * Run Method.
     *
     * Write your database seeder using this method.
     *
     * More information on writing seeders is available here:
     * https://book.cakephp.org/phinx/0/en/seeding.html
     */
    public function run(): void
    {
        $this->saveIfNotPresent('Credit');
        $this->saveIfNotPresent('Deposit');
        $this->saveIfNotPresent('Loan');
        $this->saveIfNotPresent('Mutual Funds');
        $this->saveIfNotPresent('Saving');
        $this->saveIfNotPresent('Stocks Equity');
        $this->saveIfNotPresent('Wallet');
    }

    /**
     * Saves the given name if it is not already present in the `account_types` table.
     *
     * @param string $name The name to save.
     * @return void
     */
    public function saveIfNotPresent(string $name): void {

        $foreignKeys = $this->adapter->fetchRow('SELECT id FROM `account_types` where name = "'.$name.'"');
        if (empty($foreignKeys)) {
            $data = [
                [
                    'name'    => $name
                ]
            ];

            $table = $this->table('account_types');
            $table->insert($data)
              ->saveData();
        }
    }
}
