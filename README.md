# Manga Backend API Documentation

## Base URL
`http://localhost:8080` (assuming default port 8080)

## Authentication
Authentication details are not provided in the given snippets. Implement appropriate authentication mechanisms (e.g., JWT) for protected routes.

## Endpoints

### 1. Manga

#### Get All Mangas
- **URL**: `/mangas`
- **Method**: GET
- **Description**: Retrieves a list of all mangas.
- **Response**:
    - Status Code: 200 OK
    - Body: Array of manga objects
  ```json
  [
    {
      "id": 1,
      "title": "Manga Title",
      "description": "Manga description",
      "author": "Author Name",
      "coverImage": "http://example.com/cover.jpg",
      "status": "Ongoing",
      "tags": ["Action", "Adventure"]
    },
    // ... more manga objects
  ]
  ```

### 2. Website

#### Add Website
- **URL**: `/websites`
- **Method**: POST
- **Description**: Adds a new website to scrape manga from.
- **Request Body**:
  ```json
  {
    "url": "https://example-manga-site.com",
    "name": "Example Manga Site"
  }
  ```
- **Response**:
    - Status Code: 200 OK
    - Body:
      ```json
      {
        "message": "website added"
      }
      ```

### 3. Chapters (Assumed based on repository structure)

#### Add Chapter
- **URL**: `/chapters`
- **Method**: POST
- **Description**: Adds a new chapter to a manga.
- **Request Body**:
  ```json
  {
    "mangaId": 1,
    "chapterNumber": "1.0",
    "title": "Chapter Title",
    "url": "https://example-manga-site.com/manga/1/chapter/1"
  }
  ```
- **Response**:
    - Status Code: 201 Created
    - Body: Created chapter object

### 4. Tags (Assumed based on repository structure)

#### Get All Tags
- **URL**: `/tags`
- **Method**: GET
- **Description**: Retrieves all available tags.
- **Response**:
    - Status Code: 200 OK
    - Body: Array of tag objects
  ```json
  [
    {
      "id": 1,
      "name": "Action"
    },
    {
      "id": 2,
      "name": "Adventure"
    },
    // ... more tag objects
  ]
  ```

## Error Responses

All endpoints may return the following error response:

- Status Code: 400 Bad Request or 500 Internal Server Error
- Body:
  ```json
  {
    "error": "Error message describing the issue"
  }
  ```

## Notes for Frontend Development

1. **Environment Configuration**:
    - Use the `SERVER_ADDRESS` from the backend's `.env` file to set up your API base URL in the frontend.

2. **Pagination**:
    - The current API doesn't show pagination. Consider implementing pagination for large datasets, especially for the `/mangas` endpoint.

3. **Authentication**:
    - Implement user authentication and handle token storage/management in the frontend.

4. **Error Handling**:
    - Create a global error handling mechanism to process and display error messages from the API.

5. **Real-time Updates**:
    - If real-time features are added (e.g., notifications for new chapters), consider implementing WebSocket connections.

6. **Image Handling**:
    - Ensure proper handling and caching of manga cover images and chapter images.

7. **State Management**:
    - Use a state management solution (e.g., Redux, MobX, or React Query) to manage the application state, especially for caching manga and chapter data.

8. **Responsive Design**:
    - Design your frontend to be responsive, considering various device sizes for optimal manga reading experience.

This API documentation provides a starting point for frontend development. As you expand your backend functionality, remember to update this documentation accordingly. You may also want to consider using tools like Swagger or OpenAPI for more comprehensive and interactive API documentation.