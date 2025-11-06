class MessageTemplatesController < ApplicationController
  before_action :set_message_template, only: [:show, :edit, :update, :destroy]

  def index
    @message_templates = MessageTemplate.all.order(created_at: :desc)
  end

  def show
  end

  def new
    @message_template = MessageTemplate.new
  end

  def create
    @message_template = MessageTemplate.new(message_template_params)

    if @message_template.save
      redirect_to @message_template, notice: 'Template created successfully.'
    else
      render :new
    end
  end

  def edit
  end

  def update
    if @message_template.update(message_template_params)
      redirect_to @message_template, notice: 'Template updated successfully.'
    else
      render :edit
    end
  end

  def destroy
    @message_template.destroy
    redirect_to message_templates_path, notice: 'Template deleted.'
  end

  private

  def set_message_template
    @message_template = MessageTemplate.find(params[:id])
  end

  def message_template_params
    params.require(:message_template).permit(
      :name, :content, :media_type, :media_url,
      :parse_mode_enabled, :parse_mode,
      :disable_web_page_preview, :buttons
    )
  end
end
