import type { JSX } from "react";
import "./App.css";
import { AuthProvider, useAuth } from "./auth/auth";

function RequireAuth({ children }: { children: JSX.Element }) {
  const { user, loading } = useAuth();

  if (loading) {
    return <div>loading...</div>;
  }

  return children;
}

function App() {
  return <AuthProvider></AuthProvider>;
}
export default App;
