class ScheduleChannel < ApplicationRecord
  belongs_to :schedule
  belongs_to :channel

  validates :schedule_id, uniqueness: { scope: :channel_id }
end
