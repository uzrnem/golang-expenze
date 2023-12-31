<?php

declare(strict_types=1);

use Phinx\Migration\AbstractMigration;

final class CreatePassbookTable extends AbstractMigration
{
    public function change(): void
    {
        $table = $this->table('passbooks');
        $table->addColumn('account_id', 'integer', ['null' => false, 'signed' => false])
            ->addForeignKey('account_id', 'accounts', 'id')
            ->addColumn('activity_id', 'integer', ['null' => false, 'signed' => false])
            ->addForeignKey('activity_id', 'activities', 'id')
            ->addColumn('transaction_type_id', 'integer', ['null' => false, 'signed' => false])
            ->addForeignKey('transaction_type_id', 'transaction_types', 'id')
            ->addColumn('previous_balance', 'decimal', ['null' => false, 'precision' => 20, 'scale' => 2])
            ->addColumn('balance', 'decimal', ['null' => false, 'precision' => 20, 'scale' => 2])
            ->addTimestamps()
            ->create();
    }
}