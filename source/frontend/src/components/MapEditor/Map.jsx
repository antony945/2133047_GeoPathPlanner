import { MapContainer, TileLayer } from 'react-leaflet';


import './Map.css';
import 'leaflet/dist/leaflet.css';
import 'leaflet/dist/leaflet.js';

function Map() {

    return (
        <MapContainer center={[51.505, -0.09]} zoom={13} id="geomap">
            <TileLayer url="https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png" />
        </MapContainer>     
    )
}

export default Map