'use strict';

const lumo = require('lumo');
const Stats = require('stats.js');

function generatePoint2DBuffer(numPoints, radius) {
	const buffer = new Float32Array(3 * numPoints);
	for (let i=0; i<numPoints; i++) {
		buffer[i*3] = Math.random() * 256; // x
		buffer[i*3+1] = Math.random() * 256; // y
		buffer[i*3+2] = Math.floor(Math.random() * radius) + 2; // radius
	}
	return buffer;
};

function addGeoLayer(plot) {
	const layer = new lumo.TileLayer({
		renderer: new lumo.ImageTileRenderer()
	});
	layer.requestTile = (coord, done) => {
		const SUBDOMAINS = [
			'a', 'b', 'c'
		];
		const s = SUBDOMAINS[(coord.x + coord.y + coord.z) % SUBDOMAINS.length];
		const url = `http://${s}.basemaps.cartocdn.com/light_nolabels/${coord.xyz()}.png`;
		lumo.loadImage(url, done);
	};
	plot.add(layer);
}

function addMandlebrotLayer(plot) {
	const layer = new lumo.TileLayer({
		renderer: new lumo.ImageTileRenderer(),
		opacity: 1.0
	});
	layer.requestTile = (coord, done) => {
		const url = `mandelbrot/${coord.tms()}`;
		lumo.loadBuffer(url, (err, buffer) => {
			if (err) {
				done(err);
				return;
			}
			done(null, new Uint8Array(buffer));
		});
	};
	plot.add(layer);
}

function addLineOverlay(plot) {
	const overlay = new lumo.PolylineOverlay({
		opacity: 0.5,
		renderer: new lumo.PolylineOverlayRenderer({
			lineWidth: 2
		})
	});
	for (let j=0; j<20; j++) {
		const points = [{ x: Math.random(), y: Math.random() }];
		const n = Math.floor(Math.random()*10) + 2;
		for (let i=1; i<n; i++) {
			points.push({
				x: points[i-1].x + ((Math.random() * 2.0) - 1) * 0.2,
				y: points[i-1].y + ((Math.random() * 2.0) - 1) * 0.2
			});
		}
		overlay.addPolyline(`a${j}`, points);
	}
	plot.add(overlay);
}

function addPointOverlay(plot) {
	const overlay = new lumo.PointOverlay({
		opacity: 0.5,
		pointColor: [ 0.0, 1.0, 1.0, 1.0 ],
		renderer: new lumo.PointOverlayRenderer()
	});
	for (let j=0; j<20; j++) {
		const points = [{ x: Math.random(), y: Math.random() }];
		const n = Math.floor(Math.random()*10) + 2;
		for (let i=1; i<n; i++) {
			points.push({
				x: points[i-1].x + ((Math.random() * 2.0) - 1) * 0.2,
				y: points[i-1].y + ((Math.random() * 2.0) - 1) * 0.2
			});
		}
		overlay.addPoints(`a${j}`, points);
	}
	plot.add(overlay);
}

function addPolygonOverlay(plot) {
	const overlay = new lumo.PolygonOverlay({
		opacity: 0.5,
		renderer: new lumo.PolygonOverlayRenderer()
	});
	overlay.addPolygon('test', [
		{ x: -1.67, y: -0.33 },
		{ x: -1.00, y: -0.33 },
		{ x: -0.33, y: -0.33 },
		{ x: -0.33, y: -0.67 },
		{ x: -0.33, y: -1.33 },
		{ x: -0.33, y: -1.67 },
		{ x: 0.33, y: -1.67 },
		{ x: 0.33, y: -0.33 },
		{ x: 1.67, y: -0.33 },
		{ x: 1.67, y: 0.33 },
		{ x: 0.33, y: 0.33 },
		{ x: 0.33, y: 1.67 },
		{ x: -0.33, y: 1.67 },
		{ x: -0.33, y: 0.33 },
		{ x: -0.67, y: 0.33 },
		{ x: -1.67, y: 0.33 }
	].map(p => {
		return {
			x: p.x + (Math.random() * 2 - 1) * 0.05,
			y: p.y + (Math.random() * 2 - 1) * 0.05
		};
	}));
	plot.add(overlay);
}

function addRedInteractivePointOverlay(plot) {
	const layer = new lumo.TileLayer({
		zIndex: 1,
		renderer: new lumo.InteractiveTileRenderer({
			color: [ 1.0, 0.2, 0.2, 0.8 ]
		})
	});
	layer.requestTile = (coord, done) => {
		setTimeout(() => {
			done(null, generatePoint2DBuffer(32, 16));
		}, Math.random() * 400);
	};
	plot.add(layer);
}

function addBlueInteractivePointOverlay(plot) {
	const layer = new lumo.TileLayer({
		zIndex: 2,
		renderer: new lumo.InteractiveTileRenderer({
			color: [ 0.2, 0.2, 0.8, 0.8 ]
		})
	});
	layer.requestTile = (coord, done) => {
		setTimeout(() => {
			done(null, generatePoint2DBuffer(4, 32));
		}, Math.random() * 400);
	};
	plot.add(layer);
}

function addStarLayer(plot) {
	const layer = new lumo.TileLayer({
		renderer: new lumo.InstancedTileRenderer({
			color: [ 0.8, 0.8, 0.2, 0.8 ]
		}),
		cacheSize: 256
	});
	layer.requestTile = (coord, done) => {
		setTimeout(() => {
			done(null, generatePoint2DBuffer(16, 16));
		}, Math.random() * 400);
	};
	plot.add(layer);
}

window.start = function() {

	const plot = new lumo.Plot('#plot', {
		continuousZoom: false,
		inertia: true,
		wraparound: false,
		zoom: 3,
		dirtyChecking: false
	});

	// Mandlebrot Fractal Layer
	addMandlebrotLayer(plot);

	// CartoDB Basemap Layer
	// addGeoLayer(plot);

	// Overlays
	// addPointOverlay(plot);
	// addLineOverlay(plot);
	// addPolygonOverlay(plot);

	// Interactive Point Layers
	// addRedInteractivePointOverlay(plot);
	// addBlueInteractivePointOverlay(plot);

	// Star Shape Layer
	// addStarLayer(plot);

	// Debug performance tracking

	const stats = new Stats();
	document.body.appendChild(stats.dom);
	stats.begin();
	plot.on('frame', () => {
		stats.end();
		stats.begin();
	});

};
