class CreateScheduleChannels < ActiveRecord::Migration[7.1]
  def change
    create_table :schedule_channels do |t|
      t.references :schedule, null: false, foreign_key: true
      t.references :channel, null: false, foreign_key: true
      t.timestamps
    end

    add_index :schedule_channels, [:schedule_id, :channel_id], unique: true
  end
end
