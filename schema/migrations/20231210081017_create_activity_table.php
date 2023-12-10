<?php

declare(strict_types=1);

use Phinx\Migration\AbstractMigration;

final class CreateActivityTable extends AbstractMigration
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
        $table = $this->table('activities');
        $table->addColumn('from_account_id', 'integer', ['null' => true, 'signed' => false])
            ->addForeignKey('from_account_id', 'accounts', 'id')
            ->addColumn('to_account_id', 'integer', ['null' => true, 'signed' => false])
            ->addForeignKey('to_account_id', 'accounts', 'id')
            ->addColumn('tag_id', 'integer', ['null' => true, 'signed' => false])
            ->addForeignKey('tag_id', 'tags', 'id')
            ->addColumn('sub_tag_id', 'integer', ['null' => true, 'signed' => false])
            ->addForeignKey('sub_tag_id', 'tags', 'id')
            ->addColumn('amount', 'decimal', ['null' => false, 'precision' => 20, 'scale' => 2])
            ->addColumn('event_date', 'datetime', ['null' => true])
            ->addColumn('remarks', 'string', ['limit' => 255, 'null' => true])
            ->addColumn('transaction_type_id', 'integer', ['null' => false, 'signed' => false])
            ->addForeignKey('transaction_type_id', 'transaction_types', 'id')
            ->addTimestamps()
            ->create();
            if ($this->isMigratingUp()) {
        }
    }
}
