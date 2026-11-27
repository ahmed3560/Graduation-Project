import { FormEvent, useState } from 'react';
import { Link, NavLink } from 'react-router-dom';
import logo from '../../assets/logo.png';
import { useCart } from '../../hooks/useCart';
import { CartIcon, SearchIcon } from './icons';

type Props = {
  onSearch: (query: string) => void;
};

export default function Navbar({ onSearch }: Props) {
  const [q, setQ] = useState('');
  const { count } = useCart();

  function handleSubmit(e: FormEvent) {
    e.preventDefault();
    onSearch(q.trim());
  }

  return (
    <nav className="navbar">
      <Link className="logo" to="/" aria-label="Wasla home">
        <img src={logo} alt="Wasla" />
      </Link>
      <ul className="links">
        <li>
          <NavLink to="/" end className={({ isActive }) => (isActive ? 'active' : '')}>
            Home
          </NavLink>
        </li>
        <li><a href="/#about">About</a></li>
        <li><a href="/#contact">Contact</a></li>
      </ul>
      <form className="search" role="search" onSubmit={handleSubmit}>
        <input
          type="search"
          placeholder="Search..."
          aria-label="Search"
          value={q}
          onChange={(e) => setQ(e.target.value)}
        />
        <button type="submit" aria-label="Search">
          <SearchIcon />
        </button>
      </form>
      <div className="actions">
        <Link className="btn btn-ghost" to="/signin">Sign in</Link>
        <Link className="btn btn-solid" to="/signup">Create account</Link>
        <Link className="cart" to="/cart" aria-label="Shopping cart">
          <CartIcon />
          <span className="badge" id="badge">{count}</span>
        </Link>
      </div>
    </nav>
  );
}
