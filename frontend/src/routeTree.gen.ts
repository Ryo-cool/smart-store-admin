import { Route as RootRoute } from './routes/__root';
import { Route as AuthenticatedRoute } from './routes/_authenticated/route';
import { Route as ProductsRoute } from './routes/_authenticated/products/route';
import { Route as DeliveriesRoute } from './routes/_authenticated/deliveries/index';
import { Route as InventoryRoute } from './routes/_authenticated/inventory/index';

declare module '@tanstack/react-router' {
  interface FileRoutesByPath {
    '/': {
      parentRoute: typeof RootRoute;
    };
    '/_authenticated': {
      parentRoute: typeof RootRoute;
    };
    '/_authenticated/products': {
      parentRoute: typeof AuthenticatedRoute;
    };
    '/_authenticated/products/': {
      parentRoute: typeof ProductsRoute;
    };
    '/_authenticated/products/$productId': {
      parentRoute: typeof ProductsRoute;
      params: {
        productId: string;
      };
    };
    '/_authenticated/inventory': {
      parentRoute: typeof AuthenticatedRoute;
      route: typeof InventoryRoute;
    };
    '/_authenticated/deliveries': {
      parentRoute: typeof AuthenticatedRoute;
      route: typeof DeliveriesRoute;
    };
    '/_authenticated/deliveries/$deliveryId': {
      parentRoute: typeof AuthenticatedRoute;
      params: {
        deliveryId: string;
      };
    };
  }
}
