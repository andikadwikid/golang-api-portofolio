# API Documentation - cURL Commands

Dokumentasi ini berisi daftar perintah `curl` untuk menguji semua endpoint API Portfolio.
Ganti `<TOKEN>` dengan JWT Token yang didapat setelah login, dan `<ID>` dengan ID dokumen yang relevan (seperti `user_id`, `portfolio_id`, dll).

**Base URL**: `http://localhost:8081`

---

## 1. Authentication & Users

### Health Check
```bash
curl -X GET http://localhost:8081/users/health
```

### Register User
```bash
curl -X POST http://localhost:8081/users/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Full Name",
    "username": "username123",
    "email": "user@example.com",
    "password": "password123"
  }'
```

### Login User
```bash
curl -X POST http://localhost:8081/users/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'
```

### Get All Users (Auth Required)
```bash
curl -X GET http://localhost:8081/users/ \
  -H "Authorization: Bearer <TOKEN>"
```

### Update User (Auth Required)
```bash
curl -X PUT http://localhost:8081/users/<USER_ID> \
  -H "Authorization: Bearer <TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Updated Name",
    "username": "updated_user",
    "email": "updated@example.com",
    "avatar": "https://link-to-avatar.com",
    "bio": "Software Engineer"
  }'
```

---

## 2. Social Media (Master Data)

### Create Social Media (Auth Required)
```bash
curl -X POST http://localhost:8081/social-media/ \
  -H "Authorization: Bearer <TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Instagram",
    "icon": "instagram-icon.png"
  }'
```

### Get All Social Media (Auth Required)
```bash
curl -X GET http://localhost:8081/social-media/ \
  -H "Authorization: Bearer <TOKEN>"
```

---

## 3. Social Media User (User's Links)

### Link User to Social Media (Auth Required)
```bash
curl -X POST http://localhost:8081/social-media-user/ \
  -H "Authorization: Bearer <TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{
    "link": "https://instagram.com/myusername",
    "social_media_id": "<SOCIAL_MEDIA_ID>"
  }'
```

### Get My Social Media Links (Auth Required)
```bash
curl -X GET http://localhost:8081/social-media-user/me \
  -H "Authorization: Bearer <TOKEN>"
```

---

## 4. Portfolio

### Create Portfolio (Auth Required)
```bash
curl -X POST http://localhost:8081/portofolio/ \
  -H "Authorization: Bearer <TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Web Development Portfolio"
  }'
```

### Get My Portfolios (Auth Required)
```bash
curl -X GET http://localhost:8081/portofolio/ \
  -H "Authorization: Bearer <TOKEN>"
```

### Get All Active Portfolios (Public)
```bash
curl -X GET http://localhost:8081/portofolio/all
```

### Update Portfolio Status (Auth Required)
```bash
curl -X PATCH http://localhost:8081/portofolio/<PORTFOLIO_ID>/status \
  -H "Authorization: Bearer <TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{
    "is_active": false
  }'
```

---

## 5. Project

### Create Project inside Portfolio (Auth Required)
```bash
curl -X POST http://localhost:8081/project/ \
  -H "Authorization: Bearer <TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "E-Commerce App",
    "description": "Building a fullstack e-commerce",
    "link": "https://github.com/myproject",
    "portofolio_id": "<PORTFOLIO_ID>"
  }'
```

### Get Projects by Portfolio ID
```bash
curl -X GET http://localhost:8081/project/portofolio/<PORTFOLIO_ID>
```

---

## 6. Project Images

### Upload Multiple Images (Auth Required)
```bash
curl -X POST http://localhost:8081/project/<PROJECT_ID>/images \
  -H "Authorization: Bearer <TOKEN>" \
  -F "images=@/path/to/image1.png" \
  -F "images=@/path/to/image2.jpg"
```

### Get Images by Project ID
```bash
curl -X GET http://localhost:8081/project/<PROJECT_ID>/images
```

### Update Single Image (Auth Required)
```bash
curl -X PUT http://localhost:8081/project/images/<IMAGE_ID> \
  -H "Authorization: Bearer <TOKEN>" \
  -F "image=@/path/to/new_image.png"
```

### Update Image Status (Auth Required)
```bash
curl -X PATCH http://localhost:8081/project/images/<IMAGE_ID>/status \
  -H "Authorization: Bearer <TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{
    "is_active": true
  }'
```

### Delete Image (Auth Required)
```bash
curl -X DELETE http://localhost:8081/project/images/<IMAGE_ID> \
  -H "Authorization: Bearer <TOKEN>"
```
