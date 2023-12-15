<?php
declare(strict_types=1);

use Phinx\Migration\AbstractMigration;

final class CreateTransactionTypeTable extends AbstractMigration
{
    public function change(): void
    {
        $table = $this->table('transaction_types');
        $table->addColumn('name', 'string', ['limit' => 63, 'null' => false])
            ->addIndex(['name'], ['unique' => true])
            ->create();
    }
}
