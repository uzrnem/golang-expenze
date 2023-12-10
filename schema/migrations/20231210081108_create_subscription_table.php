<?php

declare(strict_types=1);

use Phinx\Migration\AbstractMigration;

final class CreateSubscriptionTable extends AbstractMigration
{
    /**
     * Change Method.
     *
     * Write your reversible migrations using this method.
     *
     * More information on writing migrations is available here:
     * https://book.cakephp.org/phinx/0/en/migrations.html#the-change-method
     *
     * Remember to call "create()" or "update()" and NOT "save()" when working
     * with the Table class.
     */
    public function change(): void
    {
        $table = $this->table('subscriptions');
        $table->addColumn('title', 'string', ['limit' => 255, 'null' => false])
            ->addColumn('description', 'string', ['limit' => 255, 'null' => true])
            ->addColumn('amount', 'float', ['null' => false])
            ->addColumn('actual_amount', 'float', ['null' => true])
            ->addColumn('start_date', 'datetime', ['null' => true])
            ->addColumn('end_date', 'datetime', ['null' => false])
            ->addColumn('due_date', 'datetime', ['null' => true])
            ->addColumn('type', 'enum', ['values' => ['prepaid', 'postpaid'], 'default' => 'prepaid'])
            ->addColumn('status', 'boolean', ['limit' => 1, 'null' => false, 'default' => 0])
            ->addTimestamps()
            ->create();
    }
}
