import { FormEvent, useState } from 'react';
import { Link } from 'react-router-dom';
import { EMAIL_RE } from '../../lib/format';

export default function SignInPage() {
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [emailError, setEmailError] = useState('');
  const [passwordError, setPasswordError] = useState('');
  const [note, setNote] = useState('');

  function handleSubmit(e: FormEvent) {
    e.preventDefault();
    const emailOk = EMAIL_RE.test(email.trim());
    const passOk = password.length >= 6;
    setEmailError(emailOk ? '' : 'Enter a valid email address.');
    setPasswordError(passOk ? '' : 'Password must be at least 6 characters.');
    // Demo only: there is no server yet, so nothing is really checked
    setNote(emailOk && passOk ? `Signed in (demo) as ${email.trim()}` : '');
  }

  return (
    <div className="auth-page">
      <main className="auth-card">
        <Link className="auth-back" to="/">&larr; Back to home</Link>
        <h1>Sign in</h1>
        <p className="auth-sub">Welcome back to Wasla</p>

        <form onSubmit={handleSubmit} noValidate>
          <div className="field">
            <label htmlFor="email">Email</label>
            <input
              type="email"
              id="email"
              placeholder="you@example.com"
              autoComplete="email"
              className={emailError ? 'invalid' : ''}
              value={email}
              onChange={(e) => setEmail(e.target.value)}
            />
            <span className="error">{emailError}</span>
          </div>
          <div className="field">
            <label htmlFor="password">Password</label>
            <input
              type="password"
              id="password"
              placeholder="Your password"
              autoComplete="current-password"
              className={passwordError ? 'invalid' : ''}
              value={password}
              onChange={(e) => setPassword(e.target.value)}
            />
            <span className="error">{passwordError}</span>
          </div>
          <button type="submit" className="btn btn-solid btn-full">Sign in</button>
        </form>
        <p className="auth-alt">Don&apos;t have an account? <Link to="/signup">Create one</Link></p>
        <p className="auth-note">{note}</p>
      </main>
    </div>
  );
}
