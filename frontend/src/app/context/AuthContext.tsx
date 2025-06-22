// app/context/AuthContext.tsx
"use client";

import { createContext, useContext, useState, useEffect, ReactNode } from 'react';

// Define the User interface, including the role and new profile fields
interface User {
  id: number;
  email: string;
  role: 'user' | 'admin'; // Crucial for distinguishing admins
  first_name?: string;
  last_name?: string;
  // NEW PROFILE FIELDS:
  profile_picture_url?: string;
  bio?: string;
  location?: string;
  website?: string;
}

// Define the shape of your AuthContext
interface AuthContextType {
  user: User | null;         // Stores the user object
  isAuthenticated: boolean;  // Derived: true if user object exists AND JWT cookie is present
  isAdmin: boolean;          // Derived: true if user.role is 'admin'
  loading: boolean;          // To indicate if initial load from localStorage/cookie is complete
  login: (userData: User) => void; // Now accepts user data
  logout: () => void;
}

const AuthContext = createContext<AuthContextType | undefined>(undefined);

export function AuthProvider({ children }: { children: ReactNode }) {
  const [user, setUser] = useState<User | null>(null);
  const [loading, setLoading] = useState(true); // Tracks if initial load is done

  useEffect(() => {
    let storedUser: User | null = null;
    try {
      // Try to get user data from localStorage
      const userJson = localStorage.getItem('user');
      if (userJson) {
        storedUser = JSON.parse(userJson);
      }
    } catch (e) {
      console.error("Failed to parse user from localStorage:", e);
      localStorage.removeItem('user'); // Clear corrupted data
    }

    // Check for JWT cookie. The presence of this cookie still indicates login
    const jwtCookie = document.cookie.split('; ').find(row => row.startsWith('jwt='));

    // If both cookie and user data are present, set user state
    if (jwtCookie && storedUser) {
      setUser(storedUser);
    } else {
      // If either is missing, ensure no user is logged in and clean up stale data
      setUser(null);
      localStorage.removeItem('user'); // Clean up potentially stale user data
      // Optionally, clear the cookie if user data is missing but cookie is present (might be an edge case)
      if (jwtCookie) {
         document.cookie = "jwt=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;";
      }
    }
    setLoading(false); // Finished initial load
  }, []); // Run once on component mount

  const login = (userData: User) => {
    // This function is called after a successful login API call.
    // The JWT cookie is assumed to be set by the backend as an HTTP-only cookie.
    setUser(userData);
    localStorage.setItem('user', JSON.stringify(userData)); // Store user data for client-side access
  };

  const logout = () => {
    // Clear user state and remove from localStorage
    setUser(null);
    localStorage.removeItem('user');
    // Clear the JWT cookie explicitly
    document.cookie = "jwt=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;";
  };

  // Derived states for convenience
  const isAuthenticated = !!user; // True if we have a user object
  const isAdmin = user?.role === 'admin'; // True if the user's role is 'admin'

  return (
    <AuthContext.Provider value={{ user, isAuthenticated, isAdmin, loading, login, logout }}>
      {children}
    </AuthContext.Provider>
  );
}

export function useAuth() {
  const context = useContext(AuthContext);
  if (context === undefined) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
}