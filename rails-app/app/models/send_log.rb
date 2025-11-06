class SendLog < ApplicationRecord
  belongs_to :schedule
  belongs_to :telegram_account
  belongs_to :channel

  validates :status, presence: true, inclusion: {
    in: %w[pending sending sent failed]
  }

  scope :pending, -> { where(status: 'pending') }
  scope :sent, -> { where(status: 'sent') }
  scope :failed, -> { where(status: 'failed') }
  scope :recent, -> { order(created_at: :desc) }

  def success?
    status == 'sent'
  end

  def failed?
    status == 'failed'
  end

  def can_retry?
    failed? && retry_count < 3
  end
end
