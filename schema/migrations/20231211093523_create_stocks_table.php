<?php

declare(strict_types=1);

use Phinx\Migration\AbstractMigration;

final class CreateStocksTable extends AbstractMigration
{
    public function change(): void
    {
        $table = $this->table('stocks');
        $table->addColumn('code', 'string', ['limit' => 255, 'null' => false])
            ->addColumn('description', 'string', ['limit' => 255, 'null' => true])
            ->addColumn('quantity', 'integer', ['null' => false, 'signed' => false])
            ->addColumn('current_quantity', 'integer', ['null' => false, 'signed' => false])
            ->addColumn('buy_price', 'decimal', ['null' => false, 'precision' => 20, 'scale' => 2])
            ->addColumn('buy_charges', 'decimal', ['default' => 0, 'precision' => 20, 'scale' => 2])
            ->addColumn('buy_date', 'datetime', ['null' => false])
            ->addColumn('target', 'decimal', ['null' => true, 'precision' => 20, 'scale' => 2])
            ->addColumn('sell_price', 'decimal', ['null' => true, 'precision' => 20, 'scale' => 2])
            ->addColumn('sell_charges', 'decimal', ['null' => true, 'precision' => 20, 'scale' => 2])
            ->addColumn('sell_date', 'datetime', ['null' => true])
            ->addColumn('status', 'enum', ['null' => false, 'values' => ['open', 'close'], 'default' => 'open'])
            ->addColumn('parent_id', 'integer', ['null' => true, 'signed' => false])
            ->addTimestamps()
            ->create();
    }
}
