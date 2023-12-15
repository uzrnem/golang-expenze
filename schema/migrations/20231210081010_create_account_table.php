<?php

declare(strict_types=1);

use Phinx\Migration\AbstractMigration;

final class CreateAccountTable extends AbstractMigration
{
    public function change(): void
    {
        $table = $this->table('accounts');
        $table->addColumn('name', 'string', ['limit' => 255, 'null' => false])
            ->addColumn('account_type_id', 'integer', ['null' => false, 'signed' => false])
            ->addForeignKey('account_type_id', 'account_types', 'id')
            ->addColumn('amount', 'decimal', ['null' => false, 'precision' => 20, 'scale' => 2])
            ->addColumn('is_closed', 'boolean', ['limit' => 1, 'null' => false, 'default' => 0])
            ->addColumn('is_frequent', 'boolean', ['limit' => 1, 'null' => false, 'default' => 0])
            ->addColumn('is_snapshot_disable', 'boolean', ['limit' => 1, 'null' => false, 'default' => 0])
            ->addIndex(['name'], ['unique' => true])
            ->addTimestamps()
            ->create();
    }
}