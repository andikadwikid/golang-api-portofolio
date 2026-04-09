# Postman Testing Guide - Portfolio API

This guide provides the details needed to test the Portfolio API using Postman.

## 1. Setup Environment
In Postman, create a new **Environment** and add the following variables:
- `base_url`: `http://localhost:8081`
- `token`: (Leave empty, will be filled after Login)

## 2. Authentication
All protected routes require a **Bearer Token** in the `Authorization` header.

### Register User
- **Method**: `POST`
- **URL**: `{{base_url}}/users/register`
- **Body** (JSON):
```json
{
    "name": "Andika Dwi",
    "username": "andika",
    "email": "andika@example.com",
    "password": "password123"
}
```

### Login User
- **Method**: `POST`
- **URL**: `{{base_url}}/users/login`
- **Body** (JSON):
```json
{
    "email": "andika@example.com",
    "password": "password123"
}
```
> **Tip**: Copy the `token` from the response and paste it into your `token` environment variable.

---

## 3. Portofolio API

### Create Portofolio
- **Method**: `POST`
- **URL**: `{{base_url}}/portofolio/`
- **Auth**: Bearer Token `{{token}}`
- **Body** (JSON):
```json
{
    "name": "My Web Development Portfolio"
}
```

### Get My Portofolios
- **Method**: `GET`
- **URL**: `{{base_url}}/portofolio/`
- **Auth**: Bearer Token `{{token}}`

### Update Portofolio
- **Method**: `PUT`
- **URL**: `{{base_url}}/portofolio/:id`
- **Auth**: Bearer Token `{{token}}`
- **Body** (JSON):
```json
{
    "name": "Updated Portfolio Name"
}
```

### Toggle Portofolio Status
- **Method**: `PATCH`
- **URL**: `{{base_url}}/portofolio/:id/status`
- **Body** (JSON):
```json
{
    "is_active": false
}
```

---

## 4. Project API

### Create Project
- **Method**: `POST`
- **URL**: `{{base_url}}/project/`
- **Auth**: Bearer Token `{{token}}`
- **Body** (JSON):
```json
{
    "name": "E-Commerce Website",
    "description": "Built with Next.js and Go",
    "link": "https://github.com/example/ecommerce",
    "portofolio_id": "PASTE_PORTOFOLIO_ID_HERE"
}
```

### Get Projects by Portofolio ID
- **Method**: `GET`
- **URL**: `{{base_url}}/project/portofolio/:portofolio_id`

### Update Project
- **Method**: `PUT`
- **URL**: `{{base_url}}/project/:project_id`
- **Body** (JSON):
```json
{
    "name": "Updated Project Name",
    "description": "Updated Description",
    "link": "https://new-link.com"
}
```

---

## 5. Project Image API

### Upload Multiple Images (Multipart Form)
- **Method**: `POST`
- **URL**: `{{base_url}}/project/:project_id/images`
- **Auth**: Bearer Token `{{token}}`
- **Body** (form-data):
    - `images`: (Select multiple files in Postman)
    - `project_id`: (Paste Project ID)

### Get My Project Images
- **Method**: `GET`
- **URL**: `{{base_url}}/project/images/my`
- **Auth**: Bearer Token `{{token}}`

### Update Project Image (Multipart Form)
- **Method**: `PUT`
- **URL**: `{{base_url}}/project/images/:image_id`
- **Auth**: Bearer Token `{{token}}`
- **Body** (form-data):
    - `image`: (Select single file)
    - `project_id`: (Optional - to move image to another project)

### Delete Project Image
- **Method**: `DELETE`
- **URL**: `{{base_url}}/project/images/:image_id`
- **Auth**: Bearer Token `{{token}}`
