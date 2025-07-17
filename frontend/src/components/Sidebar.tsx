import { Globe, ShoppingCart, List, Settings } from 'lucide-react';
import { Link, useLocation } from 'react-router-dom';
import clsx from 'clsx';

const navItems = [
  { name: "Check Domain", path: "/check", icon: <Globe /> },
  { name: "Buy Domain", path: "/buy", icon: <ShoppingCart /> },
  { name: "Purchased Domains", path: "/domains", icon: <List /> },
  { name: "Settings", path: "/settings", icon: <Settings /> },
];

export default function Sidebar() {
  const location = useLocation();

  return (
    <aside className="w-64 h-screen bg-gray-900 text-white p-4 space-y-4">
      <h1 className="text-2xl font-bold mb-6">📦 Domain Panel</h1>
      {navItems.map(item => (
        <Link
          key={item.name}
          to={item.path}
          className={clsx(
            "flex items-center gap-3 p-2 rounded hover:bg-gray-700 transition",
            location.pathname === item.path && "bg-gray-700"
          )}
        >
          {item.icon}
          {item.name}
        </Link>
      ))}
    </aside>
  );
}
