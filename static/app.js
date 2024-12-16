// Function to dynamically create cards for each artist
function createCards(artists) {
  const bandContainer = document.getElementById('band-container');
  bandContainer.innerHTML = ''; // Clear any existing content

  artists.forEach((artist) => {
    // Ensure artist has valid data for rendering
    if (!artist.name || !artist.image || !artist.creationDate || !artist.firstAlbum || !artist.members) {
      console.warn(`Missing essential data for artist: ${artist.name}`);
      return;
    }

    const card = document.createElement('div');
    card.className = 'card';
    card.innerHTML = `
      <img src="${artist.image}" alt="${artist.name}">
      <div class="name">${artist.name}</div>
      <div class="year">${artist.creationDate}</div>
    `;
    card.addEventListener('click', () => expandCard(artist));
    bandContainer.appendChild(card);
  });
}

// Function to expand a card and show artist details
function expandCard(artist) {
  if (!artist.dates || !artist.locations || !artist.events) {
    console.warn(`Incomplete data for artist: ${artist.name}`);
    return;
  }

  // Log events to check its structure
  console.log("Artist Events:", artist.events);

  // Create an overlay to dim the background
  const overlay = document.createElement('div');
  overlay.className = 'overlay';
  overlay.addEventListener('click', closeExpandedView);

  // Create the expanded card
  const expandedCard = document.createElement('div');
  expandedCard.className = 'expanded-card';
  expandedCard.innerHTML = `
    <div class="sidebar-left" style="background-image: url('${artist.image}');"></div>
    <div class="content">
      <h2>${artist.name}</h2>
      <p><strong>Year of Formation:</strong> ${artist.creationDate}</p>
      <p><strong>First Album:</strong> ${artist.firstAlbum}</p>
      <p><strong>Members:</strong> ${artist.members.join(', ')}</p>
      <h3>Upcoming Events</h3>
      <ul>
        ${Object.entries(artist.events).map(([event, location]) => `<li><strong>${event}:</strong> ${location}</li>`).join('')}
      </ul>
    </div>
    <div class="sidebar-right" style="background-image: url('${artist.image}');"></div>
    <button class="close-btn">Close</button>
  `;

  // Add close button functionality
  expandedCard.querySelector('.close-btn').addEventListener('click', closeExpandedView);

  // Append overlay and expanded card to the document
  document.body.appendChild(overlay);
  document.body.appendChild(expandedCard);

  // Add animation to expand the card
  setTimeout(() => expandedCard.classList.add('expanded'), 10); // Add animation after a slight delay
}

// Function to close the expanded view
function closeExpandedView() {
  const overlay = document.querySelector('.overlay');
  const expandedCard = document.querySelector('.expanded-card');
  if (overlay) overlay.remove();
  if (expandedCard) expandedCard.remove();
}

// Initialize the app
function init() {
  // artistData is embedded in the HTML by Go
  if (window.artistData && Array.isArray(window.artistData)) {
    createCards(window.artistData);
  } else {
    console.error('Artist data is not available.');
  }
}

// Run the app on load
init();