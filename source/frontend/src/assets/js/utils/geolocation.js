export function getCurrentPosition() {
  return new Promise((resolve, reject) => {
    if (!navigator.geolocation) {
      return reject(new Error('Geolocation not supported'));
    }

    navigator.geolocation.getCurrentPosition(
      pos => {
        resolve(pos);
      },
      err => {
        reject(err);
      },
      { enableHighAccuracy: false, maximumAge: 60000 }
    );
  });
}

export function geocodeNominatim(q) {
  const url = `https://nominatim.openstreetmap.org/search?format=json&q=${encodeURIComponent(q)}&limit=5`;
  return fetch(url, { headers: { 'Accept-Language': 'en' } }).then(r => r.json());
}