<?php

declare(strict_types=1);

use Phinx\Migration\AbstractMigration;

final class CreateSubscriptionTable extends AbstractMigration
{
    public function change(): void
    {
        $table = $this->table('subscriptions');
        $table->addColumn('title', 'string', ['limit' => 255, 'null' => false])
            ->addColumn('description', 'string', ['limit' => 255, 'null' => true])
            ->addColumn('amount', 'float', ['null' => false])
            ->addColumn('start_date', 'datetime', ['null' => true])
            ->addColumn('end_date', 'datetime', ['null' => false])
            ->addColumn('due_date', 'datetime', ['null' => true])
            ->addColumn('type', 'enum', ['values' => ['prepaid', 'postpaid'], 'default' => 'prepaid'])
            ->addColumn('status', 'boolean', ['limit' => 1, 'null' => false, 'default' => 0])
            ->addTimestamps()
            ->create();
    }
}
