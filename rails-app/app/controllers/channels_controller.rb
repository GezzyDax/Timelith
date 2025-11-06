class ChannelsController < ApplicationController
  before_action :set_channel, only: [:show, :edit, :update, :destroy]

  def index
    @channels = Channel.all.order(created_at: :desc)
  end

  def show
  end

  def new
    @channel = Channel.new
  end

  def create
    @channel = Channel.new(channel_params)

    if @channel.save
      redirect_to @channel, notice: 'Channel created successfully.'
    else
      render :new
    end
  end

  def edit
  end

  def update
    if @channel.update(channel_params)
      redirect_to @channel, notice: 'Channel updated successfully.'
    else
      render :edit
    end
  end

  def destroy
    @channel.destroy
    redirect_to channels_path, notice: 'Channel deleted.'
  end

  def sync_from_telegram
    account = TelegramAccount.active.first

    if account.nil?
      redirect_to channels_path, alert: 'No active Telegram account found.'
      return
    end

    result = TelegramService.sync_channels(account)

    if result[:success]
      synced_count = 0
      result[:channels].each do |channel_data|
        channel = Channel.find_or_initialize_by(telegram_id: channel_data[:telegram_id])
        channel.update(
          name: channel_data[:name],
          channel_type: channel_data[:type],
          username: channel_data[:username],
          title: channel_data[:title],
          members_count: channel_data[:members_count],
          last_synced_at: Time.current
        )
        synced_count += 1
      end

      redirect_to channels_path, notice: "Synced #{synced_count} channels from Telegram."
    else
      redirect_to channels_path, alert: "Failed to sync channels: #{result[:error]}"
    end
  end

  private

  def set_channel
    @channel = Channel.find(params[:id])
  end

  def channel_params
    params.require(:channel).permit(:name, :channel_type, :telegram_id, :username, :title, :description)
  end
end
