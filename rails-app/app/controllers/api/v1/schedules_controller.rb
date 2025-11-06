module Api
  module V1
    class SchedulesController < BaseController
      def index
        schedules = Schedule.active.includes(:telegram_account, :message_template, :channels)

        render json: {
          success: true,
          schedules: schedules.map { |s| serialize_schedule(s) }
        }
      end

      private

      def serialize_schedule(schedule)
        {
          id: schedule.id,
          name: schedule.name,
          telegram_account_id: schedule.telegram_account_id,
          message_template: {
            content: schedule.message_template.content,
            media_type: schedule.message_template.media_type,
            media_url: schedule.message_template.media_url,
            parse_mode: schedule.message_template.parse_mode,
            buttons: schedule.message_template.buttons
          },
          channels: schedule.channels.map { |c| { id: c.id, telegram_id: c.telegram_id, name: c.name } },
          schedule_type: schedule.schedule_type,
          cron_expression: schedule.cron_expression,
          interval_minutes: schedule.interval_minutes,
          scheduled_at: schedule.scheduled_at,
          next_run_at: schedule.next_run_at,
          timezone: schedule.timezone
        }
      end
    end
  end
end
