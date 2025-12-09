window.addAwesomeMarker = (mapDiv, lat, lng, options) => {
    const map = mapDiv._leaflet_map;
    if (!map) return;

    const icon = L.AwesomeMarkers.icon({
        icon: options.icon || 'home',        // FontAwesome icon
        prefix: options.prefix || 'fa',      // 'fa' for FontAwesome
        markerColor: options.markerColor || 'blue',
        iconColor: options.iconColor || 'white',
        extraClasses: options.extraClasses || ''
    });

    const marker = L.marker([lat, lng], { icon }).addTo(map);
    return marker;
};
