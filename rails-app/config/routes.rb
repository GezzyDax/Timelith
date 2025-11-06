Rails.application.routes.draw do
  # Health check endpoint
  get "up" => "rails/health#show", as: :rails_health_check

  # Authentication
  get    'login',  to: 'sessions#new'
  post   'login',  to: 'sessions#create'
  delete 'logout', to: 'sessions#destroy'

  # Dashboard
  root 'dashboard#index'

  # Telegram Accounts
  resources :telegram_accounts do
    member do
      post :send_code
      post :verify_code
      post :disconnect
      get :status
    end
  end

  # Message Templates
  resources :message_templates

  # Schedules
  resources :schedules do
    member do
      post :activate
      post :deactivate
    end
  end

  # Channels
  resources :channels do
    collection do
      post :sync_from_telegram
    end
  end

  # Send History / Logs
  resources :send_logs, only: [:index, :show]

  # API for Go Backend
  namespace :api do
    namespace :v1 do
      resources :telegram_accounts, only: [] do
        member do
          get :session_data
          post :update_status
        end
      end
      resources :schedules, only: [:index]
      resources :send_logs, only: [:create]
    end
  end
end
