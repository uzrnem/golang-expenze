<?php


use Phinx\Seed\AbstractSeed;

class TransactionTypeSeeder extends AbstractSeed
{
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
