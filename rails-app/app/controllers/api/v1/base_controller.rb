module Api
  module V1
    class BaseController < ActionController::API
      before_action :verify_api_key

      private

      def verify_api_key
        api_key = request.headers['X-API-Key']
        expected_key = ENV['GO_API_KEY']

        unless api_key.present? && api_key == expected_key
          render json: { error: 'Unauthorized' }, status: :unauthorized
        end
      end
    end
  end
end
