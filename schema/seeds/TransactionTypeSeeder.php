<?php


use Phinx\Seed\AbstractSeed;

class TransactionTypeSeeder extends AbstractSeed
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
        $this->saveIfNotPresent('Expense');
        $this->saveIfNotPresent('Income');
        $this->saveIfNotPresent('Transfer');
    }

    public function saveIfNotPresent(string $name): void {

        $foreignKeys = $this->adapter->fetchRow('SELECT id FROM `transaction_types` where name = "'.$name.'"');
        if (empty($foreignKeys)) {
            $data = [
                [
                    'name'    => $name
                ]
            ];

            $table = $this->table('transaction_types');
            $table->insert($data)
              ->saveData();
        }
    }
}
