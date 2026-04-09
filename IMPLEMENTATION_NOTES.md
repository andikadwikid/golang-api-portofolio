# Implementation Notes - Portfolio API

This document provides a technical overview of the implementation patterns used in this project, specifically for those coming from a SQL background.

## 1. Database Pattern: Soft Delete
Unlike hard deletes (`DELETE FROM table`), this project uses **Soft Deletes**.
- **Field**: `is_deleted` (boolean) and `deleted_at` (timestamp).
- **Behavior**: When a user "deletes" a record, we only set `is_deleted` to `true`.
- **Querying**: Every `Find` or `FindOne` operation should include `{ "is_deleted": false }` in the filter to ensure deleted data is hidden.

## 2. MongoDB Aggregation (The "JOIN" Equivalent)
Since MongoDB is a NoSQL database, we use the **Aggregation Pipeline** to perform operations that involve multiple collections.

### Example: Get My Project Images
To find a project image belonging to a user, we must navigate the relationship: `Project Image` -> `Project` -> `Portofolio` -> `User`.

```go
pipeline := mongo.Pipeline{
    // Step 1: JOIN project_image with project
    {{Key: "$lookup", Value: bson.D{
        {Key: "from", Value: "project"},
        {Key: "localField", Value: "project_id"},
        {Key: "foreignField", Value: "_id"},
        {Key: "as", Value: "project"},
    }}},
    // Step 2: Flatten the 'project' array into an object
    {{Key: "$unwind", Value: "$project"}},
    
    // Step 3: JOIN project with portofolio
    {{Key: "$lookup", Value: bson.D{
        {Key: "from", Value: "portofolio"},
        {Key: "localField", Value: "project.portofolio_id"},
        {Key: "foreignField", Value: "_id"},
        {Key: "as", Value: "portofolio"},
    }}},
    // Step 4: Flatten the 'portofolio' array
    {{Key: "$unwind", Value: "$portofolio"}},
    
    // Step 5: WHERE filter
    {{Key: "$match", Value: bson.D{
        {Key: "portofolio.user_id", Value: userID},
        {Key: "is_deleted", Value: false},
    }}},
    
    // Step 6: SELECT (clean up temporary join fields)
    {{Key: "$project", Value: bson.D{
        {Key: "project", Value: 0},
        {Key: "portofolio", Value: 0},
    }}},
}
```

### SQL Comparison
The above pipeline is logically equivalent to this SQL:
```sql
SELECT pi.*
FROM project_image pi
JOIN project p ON pi.project_id = p._id
JOIN portofolio pr ON p.portofolio_id = pr._id
WHERE pr.user_id = 'CURRENT_USER_ID'
  AND pi.is_deleted = false;
```

## 3. Security: Ownership Validation
Before any **Update** or **Delete** operation, the API performs a nested check:
1. Find the target object (e.g., `Project`).
2. Fetch its parent (e.g., `Portofolio`).
3. Compare the parent's `user_id` with the `user_id` extracted from the JWT token.
4. Only proceed if they match.

## 4. Static Files
Uploaded images are stored in `./public/uploads/` and served via the `/public` static route configured in `main.go`:
```go
r.Static("/public", "./public")
```
Access URL: `http://<host>:<port>/public/uploads/<category>/<filename>`
