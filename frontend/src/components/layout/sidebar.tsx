import { Link } from '@tanstack/react-router';
import {
  IconBox,
  IconTruck,
  IconChartBar,
  IconShoppingCart,
  IconDashboard,
} from '@tabler/icons-react';

const navigation = [
  {
    name: 'ダッシュボード',
    href: '/',
    icon: IconDashboard,
  },
  {
    name: '商品管理',
    href: '/products',
    icon: IconBox,
  },
  {
    name: '在庫管理',
    href: '/inventory',
    icon: IconShoppingCart,
  },
  {
    name: '配送管理',
    href: '/deliveries',
    icon: IconTruck,
  },
  {
    name: '売上管理',
    href: '/sales',
    icon: IconChartBar,
  },
];

export function Sidebar() {
  return (
    <div className="flex h-full w-64 flex-col border-r bg-white">
      <div className="flex h-16 items-center border-b px-6">
        <h1 className="text-lg font-bold">スマートストア</h1>
      </div>
      <nav className="flex-1 space-y-1 px-2 py-4">
        {navigation.map((item) => (
          <Link
            key={item.href}
            to={item.href}
            className="group flex items-center rounded-md px-3 py-2 text-sm font-medium text-gray-600 hover:bg-gray-50 hover:text-gray-900"
          >
            <item.icon className="mr-3 h-5 w-5" />
            {item.name}
          </Link>
        ))}
      </nav>
    </div>
  );
} 