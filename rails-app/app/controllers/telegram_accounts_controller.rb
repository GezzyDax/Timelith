class TelegramAccountsController < ApplicationController
  before_action :set_telegram_account, only: [:show, :edit, :update, :destroy, :send_code, :verify_code, :disconnect, :status]

  def index
    @telegram_accounts = TelegramAccount.all.order(created_at: :desc)
  end

  def show
  end

  def new
    @telegram_account = TelegramAccount.new
  end

  def create
    @telegram_account = TelegramAccount.new(telegram_account_params)
    @telegram_account.status = 'pending'

    if @telegram_account.save
      # Send request to Go backend to initiate authorization
      result = TelegramService.send_auth_code(@telegram_account)

      if result[:success]
        redirect_to @telegram_account, notice: 'Please check your Telegram for the verification code.'
      else
        @telegram_account.update(status: 'error', error_message: result[:error])
        redirect_to @telegram_account, alert: "Error: #{result[:error]}"
      end
    else
      render :new
    end
  end

  def send_code
    result = TelegramService.send_auth_code(@telegram_account)

    if result[:success]
      redirect_to @telegram_account, notice: 'Verification code sent to your Telegram.'
    else
      redirect_to @telegram_account, alert: "Error: #{result[:error]}"
    end
  end

  def verify_code
    code = params[:code]
    result = TelegramService.verify_code(@telegram_account, code)

    if result[:success]
      @telegram_account.update(
        status: 'authorized',
        telegram_user_id: result[:user_id],
        first_name: result[:first_name],
        last_name: result[:last_name],
        username: result[:username],
        last_active_at: Time.current
      )
      redirect_to @telegram_account, notice: 'Successfully authorized!'
    else
      redirect_to @telegram_account, alert: "Verification failed: #{result[:error]}"
    end
  end

  def disconnect
    result = TelegramService.disconnect(@telegram_account)

    if result[:success]
      @telegram_account.update(status: 'disconnected', session_data: nil)
      redirect_to telegram_accounts_path, notice: 'Account disconnected.'
    else
      redirect_to @telegram_account, alert: "Error: #{result[:error]}"
    end
  end

  def status
    result = TelegramService.check_status(@telegram_account)

    if result[:success]
      @telegram_account.update(
        status: result[:status],
        last_active_at: Time.current
      )
    end

    redirect_to @telegram_account
  end

  def edit
  end

  def update
    if @telegram_account.update(telegram_account_params)
      redirect_to @telegram_account, notice: 'Account updated successfully.'
    else
      render :edit
    end
  end

  def destroy
    @telegram_account.destroy
    redirect_to telegram_accounts_path, notice: 'Account deleted.'
  end

  private

  def set_telegram_account
    @telegram_account = TelegramAccount.find(params[:id])
  end

  def telegram_account_params
    params.require(:telegram_account).permit(:phone_number)
  end
end
