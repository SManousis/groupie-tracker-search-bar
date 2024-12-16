# Groupie Tracker

Groupie Tracker is a web application that retrieves data from an external API containing information about music artists, their concert dates, locations, and relations. The app displays this data on a user-friendly website and allows users to interact with it through dynamic visualizations and event-based actions.

## Features

- Displays information about artists, including their name, image, members, first album, and creation date.
- Shows the locations and dates of upcoming and past concerts for each artist.
- Visualizes relationships between artists, concert dates, and locations.
- Provides an interactive web interface with data visualizations, such as blocks, cards, and lists.
- Handles client-server communication with real-time updates of concert data.

## Tech Stack

- **Backend**: Go (Golang)
- **Frontend**: HTML, CSS, JavaScript
- **API**: RESTful API with data on artists, concert dates, locations, and relations.
- **Testing**: Unit tests in Go (using `testing` package)

## Project Structure

    project/ │ ├── main.go // Go backend code for fetching and serving data ├── main_test.go // Unit tests for backend functionality ├── templates/ // Template files for frontend │ ├── home.html // Main HTML file for rendering artist data │ └── error.html // Error page (if necessary) │ └── static/ // Static assets (CSS, JS, images) ├── assets/ ├── app.js // Frontend JavaScript for rendering dynamic content └── styles.css // Stylesheet for the website


## Prerequisites

- [Go](https://golang.org/dl/) (1.18 or later) installed.
- A web browser to view the frontend.
- Internet connection to fetch data from the external API.

## Setup

1. Clone the repository to your local machine:
   
    git clone https://platform.zone01.gr/git/smanousi/groupie-tracker
    cd groupie-tracker
2. Install Go dependencies (if any). There are no external dependencies in this project as it uses only the standard Go packages.


3. Run the Go server:

    go run main.go

4. The server will start on the specified port (default: 8080). Open your web browser and visit:

    http://localhost:8080

    The application will fetch and display artist information, concert dates, and locations dynamically.

## Features and Interactions

    Dynamic Artist Data:
        The application fetches artist data from an external API and dynamically populates the frontend with the data.
        Artists' names, images, members, first albums, and creation dates are displayed in cards on the homepage.

    Concert Locations and Dates:
        Each artist card includes details about their concert locations and dates, fetched from the API.
        This data is displayed as lists on the artist cards.

    Event Handling:
        Data is loaded dynamically upon visiting the homepage, and events (like clicking or scrolling) trigger further actions, such as fetching more data or filtering the displayed results.

    Frontend Visualizations:
        The artist data is visually presented as cards, which can include additional details such as concert dates and locations.
        You can customize how you display data by modifying the frontend JavaScript and HTML templates.

## Testing

    The project includes unit tests to ensure the backend is functioning correctly. You can run the tests with the go test command.

    To run the tests, execute:

    go test -v

    This will run the tests in main_test.go and display the results.

    The tests include checking the functionality of the following:
        Fetching and decoding artist data from the external API.
        Fetching and processing concert locations, dates, and relations.
        Handling errors during the data-fetching process.

## API

    The application interacts with the following API endpoints to fetch data:

    GET /api/artists: Returns a list of artists with their details.
    GET /api/locations/{id}: Returns the concert locations for a specific artist.
    GET /api/dates/{id}: Returns the concert dates for a specific artist.
    GET /api/relation/{id}: Returns the relations between an artist, their concert dates, and locations.

## Deployment

    This project is deployed on Render. You can visit the live website here:

    https://groupie-tracker-gsdt.onrender.com/

## Deploy to Your Own Server

    To deploy this project to your own server, follow these steps:

        Ensure you have Go installed on your server.
        Clone the repository to the server.
        Set up a web server to serve the Go application and static assets (CSS, JS).
        Point the server to the appropriate port and ensure that it is accessible.

## License

    This project is licensed under the MIT License 
    Project members: 
    - Avgoustinos Andris
    - Kimon Dimitriou
    - Stamatis Manouis
