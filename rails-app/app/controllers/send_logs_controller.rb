class SendLogsController < ApplicationController
  def index
    @send_logs = SendLog.recent
                        .includes(:schedule, :telegram_account, :channel)
                        .page(params[:page])
                        .per(50)

    @send_logs = @send_logs.where(status: params[:status]) if params[:status].present?
  end

  def show
    @send_log = SendLog.find(params[:id])
  end
end
