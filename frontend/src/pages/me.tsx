import { useAuth } from "../auth/auth";

export function MePage() {
  const { user } = useAuth();

  if (!user) {
    return <div>no user</div>;
  }

  return (
    <div>
      <h1>me</h1>
      <p>id: {user.id}</p>
      <p>email: {user.email}</p>
      <p>role: {user.role}</p>
      <p>verified: {String(user.email_verified)}</p>
    </div>
  );
}
