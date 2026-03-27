const CACHE_NAME = 'tokkatot-v3';
const ASSETS_TO_CACHE = [
  '/',
  '/index.html',
  '/css/theme.css',
  '/css/components.css',
  '/css/layout.css',
  '/css/dashboard.css',
  '/css/monitoring.css',
  '/css/schedules.css',
  '/css/settings.css',
  '/css/auth.css',
  '/js/utils/i18n.js',
  '/js/utils/api.js',
  '/js/utils/components.js',
  '/js/index.js',
  '/assets/images/tokkatot logo-02.png',
  '/assets/images/farmer-avatar.png',
  '/assets/images/viewer-avatar.png',
  '/assets/images/admin-avatar.png',
];

// External CDN assets — cache on first use
const CDN_ASSETS = [
  'https://unpkg.com/vue@3/dist/vue.global.prod.js',
  'https://fonts.googleapis.com/css2?family=Material+Symbols+Outlined:opsz,wght,FILL,GRAD@20..48,100..700,0..1,-50..200',
  'https://fonts.googleapis.com/css2?family=Noto+Sans+Khmer:wght@400;500;600;700&display=swap',
];

// =====================================================================
// INSTALL: Pre-cache all app shell assets + CDN dependencies
// =====================================================================
self.addEventListener('install', event => {
  event.waitUntil(
    caches.open(CACHE_NAME).then(cache => {
      // Cache local assets synchronously
      return cache.addAll(ASSETS_TO_CACHE).then(() => {
        // Cache CDN assets — ignore individual failures so install still succeeds
        return Promise.allSettled(CDN_ASSETS.map(url => cache.add(url)));
      });
    }).then(() => self.skipWaiting())
  );
});

// =====================================================================
// ACTIVATE: Remove any old caches from previous versions
// =====================================================================
self.addEventListener('activate', event => {
  event.waitUntil(
    caches.keys().then(cacheNames => {
      return Promise.all(
        cacheNames
          .filter(name => name !== CACHE_NAME)
          .map(name => caches.delete(name))
      );
    }).then(() => self.clients.claim())
  );
});

// =====================================================================
// FETCH: Network-first for API, Cache-first for assets
// =====================================================================
self.addEventListener('fetch', event => {
  if (event.request.method !== 'GET') return;

  // API calls: always network, never cache
  if (event.request.url.includes('/v1/')) return;

  // For CDN and local assets: try network, fallback to cache
  event.respondWith(
    fetch(event.request)
      .then(response => {
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

// =====================================================================
// PUSH: Show notification when a push message is received from the server
// =====================================================================
self.addEventListener('push', event => {
  let data = {
    title: 'Tokkatot',
    body: 'New alert from your farm.',
    url: '/',
    icon: '/assets/images/tokkatot logo-02.png',
    badge: '/assets/images/tokkatot logo-02.png',
  };

  if (event.data) {
    try {
      const parsed = event.data.json();
      data = { ...data, ...parsed };
    } catch (e) {
      data.body = event.data.text();
    }
  }

  const options = {
    body: data.body,
    icon: data.icon,
    badge: data.badge,
    data: { url: data.url },
    // Vibration pattern for elderly farmers — two firm pulses
    vibrate: [300, 100, 300],
    requireInteraction: true, // Stay visible until farmer taps it
    actions: [
      { action: 'view', title: 'View Alert' },
      { action: 'dismiss', title: 'Dismiss' },
    ],
  };

  event.waitUntil(
    self.registration.showNotification(data.title, options)
  );
});

// =====================================================================
// NOTIFICATION CLICK: Navigate to the relevant page when farmer taps
// =====================================================================
self.addEventListener('notificationclick', event => {
  event.notification.close();

  if (event.action === 'dismiss') return;

  const targetUrl = (event.notification.data && event.notification.data.url)
    ? event.notification.data.url
    : '/alerts';

  event.waitUntil(
    clients.matchAll({ type: 'window', includeUncontrolled: true }).then(clientList => {
      // If app is already open, focus it
      for (const client of clientList) {
        if (client.url.includes(self.location.origin) && 'focus' in client) {
          client.navigate(targetUrl);
          return client.focus();
        }
      }
      // Otherwise open a new window
      if (clients.openWindow) {
        return clients.openWindow(targetUrl);
      }
    })
  );
});
