class Schedule < ApplicationRecord
  belongs_to :telegram_account
  belongs_to :message_template
  has_many :schedule_channels, dependent: :destroy
  has_many :channels, through: :schedule_channels
  has_many :send_logs, dependent: :destroy

  validates :name, presence: true
  validates :schedule_type, presence: true, inclusion: {
    in: %w[cron interval once]
  }
  validates :cron_expression, presence: true, if: -> { schedule_type == 'cron' }
  validates :interval_minutes, presence: true, numericality: { greater_than: 0 },
            if: -> { schedule_type == 'interval' }
  validates :scheduled_at, presence: true, if: -> { schedule_type == 'once' }

  validate :validate_channels_presence

  scope :active, -> { where(active: true) }
  scope :inactive, -> { where(active: false) }
  scope :due_for_execution, -> {
    active.where('next_run_at <= ?', Time.current)
  }

  before_save :calculate_next_run_at, if: :should_calculate_next_run?

  def activate!
    update(active: true)
    calculate_next_run_at
    save
  end

  def deactivate!
    update(active: false, next_run_at: nil)
  end

  def calculate_next_run_at
    self.next_run_at = case schedule_type
    when 'cron'
      calculate_next_cron_run
    when 'interval'
      (last_run_at || Time.current) + interval_minutes.minutes
    when 'once'
      scheduled_at
    end
  end

  def record_execution(success:, error: nil)
    self.last_run_at = Time.current
    self.total_runs += 1

    if success
      self.successful_runs += 1
    else
      self.failed_runs += 1
    end

    calculate_next_run_at if schedule_type != 'once'
    self.active = false if schedule_type == 'once' # Disable one-time schedules

    save
  end

  private

  def validate_channels_presence
    if channels.empty?
      errors.add(:channels, "must have at least one channel")
    end
  end

  def should_calculate_next_run?
    active? && (schedule_type_changed? || cron_expression_changed? ||
                interval_minutes_changed? || scheduled_at_changed?)
  end

  def calculate_next_cron_run
    # Simple cron parser - in production, use a gem like 'fugit'
    # For now, return a placeholder
    Time.current + 1.hour
  end
end
