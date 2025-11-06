module Api
  module V1
    class TelegramAccountsController < BaseController
      def session_data
        account = TelegramAccount.find(params[:id])

        render json: {
          success: true,
          account_id: account.id,
          phone_number: account.phone_number,
          session_data: account.session_data,
          status: account.status
        }
      rescue ActiveRecord::RecordNotFound
        render json: { success: false, error: 'Account not found' }, status: :not_found
      end

      def update_status
        account = TelegramAccount.find(params[:id])
        account.update(
          status: params[:status],
          last_active_at: Time.current,
          error_message: params[:error_message]
        )

        render json: { success: true }
      rescue ActiveRecord::RecordNotFound
        render json: { success: false, error: 'Account not found' }, status: :not_found
      end
    end
  end
end
