import { useState } from 'react';
import { Outlet } from 'react-router-dom';
import Navbar from '../../components/common/Navbar';

export default function MainLayout() {
  // Temporary message under the navbar (until real search is connected)
  const [msg, setMsg] = useState('');

  function handleSearch(q: string) {
    setMsg(q ? `You searched for: ${q}` : 'Type something to search.');
  }

  return (
    <>
      <Navbar onSearch={handleSearch} />
      <div id="msg">{msg}</div>
      <Outlet />
    </>
  );
}
