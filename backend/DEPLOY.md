# Backend Deployment Guide

This guide provides step-by-step instructions for deploying the Go backend of the P2P Storage System. The recommended stack is Render for hosting, Supabase for the PostgreSQL database, and Cloudflare R2 for object storage.

## Prerequisites

- A [Render](https://render.com/) account.
- A [Supabase](https://supabase.com/) account.
- A [Cloudflare](https://www.cloudflare.com/) account.

## Step 1: Set Up the Database (Supabase)

1.  **Create a New Project:** Log in to your Supabase account and create a new project.
2.  **Get the Connection String:** Navigate to your project's **Settings > Database**. Under **Connection string**, copy the `URI` that starts with `postgresql://`. You will use this as your `DATABASE_URL`.
3.  **Run Migrations:** In the Supabase dashboard, go to the **SQL Editor**. Create a **New query** and paste the contents of `backend/migrations/001_initial_schema.sql` into the editor. Click **Run** to create the necessary tables.

## Step 2: Set Up Object Storage (Cloudflare R2)

1.  **Create an R2 Bucket:** Log in to your Cloudflare account and navigate to **R2**. Create a new bucket. The recommended name is `peer-storage-anchor`.
2.  **Generate API Tokens:** Go to **R2 > Manage R2 API Tokens**. Create a new API token with **Object Read & Write** permissions.
3.  **Copy Credentials:** Safely copy the **Account ID**, **Access Key ID**, and **Secret Access Key**. You will need these for the backend configuration.

## Step 3: Deploy the Backend (Render)

1.  **Create `render.yaml`:** Create a file named `render.yaml` in the root of your forked repository with the following content:
    ```yaml
    services:
      - type: web
        name: p2p-storage-api
        env: go
        buildCommand: go build -o bin/server cmd/server/main.go
        startCommand: ./bin/server
        envVars:
          - key: DATABASE_URL
            sync: false
          - key: JWT_SECRET
            generateValue: true # Render will generate a secure secret
          - key: R2_ACCOUNT_ID
            sync: false
          - key: R2_ACCESS_KEY_ID
            sync: false
          - key: R2_SECRET_ACCESS_KEY
            sync: false
    ```
    *Note: This file is for setting up a new service from scratch. If you are connecting an existing repository, you can set these values in the Render dashboard instead.*

2.  **Create a New Web Service:** Log in to your Render account and create a new **Web Service**.
3.  **Connect Your Repository:** Connect the GitHub repository where your backend code is located.
4.  **Configure Environment Variables:** In the Render dashboard for your new service, go to the **Environment** tab and add the following secrets:
    *   `DATABASE_URL`: The connection string you copied from Supabase.
    *   `JWT_SECRET`: A long, random, and secure string that you generate.
    *   `R2_ACCOUNT_ID`: Your Cloudflare R2 Account ID.
    *   `R2_ACCESS_KEY_ID`: Your Cloudflare R2 Access Key ID.
    *   `R2_SECRET_ACCESS_KEY`: Your Cloudflare R2 Secret Access Key.

5.  **Deploy:** Render will automatically build and deploy your application. After the first deployment, it will automatically redeploy on every push to your connected branch.

Your backend API is now live. You can find its public URL in the Render dashboard.
