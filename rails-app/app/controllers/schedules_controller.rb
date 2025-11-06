class SchedulesController < ApplicationController
  before_action :set_schedule, only: [:show, :edit, :update, :destroy, :activate, :deactivate]

  def index
    @schedules = Schedule.all.includes(:telegram_account, :message_template, :channels).order(created_at: :desc)
  end

  def show
    @recent_logs = @schedule.send_logs.recent.limit(20)
  end

  def new
    @schedule = Schedule.new
    @telegram_accounts = TelegramAccount.active
    @message_templates = MessageTemplate.all
    @channels = Channel.all
  end

  def create
    @schedule = Schedule.new(schedule_params)

    if @schedule.save
      redirect_to @schedule, notice: 'Schedule created successfully.'
    else
      @telegram_accounts = TelegramAccount.active
      @message_templates = MessageTemplate.all
      @channels = Channel.all
      render :new
    end
  end

  def edit
    @telegram_accounts = TelegramAccount.active
    @message_templates = MessageTemplate.all
    @channels = Channel.all
  end

  def update
    if @schedule.update(schedule_params)
      redirect_to @schedule, notice: 'Schedule updated successfully.'
    else
      @telegram_accounts = TelegramAccount.active
      @message_templates = MessageTemplate.all
      @channels = Channel.all
      render :edit
    end
  end

  def destroy
    @schedule.destroy
    redirect_to schedules_path, notice: 'Schedule deleted.'
  end

  def activate
    if @schedule.activate!
      redirect_to @schedule, notice: 'Schedule activated.'
    else
      redirect_to @schedule, alert: 'Failed to activate schedule.'
    end
  end

  def deactivate
    if @schedule.deactivate!
      redirect_to @schedule, notice: 'Schedule deactivated.'
    else
      redirect_to @schedule, alert: 'Failed to deactivate schedule.'
    end
  end

  private

  def set_schedule
    @schedule = Schedule.find(params[:id])
  end

  def schedule_params
    params.require(:schedule).permit(
      :name, :telegram_account_id, :message_template_id,
      :schedule_type, :cron_expression, :interval_minutes,
      :scheduled_at, :timezone, :active,
      channel_ids: []
    )
  end
end
