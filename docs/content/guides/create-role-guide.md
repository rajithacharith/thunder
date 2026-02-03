# Guide: Creating a Role in Thunder

## Purpose
This document provides step-by-step instructions for creating a role in Thunder. Roles are used to define permissions and access levels for users in the system.

## When to Use
Use this guide when you need to create a new role in Thunder to manage user permissions and access control.

## Prerequisites

Before creating a role, ensure the following:

1. **Thunder Server**:
   - Verify that the Thunder server is running. Refer to the [Setup Guide](../setup/README.md) for details.

2. **Authentication**:
   - Obtain an access token or API key for authentication.

3. **Permissions**:
   - Ensure you have the necessary permissions to create roles.

4. **Role Details**:
   - Prepare the required role details, such as the role name and associated permissions.

## Configuration or Steps

### Method 1: Using the API

#### Step 1: Prepare the API Request

- **Endpoint**: `/api/v1/roles`
- **Method**: `POST`
- **Headers**:
  - `Authorization: Bearer <access_token>`
  - `Content-Type: application/json`

- **Request Body**:

  ```json
  {
    "name": "admin",
    "description": "Administrator role with full access",
    "permissions": ["read", "write", "delete"]
  }
  ```

#### Step 2: Send the Request

Use a tool like `curl`, Postman, or any HTTP client to send the request. Example using `curl`:

```bash
curl -X POST \
  https://<thunder-server-url>/api/v1/roles \
  -H "Authorization: Bearer <access_token>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "admin",
    "description": "Administrator role with full access",
    "permissions": ["read", "write", "delete"]
  }'
```

#### Step 3: Verify the Response

A successful response will return a `201 Created` status code with the created role's details:

```json
{
  "id": "67890",
  "name": "admin",
  "description": "Administrator role with full access",
  "permissions": ["read", "write", "delete"]
}
```

#### Step 4: Handle Errors

If the API returns an error, check the response for details. Common issues include:

- Missing required fields
- Invalid token
- Insufficient permissions

### Method 2: Using the Thunder Console

#### Step 1: Sign In
   - Open the Thunder console and sign in with an admin account.

#### Step 2: Navigate to Role Management
   - Go to the "Roles" section in the console.

#### Step 3: Create a New Role
   - Click the "Create Role" button.
   - Fill in the required fields (e.g., role name, description, permissions).
   - Click "Save" to create the role.

#### Step 4: Verify the Role
   - Ensure the new role appears in the role list.

#### Step 5: Edit or Delete Roles
   - Use the console to manage existing roles, including editing details or removing roles.

## Verify

After creating a role, verify the following:

1. The role appears in the role list (console) or the API response.
2. The role has the correct permissions assigned.
3. Users assigned to the role have the expected access levels.

## Troubleshooting

- **Issue**: API request fails with a `401 Unauthorized` error.
  - **Solution**: Verify the access token and ensure it hasn't expired.

- **Issue**: Role creation fails with a `400 Bad Request` error.
  - **Solution**: Check the request body for missing or invalid fields.

- **Issue**: Role permissions are not applied correctly.
  - **Solution**: Verify the permissions assigned to the role and ensure they are valid.

## Related References

- [API Documentation](../content/apis/role.yaml)
- [Thunder Console Guide](../guides/thunder-console-guide.md)
- [Setup Guide](../setup/README.md)