// components/Footer.tsx
import Link from 'next/link';

export default function Footer() {
  const currentYear = new Date().getFullYear();

  return (
    <footer className="bg-gray-800 text-white p-6 mt-12">
      <div className="container mx-auto text-center">
        <div className="flex flex-wrap justify-center gap-x-6 gap-y-2 mb-4">
          <Link href="/about" className="hover:text-blue-400 transition-colors">
            About Us
          </Link>
          <Link href="/contact" className="hover:text-blue-400 transition-colors">
            Contact
          </Link>
          <Link href="/privacy-policy" className="hover:text-blue-400 transition-colors">
            Privacy Policy
          </Link>
          <Link href="/terms-of-service" className="hover:text-blue-400 transition-colors">
            Terms of Service
          </Link>
        </div>
        <p className="text-sm text-gray-400">
          © {currentYear} Gin Blog. All rights reserved.
        </p>
        <p className="text-xs text-gray-500 mt-2">
          Built with Next.js and Tailwind CSS
        </p>
      </div>
    </footer>
  );
}