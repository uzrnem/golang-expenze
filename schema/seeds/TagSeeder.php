<?php


use Phinx\Seed\AbstractSeed;

class TagSeeder extends AbstractSeed
{
    public function run(): void
    {
        $expense = $this->adapter->fetchRow('SELECT id FROM `transaction_types` where name = "Expense"');
        $income = $this->adapter->fetchRow('SELECT id FROM `transaction_types` where name = "Income"');
        $transfer = $this->adapter->fetchRow('SELECT id FROM `transaction_types` where name = "Transfer"');
        $this->saveIfNotPresent('Transfer', 'Credit Card Bill', $transfer['id']);

        $this->saveIfNotPresent('Food', 'Restaurant', $expense['id']);
        $this->saveIfNotPresent('Food', 'Fruits', $expense['id']);
        $this->saveIfNotPresent('Food', 'Online', $expense['id']);

        $this->saveIfNotPresent('House Hold', 'Milk', $expense['id']);
        $this->saveIfNotPresent('House Hold', 'Online', $expense['id']);
        $this->saveIfNotPresent('House Hold', 'Mart', $expense['id']);
        $this->saveIfNotPresent('House Hold', 'Gas Cylinder', $expense['id']);
        $this->saveIfNotPresent('House Hold', 'Veggies', $expense['id']);

        $this->saveIfNotPresent('Earning', 'Cashback', $income['id']);
        $this->saveIfNotPresent('Earning', 'Salary', $income['id']);

        $this->saveIfNotPresent('Purchase', 'Fuel', $expense['id']);
        $this->saveIfNotPresent('Purchase', 'Cloths', $expense['id']);
        $this->saveIfNotPresent('Purchase', 'Medician', $expense['id']);
        $this->saveIfNotPresent('Purchase', 'Offline', $expense['id']);
        $this->saveIfNotPresent('Purchase', 'Online', $expense['id']);

        $this->saveIfNotPresent('Bill', 'Rent', $expense['id']);
        $this->saveIfNotPresent('Bill', 'Service', $expense['id']);
        $this->saveIfNotPresent('Bill', 'Recharge', $expense['id']);
        $this->saveIfNotPresent('Bill', 'Travelling', $expense['id']);
        $this->saveIfNotPresent('Bill', 'Electricity', $expense['id']);
        $this->saveIfNotPresent('Bill', 'Hospital', $expense['id']);
    }

    public function saveIfNotPresent(string $parent, string $child, string $transactionId): void {

        $parentId = $this->adapter->fetchRow('SELECT id FROM `tags` where name = "'.$parent.'"');
        if (empty($parentId)) {
            $parentId = [];
            $table = $this->table('tags');
            $data = [
                [
                    'name'    => $parent,
                    'transaction_type_id' => $transactionId
                ]
            ];
            $var = $table->insert($data)
              ->saveData();
            $parentId["id"] = $this->getAdapter()->getConnection()->lastInsertId();
        }
        $childId = $this->adapter->fetchRow('SELECT id FROM `tags` where name = "'.$child.'"');
        if (empty($childId)) {
            $table = $this->table('tags');
            $data = [
                [
                    'name'    => $child,
                    'tag_id' => $parentId["id"],
                    'transaction_type_id' => $transactionId
                ]
            ];
            $var = $table->insert($data)
              ->saveData();
        }
    }
}
