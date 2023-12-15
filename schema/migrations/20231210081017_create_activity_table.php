<?php

declare(strict_types=1);

use Phinx\Migration\AbstractMigration;

final class CreateActivityTable extends AbstractMigration
{
    public function change(): void
    {
        if (!$this->isMigratingUp()) {
            $this->query('DROP TRIGGER IF EXISTS `a1_before_insert_activity`');
            $this->query('DROP TRIGGER IF EXISTS `a2_after_insert_activity`');
            $this->query('DROP TRIGGER IF EXISTS `b1_before_update_activity`');
            $this->query('DROP TRIGGER IF EXISTS `b2_after_udpate_activity`');
            $this->query('DROP TRIGGER IF EXISTS `c1_before_delete_activity`');
        }
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
            $this->execute($this->getBeforeInsertActivityTrigger());
            $this->execute($this->getAfterInsertActivityTrigger());
            $this->execute($this->getBeforeUpdateActivityTrigger());
            $this->execute($this->getAfterUpdateActivityTrigger());
            $this->execute($this->getBeforeDeleteActivityTrigger());
        }
    }

    public function getBeforeInsertActivityTrigger():string {
        return "CREATE TRIGGER `a1_before_insert_activity` BEFORE INSERT ON `activities` FOR EACH ROW BEGIN
        IF NEW.from_account_id IS NOT NULL AND NEW.from_account_id in ('', 0, '0') THEN
          SET NEW.from_account_id = NULL;
        END IF;
        IF NEW.to_account_id IS NOT NULL AND NEW.to_account_id in ('', 0, '0') THEN
          SET NEW.to_account_id = NULL;
        END IF;
        IF NEW.sub_tag_id IS NOT NULL AND NEW.sub_tag_id in ('', 0, '0') THEN
          SET NEW.sub_tag_id = NULL;
        END IF;
        IF NEW.tag_id IS NOT NULL AND NEW.tag_id in ('', 0, '0') THEN
          SET NEW.tag_id = NULL;
        END IF;
        IF NEW.amount IS NOT NULL AND NEW.amount in ('', 0, '0') THEN
          SIGNAL SQLSTATE '45000' 
          SET MESSAGE_TEXT = 'amount can not be zero';
        ELSEIF NEW.from_account_id IS NULL AND NEW.to_account_id IS NULL THEN
          SIGNAL SQLSTATE '45000' 
          SET MESSAGE_TEXT = 'transaction not allowed';
        ELSEIF NEW.from_account_id IS NULL THEN
          SET NEW.transaction_type_id = ( select id from transaction_types where name = 'Income' );
        ELSEIF NEW.to_account_id IS NULL THEN
          SET NEW.transaction_type_id = ( select id from transaction_types where name = 'Expense' );
        ELSE
          SET NEW.transaction_type_id = ( select id from transaction_types where name = 'Transfer' );
        END IF;
      END";
    }

    public function getAfterInsertActivityTrigger():string {
        return "CREATE TRIGGER `a2_after_insert_activity` AFTER INSERT ON `activities` FOR EACH ROW BEGIN
        DECLARE old_balance decimal(20,2);
        DECLARE transaction_type varchar(20);
      
        IF NEW.to_account_id IS NOT NULL THEN
          set old_balance = ( select amount from accounts where id = NEW.to_account_id );
          set transaction_type = ( select id from transaction_types where name = 'Income' );
          INSERT INTO passbooks(account_id, activity_id, previous_balance, transaction_type_id, balance, created_at, updated_at)
          VALUES(NEW.to_account_id, NEW.id, old_balance, transaction_type, old_balance + NEW.amount, now(), now());
          update accounts set amount = old_balance + NEW.amount, updated_at = now() where id = NEW.to_account_id;
        END IF;
      
        IF NEW.from_account_id IS NOT NULL THEN
          set old_balance = ( select amount from accounts where id = NEW.from_account_id );
          set transaction_type = ( select id from transaction_types where name = 'Expense' );
          INSERT INTO passbooks(account_id, activity_id, previous_balance, transaction_type_id, balance, created_at, updated_at)
          VALUES(NEW.from_account_id, NEW.id, old_balance, transaction_type, old_balance - NEW.amount, now(), now());
          update accounts set amount = old_balance - NEW.amount, updated_at = now() where id = NEW.from_account_id;
        END IF;
      END";
    }

    public function getBeforeUpdateActivityTrigger():string {
        return "CREATE TRIGGER `b1_before_update_activity` BEFORE UPDATE ON `activities` FOR EACH ROW BEGIN
        IF NEW.from_account_id IS NOT NULL AND NEW.from_account_id in ('', 0, '0') THEN
          SET NEW.from_account_id = NULL;
        END IF;
        IF NEW.to_account_id IS NOT NULL AND NEW.to_account_id in ('', 0, '0') THEN
          SET NEW.to_account_id = NULL;
        END IF;
        IF NEW.sub_tag_id IS NOT NULL AND NEW.sub_tag_id in ('', 0, '0') THEN
          SET NEW.sub_tag_id = NULL;
        END IF;
        IF NEW.tag_id IS NOT NULL AND NEW.tag_id in ('', 0, '0') THEN
          SET NEW.tag_id = NULL;
        END IF;
        IF NEW.amount IS NOT NULL AND NEW.amount in ('', 0, '0') THEN
          SIGNAL SQLSTATE '45000' 
          SET MESSAGE_TEXT = 'amount can not be zero';
        ELSEIF NEW.from_account_id IS NULL AND NEW.to_account_id IS NULL THEN
          SIGNAL SQLSTATE '45000' 
          SET MESSAGE_TEXT = 'transaction not allowed';
        ELSEIF NEW.from_account_id IS NULL THEN
          SET NEW.transaction_type_id = ( select id from transaction_types where name = 'Income' );
        ELSEIF NEW.to_account_id IS NULL THEN
          SET NEW.transaction_type_id = ( select id from transaction_types where name = 'Expense' );
        ELSE
          SET NEW.transaction_type_id = ( select id from transaction_types where name = 'Transfer' );
        END IF;
      END";
    }

    public function getAfterUpdateActivityTrigger():string {
        return "CREATE TRIGGER `b2_after_udpate_activity` AFTER UPDATE ON `activities` FOR EACH ROW BEGIN
        DECLARE old_balance decimal(20,2);
        DECLARE transaction_type varchar(20);
      
      
        IF NOT (OLD.from_account_id <=> NEW.from_account_id) THEN
          IF OLD.from_account_id IS NOT NULL THEN
            DELETE FROM passbooks where activity_id = OLD.id and account_id = OLD.from_account_id;
            update accounts set amount = amount + OLD.amount, updated_at = now() where id = OLD.from_account_id;
          END IF;
          
          IF NEW.from_account_id IS NOT NULL THEN
            set old_balance = ( select amount from accounts where id = NEW.from_account_id );
            set transaction_type = ( select id from transaction_types where name = 'Expense' );
            INSERT INTO passbooks(account_id, activity_id, previous_balance, transaction_type_id, balance, created_at, updated_at)
            VALUES(NEW.from_account_id, NEW.id, old_balance, transaction_type, old_balance - NEW.amount, now(), now());
            update accounts set amount = old_balance - NEW.amount, updated_at = now() where id = NEW.from_account_id;
          END IF;
        END IF;
      
        IF NOT (OLD.to_account_id <=> NEW.to_account_id) THEN
          IF OLD.to_account_id IS NOT NULL THEN
            DELETE FROM passbooks where activity_id = OLD.id and account_id = OLD.to_account_id;
            update accounts a set amount = amount - OLD.amount, updated_at = now() where id = OLD.to_account_id;
          END IF;
      
          IF NEW.to_account_id IS NOT NULL THEN
            set old_balance = ( select amount from accounts where id = NEW.to_account_id );
            set transaction_type = ( select id from transaction_types where name = 'Income' );
            INSERT INTO passbooks(account_id, activity_id, previous_balance, transaction_type_id, balance, created_at, updated_at)
            VALUES(NEW.to_account_id, NEW.id, old_balance, transaction_type, old_balance + NEW.amount, now(), now());
            update accounts set amount = old_balance + NEW.amount, updated_at = now() where id = NEW.to_account_id;
          END IF;
        END IF;
      
        IF (OLD.from_account_id <=> NEW.from_account_id) AND
          (OLD.to_account_id <=> NEW.to_account_id) AND
          NOT (OLD.amount <=> NEW.amount) THEN
          IF OLD.from_account_id IS NOT NULL THEN
            update passbooks set balance = balance + OLD.amount - NEW.amount, updated_at = now()
            where activity_id = OLD.id and account_id = OLD.from_account_id;
      
            update accounts a set amount = amount + OLD.amount - NEW.amount, updated_at = now() where id = OLD.from_account_id;
          END IF;
      
          IF OLD.to_account_id IS NOT NULL THEN
            update passbooks set balance = balance - OLD.amount + NEW.amount, updated_at = now()
            where activity_id = OLD.id and account_id = OLD.to_account_id;
      
            update accounts a set amount = amount - OLD.amount + NEW.amount, updated_at = now() where id = OLD.to_account_id;
          END IF;
        END IF;
      END";
    }

    public function getBeforeDeleteActivityTrigger():string {
        return "CREATE TRIGGER `c1_before_delete_activity` BEFORE DELETE ON `activities` FOR EACH ROW BEGIN

        DELETE FROM passbooks where activity_id = OLD.id;
        IF OLD.from_account_id IS NOT NULL THEN
          update accounts set amount = amount + OLD.amount, updated_at = now() where id = OLD.from_account_id;
        END IF;
        
        IF OLD.to_account_id IS NOT NULL THEN
          update accounts a set amount = amount - OLD.amount, updated_at = now() where id = OLD.to_account_id;
        END IF;
      END";
    }
}
