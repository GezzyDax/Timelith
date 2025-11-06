module Api
  module V1
    class SendLogsController < BaseController
      def create
        log = SendLog.new(send_log_params)

        if log.save
          # Update schedule statistics
          schedule = log.schedule
          schedule.record_execution(
            success: log.status == 'sent',
            error: log.error_message
          )

          render json: { success: true, log_id: log.id }, status: :created
        else
          render json: { success: false, errors: log.errors.full_messages }, status: :unprocessable_entity
        end
      end

      private

      def send_log_params
        params.require(:send_log).permit(
          :schedule_id, :telegram_account_id, :channel_id,
          :status, :message_content, :telegram_message_id,
          :error_message, :sent_at, :retry_count
        )
      end
    end
  end
end
