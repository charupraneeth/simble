# simble 🟢

**simble** is a lightweight, privacy-friendly, open-source web analytics platform. Built for developers who want a simple, clean, and fast dashboard to track their website traffic without compromising user privacy.

![simble Dashboard Preview](web/src/assets/dashboard-preview.png)

## Features

- **Privacy-First:** Lightweight tracking that respects users.
- **Fast and Simple Dashboard:** Beautiful, responsive UI built with Vue 3, Tailwind CSS, and Chart.js.
- **GitHub Authentication:** Secure, one-click OAuth login.
- **Self-Hostable:** Easy deployment to platforms like [Railway](https://railway.app/).
- **Tiny Tracking Script:** Drop-in `<script>` tag under 1KB that stays out of your way.

## Tech Stack

- **Backend:** Go (`net/http`), pgx (Postgres driver)
- **Frontend:** Vue 3 (Composition API), Vite, Tailwind CSS
- **Database:** PostgreSQL
- **GeoIP:** MaxMind GeoLite2

---

## Local Development Setup

To run **simble** locally, you will need Go 1.23+, Node.js, and a running PostgreSQL instance.

### 1. Database & Environment Variables

Create a PostgreSQL database for the application. Then, set up your local environment variables:

```bash
# Database connection
export DATABASE_URL="postgres://user:password@localhost:5432/simble_dev"
export AUTO_MIGRATE_DB="true"

# GitHub OAuth App Credentials (create one at https://github.com/settings/developers)
export GITHUB_CLIENT_ID="your_client_id"
export GITHUB_CLIENT_SECRET="your_client_secret"
export GITHUB_REDIRECT_URL="http://localhost:8080/auth/github/callback"

# Server Port
export PORT="8080"
```

### 2. Run the Go Backend

The backend serves the API, handles OAuth, and ingests analytics events.

```bash
# From the project root
go run main.go
```
*Note: Make sure your `GeoLite2-City.mmdb` database is downloaded and available at `./GeoLite2-City/GeoLite2-City.mmdb` for the Geo IP features to work!*

### 3. Run the Frontend Vue App

The frontend provides the main dashboard and onboarding interfaces.

```bash
cd web

# Install dependencies
pnpm install

# Start the development server
pnpm run dev
```

The frontend will run on `http://localhost:5173` but API requests are proxied over to the Go backend running on port `8080`.

## Deployment

simble is designed to be easily deployed to Railway:
1. Provision a PostgreSQL container.
2. Link your GitHub repository.
3. Configure the environment variables (make sure to set `GITHUB_REDIRECT_URL` to your production domain).
4. The Railpack build process will automatically bundle the Vue frontend and compile the Go binary for production.

## Tracking Snippet Integration

Once deployed, simply copy the tracking snippet provided during the site onboarding flow into your website's `<head>`:

```html
<script defer data-domain="yourdomain.com" src="https://simble.dev/script.js"></script>
```

---
*simble is built keeping simplicity in mind.*
