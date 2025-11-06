class Channel < ApplicationRecord
  has_many :schedule_channels, dependent: :destroy
  has_many :schedules, through: :schedule_channels
  has_many :send_logs, dependent: :destroy

  validates :name, presence: true
  validates :telegram_id, presence: true, uniqueness: true
  validates :channel_type, presence: true, inclusion: {
    in: %w[channel group supergroup user]
  }

  scope :channels, -> { where(channel_type: 'channel') }
  scope :groups, -> { where(channel_type: ['group', 'supergroup']) }
  scope :users, -> { where(channel_type: 'user') }

  def display_name
    title.presence || username.presence || "ID: #{telegram_id}"
  end

  def formatted_type
    channel_type.titleize
  end
end
