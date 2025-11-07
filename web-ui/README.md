# Timelith Web UI

Next.js web dashboard for Timelith.

## Tech Stack

- Next.js 14 (App Router)
- TypeScript
- TailwindCSS
- TanStack Query
- Axios

## Development

```bash
# Install dependencies
npm install

# Run dev server
npm run dev

# Build
npm run build

# Start production
npm start
```

## Environment Variables

- `NEXT_PUBLIC_API_URL` - Backend API URL (default: http://localhost:8080)

## Pages

- `/` - Dashboard
- `/login` - Login page
- `/accounts` - Manage Telegram accounts
- `/templates` - Message templates
- `/channels` - Target channels
- `/schedules` - Scheduled jobs
- `/logs` - Job execution logs
