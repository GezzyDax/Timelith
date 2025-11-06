class TelegramService
  class << self
    def send_auth_code(account)
      response = backend_request(:post, "/api/v1/telegram/auth/send_code", {
        phone_number: account.phone_number,
        account_id: account.id
      })

      if response[:success]
        { success: true }
      else
        { success: false, error: response[:error] || 'Failed to send code' }
      end
    rescue StandardError => e
      { success: false, error: e.message }
    end

    def verify_code(account, code)
      response = backend_request(:post, "/api/v1/telegram/auth/verify_code", {
        account_id: account.id,
        code: code
      })

      if response[:success]
        {
          success: true,
          user_id: response[:user_id],
          first_name: response[:first_name],
          last_name: response[:last_name],
          username: response[:username]
        }
      else
        { success: false, error: response[:error] || 'Verification failed' }
      end
    rescue StandardError => e
      { success: false, error: e.message }
    end

    def disconnect(account)
      response = backend_request(:post, "/api/v1/telegram/disconnect", {
        account_id: account.id
      })

      if response[:success]
        { success: true }
      else
        { success: false, error: response[:error] || 'Failed to disconnect' }
      end
    rescue StandardError => e
      { success: false, error: e.message }
    end

    def check_status(account)
      response = backend_request(:get, "/api/v1/telegram/status/#{account.id}")

      if response[:success]
        { success: true, status: response[:status] }
      else
        { success: false, error: response[:error] }
      end
    rescue StandardError => e
      { success: false, error: e.message }
    end

    def sync_channels(account)
      response = backend_request(:post, "/api/v1/telegram/channels/sync", {
        account_id: account.id
      })

      if response[:success]
        { success: true, channels: response[:channels] }
      else
        { success: false, error: response[:error] }
      end
    rescue StandardError => e
      { success: false, error: e.message }
    end

    private

    def backend_request(method, path, body = nil)
      url = "#{ENV['GO_BACKEND_URL']}#{path}"
      api_key = ENV['GO_API_KEY']

      connection = Faraday.new do |conn|
        conn.request :json
        conn.response :json
        conn.adapter Faraday.default_adapter
      end

      response = case method
      when :get
        connection.get(url) do |req|
          req.headers['X-API-Key'] = api_key
        end
      when :post
        connection.post(url) do |req|
          req.headers['X-API-Key'] = api_key
          req.body = body
        end
      when :put
        connection.put(url) do |req|
          req.headers['X-API-Key'] = api_key
          req.body = body
        end
      when :delete
        connection.delete(url) do |req|
          req.headers['X-API-Key'] = api_key
        end
      end

      if response.success?
        JSON.parse(response.body).deep_symbolize_keys
      else
        { success: false, error: "HTTP #{response.status}: #{response.body}" }
      end
    rescue Faraday::Error => e
      { success: false, error: "Connection error: #{e.message}" }
    end
  end
end
