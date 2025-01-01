import { Route as RootRoute } from './routes/__root';
import { Route as AuthenticatedRoute } from './routes/_authenticated/route';
import { Route as ProductsRoute } from './routes/_authenticated/products/route';
import { Route as ProductIndexRoute } from './routes/_authenticated/products/index';
import { Route as ProductDetailRoute } from './routes/_authenticated/products/$productId';

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
  }
}
