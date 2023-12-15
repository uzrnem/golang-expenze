<?php

declare(strict_types=1);

use Phinx\Migration\AbstractMigration;

final class CreateTagTable extends AbstractMigration
{
    public function change(): void
    {
        $table = $this->table('tags');
        $table->addColumn('name', 'string', ['limit' => 63, 'null' => false])
            ->addColumn('transaction_type_id', 'integer', ['null' => false, 'signed' => false])
            ->addForeignKey('transaction_type_id', 'transaction_types', 'id')
            ->addColumn('tag_id', 'integer', ['null' => true, 'signed' => false])
            ->addForeignKey('tag_id', 'tags', 'id', ['delete'=> 'SET_NULL', 'update'=> 'NO_ACTION'])
            ->addIndex(['name', 'transaction_type_id', 'tag_id'], ['unique' => true])
            ->create();
    }
}
