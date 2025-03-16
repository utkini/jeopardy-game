# Jeopardy Game

A web-based implementation of the classic Jeopardy game show, built with Go and modern web technologies. This application allows you to host your own Jeopardy-style quiz game with customizable questions, multiple players, and multimedia support.

## Features

- **Interactive Game Board**: Classic Jeopardy-style board with categories and point values
- **Turn-Based Gameplay**: Players take turns answering questions in a circular order
- **Multimedia Questions**: Support for various types of media in questions:
  - Text-based questions
  - Image questions
  - YouTube video questions
  - Audio questions
- **Real-time Updates**: Using HTMX for smooth, dynamic updates without page reloads
- **Responsive Design**: Works on both desktop and mobile devices
- **Score Tracking**: Automatic tracking of player scores
- **Answer Validation**: Easy-to-use correct/incorrect buttons for scoring
- **Custom Categories**: Configurable questions and categories via YAML files

## Technology Stack

- Backend: Go
- Frontend: HTML, CSS (Bulma), JavaScript
- Dynamic Updates: HTMX
- Styling: Bulma CSS Framework
- Icons: Font Awesome
- Media Support: YouTube IFrame API
- Containerization: Docker

## Setup

### Local Setup

1. Configure your questions in `configs/questions.yaml`:
   ```yaml
   categories:
     - name: "Your Category"
       questions:
         - points: 100
           question: "Your question"
           answer: "The answer"
           # Media configuration:
           mediaType: "video"  # Options: "video", "image", "audio", or "" for no media
           mediaURL: "https://youtube.com/..."  # URL format depends on mediaType
   ```

   Media URL formats:
   - For YouTube videos: Use the standard YouTube URL (e.g., "https://youtube.com/watch?v=...")
   - For images: Direct URL to the image file (JPG, PNG, GIF, WebP)
   - For audio: Direct URL to the MP3 file
   - For no media: Leave mediaType and mediaURL empty

2. Run the application:
   ```shell
   go run cmd/main.go
   ```

### Docker Setup

1. Using Docker Compose (recommended):
   ```shell
   docker-compose up -d
   ```

2. Using Docker directly:
   ```shell
   # Build the image
   docker build -t jeopardy-game .

   # Run the container
   docker run -d -p 8080:8080 -v $(pwd)/configs:/app/configs:ro jeopardy-game
   ```

3. Access the application:
   Open your browser and navigate to `http://localhost:8080`

## Game Rules

1. Add players at the start of the game
2. Players take turns selecting questions from the board
3. When a question is selected:
   - The question and any associated media are displayed
   - The current player can attempt to answer
   - Use the "Correct" or "Incorrect" button to score the answer
   - Points are automatically added/subtracted from the player's score
4. The game continues until all questions have been answered

## Configuration

The game can be configured through two main files:

- `configs/base.yaml`: Basic application settings
- `configs/questions.yaml`: Game questions and categories

When running with Docker, you can modify these files in your local directory, and they will be automatically available to the container through volume mounting.

### Media Support Details

1. **Images**:
   - Supported formats: JPG, PNG, GIF, WebP
   - Images are displayed responsively (max-width: 100%, maintaining aspect ratio)
   - Example:
     ```yaml
     mediaType: "image"
     mediaURL: "https://example.com/path/to/image.jpg"
     ```

2. **Audio**:
   - Supported format: MP3
   - Displays with standard HTML5 audio controls
   - Example:
     ```yaml
     mediaType: "audio"
     mediaURL: "https://example.com/path/to/audio.mp3"
     ```

3. **YouTube Videos**:
   - Use standard YouTube URLs
   - Videos are embedded responsively (16:9 aspect ratio)
   - Example:
     ```yaml
     mediaType: "video"
     mediaURL: "https://youtube.com/watch?v=VIDEO_ID"
     ```

## Contributing

Feel free to submit issues, fork the repository, and create pull requests for any improvements.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
