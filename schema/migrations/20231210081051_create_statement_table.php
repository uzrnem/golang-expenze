<?php

declare(strict_types=1);

use Phinx\Migration\AbstractMigration;

final class CreateStatementTable extends AbstractMigration
{
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