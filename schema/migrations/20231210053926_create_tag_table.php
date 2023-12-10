<?php

declare(strict_types=1);

use Phinx\Migration\AbstractMigration;

final class CreateTagTable extends AbstractMigration
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
