class CreateChannels < ActiveRecord::Migration[7.1]
  def change
    create_table :channels do |t|
      t.string :name, null: false
      t.string :channel_type, null: false # 'channel', 'group', 'user'
      t.bigint :telegram_id, null: false
      t.string :username
      t.string :title
      t.integer :members_count
      t.text :description
      t.datetime :last_synced_at
      t.timestamps
    end

    add_index :channels, :telegram_id, unique: true
    add_index :channels, :username
    add_index :channels, :channel_type
  end
end
