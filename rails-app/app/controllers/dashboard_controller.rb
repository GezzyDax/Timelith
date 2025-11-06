class DashboardController < ApplicationController
  def index
    @telegram_accounts_count = TelegramAccount.count
    @active_accounts_count = TelegramAccount.active.count
    @schedules_count = Schedule.count
    @active_schedules_count = Schedule.active.count
    @total_sends = SendLog.count
    @successful_sends = SendLog.sent.count
    @failed_sends = SendLog.failed.count
    @recent_logs = SendLog.recent.limit(10).includes(:schedule, :channel, :telegram_account)
  end
end
