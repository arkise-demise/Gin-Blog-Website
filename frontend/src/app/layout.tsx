// app/layout.tsx
import type { Metadata } from "next";
import "./globals.css";
import Navbar from "../components/Navbar";
import Footer from "../components/Footer"; 
import { AuthProvider } from './context/AuthContext';

export const metadata: Metadata = {
  title: "Gin Blog",
  description: "A blog built with Next.js and GoLang",
};

export default function RootLayout({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) {
  return (
    <html lang="en">
      <body>
        <AuthProvider>
          <Navbar />
          <main className="flex-grow">
            {children} 
          </main>
          <Footer />
        </AuthProvider>
      </body>
    </html>
  );
}