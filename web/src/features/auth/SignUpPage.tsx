import { FormEvent, useState } from 'react';
import { Link } from 'react-router-dom';
import signupImage from '../../assets/signup-image.jpg';
import { EMAIL_RE } from '../../lib/format';
import { EyeIcon, FacebookIcon, GoogleIcon } from '../../components/common/icons';

type Errors = {
  first?: string;
  last?: string;
  email?: string;
  password?: string;
  terms?: string;
};

export default function SignUpPage() {
  const [firstName, setFirstName] = useState('');
  const [lastName, setLastName] = useState('');
  const [email, setEmail] = useState('');
  const [password, setPassword] = useState('');
  const [showPassword, setShowPassword] = useState(false);
  const [terms, setTerms] = useState(false);
  const [errors, setErrors] = useState<Errors>({});
  const [note, setNote] = useState('');

  function handleSubmit(e: FormEvent) {
    e.preventDefault();

    const firstOk = firstName.trim() !== '';
    const lastOk = lastName.trim() !== '';
    const emailOk = EMAIL_RE.test(email.trim());
    const passOk = password.length >= 6;

    setErrors({
      first: firstOk ? '' : 'Required.',
      last: lastOk ? '' : 'Required.',
      email: emailOk ? '' : 'Enter a valid email address.',
      password: passOk ? '' : 'At least 6 characters.',
      terms: terms ? '' : 'Please accept the Terms & Conditions.',
    });

    // Demo only: there is no server yet, so nothing is really saved
    setNote(
      firstOk && lastOk && emailOk && passOk && terms
        ? `Account created (demo) for ${firstName.trim()}`
        : '',
    );
  }

  return (
    <div className="auth-page">
      <main className="signup-card">
        <div className="su-image">
          <img src={signupImage} alt="Wasla delivery" />
          <Link className="su-back" to="/">Back to website &rarr;</Link>
        </div>

        <section className="su-form">
          <h1>Create an account</h1>
          <p className="su-sub">Already have an account? <Link to="/signin">Sign in</Link></p>

          <form onSubmit={handleSubmit} noValidate>
            <div className="su-row">
              <div className="field">
                <input
                  type="text"
                  placeholder="First name"
                  aria-label="First name"
                  autoComplete="given-name"
                  className={errors.first ? 'invalid' : ''}
                  value={firstName}
                  onChange={(e) => setFirstName(e.target.value)}
                />
                <span className="error">{errors.first}</span>
              </div>
              <div className="field">
                <input
                  type="text"
                  placeholder="Last name"
                  aria-label="Last name"
                  autoComplete="family-name"
                  className={errors.last ? 'invalid' : ''}
                  value={lastName}
                  onChange={(e) => setLastName(e.target.value)}
                />
                <span className="error">{errors.last}</span>
              </div>
            </div>
            <div className="field">
              <input
                type="email"
                placeholder="Email"
                aria-label="Email"
                autoComplete="email"
                className={errors.email ? 'invalid' : ''}
                value={email}
                onChange={(e) => setEmail(e.target.value)}
              />
              <span className="error">{errors.email}</span>
            </div>
            <div className="field">
              <div className="pw-wrap">
                <input
                  type={showPassword ? 'text' : 'password'}
                  placeholder="Enter your password"
                  aria-label="Password"
                  autoComplete="new-password"
                  className={errors.password ? 'invalid' : ''}
                  value={password}
                  onChange={(e) => setPassword(e.target.value)}
                />
                <button
                  type="button"
                  className="eye"
                  aria-label={showPassword ? 'Hide password' : 'Show password'}
                  onClick={() => setShowPassword((s) => !s)}
                >
                  <EyeIcon />
                </button>
              </div>
              <span className="error">{errors.password}</span>
            </div>

            <label className="terms">
              <input type="checkbox" checked={terms} onChange={(e) => setTerms(e.target.checked)} />
              I agree to the <a href="#">Terms &amp; Conditions</a>
            </label>
            <span className="error">{errors.terms}</span>

            <button type="submit" className="btn btn-solid btn-full">Create account</button>
          </form>

          <div className="divider">Or register with</div>
          <div className="social">
            <button type="button" className="btn-social" onClick={() => setNote('Sign in with Google (demo).')}>
              <GoogleIcon /> Google
            </button>
            <button type="button" className="btn-social" onClick={() => setNote('Sign in with Facebook (demo).')}>
              <FacebookIcon /> Facebook
            </button>
          </div>
          <p className="auth-note">{note}</p>
        </section>
      </main>
    </div>
  );
}
