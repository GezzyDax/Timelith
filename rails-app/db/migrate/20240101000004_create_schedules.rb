class CreateSchedules < ActiveRecord::Migration[7.1]
  def change
    create_table :schedules do |t|
      t.string :name, null: false
      t.references :telegram_account, null: false, foreign_key: true
      t.references :message_template, null: false, foreign_key: true
      t.string :schedule_type, null: false # 'cron', 'interval', 'once'
      t.string :cron_expression
      t.integer :interval_minutes
      t.datetime :scheduled_at # for 'once' type
      t.string :timezone, default: 'UTC'
      t.boolean :active, default: false
      t.datetime :next_run_at
      t.datetime :last_run_at
      t.integer :total_runs, default: 0
      t.integer :successful_runs, default: 0
      t.integer :failed_runs, default: 0
      t.timestamps
    end

    add_index :schedules, :active
    add_index :schedules, :next_run_at
    add_index :schedules, :schedule_type
  end
end
