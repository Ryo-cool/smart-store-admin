import { createFileRoute, Outlet } from '@tanstack/react-router';

export const Route = createFileRoute('/_authenticated/products')({
  component: ProductsLayout,
});

function ProductsLayout() {
  return (
    <div className="container py-6">
      <Outlet />
    </div>
  );
} 