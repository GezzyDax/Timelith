class MessageTemplate < ApplicationRecord
  has_many :schedules, dependent: :destroy

  validates :name, presence: true
  validates :content, presence: true
  validates :parse_mode, inclusion: { in: %w[markdown html], allow_nil: true }
  validates :media_type, inclusion: {
    in: %w[photo video document audio voice],
    allow_nil: true
  }

  scope :with_media, -> { where.not(media_url: nil) }
  scope :text_only, -> { where(media_url: nil) }

  def has_media?
    media_url.present?
  end

  def formatted_buttons
    return [] if buttons.blank?

    # Parse JSON buttons for display
    JSON.parse(buttons) rescue []
  end
end
