class CreateTelegramAccounts < ActiveRecord::Migration[7.1]
  def change
    create_table :telegram_accounts do |t|
      t.string :phone_number, null: false
      t.string :session_data, limit: 10000
      t.string :status, default: 'pending', null: false
      t.string :first_name
      t.string :last_name
      t.string :username
      t.bigint :telegram_user_id
      t.datetime :last_active_at
      t.text :error_message
      t.timestamps
    end

    add_index :telegram_accounts, :phone_number, unique: true
    add_index :telegram_accounts, :status
    add_index :telegram_accounts, :telegram_user_id, unique: true
  end
end
