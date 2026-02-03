HELLO

# Guide: Creating a User in Thunder

## Purpose

This document provides step-by-step instructions for creating a user in Thunder. It includes both API-based and console-based methods to ensure flexibility.

## When to Use

Use this guide when you need to add a new user to your Thunder instance, either programmatically via the API or manually through the console.

## Before You Begin

Before creating a user, ensure the following:

1. **Thunder Server**: Verify that the Thunder server is running.
2. **Authentication**: Get an access token or API key for authentication.
3. **Permissions**: Ensure you have the necessary permissions to create users.
4. **User Data**: Prepare the required user details, such as username, email, password, and roles.

## Configuration Steps

### Using the API

The API method allows you to programmatically create users by sending HTTP requests to the Thunder server. Follow the steps below to complete the process.

#### Step 1: Prepare the API Request

- **Endpoint**: `/api/v1/users`
- **Method**: `POST`
- **Headers**:
  - `Authorization: Bearer <access_token>`
  - `Content-Type: application/json`

- **Request Body**:

  ```json
  {
    "username": "johndoe",
    "email": "johndoe@example.com",
    "password": "securepassword",
    "firstName": "John",
    "lastName": "Doe",
    "roles": ["user"]
  }
  ```

#### Step 2: Send the Request

Use a tool like `curl`, Postman, or any HTTP client to send the request. Example using `curl`:

```bash
curl -X POST \
  https://<thunder-server-url>/api/v1/users \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "username": "johndoe",
    "email": "johndoe@example.com",
    "password": "securepassword",
    "firstName": "John",
    "lastName": "Doe",
    "roles": ["user"]
  }'
```

#### Step 3: Verify the Response

A successful response will return a `201 Created` status code with the created user's details:

```json
{
  "id": "12345",
  "username": "johndoe",
  "email": "johndoe@example.com",
  "firstName": "John",
  "lastName": "Doe",
  "roles": ["user"]
}
```

### Using the Thunder Console

The Thunder Console provides a graphical interface for managing users. Follow the steps below to create a user through the console.

#### Step 1: Sign In

- Open the Thunder console and sign in with an admin account.

#### Step 2: Navigate to User Management

- Go to the "Users" section in the console.

#### Step 3: Create a New User

- Click the "Create User" button.
- Fill in the required fields (e.g., username, email, password, roles).
- Click "Save" to create the user.

#### Step 4: Verify the User

- Ensure the new user appears in the user list.

## How to Verify

After creating a user, verify the following:

1. The user appears in the user list (console) or the API response.
2. The user can sign in with the provided credentials.
3. Assigned roles and permissions are functioning as expected.

## How to Troubleshoot

- **Issue**: API request fails with a `401 Unauthorized` error.
  - **Solution**: Verify the access token and ensure it hasn't expired.

- **Issue**: User creation fails with a `400 Bad Request` error.
  - **Solution**: Check the request body for missing or invalid fields.

- **Issue**: User can't sign in after creation.
  - **Solution**: Verify the password and assigned roles.

## Related References

- [API Documentation](../content/apis/user.yaml)
- [Setup Guide](../setup/README.md)