class TelegramAccount < ApplicationRecord
  has_many :schedules, dependent: :destroy
  has_many :send_logs, dependent: :destroy

  validates :phone_number, presence: true, uniqueness: true
  validates :status, presence: true, inclusion: {
    in: %w[pending authorized disconnected error]
  }

  scope :active, -> { where(status: 'authorized') }
  scope :pending, -> { where(status: 'pending') }

  # Encrypt session data before saving
  before_save :encrypt_session_data, if: :session_data_changed?

  def authorized?
    status == 'authorized'
  end

  def display_name
    if first_name.present?
      [first_name, last_name].compact.join(' ')
    elsif username.present?
      "@#{username}"
    else
      phone_number
    end
  end

  private

  def encrypt_session_data
    # Session data encryption logic will be handled by Go backend
    # This is just a placeholder for Rails validation
  end
end
