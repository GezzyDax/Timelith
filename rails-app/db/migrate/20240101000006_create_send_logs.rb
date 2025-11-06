class CreateSendLogs < ActiveRecord::Migration[7.1]
  def change
    create_table :send_logs do |t|
      t.references :schedule, null: false, foreign_key: true
      t.references :telegram_account, null: false, foreign_key: true
      t.references :channel, null: false, foreign_key: true
      t.string :status, null: false # 'pending', 'sending', 'sent', 'failed'
      t.text :message_content
      t.bigint :telegram_message_id
      t.text :error_message
      t.datetime :sent_at
      t.integer :retry_count, default: 0
      t.timestamps
    end

    add_index :send_logs, :status
    add_index :send_logs, :sent_at
    add_index :send_logs, :created_at
  end
end
