const CACHE_NAME = 'tokkatot-v2';
const ASSETS_TO_CACHE = [
  '/',
  '/index.html',
  '/css/theme.css',
  '/css/components.css',
  '/css/layout.css',
  '/css/dashboard.css',
  '/js/utils/i18n.js',
  '/js/utils/api.js',
  '/js/utils/components.js',
  '/js/index.js',
  '/assets/images/tokkatot logo-02.png',
  '/assets/images/farmer-avatar.png',
  '/assets/images/viewer-avatar.png',
  '/assets/images/admin-avatar.png',
  'https://unpkg.com/vue@3/dist/vue.global.prod.js',
  'https://fonts.googleapis.com/css2?family=Material+Symbols+Outlined:opsz,wght,FILL,GRAD@20..48,100..700,0..1,-50..200',
  'https://fonts.googleapis.com/css2?family=Noto+Sans+Khmer:wght@400;500;600;700&display=swap'
];

// Install event: cache core assets
self.addEventListener('install', event => {
  event.waitUntil(
    caches.open(CACHE_NAME)
      .then(cache => cache.addAll(ASSETS_TO_CACHE))
      .then(() => self.skipWaiting())
  );
});

// Activate event: cleanup old caches
self.addEventListener('activate', event => {
  event.waitUntil(
    caches.keys().then(cacheNames => {
      return Promise.all(
        cacheNames.filter(name => name !== CACHE_NAME)
          .map(name => caches.delete(name))
      );
    }).then(() => self.clients.claim())
  );
});

// Fetch event: network first, then cache
self.addEventListener('fetch', event => {
  // Only cache GET requests
  if (event.request.method !== 'GET') return;

  // For API calls, always go to network (no cache)
  if (event.request.url.includes('/v1/')) {
    return;
  }

  event.respondWith(
    fetch(event.request)
      .then(response => {
        // Cache valid responses
        if (response && response.status === 200) {
          const responseToCache = response.clone();
          caches.open(CACHE_NAME).then(cache => {
            cache.put(event.request, responseToCache);
          });
        }
        return response;
      })
      .catch(() => caches.match(event.request))
  );
});
