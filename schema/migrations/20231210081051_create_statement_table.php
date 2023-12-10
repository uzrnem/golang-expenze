<?php

declare(strict_types=1);

use Phinx\Migration\AbstractMigration;

final class CreateStatementTable extends AbstractMigration
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
      $table = $this->table('statements');
      $table->addColumn('account_id', 'integer', ['null' => false, 'signed' => false])
          ->addForeignKey('account_id', 'accounts', 'id')
          ->addColumn('amount', 'decimal', ['null' => false, 'precision' => 20, 'scale' => 2])
          ->addColumn('event_date', 'datetime', ['null' => true])
          ->addColumn('remarks', 'string', ['limit' => 255, 'null' => true])
          ->addTimestamps()
          ->create();
  }
}