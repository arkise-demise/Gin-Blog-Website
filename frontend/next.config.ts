/** @type {import('next').NextConfig} */
const nextConfig = {
  images: {
    // Add the domains where your images are hosted.
    // Since your backend serves images from localhost:8080, 'localhost' is needed.
    // If you access via IP (like 192.168.0.199), add that IP as well.
    domains: ['localhost', '192.168.0.199'], // Add any other domains if you use a CDN later
  },
};

module.exports = nextConfig;