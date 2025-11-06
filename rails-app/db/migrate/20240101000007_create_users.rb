class CreateUsers < ActiveRecord::Migration[7.1]
  def change
    create_table :users do |t|
      t.string :email, null: false
      t.string :password_digest, null: false
      t.string :role, default: 'admin', null: false
      t.datetime :last_login_at
      t.timestamps
    end

    add_index :users, :email, unique: true
  end
end
